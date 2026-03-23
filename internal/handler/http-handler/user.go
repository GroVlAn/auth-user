package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GroVlAn/auth-user/internal/core"
	"github.com/GroVlAn/auth-user/internal/core/e"
	"github.com/go-chi/chi"
)

const (
	registerEndpoint       = "/user/register"
	userEndpoint           = "/user"
	inactivateUserEndpoint = "/user/inactivate"
	restoreUserEndpoint    = "/user/restore"
	banUserEndpoint        = "/user/ban"
	unbanUserEndpoint      = "/user/unban"
)

func (h *HTTPHandler) userRoute(r chi.Router) {
	r.Post(registerEndpoint, h.register)
	r.Get(userEndpoint, h.user)
	r.Patch(inactivateUserEndpoint, h.inactivateUser)
	r.Patch(restoreUserEndpoint, h.restoreUser)
	r.Patch(banUserEndpoint, h.banUser)
	r.Patch(unbanUserEndpoint, h.unbanUser)
}

func (h *HTTPHandler) register(w http.ResponseWriter, r *http.Request) {
	h.withBodyClose(r.Body, func(body io.ReadCloser) {
		var user core.User
		err := json.NewDecoder(body).Decode(&user)
		if err != nil {
			h.handleDecodeBody(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), h.DefaultTimeout)
		defer cancel()

		err = h.s.Create(ctx, user)
		if err != nil {
			status, res := h.handleError(err)

			h.sendResponse(w, res, status)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("user created"))
	})
}

func (h *HTTPHandler) user(w http.ResponseWriter, r *http.Request) {
	h.withBodyClose(r.Body, func(body io.ReadCloser) {
		var userQuery core.UserQuery
		err := json.NewDecoder(body).Decode(&userQuery)
		if err != nil {
			h.handleDecodeBody(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), h.DefaultTimeout)
		defer cancel()

		user, err := h.s.User(ctx, userQuery)
		if err != nil {
			status, res := h.handleError(err)

			h.sendResponse(w, res, status)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			status, res := h.handleError(
				e.NewErrInternal(fmt.Errorf("failed to encode response body: %w", err)),
			)

			h.sendResponse(w, res, status)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (h *HTTPHandler) inactivateUser(w http.ResponseWriter, r *http.Request) {
	h.withBodyClose(r.Body, func(body io.ReadCloser) {
		h.changeUserStatus(w, r, body, "user inactivated", h.s.InactivateUser)
	})
}

func (h *HTTPHandler) restoreUser(w http.ResponseWriter, r *http.Request) {
	h.withBodyClose(r.Body, func(body io.ReadCloser) {
		h.changeUserStatus(w, r, body, "user restored", h.s.RestoreUser)
	})
}

func (h *HTTPHandler) banUser(w http.ResponseWriter, r *http.Request) {
	h.withBodyClose(r.Body, func(body io.ReadCloser) {
		h.changeUserStatus(w, r, body, "user banned", h.s.BanUser)
	})
}

func (h *HTTPHandler) unbanUser(w http.ResponseWriter, r *http.Request) {
	h.withBodyClose(r.Body, func(body io.ReadCloser) {
		h.changeUserStatus(w, r, body, "user unbanned", h.s.UnbanUser)
	})
}

func (h *HTTPHandler) changeUserStatus(
	w http.ResponseWriter,
	r *http.Request,
	body io.ReadCloser,
	successMessage string,
	fn func(context.Context, core.UserQuery) error,
) {
	var userQuery core.UserQuery
	if err := json.NewDecoder(body).Decode(&userQuery); err != nil {
		h.handleDecodeBody(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.DefaultTimeout)
	defer cancel()

	if err := fn(ctx, userQuery); err != nil {
		status, res := h.handleError(err)

		h.sendResponse(w, res, status)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(successMessage))
}
