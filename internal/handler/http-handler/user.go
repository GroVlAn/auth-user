package http_handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/GroVlAn/auth-user/internal/domain"
	"github.com/go-chi/chi"
)

const (
	registerEndpoint       = "/register"
	userEndPoint           = "/"
	userInfoEndpoint       = "/info"
	changePasswordEndpoint = "/change-password"
	inactivateUserEndpoint = "/inactivate"
	restoreUserEndpoint    = "/restore"
	banUserEndpoint        = "/ban"
	unbanUserEndpoint      = "/unban"
)

func (h *HTTPHandler) userRoute(r chi.Router) {
	r.Post(registerEndpoint, h.register)
	r.Get(userEndPoint, h.user)
	r.Get(userInfoEndpoint, h.userInfo)
	r.Patch(changePasswordEndpoint, h.changePassword)
	r.Patch(inactivateUserEndpoint, h.inactivateUser)
	r.Patch(restoreUserEndpoint, h.restoreUser)
	r.Patch(banUserEndpoint, h.banUser)
	r.Patch(unbanUserEndpoint, h.unbanUser)
}

func (h *HTTPHandler) register(w http.ResponseWriter, r *http.Request) {
	h.withBodyClose(r.Body, func(body io.ReadCloser) {
		var user domain.User
		err := json.NewDecoder(body).Decode(&user)
		if err != nil {
			h.handleDecodeBody(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), h.DefaultTimeout)
		defer cancel()

		if err = h.s.Create(ctx, user); err != nil {
			h.handleError(w, err)
			return
		}

		h.sendSuccessResponse(w, "user created", http.StatusCreated)
	})
}

func (h *HTTPHandler) user(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	req := domain.UserQuery{}

	if id := query.Get("id"); len(id) > 0 {
		req.ID = &id
	}

	if username := query.Get("username"); len(username) > 0 {
		req.Username = &username
	}

	if email := query.Get("email"); len(email) > 0 {
		req.Email = &email
	}

	h.userByQuery(w, r, req)
}

func (h *HTTPHandler) userInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	req := domain.UserQuery{}

	if id := query.Get("id"); len(id) > 0 {
		req.ID = &id
	}

	if username := query.Get("username"); len(username) > 0 {
		req.Username = &username
	}

	if email := query.Get("email"); len(email) > 0 {
		req.Email = &email
	}

	h.userInfoByQuery(w, r, req)
}

func (h *HTTPHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	h.withBodyClose(r.Body, func(body io.ReadCloser) {
		var userQueryNewPassword domain.UserQueryNewPassword

		err := json.NewDecoder(body).Decode(&userQueryNewPassword)
		if err != nil {
			h.handleDecodeBody(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), h.DefaultTimeout)
		defer cancel()

		if err := h.s.UpdatePassword(ctx, userQueryNewPassword); err != nil {
			h.handleError(w, err)

			return
		}

		h.sendSuccessResponse(w, "password changed", http.StatusOK)
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

func (h *HTTPHandler) userByQuery(w http.ResponseWriter, r *http.Request, userQuery domain.UserQuery) {
	ctx, cancel := context.WithTimeout(r.Context(), h.DefaultTimeout)
	defer cancel()

	user, err := h.s.User(ctx, userQuery)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponseWithData(w, user, http.StatusOK)
}

func (h *HTTPHandler) userInfoByQuery(w http.ResponseWriter, r *http.Request, userQuery domain.UserQuery) {
	ctx, cancel := context.WithTimeout(r.Context(), h.DefaultTimeout)
	defer cancel()

	userInfo, err := h.s.UserInfo(ctx, userQuery)
	if err != nil {
		h.handleError(w, err)

		return
	}

	h.sendResponseWithData(w, userInfo, http.StatusOK)
}

func (h *HTTPHandler) changeUserStatus(
	w http.ResponseWriter,
	r *http.Request,
	body io.ReadCloser,
	successMessage string,
	fn func(context.Context, domain.UserQuery) error,
) {
	var userQuery domain.UserQuery
	if err := json.NewDecoder(body).Decode(&userQuery); err != nil {
		h.handleDecodeBody(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.DefaultTimeout)
	defer cancel()

	if err := fn(ctx, userQuery); err != nil {
		h.handleError(w, err)

		return
	}

	h.sendSuccessResponse(w, successMessage, http.StatusOK)
}
