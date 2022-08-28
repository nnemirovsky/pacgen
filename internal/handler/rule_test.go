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
	"strconv"
	"strings"
	"testing"
)

func testPrepareRuleHandler(t *testing.T) (*RuleHandler, *mock.RuleService) {
	ctrl := gomock.NewController(t)
	ruleSrvcMock := mock.NewRuleService(ctrl)

	return NewRuleHandler(ruleSrvcMock, logutil.DiscardLogger), ruleSrvcMock
}

func TestRuleHandler_GetAll_OK(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	rules := []model.Rule{
		{
			ID:           1,
			Regex:        `^www\.google\.com$`,
			ProxyProfile: &model.ProxyProfile{ID: 1},
		},
		{
			ID:           2,
			Regex:        `(?:^|\.)facebook\.com$`,
			ProxyProfile: &model.ProxyProfile{ID: 2},
		},
	}

	want := `[{"id":1,"regexp":"^www\\.google\\.com$","proxy_profile_id":1},` +
		`{"id":2,"regexp":"(?:^|\\.)facebook\\.com$","proxy_profile_id":2}]`

	ruleSrvcMock.EXPECT().GetAll(gomock.Any()).Return(rules, nil)

	req, err := http.NewRequest(http.MethodGet, "/rules", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.GetAll)

	handler.ServeHTTP(rr, req)

	got := strings.TrimSuffix(rr.Body.String(), "\n")

	assert.Equal(t, rr.Code, http.StatusOK)

	assert.Equal(t, got, want)
}

func TestRuleHandler_GetAll_InternalServerError(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	ruleSrvcMock.EXPECT().GetAll(gomock.Any()).Return(nil, errs.ServiceUnknownError)

	req, err := http.NewRequest(http.MethodGet, "/rules", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.GetAll)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestRuleHandler_GetById_OK(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	rule := model.Rule{
		ID:           1,
		Regex:        `^www\.google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 14},
	}

	want := `{"id":1,"regexp":"^www\\.google\\.com$","proxy_profile_id":14}`
	ruleSrvcMock.EXPECT().GetByID(gomock.Any(), 1).Return(rule, nil)

	req, err := http.NewRequest(http.MethodGet, "/rules/1", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.GetByID)

	handler.ServeHTTP(rr, req)

	got := strings.TrimSuffix(rr.Body.String(), "\n")

	assert.Equal(t, rr.Code, http.StatusOK)

	assert.Equal(t, got, want)
}

func TestRuleHandler_GetById_NotFound(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	ruleSrvcMock.EXPECT().GetByID(gomock.Any(), 19).Return(model.Rule{}, &errs.EntityNotFoundError{})

	req, err := http.NewRequest(http.MethodGet, "/rules/19", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "19")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.GetByID)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound)
}

func TestRuleHandler_GetById_InternalServerError(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	ruleSrvcMock.EXPECT().GetByID(gomock.Any(), 1).Return(model.Rule{}, errs.ServiceUnknownError)

	req, err := http.NewRequest(http.MethodGet, "/rules/1", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.GetByID)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestRuleHandler_Create_OK(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	rule := model.Rule{
		Regex:        `^google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	const insertedID = 15

	ruleSrvcMock.EXPECT().Create(gomock.Any(), &rule).DoAndReturn(
		func(ctx context.Context, rule *model.Rule) error {
			rule.ID = insertedID
			return nil
		},
	)

	body := `{"domain":"google.com","mode":"domain","proxy_profile_id":1}`

	req, err := http.NewRequest(http.MethodPost, "http://localhost/rules", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Create)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusCreated)

	assert.Equal(t, rr.Header().Get("Location"), "http://localhost/rules/"+strconv.Itoa(insertedID))
}

func TestRuleHandler_Create_UnprocessableEntity(t *testing.T) {
	t.Parallel()

	ruleHandler, _ := testPrepareRuleHandler(t)

	cases := map[string]string{
		"missing proxy_profile_id": `{"domain":"google.com","mode":"domain"}`,
		"invalid mode":             `{"domain":"google.com","proxy_profile_id":1,"mode":"just_domain"}`,
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
			handler := http.HandlerFunc(ruleHandler.Create)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
		})
	}
}

