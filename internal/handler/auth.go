package handler

import (
	"errors"
	"fmt"

	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/auth"
	"github.com/go-openapi/runtime/middleware"
)

func (h *Handler) SingUp(d auth.PostDriverSingUpParams) middleware.Responder {
	driver := model.Driver{
		Name:        d.Input.Name,
		PhoneNumber: d.Input.PhoneNumber,
		Email:       d.Input.Email,
		Password:    d.Input.Password,
	}

	err := h.s.CreateDriver(driver)
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
		Status: model.StatusCreated,
	}
	return auth.NewPostDriverSingUpCreated().WithPayload(&body)
}
