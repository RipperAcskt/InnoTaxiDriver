package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/RipperAcskt/innotaxidriver/internal/model"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/auth"
	"github.com/go-openapi/runtime/middleware"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func (h *Handler) SingUp(d auth.PostDriverSingUpParams) middleware.Responder {
	driver := model.Driver{
		Name:        d.Input.Name,
		PhoneNumber: d.Input.PhoneNumber,
		Email:       d.Input.Email,
		Password:    d.Input.Password,
	}

	err := h.s.SingUp(driver)
	if err != nil {
		if errors.Is(err, service.ErrDriverDoesNotExists) {
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

func (h *Handler) SingIn(d auth.PostDriverSingInParams) middleware.Responder {
	driver := model.Driver{
		PhoneNumber: d.Input.PhoneNumber,
		Password:    d.Input.Password,
	}

	token, err := h.s.SingIn(driver)
	if err != nil {
		if errors.Is(err, service.ErrIncorrectPassword) {
			body := auth.PostDriverSingInForbiddenBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverSingInForbidden().WithPayload(&body)
		}

		body := auth.PostDriverSingInInternalServerErrorBody{
			Error: fmt.Errorf("sing in failed: %v", err).Error(),
		}
		return auth.NewPostDriverSingInInternalServerError().WithPayload(&body)
	}

	body := auth.PostDriverSingInOKBody{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return auth.NewPostDriverSingInOK().WithPayload(&body)
}

func (h *Handler) VerifyToken(handler http.Handler) http.Handler {
	resp := make(map[string]string)
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("Authorization"), " ")
		if len(token) < 2 {
			rw.WriteHeader(http.StatusUnauthorized)
			resp["error"] = fmt.Errorf("access token required").Error()
			jsonResp, _ := json.Marshal(resp)
			rw.Write(jsonResp)
			return
		}
		accessToken := token[1]

		_, err := service.Verify(accessToken, h.Cfg)
		if err != nil {
			if strings.Contains(err.Error(), "Token is expired") {
				rw.WriteHeader(http.StatusUnauthorized)
				resp["error"] = err.Error()
				jsonResp, _ := json.Marshal(resp)
				rw.Write(jsonResp)
				return
			}
			if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
				rw.WriteHeader(http.StatusForbidden)
				resp["error"] = fmt.Errorf("wrong signature").Error()
				jsonResp, _ := json.Marshal(resp)
				rw.Write(jsonResp)
				return
			}

			rw.WriteHeader(http.StatusInternalServerError)
			resp["error"] = fmt.Errorf("verify failed: %w", err).Error()
			jsonResp, _ := json.Marshal(resp)
			rw.Write(jsonResp)
			return
		}

		handler.ServeHTTP(rw, r)

	})
}

func (h *Handler) Refresh(token auth.PostDriverRefreshParams) middleware.Responder {

	id, err := service.Verify(token.Input.RefreshToken, h.Cfg)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			body := auth.PostDriverRefreshUnauthorizedBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverRefreshUnauthorized().WithPayload(&body)
		}
		if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
			body := auth.PostDriverRefreshForbiddenBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverRefreshForbidden().WithPayload(&body)
		}

		body := auth.PostDriverRefreshInternalServerErrorBody{
			Error: fmt.Errorf("verify rt failed: %w", err).Error(),
		}
		return auth.NewPostDriverRefreshInternalServerError().WithPayload(&body)
	}
	uuid, _ := uuid.Parse(id)
	driver := model.Driver{
		ID: uuid,
	}
	t, err := h.s.Refresh(driver)
	if err != nil {
		if errors.Is(err, service.ErrIncorrectPassword) {
			body := auth.PostDriverSingInForbiddenBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverSingInForbidden().WithPayload(&body)
		}

		body := auth.PostDriverSingInInternalServerErrorBody{
			Error: fmt.Errorf("sing in failed: %v", err).Error(),
		}
		return auth.NewPostDriverSingInInternalServerError().WithPayload(&body)
	}

	body := auth.PostDriverRefreshOKBody{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}

	return auth.NewPostDriverRefreshOK().WithPayload(&body)
}