func TestRuleHandler_Create_Conflict(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	rule := model.Rule{
		Regex:        `^google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	ruleSrvcMock.EXPECT().Create(gomock.Any(), &rule).Return(errs.InvalidReferenceError)

	body := `{"domain":"google.com","mode":"domain","proxy_profile_id":1}`

	req, err := http.NewRequest(http.MethodPost, "/rules", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Create)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusConflict)
}

func TestRuleHandler_Create_InternalServerError(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	rule := model.Rule{
		Regex:        `^google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	ruleSrvcMock.EXPECT().Create(gomock.Any(), &rule).Return(errs.ServiceUnknownError)

	body := `{"domain":"google.com","mode":"domain","proxy_profile_id":1}`

	req, err := http.NewRequest(http.MethodPost, "/rules", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Create)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestRuleHandler_Update_OK(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	rule := model.Rule{
		ID:           12,
		Regex:        `^google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	ruleSrvcMock.EXPECT().Update(gomock.Any(), rule).Return(nil)

	body := `{"domain":"google.com","mode":"domain","proxy_profile_id":1}`

	req, err := http.NewRequest(http.MethodPut, "/rules/12", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "12")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNoContent)
}

func TestRuleHandler_Update_UnprocessableEntity(t *testing.T) {
	t.Parallel()

	ruleHandler, _ := testPrepareRuleHandler(t)

	cases := map[string]string{
		"missing proxy_profile_id": `{"domain":"google.com","mode":"domain"}`,
		"invalid mode":             `{"domain":"google.com","proxy_profile_id":1,"mode":"just_domain"}`,
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodGet, "/rules/1", strings.NewReader(body))
			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ruleHandler.Update)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
		})
	}
}

func TestRuleHandler_Update_BadRequest(t *testing.T) {
	t.Parallel()

	ruleHandler, _ := testPrepareRuleHandler(t)

	body := `{"domain":"google.com","mode":"domain","proxy_profile_id":1}`

	req, err := http.NewRequest(http.MethodPut, "/rules/abcd", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abcd")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
}

func TestRuleHandler_Update_NotFound(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	rule := model.Rule{
		ID:           12,
		Regex:        `^google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	ruleSrvcMock.EXPECT().Update(gomock.Any(), rule).Return(&errs.EntityNotFoundError{})
	body := `{"domain":"google.com","mode":"domain","proxy_profile_id":1}`

	req, err := http.NewRequest(http.MethodPut, "/rules/12", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "12")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound)
}

func TestRuleHandler_Update_Conflict(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	rule := model.Rule{
		ID:           12,
		Regex:        `^google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	ruleSrvcMock.EXPECT().Update(gomock.Any(), rule).Return(errs.InvalidReferenceError)
	body := `{"domain":"google.com","mode":"domain","proxy_profile_id":1}`

	req, err := http.NewRequest(http.MethodPut, "/rules/12", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "12")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusConflict)
}

func TestRuleHandler_Update_InternalServerError(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	rule := model.Rule{
		ID:           12,
		Regex:        `^google\.com$`,
		ProxyProfile: &model.ProxyProfile{ID: 1},
	}

	ruleSrvcMock.EXPECT().Update(gomock.Any(), rule).Return(errs.ServiceUnknownError)
	body := `{"domain":"google.com","mode":"domain","proxy_profile_id":1}`

	req, err := http.NewRequest(http.MethodPut, "/rules/12", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "12")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Update)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestRuleHandler_Delete_OK(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	ruleSrvcMock.EXPECT().Delete(gomock.Any(), 12).Return(nil)

	req, err := http.NewRequest(http.MethodDelete, "/rules/12", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "12")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Delete)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNoContent)
}

func TestRuleHandler_Delete_NotFound(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	ruleSrvcMock.EXPECT().Delete(gomock.Any(), 12).Return(&errs.EntityNotFoundError{})

	req, err := http.NewRequest(http.MethodDelete, "/rules/12", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "12")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Delete)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound)
}

func TestRuleHandler_Delete_InternalServerError(t *testing.T) {
	t.Parallel()

	ruleHandler, ruleSrvcMock := testPrepareRuleHandler(t)

	ruleSrvcMock.EXPECT().Delete(gomock.Any(), 12).Return(errs.ServiceUnknownError)

	req, err := http.NewRequest(http.MethodDelete, "/rules/12", nil)
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "12")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ruleHandler.Delete)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}
