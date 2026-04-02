package identity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
)

func toCreateUserDomain(req CreateUserRequest, actorID string) (domain.UserCreateInput, error) {
	roleIDs, err := parseUUIDList(req.RoleIDs, "role_ids")
	if err != nil {
		return domain.UserCreateInput{}, err
	}

	avatarID, err := parseOptionalUUID(req.AvatarID, "avatar_id")
	if err != nil {
		return domain.UserCreateInput{}, err
	}

	return domain.UserCreateInput{
		Email:       req.Email,
		Username:    req.Username,
		DisplayName: req.DisplayName,
		AvatarID:    avatarID,
		Status:      domain.UserStatus(req.Status),
		RoleIDs:     roleIDs,
		ActorID:     actorID,
	}, nil
}

func toUpdateUserDomain(id uuid.UUID, req UpdateUserRequest, actorID string) (domain.UserUpdateInput, error) {
	input := domain.UserUpdateInput{
		ID:      id,
		ActorID: actorID,
	}

	if req.Email != nil {
		input.Email = domain.NewEntityField(strings.TrimSpace(*req.Email))
	}
	if req.Username != nil {
		input.Username = domain.NewEntityField(strings.TrimSpace(*req.Username))
	}
	if req.DisplayName != nil {
		input.DisplayName = domain.NewEntityField(strings.TrimSpace(*req.DisplayName))
	}
	if req.Status != nil {
		input.Status = domain.NewEntityField(domain.UserStatus(strings.TrimSpace(*req.Status)))
	}
	if req.AvatarID != nil {
		avatarID, err := parseOptionalUUID(req.AvatarID, "avatar_id")
		if err != nil {
			return domain.UserUpdateInput{}, err
		}
		input.AvatarID = domain.NewEntityField(avatarID)
	}
	if req.RoleIDs != nil {
		roleIDs, err := parseUUIDList(req.RoleIDs, "role_ids")
		if err != nil {
			return domain.UserUpdateInput{}, err
		}
		input.RoleIDs = domain.NewEntityField(roleIDs)
	}

	return input, nil
}

func toUserResponse(user domain.User) UserResponse {
	response := UserResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Status:      string(user.Status),
		Roles:       make([]RoleResponse, 0, len(user.Roles)),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		CreatedBy:   user.CreatedBy,
		UpdatedBy:   user.UpdatedBy,
	}

	if user.AvatarID != nil {
		value := user.AvatarID.String()
		response.AvatarID = &value
	}

	for _, role := range user.Roles {
		response.Roles = append(response.Roles, toRoleResponse(role))
	}

	if len(user.Identities) > 0 {
		response.Identities = make([]AuthIdentityResponse, 0, len(user.Identities))
		for _, identity := range user.Identities {
			response.Identities = append(response.Identities, toAuthIdentityResponse(identity))
		}
	}

	return response
}

func toUsersListResponse(result domain.PageResult[domain.User], query domain.UserQuery) UsersListResponse {
	items := make([]UserResponse, 0, len(result.Items))
	for _, user := range result.Items {
		items = append(items, toUserResponse(user))
	}

	payload := UserQueryPayload{
		Search:      query.Search,
		SortBy:      query.SortBy,
		SortOrder:   query.SortOrder,
		IncludeAuth: query.IncludeAuth,
	}
	if query.RoleID != nil {
		payload.RoleID = query.RoleID.String()
	}
	if query.Status != nil {
		payload.Status = string(*query.Status)
	}

	return UsersListResponse{
		Items: items,
		Paging: PagingResponse{
			Page:  result.Paging.Page,
			Limit: result.Paging.Limit,
			Total: result.Paging.Total,
		},
		Query: payload,
	}
}

func toCreateRoleDomain(req CreateRoleRequest, actorID string) domain.RoleCreateInput {
	return domain.RoleCreateInput{
		Name:        domain.RoleName(req.Name),
		Description: req.Description,
		Permissions: toRolePermissionInputs(req.Permissions),
		ActorID:     actorID,
	}
}

func toUpdateRoleDomain(id uuid.UUID, req UpdateRoleRequest, actorID string) domain.RoleUpdateInput {
	input := domain.RoleUpdateInput{
		ID:      id,
		ActorID: actorID,
	}

	if req.Name != nil {
		input.Name = domain.NewEntityField(domain.RoleName(strings.TrimSpace(*req.Name)))
	}
	if req.Description != nil {
		input.Description = domain.NewEntityField(strings.TrimSpace(*req.Description))
	}
	if req.Permissions != nil {
		input.Permissions = domain.NewEntityField(toRolePermissionInputs(req.Permissions))
	}

	return input
}

