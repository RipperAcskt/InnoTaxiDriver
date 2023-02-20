package handler

import (
	"github.com/RipperAcskt/innotaxidriver/config"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
)

type Handler struct {
	s   *service.Service
	Cfg *config.Config
}

func New(s *service.Service, cfg *config.Config) *Handler {
	return &Handler{s, cfg}
}
