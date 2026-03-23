package httphandler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type service interface{}

type Deps struct {
	BasePath       string
	DefaultTimeout time.Duration
}

type HTTPHandler struct {
	s service
	l zerolog.Logger
	Deps
}

func New(s service, l zerolog.Logger, deps Deps) *HTTPHandler {
	return &HTTPHandler{
		s:    s,
		l:    l,
		Deps: deps,
	}
}

func (h *HTTPHandler) Handler() http.Handler {
	r := chi.NewRouter()

	h.useMiddleware(r)

	r.Route("/", func(r chi.Router) {
		r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome to the Home Page!"))
		})
	})

	return r
}
