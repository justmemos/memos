package v2

import (
	"context"
	"fmt"

	grpcRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/usememos/memos/api/v2/auth"
	apiv2pb "github.com/usememos/memos/proto/gen/api/v2"
	"github.com/usememos/memos/server/profile"
	"github.com/usememos/memos/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type APIV2Service struct {
	Secret  string
	Profile *profile.Profile
	Store   *store.Store

	grpcServer     *grpc.Server
	grpcServerPort int
}

func NewAPIV2Service(secret string, profile *profile.Profile, store *store.Store, grpcServerPort int) *APIV2Service {
	authProvider := auth.NewAuthInterceptor(store, secret)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authProvider.AuthenticationInterceptor,
		),
	)
	apiv2pb.RegisterUserServiceServer(grpcServer, NewUserService(store))
	apiv2pb.RegisterTagServiceServer(grpcServer, NewTagService(store))

	return &APIV2Service{
		Secret:         secret,
		Profile:        profile,
		Store:          store,
		grpcServer:     grpcServer,
		grpcServerPort: grpcServerPort,
	}
}

func (s *APIV2Service) GetGRPCServer() *grpc.Server {
	return s.grpcServer
}

// RegisterGateway registers the gRPC-Gateway with the given Echo instance.
func (s *APIV2Service) RegisterGateway(ctx context.Context, e *echo.Echo) error {
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf(":%d", s.grpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	gwMux := grpcRuntime.NewServeMux()
	if err := apiv2pb.RegisterUserServiceHandler(context.Background(), gwMux, conn); err != nil {
		return err
	}
	if err := apiv2pb.RegisterTagServiceHandler(context.Background(), gwMux, conn); err != nil {
		return err
	}
	e.Any("/api/v2/*", echo.WrapHandler(gwMux))

	return nil
}
