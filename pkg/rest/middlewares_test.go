package rest

import (
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var fakeHandler = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

func TestURLFixer(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "/rules", nil)
	req.Host = "localhost:8080"
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	req2, err := http.NewRequest(http.MethodGet, "/rules", nil)
	req2.Host = "localhost:8080"
	req2.Header.Set("X-Forwarded-Proto", "https")
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rr := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()

	fixer := URLFixer(fakeHandler)

	fixer.ServeHTTP(rr, req)
	fixer.ServeHTTP(rr2, req2)

	assert.Equal(t, req.URL.Host, req.Host)
	assert.Equal(t, req.URL.Scheme, "http")

	assert.Equal(t, req2.URL.Host, req2.Host)
	assert.Equal(t, req2.URL.Scheme, "https")
}

func TestValidateJSONBody_OK(t *testing.T) {
	t.Parallel()

	body := `{"vars": [124, 8764, 435, 1e6], "name": "test", "object": {"key": "value"}}`
	req, err := http.NewRequest(http.MethodPost, "/rules", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rr := httptest.NewRecorder()

	validator := ValidateJSONBody(fakeHandler)

	validator.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestValidateJSONBody_Invalid(t *testing.T) {
	t.Parallel()

	body := `{"vars": [124, 8764, 435, abcd], "name": "test", "object": {"key": "value"}}`
	req, err := http.NewRequest(http.MethodPost, "/rules", strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}

	rr := httptest.NewRecorder()

	validator := ValidateJSONBody(fakeHandler)

	validator.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
}
