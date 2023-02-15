package handler

import (
	"errors"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/auth"
	"github.com/go-openapi/runtime/middleware"
)

func (h *Handler) SingUp(user auth.PostDriverSingUpParams) middleware.Responder {
	err := h.s.CreateDriver(user.Input)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			body := auth.PostDriverSingUpBadRequestBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverSingUpBadRequest().WithPayload(&body)
		}

		body := auth.PostDriverSingUpInternalServerErrorBody{
			Error: fmt.Errorf("create driver failed: %v", err).Error(),
		}
		return auth.NewPostDriverSingUpInternalServerError().WithPayload(&body)
	}

	body := auth.PostDriverSingUpCreatedBody{
		Status: service.StatusCreated,
	}
	return auth.NewPostDriverSingUpCreated().WithPayload(&body)
}

func (h *Handler) SingIn(user auth.PostDriverSingUpParams) middleware.Responder {
	err := h.s.CreateDriver(user.Input)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			body := auth.PostDriverSingUpBadRequestBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverSingUpBadRequest().WithPayload(&body)
		}

		body := auth.PostDriverSingUpInternalServerErrorBody{
			Error: fmt.Errorf("create driver failed: %v", err).Error(),
		}
		return auth.NewPostDriverSingUpInternalServerError().WithPayload(&body)
	}

	body := auth.PostDriverSingUpCreatedBody{
		Status: service.StatusCreated,
	}
	return auth.NewPostDriverSingUpCreated().WithPayload(&body)
}
