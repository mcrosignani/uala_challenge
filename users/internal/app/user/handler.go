package user

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mcrosignani/uala_challenge/users/internal/entities"
	"github.com/mcrosignani/uala_challenge/users/internal/entities/customerrors"
)

type (
	Servicier interface {
		GetUsers(ctx context.Context, req entities.GetUserRequest, pagination entities.Pagination) (entities.PagedResponse, error)
		Create(ctx context.Context, user entities.User) (entities.User, error)
		Delete(ctx context.Context, id int64) error
	}

	Handler struct {
		service Servicier
	}
)

func NewHandler(service Servicier) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetUsers(c echo.Context) error {
	var req entities.GetUserRequest
	pagination := entities.NewDefaultPagination()
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	users, err := h.service.GetUsers(c.Request().Context(), req, pagination)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) Create(c echo.Context) error {
	var user entities.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	user, err := h.service.Create(c.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ID is required")
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	if err := h.service.Delete(c.Request().Context(), userID); err != nil {
		var notFoundErr *customerrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
