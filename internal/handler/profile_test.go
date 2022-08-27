package handler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/nnemirovsky/pacgen/internal/errs"
	"github.com/nnemirovsky/pacgen/internal/handler/mock"
	"github.com/nnemirovsky/pacgen/internal/model"
	"github.com/nnemirovsky/pacgen/pkg/logutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testPrepareProfileHandler(t *testing.T) (*ProxyProfileHandler, *mock.ProxyProfileService) {
	ctrl := gomock.NewController(t)
	profileSrvcMock := mock.NewProxyProfileService(ctrl)

	return NewProxyProfileHandler(profileSrvcMock, logutil.DiscardLogger), profileSrvcMock
}

func TestProxyProfileHandler_GetAll_OK(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profiles := []model.ProxyProfile{
		{
			ID:      1,
			Name:    "shadowsocks",
			Type:    model.Socks5,
			Address: "localhost:1080",
		},
		{
			ID:      2,
			Name:    "some http proxy",
			Type:    model.Http,
			Address: "::1:8080",
		},
	}

	want := `[{"id":1,"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"},` +
		`{"id":2,"name":"some http proxy","type":"HTTP","address":"::1:8080"}]`

	profileSrvcMock.EXPECT().GetAll(gomock.Any()).Return(profiles, nil)

	req, err := http.NewRequest(http.MethodGet, "/rules", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.GetAll)

	handler.ServeHTTP(rr, req)

	got := strings.TrimSuffix(rr.Body.String(), "\n")

	assert.Equal(t, rr.Code, http.StatusOK)

	assert.Equal(t, got, want)
}

func TestProxyProfileHandler_GetAll_InternalServerError(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profileSrvcMock.EXPECT().GetAll(gomock.Any()).Return(nil, errs.ServiceUnknownError)

	req, err := http.NewRequest(http.MethodGet, "/rules", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.GetAll)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestProxyProfileHandler_GetByID_OK(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profile := model.ProxyProfile{
		ID:      1,
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	want := `{"id":1,"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"}`

	profileSrvcMock.EXPECT().GetByID(gomock.Any(), 1).Return(profile, nil)

	req, err := http.NewRequest(http.MethodGet, "/rules/1", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.GetByID)

	handler.ServeHTTP(rr, req)

	got := strings.TrimSuffix(rr.Body.String(), "\n")

	assert.Equal(t, rr.Code, http.StatusOK)

	assert.Equal(t, got, want)
}

func TestProxyProfileHandler_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profileSrvcMock.EXPECT().GetByID(gomock.Any(), 1).Return(model.ProxyProfile{}, &errs.EntityNotFoundError{})

	req, err := http.NewRequest(http.MethodGet, "/rules/1", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.GetByID)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound)
}

func TestProxyProfileHandler_GetByID_InternalServerError(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profileSrvcMock.EXPECT().GetByID(gomock.Any(), 1).Return(model.ProxyProfile{}, errs.ServiceUnknownError)

	req, err := http.NewRequest(http.MethodGet, "/rules/1", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.GetByID)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestProxyProfileHandler_Create_OK(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profile := model.ProxyProfile{
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	profileSrvcMock.EXPECT().Create(gomock.Any(), &profile).DoAndReturn(
		func(ctx context.Context, p *model.ProxyProfile) error {
			p.ID = 17
			return nil
		},
	)

	body := `{"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"}`

	req, err := http.NewRequest(http.MethodPost, "/rules", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Create)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusCreated)

	assert.Equal(t, rr.Header().Get("Location"), "/rules/17")
}

func TestProxyProfileHandler_Create_UnprocessableEntity(t *testing.T) {
	t.Parallel()

	profileHandler, _ := testPrepareProfileHandler(t)

	cases := map[string]string{
		"missing address": `{"name":"shadowsocks","type":"SOCKS5"}`,
		"invalid type":    `{"name":"shadowsocks","type":"qwerty","address":"localhost:1080"}`,
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodGet, "/rules/1", strings.NewReader(body))
			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(profileHandler.Create)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
		})
	}
}

func TestProxyProfileHandler_Create_Conflict(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profile := model.ProxyProfile{
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	profileSrvcMock.EXPECT().Create(gomock.Any(), &profile).Return(&errs.EntityAlreadyExistsError{})

	body := `{"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"}`

	req, err := http.NewRequest(http.MethodPost, "/rules", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Create)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusConflict)
}

