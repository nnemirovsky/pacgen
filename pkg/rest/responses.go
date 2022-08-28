package rest

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	url2 "net/url"
)

// Created returns an HTTP 201 Created response with Location header set to the corresponding url.
func Created(w http.ResponseWriter, r *http.Request, id any) {
	url := url2.URL{
		Scheme: GetScheme(r),
		Host:   GetHost(r),
		Path:   fmt.Sprintf("%s/%v", GetPath(r), id),
	}
	w.Header().Set("Location", url.String())
	w.WriteHeader(http.StatusCreated)
}

type ErrorResponse struct {
	StatusCode int    `json:"-"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

//func InternalErrorResponse(errorText string) *ErrorResponse {
//	return &ErrorResponse{StatusCode: http.StatusInternalServerError, ErrorText: errorText}
//}

func BadRequestResponse(errorText string) *ErrorResponse {
	return &ErrorResponse{StatusCode: http.StatusBadRequest, ErrorText: errorText}
}

func NotFoundResponse(errorText string) *ErrorResponse {
	return &ErrorResponse{StatusCode: http.StatusNotFound, ErrorText: errorText}
}

func ConflictResponse(errorText string) *ErrorResponse {
	return &ErrorResponse{StatusCode: http.StatusConflict, ErrorText: errorText}
}

//func UnprocessableEntityResponse(errorText string) *ErrorResponse {
//	return &ErrorResponse{StatusCode: http.StatusUnprocessableEntity, ErrorText: errorText}
//}
