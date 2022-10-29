package grpc

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

type GrpcServer struct {
	Grpc   *grpc.Server
	Config *config.GrpcConfig
	Log    logger.ILogger
}

func NewGrpcServer(log logger.ILogger, config *config.GrpcConfig) *GrpcServer {

	unaryServerInterceptors := []grpc.UnaryServerInterceptor{
		otelgrpc.UnaryServerInterceptor(),
	}
	streamServerInterceptors := []grpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),
	}

	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		//https://github.com/open-telemetry/opentelemetry-go-contrib/tree/00b796d0cdc204fa5d864ec690b2ee9656bb5cfc/instrumentation/google.golang.org/grpc/otelgrpc
		//github.com/grpc-ecosystem/go-grpc-middleware
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			streamServerInterceptors...,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			unaryServerInterceptors...,
		)),
	)

	return &GrpcServer{Grpc: s, Config: config, Log: log}
}

func (s *GrpcServer) RunGrpcServer(ctx context.Context, configGrpc ...func(grpcServer *grpc.Server)) error {
	listen, err := net.Listen("tcp", s.Config.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.Grpc)
		}
	}

	if s.Config.Development {
		reflection.Register(s.Grpc)
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.Grpc)
		}
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.Log.Errorf("shutting down grpc PORT: {%s}", s.Config.Port)
				s.shutdown()
				s.Log.Error("grpc exited properly")
				return
			}
		}
	}()

	s.Log.Infof("grpc server is listening on port: %s", s.Config.Port)

	err = s.Grpc.Serve(listen)

	if err != nil {
		s.Log.Error(fmt.Sprintf("[grpcServer_RunGrpcServer.Serve] grpc server serve error: %+v", err))
	}

	return err
}

func (s *GrpcServer) shutdown() {
	s.Grpc.Stop()
	s.Grpc.GracefulStop()
}
