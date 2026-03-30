package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type UserServiceDependencies struct {
	UserRepository port.UserRepository
	TxManager      port.TransactionManager
	Storage        port.FileStorage
}

type userService struct {
	userRepository port.UserRepository
	txManager      port.TransactionManager
	storage        port.FileStorage
}

func NewUserService(deps UserServiceDependencies) port.UserService {
	return &userService{
		userRepository: deps.UserRepository,
		txManager:      deps.TxManager,
		storage:        deps.Storage,
	}
}

func (s *userService) Create(ctx context.Context, input domain.UserCreateInput) (domain.User, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Code = strings.TrimSpace(strings.ToUpper(input.Code))
	input.Name = strings.TrimSpace(input.Name)

	if input.Email == "" || input.Code == "" || input.Name == "" {
		return domain.User{}, helper.BadRequest("invalid_user", "name, email and code are required", nil)
	}

	var created domain.User
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.ensureRole(txCtx, input.RoleID); err != nil {
			return err
		}

		if err := s.ensureUniqueEmail(txCtx, input.Email, nil); err != nil {
			return err
		}

		if err := s.ensureUniqueCode(txCtx, input.Code, nil); err != nil {
			return err
		}

		user, err := s.userRepository.Create(txCtx, domain.User{
			Code:   input.Code,
			Name:   input.Name,
			Email:  input.Email,
			Status: input.Status,
			RoleID: input.RoleID,
			AuditFields: domain.AuditFields{
				CreatedBy: input.ActorID,
				UpdatedBy: input.ActorID,
			},
		})
		if err != nil {
			return err
		}

		created = user
		return nil
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
			email := strings.TrimSpace(strings.ToLower(input.Email.Value))
			if email == "" {
				return helper.BadRequest("invalid_email", "email cannot be empty", nil)
			}

			input.Email.Value = email
			if !strings.EqualFold(current.Email, email) {
				if err := s.ensureUniqueEmail(txCtx, email, &input.ID); err != nil {
					return err
				}
			}
		}

		if input.Code.Set {
			code := strings.TrimSpace(strings.ToUpper(input.Code.Value))
			if code == "" {
				return helper.BadRequest("invalid_code", "code cannot be empty", nil)
			}

			input.Code.Value = code
			if !strings.EqualFold(current.Code, code) {
				if err := s.ensureUniqueCode(txCtx, code, &input.ID); err != nil {
					return err
				}
			}
		}

		if input.Name.Set {
			name := strings.TrimSpace(input.Name.Value)
			if name == "" {
				return helper.BadRequest("invalid_name", "name cannot be empty", nil)
			}

			input.Name.Value = name
		}

		if input.RoleID.Set && input.RoleID.Value != nil {
			if err := s.ensureRole(txCtx, input.RoleID.Value); err != nil {
				return err
			}
		}

		updated, err = s.userRepository.Update(txCtx, input)
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
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return domain.User{}, helper.BadRequest("invalid_email", "email is required", nil)
	}

	return s.userRepository.FindByEmail(ctx, email)
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

func (s *userService) ensureUniqueCode(ctx context.Context, code string, excludeID *uuid.UUID) error {
	exists, err := s.userRepository.ExistsByCode(ctx, code, excludeID)
	if err != nil {
		return err
	}

	if exists {
		return helper.Conflict("duplicate_code", "code already exists", nil)
	}

	return nil
}

func (s *userService) ensureRole(ctx context.Context, roleID *uuid.UUID) error {
	if roleID == nil {
		return nil
	}

	if _, err := s.userRepository.FindRoleByID(ctx, *roleID); err != nil {
		var appErr *helper.AppError
		if errors.As(err, &appErr) {
			return err
		}

		return helper.BadRequest("invalid_role", "role does not exist", err)
	}

	return nil
}
