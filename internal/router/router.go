package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/nnemirovsky/pacgen/pkg/rest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"net/http"
	"time"
)

type RuleHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type ProxyProfileHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func New(
	ruleHandler RuleHandler,
	profileHandler ProxyProfileHandler,
	logger zerolog.Logger,
	basicAuthCreds map[string]string,
) http.Handler {
	router := chi.NewRouter()
	initMiddlewares(router, logger, basicAuthCreds)
	initRoutes(router, ruleHandler, profileHandler)
	return router
}

func initMiddlewares(router *chi.Mux, logger zerolog.Logger, basicAuthCreds map[string]string) {
	router.Use(hlog.NewHandler(logger))
	router.Use(hlog.RequestIDHandler("request-id", "X-Request-Id"))
	router.Use(hlog.URLHandler("url"))
	router.Use(hlog.MethodHandler("method"))
	router.Use(hlog.RemoteAddrHandler("remote-addr"))
	router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Debug().
			Int("status", status).
			//Dur("duration", duration).
			Stringer("duration", duration).
			Msg("Request processed")
	}))
	//router.Use(middleware.Recoverer)
	router.Use(rest.Recoverer)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(rest.ValidateJSONBody)
	router.Use(middleware.BasicAuth("/", basicAuthCreds))
}

func initRoutes(router *chi.Mux, ruleHandler RuleHandler, profileHandler ProxyProfileHandler) {
	router.Route("/rules", func(r chi.Router) {
		r.Get("/", ruleHandler.GetAll)
		r.Get("/{id}", ruleHandler.GetByID)
		r.Post("/", ruleHandler.Create)
		r.Put("/{id}", ruleHandler.Update)
		r.Delete("/{id}", ruleHandler.Delete)
	})
	router.Route("/profiles", func(r chi.Router) {
		r.Get("/", profileHandler.GetAll)
		r.Get("/{id}", profileHandler.GetByID)
		r.Post("/", profileHandler.Create)
		r.Put("/{id}", profileHandler.Update)
		r.Delete("/{id}", profileHandler.Delete)
	})
}
