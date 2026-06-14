package http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/GroVlAn/auth-base/ew"
	"github.com/GroVlAn/auth-base/ew/httpx"
	"github.com/GroVlAn/auth-user/internal/domain"
)

func (h *HTTPHandler) withBodyClose(body io.ReadCloser, fn func(io.ReadCloser)) {
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			h.l.Error().Err(err).Msg("failed to close request body")
		}
	}(body)

	fn(body)
}

func (h *HTTPHandler) handleDecodeBody(w http.ResponseWriter, err error) {
	ev := ew.NewErrValidation("failed read request body")

	switch e := err.(type) {
	case *json.SyntaxError:
		ev.AddField("body", fmt.Sprintf("invalid JSON syntax at offset %d", e.Offset))
	case *json.UnmarshalTypeError:
		ev.AddField(e.Field, fmt.Sprintf("expected %v but got %v", e.Type, e.Value))
	default:
		if errors.Is(err, io.EOF) {
			ev.AddField("body", "empty request body")
		} else {
			ev.AddField("body", err.Error())
		}
	}

	h.l.Error().Err(err).Msg("failed to decode request body")

	h.handleError(w, err)
}

func (h *HTTPHandler) sendResponse(w http.ResponseWriter, resp domain.Response, status int) {
	b, err := json.Marshal(resp)
	if err != nil {
		h.l.Error().Err(err).Msg("failed marshal response")
	}

	w.WriteHeader(status)

	_, err = w.Write(b)
	if err != nil {
		h.l.Error().Err(err).Msg("failed write response")
	}
}

func (h *HTTPHandler) handleError(w http.ResponseWriter, err error) {
	respErr := httpx.HandleError(err)

	resp := domain.Response{
		Error: &domain.ErrorResponse{
			Code: respErr.Status,
			Text: respErr.Message,
		},
		Data: respErr.Fields,
	}

	h.l.Err(err).Msg(respErr.LogMsg)

	h.sendResponse(w, resp, respErr.Status)
}
