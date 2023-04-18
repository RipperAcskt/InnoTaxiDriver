package handler

import (
	"errors"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/driver"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	s   *service.Service
	Cfg *config.Config
}

func New(s *service.Service, cfg *config.Config) *Handler {
	return &Handler{s, cfg}
}

func (h *Handler) GetProfile(d driver.GetDriverParams) middleware.Responder {
	log, ok := LoggerFromContext(d.HTTPRequest.Context())
	if !ok {
		body := driver.GetDriverInternalServerErrorBody{
			Error: fmt.Errorf("can't get logger").Error(),
		}
		return driver.NewGetDriverInternalServerError().WithPayload(&body)
	}

	id, ok := IdFromContext(d.HTTPRequest.Context())
	if !ok {
		body := driver.GetDriverBadRequestBody{
			Error: fmt.Errorf("can't get id").Error(),
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

		log.Error("/driver", zap.Error(fmt.Errorf("get profile failed: %w", err)))
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
	log, ok := LoggerFromContext(d.HTTPRequest.Context())
	if !ok {
		body := driver.PutDriverInternalServerErrorBody{
			Error: fmt.Errorf("can't get logger").Error(),
		}
		return driver.NewPutDriverInternalServerError().WithPayload(&body)
	}

	id, ok := IdFromContext(d.HTTPRequest.Context())
	if !ok {
		body := driver.PutDriverBadRequestBody{
			Error: fmt.Errorf("can't get id").Error(),
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

		log.Error("/driver", zap.Error(fmt.Errorf("update profile failed: %w", err)))

		body := driver.PutDriverInternalServerErrorBody{
			Error: fmt.Errorf("update profile failed: %w", err).Error(),
		}
		return driver.NewPutDriverInternalServerError().WithPayload(&body)
	}
	return driver.NewPutDriverOK()
}

func (h *Handler) DeleteProfile(d driver.DeleteDriverParams) middleware.Responder {
	log, ok := LoggerFromContext(d.HTTPRequest.Context())
	if !ok {
		body := driver.DeleteDriverInternalServerErrorBody{
			Error: fmt.Errorf("bad access token").Error(),
		}
		return driver.NewDeleteDriverInternalServerError().WithPayload(&body)
	}

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

		log.Error("/driver", zap.Error(fmt.Errorf("delete profile failed: %w", err)))
		body := driver.DeleteDriverInternalServerErrorBody{
			Error: fmt.Errorf("delete profile failed: %w", err).Error(),
		}
		return driver.NewDeleteDriverInternalServerError().WithPayload(&body)
	}
	return driver.NewDeleteDriverOK()
}

func (h *Handler) Recovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log, ok := LoggerFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(fmt.Errorf("bad access token").Error()))
			if err != nil {
				log.Error("Recovery", zap.Error(fmt.Errorf("write  failed: %w", err)))
			}
			return
		}

		defer func() {
			err := recover()
			if err != nil {
				log.Error("recovery", zap.Any("error", err))

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, erro := w.Write(jsonBody)
				if err != nil {
					log.Error("Recovery", zap.Error(fmt.Errorf("write  failed: %w", erro)))
				}
			}
		}()

		handler.ServeHTTP(w, r)

	})
}
