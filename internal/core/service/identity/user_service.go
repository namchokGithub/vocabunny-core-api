package identity

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type UserServiceDependencies struct {
	UserRepository port.UserRepository
	RoleRepository port.RoleRepository
	TxManager      port.TransactionManager
}

type userService struct {
	userRepository port.UserRepository
	roleRepository port.RoleRepository
	txManager      port.TransactionManager
}

func NewUserService(deps UserServiceDependencies) port.UserService {
	return &userService{
		userRepository: deps.UserRepository,
		roleRepository: deps.RoleRepository,
		txManager:      deps.TxManager,
	}
}

func (s *userService) Create(ctx context.Context, input domain.UserCreateInput) (domain.User, error) {
	input.Email = normalizeEmail(input.Email)
	input.Username = normalizeUsername(input.Username)
	input.DisplayName = strings.TrimSpace(input.DisplayName)

	if input.Email == "" || input.Username == "" || input.DisplayName == "" {
		return domain.User{}, helper.BadRequest("invalid_user", "email, username and display_name are required", nil)
	}

	if !isValidUserStatus(input.Status) {
		return domain.User{}, helper.BadRequest("invalid_status", "status is invalid", nil)
	}

	var created domain.User
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.ensureRolesExist(txCtx, input.RoleIDs); err != nil {
			return err
		}

		if err := s.ensureUniqueEmail(txCtx, input.Email, nil); err != nil {
			return err
		}

		if err := s.ensureUniqueUsername(txCtx, input.Username, nil); err != nil {
			return err
		}

		user, err := s.userRepository.Create(txCtx, domain.User{
			Email:       input.Email,
			Username:    input.Username,
			DisplayName: input.DisplayName,
			AvatarID:    input.AvatarID,
			Status:      input.Status,
			AuditFields: domain.AuditFields{
				CreatedBy: input.ActorID,
				UpdatedBy: input.ActorID,
			},
		})
		if err != nil {
			return err
		}

		if err := s.userRepository.ReplaceRoles(txCtx, user.ID, input.RoleIDs, input.ActorID); err != nil {
			return err
		}

		created, err = s.userRepository.FindByID(txCtx, user.ID)
		return err
	})
	if err != nil {
		return domain.User{}, err
	}

	return created, nil
}

func (s *userService) Update(ctx context.Context, input domain.UserUpdateInput) (domain.User, error) {
	var updated domain.User
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		current, err := s.userRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}

		if input.Email.Set {
			input.Email.Value = normalizeEmail(input.Email.Value)
			if input.Email.Value == "" {
				return helper.BadRequest("invalid_email", "email cannot be empty", nil)
			}
			if !strings.EqualFold(current.Email, input.Email.Value) {
				if err := s.ensureUniqueEmail(txCtx, input.Email.Value, &input.ID); err != nil {
					return err
				}
			}
		}

		if input.Username.Set {
			input.Username.Value = normalizeUsername(input.Username.Value)
			if input.Username.Value == "" {
				return helper.BadRequest("invalid_username", "username cannot be empty", nil)
			}
			if !strings.EqualFold(current.Username, input.Username.Value) {
				if err := s.ensureUniqueUsername(txCtx, input.Username.Value, &input.ID); err != nil {
					return err
				}
			}
		}

		if input.DisplayName.Set {
			input.DisplayName.Value = strings.TrimSpace(input.DisplayName.Value)
			if input.DisplayName.Value == "" {
				return helper.BadRequest("invalid_display_name", "display_name cannot be empty", nil)
			}
		}

		if input.Status.Set && !isValidUserStatus(input.Status.Value) {
			return helper.BadRequest("invalid_status", "status is invalid", nil)
		}

		if input.RoleIDs.Set {
			if err := s.ensureRolesExist(txCtx, input.RoleIDs.Value); err != nil {
				return err
			}
		}

		if _, err := s.userRepository.Update(txCtx, input); err != nil {
			return err
		}

		if input.RoleIDs.Set {
			if err := s.userRepository.ReplaceRoles(txCtx, input.ID, input.RoleIDs.Value, input.ActorID); err != nil {
				return err
			}
		}

		updated, err = s.userRepository.FindByID(txCtx, input.ID)
		return err
	})
	if err != nil {
		return domain.User{}, err
	}

	return updated, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.userRepository.FindByID(txCtx, id); err != nil {
			return err
		}

		return s.userRepository.Delete(txCtx, id, actorID)
	})
}

func (s *userService) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.userRepository.FindByID(ctx, id)
}

func (s *userService) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	email = normalizeEmail(email)
	if email == "" {
		return domain.User{}, helper.BadRequest("invalid_email", "email is required", nil)
	}

	return s.userRepository.FindByEmail(ctx, email)
}

func (s *userService) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	username = normalizeUsername(username)
	if username == "" {
		return domain.User{}, helper.BadRequest("invalid_username", "username is required", nil)
	}

	return s.userRepository.FindByUsername(ctx, username)
}

func (s *userService) FindBySubject(ctx context.Context, subject string) (domain.User, error) {
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return domain.User{}, helper.Unauthorized("invalid_subject", "jwt subject is required", nil)
	}

	return s.userRepository.FindBySubject(ctx, subject)
}

func (s *userService) FindAll(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error) {
	return s.userRepository.FindAll(ctx, query)
}

func (s *userService) ensureUniqueEmail(ctx context.Context, email string, excludeID *uuid.UUID) error {
	exists, err := s.userRepository.ExistsByEmail(ctx, email, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return helper.Conflict("duplicate_email", "email already exists", nil)
	}
	return nil
}

func (s *userService) ensureUniqueUsername(ctx context.Context, username string, excludeID *uuid.UUID) error {
	exists, err := s.userRepository.ExistsByUsername(ctx, username, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return helper.Conflict("duplicate_username", "username already exists", nil)
	}
	return nil
}

func (s *userService) ensureRolesExist(ctx context.Context, roleIDs []uuid.UUID) error {
	for _, roleID := range roleIDs {
		if _, err := s.roleRepository.FindByID(ctx, roleID); err != nil {
			return helper.BadRequest("invalid_role_id", "one or more roles do not exist", err)
		}
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func normalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func isValidUserStatus(status domain.UserStatus) bool {
	switch status {
	case domain.UserStatusInactive, domain.UserStatusActive, domain.UserStatusBanned, domain.UserStatusDeleted:
		return true
	default:
		return false
	}
}