func TestProxyProfileHandler_Create_InternalServerError(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profile := model.ProxyProfile{
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	profileSrvcMock.EXPECT().Create(gomock.Any(), &profile).Return(errs.ServiceUnknownError)

	body := `{"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"}`

	req, err := http.NewRequest(http.MethodPost, "/rules", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Create)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestProxyProfileHandler_Update_OK(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profile := model.ProxyProfile{
		ID:      1,
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	profileSrvcMock.EXPECT().Update(gomock.Any(), profile).Return(nil)

	body := `{"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"}`

	req, err := http.NewRequest(http.MethodPut, "/rules/1", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNoContent)
}

func TestProxyProfileHandler_Update_BadRequest(t *testing.T) {
	t.Parallel()

	profileHandler, _ := testPrepareProfileHandler(t)

	body := `{"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"}`

	req, err := http.NewRequest(http.MethodPut, "/rules/abcd", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abcd")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
}

func TestProxyProfileHandler_Update_UnprocessableEntity(t *testing.T) {
	t.Parallel()

	profileHandler, _ := testPrepareProfileHandler(t)

	cases := map[string]string{
		"missing address": `{"name":"shadowsocks","type":"SOCKS5"}`,
		"invalid type":    `{"name":"shadowsocks","type":"qwerty","address":"localhost:1080"}`,
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodPut, "/rules/1", strings.NewReader(body))
			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(profileHandler.Update)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
		})
	}
}

func TestProxyProfileHandler_Update_NotFound(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profile := model.ProxyProfile{
		ID:      1,
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	profileSrvcMock.EXPECT().Update(gomock.Any(), profile).Return(&errs.EntityNotFoundError{})

	body := `{"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"}`

	req, err := http.NewRequest(http.MethodPut, "/rules/1", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound)
}

// update already exists
func TestProxyProfileHandler_Update_AlreadyExists(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profile := model.ProxyProfile{
		ID:      1,
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	profileSrvcMock.EXPECT().Update(gomock.Any(), profile).Return(&errs.EntityAlreadyExistsError{})

	body := `{"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"}`

	req, err := http.NewRequest(http.MethodPut, "/rules/1", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusConflict)
}

func TestProxyProfileHandler_Update_InternalServerError(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profile := model.ProxyProfile{
		ID:      1,
		Name:    "shadowsocks",
		Type:    model.Socks5,
		Address: "localhost:1080",
	}

	profileSrvcMock.EXPECT().Update(gomock.Any(), profile).Return(errs.ServiceUnknownError)

	body := `{"name":"shadowsocks","type":"SOCKS5","address":"localhost:1080"}`

	req, err := http.NewRequest(http.MethodPut, "/rules/1", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestProxyProfileHandler_Delete_OK(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profileSrvcMock.EXPECT().Delete(gomock.Any(), 1).Return(nil)

	req, err := http.NewRequest(http.MethodDelete, "/rules/1", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Delete)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNoContent)
}

func TestProxyProfileHandler_Delete_BadRequest(t *testing.T) {
	t.Parallel()

	profileHandler, _ := testPrepareProfileHandler(t)

	req, err := http.NewRequest(http.MethodDelete, "/rules/abcd", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abcd")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Delete)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
}

func TestProxyProfileHandler_Delete_NotFound(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profileSrvcMock.EXPECT().Delete(gomock.Any(), 1).Return(&errs.EntityNotFoundError{})

	req, err := http.NewRequest(http.MethodDelete, "/rules/1", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Delete)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound)
}

func TestProxyProfileHandler_Delete_Conflict(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profileSrvcMock.EXPECT().Delete(gomock.Any(), 1).Return(&errs.EntityStillReferencedError{})

	req, err := http.NewRequest(http.MethodDelete, "/rules/1", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Delete)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusConflict)
}

func TestProxyProfileHandler_Delete_InternalServerError(t *testing.T) {
	t.Parallel()

	profileHandler, profileSrvcMock := testPrepareProfileHandler(t)

	profileSrvcMock.EXPECT().Delete(gomock.Any(), 1).Return(errs.ServiceUnknownError)

	req, err := http.NewRequest(http.MethodDelete, "/rules/1", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler.Delete)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}
