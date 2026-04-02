package identity

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type RoleServiceDependencies struct {
	RoleRepository port.RoleRepository
	TxManager      port.TransactionManager
}

type roleService struct {
	roleRepository port.RoleRepository
	txManager      port.TransactionManager
}

func NewRoleService(deps RoleServiceDependencies) port.RoleService {
	return &roleService{
		roleRepository: deps.RoleRepository,
		txManager:      deps.TxManager,
	}
}

func (s *roleService) Create(ctx context.Context, input domain.RoleCreateInput) (domain.Role, error) {
	input.Description = strings.TrimSpace(input.Description)
	if !isValidRoleName(input.Name) {
		return domain.Role{}, helper.BadRequest("invalid_role_name", "role name is invalid", nil)
	}

	var created domain.Role
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		exists, err := s.roleRepository.ExistsByName(txCtx, input.Name, nil)
		if err != nil {
			return err
		}
		if exists {
			return helper.Conflict("duplicate_role_name", "role name already exists", nil)
		}

		role, err := s.roleRepository.Create(txCtx, domain.Role{
			Name:        input.Name,
			Description: input.Description,
			AuditFields: domain.AuditFields{
				CreatedBy: input.ActorID,
				UpdatedBy: input.ActorID,
			},
		})
		if err != nil {
			return err
		}

		if err := s.roleRepository.ReplacePermissions(txCtx, role.ID, input.Permissions, input.ActorID); err != nil {
			return err
		}

		created, err = s.roleRepository.FindByID(txCtx, role.ID)
		return err
	})
	if err != nil {
		return domain.Role{}, err
	}

	return created, nil
}

func (s *roleService) Update(ctx context.Context, input domain.RoleUpdateInput) (domain.Role, error) {
	var updated domain.Role
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.roleRepository.FindByID(txCtx, input.ID); err != nil {
			return err
		}

		if input.Name.Set {
			if !isValidRoleName(input.Name.Value) {
				return helper.BadRequest("invalid_role_name", "role name is invalid", nil)
			}
			exists, err := s.roleRepository.ExistsByName(txCtx, input.Name.Value, &input.ID)
			if err != nil {
				return err
			}
			if exists {
				return helper.Conflict("duplicate_role_name", "role name already exists", nil)
			}
		}

		if input.Description.Set {
			input.Description.Value = strings.TrimSpace(input.Description.Value)
		}

		if _, err := s.roleRepository.Update(txCtx, input); err != nil {
			return err
		}

		if input.Permissions.Set {
			if err := s.roleRepository.ReplacePermissions(txCtx, input.ID, input.Permissions.Value, input.ActorID); err != nil {
				return err
			}
		}

		role, err := s.roleRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}
		updated = role
		return nil
	})
	if err != nil {
		return domain.Role{}, err
	}

	return updated, nil
}

func (s *roleService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.roleRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.roleRepository.Delete(txCtx, id, actorID)
	})
}

func (s *roleService) FindByID(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	return s.roleRepository.FindByID(ctx, id)
}

func (s *roleService) FindAll(ctx context.Context, query domain.RoleQuery) (domain.PageResult[domain.Role], error) {
	return s.roleRepository.FindAll(ctx, query)
}

func isValidRoleName(name domain.RoleName) bool {
	switch name {
	case domain.RoleNameAdmin, domain.RoleNameContentAdmin, domain.RoleNameModerator, domain.RoleNameUser:
		return true
	default:
		return false
	}
}
