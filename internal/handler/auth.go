package handler

import (
	"context"
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
	"go.uber.org/zap"
)

type key string

const userId key = "id"

func (h *Handler) SingUp(d auth.PostDriverSingUpParams) middleware.Responder {
	log, ok := LoggerFromContext(d.HTTPRequest.Context())
	if !ok {
		body := auth.PostDriverSingUpInternalServerErrorBody{
			Error: fmt.Errorf("can't get logger").Error(),
		}
		return auth.NewPostDriverSingUpInternalServerError().WithPayload(&body)
	}

	driver := model.Driver{
		Name:        d.Input.Name,
		PhoneNumber: d.Input.PhoneNumber,
		Email:       d.Input.Email,
		Password:    d.Input.Password,
		TaxiType:    d.Input.TaxiType,
	}

	err := h.s.SingUp(d.HTTPRequest.Context(), driver)
	if err != nil {
		if errors.Is(err, service.ErrDriverDoesNotExists) {
			body := auth.PostDriverSingUpBadRequestBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverSingUpBadRequest().WithPayload(&body)
		}

		log.Error("/driver/sing-up", zap.Error(fmt.Errorf("service sing up failed: %w", err)))
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
	log, ok := LoggerFromContext(d.HTTPRequest.Context())
	if !ok {
		body := auth.PostDriverSingInInternalServerErrorBody{
			Error: fmt.Errorf("can't get logger").Error(),
		}
		return auth.NewPostDriverSingInInternalServerError().WithPayload(&body)
	}

	driver := model.Driver{
		PhoneNumber: d.Input.PhoneNumber,
		Password:    d.Input.Password,
	}

	token, err := h.s.SingIn(d.HTTPRequest.Context(), driver)
	if err != nil {
		if errors.Is(err, service.ErrIncorrectPassword) {
			body := auth.PostDriverSingInForbiddenBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverSingInForbidden().WithPayload(&body)
		}

		log.Error("/driver/sing-in", zap.Error(fmt.Errorf("service sing in failed: %w", err)))
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
		log, ok := LoggerFromContext(r.Context())
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			_, err := rw.Write([]byte(fmt.Errorf("can't get logger").Error()))
			if err != nil {
				log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
			}
			return
		}

		token := strings.Split(r.Header.Get("Authorization"), " ")
		if len(token) < 2 {
			rw.WriteHeader(http.StatusUnauthorized)
			resp["error"] = fmt.Errorf("access token required").Error()
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

				rw.WriteHeader(http.StatusInternalServerError)
				_, err := rw.Write([]byte(err.Error()))
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}
			_, err = rw.Write(jsonResp)
			if err != nil {
				log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
			}
			return
		}
		accessToken := token[1]

		id, err := service.Verify(accessToken, h.Cfg)
		if err != nil {
			if errors.Is(err, jwt.ValidationError{Errors: jwt.ValidationErrorExpired}) {
				rw.WriteHeader(http.StatusUnauthorized)
				resp["error"] = err.Error()
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

					rw.WriteHeader(http.StatusInternalServerError)
					_, err := rw.Write([]byte(err.Error()))
					if err != nil {
						log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
					}
					return
				}
				_, err = rw.Write(jsonResp)
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				rw.WriteHeader(http.StatusForbidden)
				resp["error"] = fmt.Errorf("wrong signature").Error()
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

					rw.WriteHeader(http.StatusInternalServerError)
					_, err := rw.Write([]byte(err.Error()))
					if err != nil {
						log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
					}
					return
				}
				_, err = rw.Write(jsonResp)
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}
			if errors.Is(err, service.ErrTokenId) {
				rw.WriteHeader(http.StatusForbidden)
				resp["error"] = fmt.Errorf("id failed").Error()
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

					rw.WriteHeader(http.StatusInternalServerError)
					_, err := rw.Write([]byte(err.Error()))
					if err != nil {
						log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
					}
					return
				}
				_, err = rw.Write(jsonResp)
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}

			rw.WriteHeader(http.StatusInternalServerError)
			resp["error"] = fmt.Errorf("verify failed: %w", err).Error()
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Error("verefy", zap.Error(fmt.Errorf("json marshal failed: %w", err)))

				rw.WriteHeader(http.StatusInternalServerError)
				_, err := rw.Write([]byte(err.Error()))
				if err != nil {
					log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
				}
				return
			}
			_, err = rw.Write(jsonResp)
			if err != nil {
				log.Error("verify", zap.Error(fmt.Errorf("write  failed: %w", err)))
			}
			return
		}
		ctx := ContextWithId(r.Context(), id)
		r = r.WithContext(ctx)
		handler.ServeHTTP(rw, r)

	})
}

func ContextWithId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userId, id)
}

func IdFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userId).(string)
	return id, ok
}

func (h *Handler) Refresh(token auth.PostDriverRefreshParams) middleware.Responder {
	log, ok := LoggerFromContext(token.HTTPRequest.Context())
	if !ok {
		body := auth.PostDriverRefreshInternalServerErrorBody{
			Error: fmt.Errorf("can't get logger").Error(),
		}
		return auth.NewPostDriverRefreshInternalServerError().WithPayload(&body)
	}

	id, err := service.Verify(token.Input.RefreshToken, h.Cfg)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			body := auth.PostDriverRefreshUnauthorizedBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverRefreshUnauthorized().WithPayload(&body)
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			body := auth.PostDriverRefreshForbiddenBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverRefreshForbidden().WithPayload(&body)
		}

		log.Error("/driver/refresh", zap.Error(fmt.Errorf("service verify failed: %w", err)))
		body := auth.PostDriverRefreshInternalServerErrorBody{
			Error: fmt.Errorf("verify rt failed: %w", err).Error(),
		}
		return auth.NewPostDriverRefreshInternalServerError().WithPayload(&body)
	}
	uuid, _ := uuid.Parse(id)
	driver := model.Driver{
		ID: uuid,
	}
	t, err := h.s.Refresh(token.HTTPRequest.Context(), driver)
	if err != nil {
		if errors.Is(err, service.ErrIncorrectPassword) {
			body := auth.PostDriverSingInForbiddenBody{
				Error: err.Error(),
			}
			return auth.NewPostDriverSingInForbidden().WithPayload(&body)
		}

		log.Error("/driver/refresh", zap.Error(fmt.Errorf("service refresh failed: %w", err)))
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
