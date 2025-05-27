package follower_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mcrosignani/uala_challenge/users/internal/app/follower"
	"github.com/mcrosignani/uala_challenge/users/internal/entities"
	"github.com/mcrosignani/uala_challenge/users/internal/entities/customerrors"
	"github.com/mcrosignani/uala_challenge/users/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetFollowers_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockSvc := new(mocks.FollowerServiceMock)
	mockSvc.On("GetFollowers", mock.Anything, mock.Anything, mock.Anything).Return(entities.PagedResponse{Total: 1}, nil)
	h := follower.NewHandler(mockSvc)

	err := h.GetFollowers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_GetFollowers_BindError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(`invalid-json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockSvc := new(mocks.FollowerServiceMock)
	h := follower.NewHandler(mockSvc)

	err := h.GetFollowers(c)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestHandler_Create_Success(t *testing.T) {
	e := echo.New()
	body := `{"user_id":1,"follow_id":2}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockSvc := new(mocks.FollowerServiceMock)
	mockSvc.On("InsertFollower", mock.Anything, mock.Anything).Return(nil)
	h := follower.NewHandler(mockSvc)

	err := h.Create(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestHandler_Create_AlreadyExists(t *testing.T) {
	e := echo.New()
	body := `{"user_id":1,"follow_id":2}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockSvc := new(mocks.FollowerServiceMock)
	mockSvc.On("InsertFollower", mock.Anything, mock.Anything).Return(&customerrors.AlreadyExistsError{})
	h := follower.NewHandler(mockSvc)

	err := h.Create(c)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusConflict, httpErr.Code)
}

func TestHandler_Create_NotFound(t *testing.T) {
	e := echo.New()
	body := `{"user_id":1,"follow_id":2}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockSvc := new(mocks.FollowerServiceMock)
	mockSvc.On("InsertFollower", mock.Anything, mock.Anything).Return(&customerrors.NotFoundError{})
	h := follower.NewHandler(mockSvc)

	err := h.Create(c)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/followers/1/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id", "follow_id")
	c.SetParamValues("1", "2")

	mockSvc := new(mocks.FollowerServiceMock)
	mockSvc.On("DeleteFollower", mock.Anything, int64(1), int64(2)).Return(nil)
	h := follower.NewHandler(mockSvc)

	err := h.Delete(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestHandler_Delete_NotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/followers/1/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id", "follow_id")
	c.SetParamValues("1", "2")

	mockSvc := new(mocks.FollowerServiceMock)
	mockSvc.On("DeleteFollower", mock.Anything, int64(1), int64(2)).Return(&customerrors.NotFoundError{})
	h := follower.NewHandler(mockSvc)

	err := h.Delete(c)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestHandler_Delete_BadRequest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/followers/x/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id", "follow_id")
	c.SetParamValues("x", "2")

	mockSvc := new(mocks.FollowerServiceMock)
	h := follower.NewHandler(mockSvc)

	err := h.Delete(c)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}