func toRolePermissionInputs(values []RolePermissionRequest) []domain.RolePermissionInput {
	items := make([]domain.RolePermissionInput, 0, len(values))
	for _, value := range values {
		items = append(items, domain.RolePermissionInput{
			PermissionCode: domain.PermissionCode(value.PermissionCode),
			Scope:          strings.TrimSpace(value.Scope),
		})
	}
	return items
}

func toRoleResponse(role domain.Role) RoleResponse {
	response := RoleResponse{
		ID:          role.ID.String(),
		Name:        string(role.Name),
		Description: role.Description,
		Permissions: make([]RolePermissionResponse, 0, len(role.Permissions)),
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   role.UpdatedAt.Format(time.RFC3339),
		CreatedBy:   role.CreatedBy,
		UpdatedBy:   role.UpdatedBy,
	}
	for _, permission := range role.Permissions {
		response.Permissions = append(response.Permissions, RolePermissionResponse{
			PermissionCode: string(permission.PermissionCode),
			Scope:          permission.Scope,
		})
	}
	return response
}

func toRolesListResponse(result domain.PageResult[domain.Role], query domain.RoleQuery) RolesListResponse {
	items := make([]RoleResponse, 0, len(result.Items))
	for _, role := range result.Items {
		items = append(items, toRoleResponse(role))
	}

	payload := RoleQueryPayload{
		Search:    query.Search,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
	}
	if query.Permission != nil {
		payload.Permission = string(*query.Permission)
	}

	return RolesListResponse{
		Items: items,
		Paging: PagingResponse{
			Page:  result.Paging.Page,
			Limit: result.Paging.Limit,
			Total: result.Paging.Total,
		},
		Query: payload,
	}
}

func toCreateAuthIdentityDomain(req CreateAuthIdentityRequest, actorID string) (domain.AuthIdentityCreateInput, error) {
	userID, err := uuid.Parse(strings.TrimSpace(req.UserID))
	if err != nil {
		return domain.AuthIdentityCreateInput{}, helper.BadRequest("invalid_user_id", "user_id must be a valid uuid", err)
	}

	return domain.AuthIdentityCreateInput{
		UserID:         userID,
		Provider:       domain.AuthProvider(req.Provider),
		ProviderUserID: strings.TrimSpace(req.ProviderUserID),
		Password:       req.Password,
		ActorID:        actorID,
	}, nil
}

func toUpdateAuthIdentityDomain(id uuid.UUID, req UpdateAuthIdentityRequest, actorID string) domain.AuthIdentityUpdateInput {
	input := domain.AuthIdentityUpdateInput{
		ID:      id,
		ActorID: actorID,
	}
	if req.ProviderUserID != nil {
		input.ProviderUserID = domain.NewEntityField(strings.TrimSpace(*req.ProviderUserID))
	}
	if req.Password != nil {
		input.Password = domain.NewEntityField(*req.Password)
	}
	return input
}

func toAuthIdentityResponse(identity domain.AuthIdentity) AuthIdentityResponse {
	return AuthIdentityResponse{
		ID:             identity.ID.String(),
		UserID:         identity.UserID.String(),
		Provider:       string(identity.Provider),
		ProviderUserID: identity.ProviderUserID,
		HasPassword:    identity.PasswordHash != "",
		CreatedAt:      identity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      identity.UpdatedAt.Format(time.RFC3339),
		CreatedBy:      identity.CreatedBy,
		UpdatedBy:      identity.UpdatedBy,
	}
}

func toAuthIdentitiesListResponse(result domain.PageResult[domain.AuthIdentity], query domain.AuthIdentityQuery) AuthIdentitiesListResponse {
	items := make([]AuthIdentityResponse, 0, len(result.Items))
	for _, identity := range result.Items {
		items = append(items, toAuthIdentityResponse(identity))
	}
	payload := AuthIdentityQueryPayload{
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
	}
	if query.UserID != nil {
		payload.UserID = query.UserID.String()
	}
	if query.Provider != nil {
		payload.Provider = string(*query.Provider)
	}
	return AuthIdentitiesListResponse{
		Items: items,
		Paging: PagingResponse{
			Page:  result.Paging.Page,
			Limit: result.Paging.Limit,
			Total: result.Paging.Total,
		},
		Query: payload,
	}
}

func toLoginResponse(token domain.AuthToken) LoginResponse {
	return LoginResponse{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		TokenType:        string(token.TokenType),
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
		User:             toUserResponse(token.User),
	}
}

func parseOptionalUUID(value *string, field string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(trimmed)
	if err != nil {
		return nil, helper.BadRequest("invalid_"+field, field+" must be a valid uuid", err)
	}
	return &parsed, nil
}

func parseUUIDList(values []string, field string) ([]uuid.UUID, error) {
	items := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		parsed, err := uuid.Parse(trimmed)
		if err != nil {
			return nil, helper.BadRequest("invalid_"+field, field+" must contain valid uuid values", err)
		}
		items = append(items, parsed)
	}
	return items, nil
}
