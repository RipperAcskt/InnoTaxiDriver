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

func (h *Handler) GetProfile(d driver.GetDriverParams) middleware.Responder {
	id, ok := IdFromContext(d.HTTPRequest.Context())
	if !ok {
		body := driver.GetDriverBadRequestBody{
			Error: fmt.Errorf("bad access token").Error(),
		}
		return driver.NewGetDriverBadRequest().WithPayload(&body)
	}

	dr, err := h.s.GetProfile(d.HTTPRequest.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrDriverDoesNotExists) {
			body := driver.GetDriverBadRequestBody{
				Error: err.Error(),
			}
			return driver.NewGetDriverBadRequest().WithPayload(&body)
		}

		body := driver.GetDriverInternalServerErrorBody{
			Error: fmt.Errorf("get profile failed: %w", err).Error(),
		}
		return driver.NewGetDriverInternalServerError().WithPayload(&body)
	}

	body := driver.GetDriverOKBody{
		ID:      fmt.Sprint(dr.ID),
		Name:    dr.Name,
		Email:   dr.Email,
		Raiting: float64(dr.Raiting),
	}
	return driver.NewGetDriverOK().WithPayload(&body)
}
