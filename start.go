package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"vpn/src/api/config"
	pbConfig "vpn/src/api/config/proto"
	"vpn/src/api/monitoring"
	pbMonitoring "vpn/src/api/monitoring/proto"
	"vpn/src/core/interceptor"
	"vpn/src/core/logger"
	"vpn/src/infra/hasher"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// Инициализируем глобальный логгер
	logger.InitLogger(file)

	configHandler := config.NewHandler()
	h := hasher.NewHasher()
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthUnaryInterceptor(h)),
	)
	monitoringHandler := monitoring.NewHandler()
	pbConfig.RegisterConfigServiceServer(server, configHandler)
	pbMonitoring.RegisterMonitoringServiceServer(server, monitoringHandler)

	log.Println("gRPC server is running on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
