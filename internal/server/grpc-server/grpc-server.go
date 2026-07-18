package grpc_server

import (
	"fmt"
	"net"

	api "github.com/GroVlAn/auth-api/user"
	"google.golang.org/grpc"
)

type Server struct {
	srv     *grpc.Server
	handler api.UserServiceServer
}

func New(handler api.UserServiceServer) *Server {
	return &Server{
		srv:     grpc.NewServer(),
		handler: handler,
	}
}

func (s *Server) ListenAndServe(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listening tcp server: %w", err)
	}

	api.RegisterUserServiceServer(s.srv, s.handler)

	if err = s.srv.Serve(lis); err != nil {
		return fmt.Errorf("serving grpc server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
