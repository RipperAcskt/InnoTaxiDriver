package handler

import (
	"errors"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/driver"
	"github.com/go-openapi/runtime/middleware"
)

type Handler struct {
	s   *service.Service
	Cfg *config.Config
}

func New(s *service.Service, cfg *config.Config) *Handler {
	return &Handler{s, cfg}
}

func (h *Handler) DeleteProfile(d driver.DeleteDriverParams) middleware.Responder {
	id, ok := IdFromContext(d.HTTPRequest.Context())
	if !ok {
		body := driver.DeleteDriverBadRequestBody{
			Error: fmt.Errorf("bad access token").Error(),
		}
		return driver.NewDeleteDriverBadRequest().WithPayload(&body)
	}

	err := h.s.DeleteProfile(d.HTTPRequest.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrDriverDoesNotExists) {
			body := driver.DeleteDriverBadRequestBody{
				Error: err.Error(),
			}
			return driver.NewDeleteDriverBadRequest().WithPayload(&body)
		}

		body := driver.DeleteDriverInternalServerErrorBody{
			Error: fmt.Errorf("delete profile failed: %w", err).Error(),
		}
		return driver.NewDeleteDriverInternalServerError().WithPayload(&body)
	}
	return driver.NewDeleteDriverOK()
}
