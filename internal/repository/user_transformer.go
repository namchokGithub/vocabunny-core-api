package repository

import (
	"time"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func toDomainRole(model RoleModel) domain.Role {
	role := domain.Role{
		ID:   model.ID,
		Code: model.Code,
		Name: model.Name,
		AuditFields: domain.AuditFields{
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
		},
	}

	if model.DeletedAt.Valid {
		deletedAt := model.DeletedAt.Time
		role.DeletedAt = &deletedAt
	}

	return role
}

func toDomainUser(model UserModel) domain.User {
	user := domain.User{
		ID:     model.ID,
		Code:   model.Code,
		Name:   model.Name,
		Email:  model.Email,
		Status: domain.UserStatus(model.Status),
		RoleID: model.RoleID,
		AuditFields: domain.AuditFields{
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
		},
	}

	if model.DeletedAt.Valid {
		deletedAt := model.DeletedAt.Time
		user.DeletedAt = &deletedAt
	}

	if model.Role != nil {
		role := toDomainRole(*model.Role)
		user.Role = &role
	}

	return user
}

func applyUserPatch(model *UserModel, input domain.UserUpdateInput) {
	if input.Name.Set {
		model.Name = input.Name.Value
	}

	if input.Email.Set {
		model.Email = input.Email.Value
	}

	if input.Code.Set {
		model.Code = input.Code.Value
	}

	if input.Status.Set {
		model.Status = string(input.Status.Value)
	}

	if input.RoleID.Set {
		model.RoleID = input.RoleID.Value
	}

	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()
}
