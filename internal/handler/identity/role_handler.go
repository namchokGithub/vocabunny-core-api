package identity

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type RoleHandler struct {
	roleService port.RoleService
}

func NewRoleHandler(roleService port.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) Create(c echo.Context) error {
	var req CreateRoleRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	role, err := h.roleService.Create(c.Request().Context(), toCreateRoleDomain(req, helper.ActorIDFromContext(c)))
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusCreated, toRoleResponse(role))
}

func (h *RoleHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_role_id", "role id must be a valid uuid", err))
	}
	var req UpdateRoleRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	role, err := h.roleService.Update(c.Request().Context(), toUpdateRoleDomain(id, req, helper.ActorIDFromContext(c)))
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toRoleResponse(role))
}

func (h *RoleHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_role_id", "role id must be a valid uuid", err))
	}
	if err := h.roleService.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}

func (h *RoleHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_role_id", "role id must be a valid uuid", err))
	}
	role, err := h.roleService.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toRoleResponse(role))
}

func (h *RoleHandler) FindAll(c echo.Context) error {
	query := domain.RoleQuery{
		Paging: helper.BuildPaging(c),
		Search: strings.TrimSpace(c.QueryParam("search")),
	}
	query.SortBy, query.SortOrder = helper.BuildSort(c)
	if permission := strings.TrimSpace(c.QueryParam("permission")); permission != "" {
		code := domain.PermissionCode(permission)
		query.Permission = &code
	}
	result, err := h.roleService.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toRolesListResponse(result, query))
}
