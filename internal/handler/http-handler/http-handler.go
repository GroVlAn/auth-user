package http_handler

import (
	"context"
	"net/http"
	"time"

	"github.com/GroVlAn/auth-user/internal/domain"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type service interface {
	Create(ctx context.Context, user domain.User) error
	User(ctx context.Context, userQuery domain.UserQuery) (domain.User, error)
	UserInfo(ctx context.Context, userQuery domain.UserQuery) (domain.UserInfo, error)
	UpdatePassword(ctx context.Context, userQueryNewPassword domain.UserQueryNewPassword) error
	InactivateUser(ctx context.Context, userQuery domain.UserQuery) error
	RestoreUser(ctx context.Context, userQuery domain.UserQuery) error
	BanUser(ctx context.Context, userQuery domain.UserQuery) error
	UnbanUser(ctx context.Context, userQuery domain.UserQuery) error
}

type Deps struct {
	BasePath       string
	DefaultTimeout time.Duration
}

type HTTPHandler struct {
	l zerolog.Logger
	s service
	Deps
}

func New(l zerolog.Logger, s service, deps Deps) *HTTPHandler {
	return &HTTPHandler{
		l:    l,
		s:    s,
		Deps: deps,
	}
}

func (h *HTTPHandler) Handler() http.Handler {
	r := chi.NewRouter()

	h.cors(r)

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
