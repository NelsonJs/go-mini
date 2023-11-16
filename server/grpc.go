package server

import "google.golang.org/grpc"

type Server interface {
	GrpcServer() *grpc.Server
}

type grpcServer struct {
}

func NewGrpcServer() Server {
	return &grpcServer{}
}

func (s *grpcServer) GrpcServer() *grpc.Server {
	return grpc.NewServer()
}

func (s *grpcServer) Register() {

}
