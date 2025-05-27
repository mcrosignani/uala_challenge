package follower

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
		GetFollowers(ctx context.Context, req entities.GetFollowersRequest, pagination entities.Pagination) (entities.PagedResponse, error)
		InsertFollower(ctx context.Context, req entities.InsertFollowerRequest) error
		DeleteFollower(ctx context.Context, userID, followID int64) error
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

func (h *Handler) GetFollowers(c echo.Context) error {
	var req entities.GetFollowersRequest
	pagination := entities.NewDefaultPagination()
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	followers, err := h.service.GetFollowers(c.Request().Context(), req, pagination)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, followers)
}

func (h *Handler) Create(c echo.Context) error {
	var req entities.InsertFollowerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.service.InsertFollower(c.Request().Context(), req); err != nil {
		var alreadyExistsErr *customerrors.AlreadyExistsError
		if errors.As(err, &alreadyExistsErr) {
			return echo.NewHTTPError(http.StatusConflict, err)
		}
		var notFoundErr *customerrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *Handler) Delete(c echo.Context) error {
	userID := c.Param("user_id")
	followID := c.Param("follow_id")
	if userID == "" || followID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id and follow_id are required")
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user_id format")
	}

	followIDInt, err := strconv.ParseInt(followID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid follow_id format")
	}

	if err := h.service.DeleteFollower(c.Request().Context(), userIDInt, followIDInt); err != nil {
		var notFoundErr *customerrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
