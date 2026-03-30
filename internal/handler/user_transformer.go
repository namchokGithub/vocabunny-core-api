package handler

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
)

func toCreateUserDomain(req CreateUserRequest, actorID string) (domain.UserCreateInput, error) {
	var roleID *uuid.UUID
	if req.RoleID != nil && strings.TrimSpace(*req.RoleID) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(*req.RoleID))
		if err != nil {
			return domain.UserCreateInput{}, helper.BadRequest("invalid_role_id", "role_id must be a valid uuid", err)
		}

		roleID = &parsed
	}

	return domain.UserCreateInput{
		Code:    req.Code,
		Name:    req.Name,
		Email:   req.Email,
		Status:  domain.UserStatus(req.Status),
		RoleID:  roleID,
		ActorID: actorID,
	}, nil
}

func toUpdateUserDomain(id uuid.UUID, req UpdateUserRequest, actorID string) (domain.UserUpdateInput, error) {
	input := domain.UserUpdateInput{
		ID:      id,
		ActorID: actorID,
	}

	if req.Name != nil {
		input.Name = domain.NewEntityField(strings.TrimSpace(*req.Name))
	}

	if req.Email != nil {
		input.Email = domain.NewEntityField(strings.TrimSpace(*req.Email))
	}

	if req.Code != nil {
		input.Code = domain.NewEntityField(strings.TrimSpace(*req.Code))
	}

	if req.Status != nil {
		input.Status = domain.NewEntityField(domain.UserStatus(strings.TrimSpace(*req.Status)))
	}

	if req.RoleID != nil {
		if strings.TrimSpace(*req.RoleID) == "" {
			input.RoleID = domain.NewEntityField[*uuid.UUID](nil)
			return input, nil
		}

		parsed, err := uuid.Parse(strings.TrimSpace(*req.RoleID))
		if err != nil {
			return domain.UserUpdateInput{}, helper.BadRequest("invalid_role_id", "role_id must be a valid uuid", err)
		}

		input.RoleID = domain.NewEntityField(&parsed)
	}

	return input, nil
}

func toUserResponse(user domain.User) UserResponse {
	response := UserResponse{
		ID:        user.ID.String(),
		Code:      user.Code,
		Name:      user.Name,
		Email:     user.Email,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.UpdatedBy,
	}

	if user.Role != nil {
		response.Role = &RoleResponse{
			ID:   user.Role.ID.String(),
			Code: user.Role.Code,
			Name: user.Role.Name,
		}
	}

	return response
}

func toUsersListResponse(result domain.PageResult[domain.User], query domain.UserQuery) UsersListResponse {
	items := make([]UserResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toUserResponse(item))
	}

	queryPayload := UserQueryPayload{
		Search: query.Search,
		SortBy: query.SortBy,
		Order:  query.SortOrder,
	}

	if query.RoleID != nil {
		queryPayload.RoleID = query.RoleID.String()
	}

	if query.Status != nil {
		queryPayload.Status = string(*query.Status)
	}

	return UsersListResponse{
		Items: items,
		Paging: PagingResponse{
			Page:  result.Paging.Page,
			Limit: result.Paging.Limit,
			Total: result.Paging.Total,
		},
		Query: queryPayload,
	}
}
