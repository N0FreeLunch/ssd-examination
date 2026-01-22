package service

import (
	"context"
	"examination/internal/api"
	"examination/internal/types"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetHello(ctx context.Context, request api.GetHelloRequestObject) (api.GetHelloResponseObject, error) {
	return api.GetHello200JSONResponse(types.HelloResponse{
		Message: "Hello, World!",
	}), nil
}
