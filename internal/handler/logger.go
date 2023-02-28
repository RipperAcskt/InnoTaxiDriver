package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handler) Log(handler http.Handler) http.Handler {
	resp := make(map[string]string)
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log, err := zap.NewProduction(zap.Fields(zap.String("url", r.URL.Path), zap.String("method", r.Method), zap.Any("uuid", uuid.New()), zap.String("request time", time.Now().String())))
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			resp["error"] = fmt.Errorf("create logger failed").Error()
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(err.Error()))
				return
			}
			rw.Write(jsonResp)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "logger", log))

		handler.ServeHTTP(rw, r)

		log.Info("request", zap.String("time", time.Since(start).String()))
	})
}

func GetLogger(r *http.Request) zap.Logger {
	log := r.Context().Value("logger")
	logger := log.(zap.Logger)
	return logger
}
