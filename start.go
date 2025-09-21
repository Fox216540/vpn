package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"vpn/src/api/config"
	pbConfig "vpn/src/api/config/proto"
	"vpn/src/api/traffic"
	pbTraffic "vpn/src/api/traffic/proto"
	"vpn/src/core/interceptor"
	"vpn/src/infra/hasher"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	configHandler := config.NewHandler()
	h := hasher.NewHasher()
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthUnaryInterceptor(h)),
	)
	trafficHandler := traffic.NewHandler()
	pbConfig.RegisterConfigServiceServer(server, configHandler)
	pbTraffic.RegisterTrafficServiceServer(server, trafficHandler)

	log.Println("gRPC server is running on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
