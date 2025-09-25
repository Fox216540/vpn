package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("Failed to create log dir: %v", err)
	}

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// Инициализируем глобальный логгер
	logger.InitLogger(file)

	configHandler := config.NewHandler()
	h := hasher.NewHasher()
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		//log.Fatalf("failed to load TLS keys: %v", err)
	}

	// Создаём gRPC сервер с TLS и интерцептором
	server := grpc.NewServer(
		grpc.Creds(creds), // 🔒 TLS
		grpc.UnaryInterceptor(interceptor.AuthUnaryInterceptor(h)), // 🔐 твой интерцептор
	)
	monitoringHandler := monitoring.NewHandler()
	pbConfig.RegisterConfigServiceServer(server, configHandler)
	pbMonitoring.RegisterMonitoringServiceServer(server, monitoringHandler)

	log.Println("gRPC server is running on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
