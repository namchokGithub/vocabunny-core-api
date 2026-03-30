package handler

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type UserHandler struct {
	userService port.UserService
}

func NewUserHandler(userService port.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(c echo.Context) error {
	var req CreateUserRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}

	input, err := toCreateUserDomain(req, helper.ActorIDFromContext(c))
	if err != nil {
		return helper.RespondError(c, err)
	}

	user, err := h.userService.Create(c.Request().Context(), input)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_user_id", "user id must be a valid uuid", err))
	}

	var req UpdateUserRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}

	input, err := toUpdateUserDomain(id, req, helper.ActorIDFromContext(c))
	if err != nil {
		return helper.RespondError(c, err)
	}

	user, err := h.userService.Update(c.Request().Context(), input)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_user_id", "user id must be a valid uuid", err))
	}

	if err := h.userService.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}

func (h *UserHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_user_id", "user id must be a valid uuid", err))
	}

	user, err := h.userService.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) FindByEmail(c echo.Context) error {
	email := strings.TrimSpace(c.QueryParam("email"))
	user, err := h.userService.FindByEmail(c.Request().Context(), email)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) FindAll(c echo.Context) error {
	query, err := buildUserQuery(c)
	if err != nil {
		return helper.RespondError(c, err)
	}

	result, err := h.userService.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, http.StatusOK, toUsersListResponse(result, query))
}

func buildUserQuery(c echo.Context) (domain.UserQuery, error) {
	paging := helper.BuildPaging(c)
	sortBy, sortOrder := helper.BuildSort(c)

	query := domain.UserQuery{
		Paging:    paging,
		Search:    strings.TrimSpace(c.QueryParam("search")),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	if roleID := strings.TrimSpace(c.QueryParam("role_id")); roleID != "" {
		parsed, err := uuid.Parse(roleID)
		if err != nil {
			return domain.UserQuery{}, helper.BadRequest("invalid_role_id", "role_id must be a valid uuid", err)
		}

		query.RoleID = &parsed
	}

	if status := strings.TrimSpace(c.QueryParam("status")); status != "" {
		userStatus := domain.UserStatus(status)
		if userStatus != domain.UserStatusActive && userStatus != domain.UserStatusInactive {
			return domain.UserQuery{}, helper.BadRequest("invalid_status", "status must be active or inactive", nil)
		}

		query.Status = &userStatus
	}

	return query, nil
}
