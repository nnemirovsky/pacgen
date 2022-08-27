package rest

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/hlog"
	"io"
	"net/http"
	"os"
	"runtime/debug"
)

func URLFixer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "http"
		if r.TLS != nil {
			r.URL.Scheme = "https"
		}
		if v := r.Header.Get("X-Forwarded-Proto"); v != "" {
			r.URL.Scheme = v
		}

		if r.URL.Host == "" {
			r.URL.Host = r.Host
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateJSONBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			middleware.GetLogEntry(r).Panic(err, debug.Stack())
			return
		}
		if err := r.Body.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			middleware.GetLogEntry(r).Panic(err, debug.Stack())
			return
		}

		if !json.Valid(body) {
			if err := render.Render(w, r, BadRequestResponse("Invalid JSON")); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				middleware.GetLogEntry(r).Panic(err, debug.Stack())
			}
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		next.ServeHTTP(w, r)
	})
}

//func RequestLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			entry := logger.WithField("request-id", middleware.GetReqID(r.Context()))
//			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
//			start := time.Now()
//
//			defer func() {
//				entry.WithFields(logrus.Fields{
//					"uri":      r.RequestURI,
//					"method":   r.Method,
//					"status":   ww.Status(),
//					"duration": time.Since(start),
//					"size":     strconv.Itoa(ww.BytesWritten()) + "B",
//				}).Debug("request completed")
//			}()
//
//			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
//			next.ServeHTTP(ww, r)
//		})
//	}
//}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				logger := hlog.FromRequest(r)
				if isatty.IsTerminal(os.Stderr.Fd()) {
					middleware.PrintPrettyStack(rvr)
				} else {
					//logger.Error().Stack().Err(errors.WithStack(errors.New(rvr.(string)))).Send()
					logger.Error().Stack().Err(errors.New(rvr.(string))).Send()
				}

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
