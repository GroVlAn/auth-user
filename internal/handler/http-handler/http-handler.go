package httphandler

import (
	"context"
	"net/http"
	"time"

	"github.com/GroVlAn/auth-user/internal/core"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type service interface {
	Create(ctx context.Context, user core.User) error
	User(ctx context.Context, userQuery core.UserQuery) (core.User, error)
	UserInfo(ctx context.Context, userQuery core.UserQuery) (core.UserInfo, error)
	ChangePassword(ctx context.Context, userQueryNewPassword core.UserQueryNewPassword) error
	InactivateUser(ctx context.Context, userQuery core.UserQuery) error
	RestoreUser(ctx context.Context, userQuery core.UserQuery) error
	BanUser(ctx context.Context, userQuery core.UserQuery) error
	UnbanUser(ctx context.Context, userQuery core.UserQuery) error
}

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

	r.Route(h.BasePath, func(r chi.Router) {
		h.userRoute(r)
	})

	return r
}
