package handler

import (
	"errors"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/driver"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
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

func (h *Handler) UpdateProfile(d driver.PutDriverParams) middleware.Responder {
	id, ok := IdFromContext(d.HTTPRequest.Context())
	if !ok {
		body := driver.PutDriverBadRequestBody{
			Error: fmt.Errorf("bad access token").Error(),
		}
		return driver.NewPutDriverBadRequest().WithPayload(&body)
	}

	dr := model.Driver{
		ID:          uuid.MustParse(id),
		Name:        d.Input.Name,
		PhoneNumber: d.Input.PhoneNumber,
		Email:       d.Input.Email,
	}

	err := h.s.UpdateProfile(d.HTTPRequest.Context(), dr)
	if err != nil {
		if errors.Is(err, service.ErrDriverDoesNotExists) {
			body := driver.PutDriverBadRequestBody{
				Error: err.Error(),
			}
			return driver.NewPutDriverBadRequest().WithPayload(&body)
		}

		body := driver.PutDriverInternalServerErrorBody{
			Error: fmt.Errorf("update profile failed: %w", err).Error(),
		}
		return driver.NewPutDriverInternalServerError().WithPayload(&body)
	}
	return driver.NewPutDriverOK()
}
