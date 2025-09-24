package config

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"vpn/src/api/config/di"
	pb "vpn/src/api/config/proto"
	"vpn/src/app/config"
	"vpn/src/core/mapError"
)

type Handler struct {
	service config.UseCase
	pb.UnimplementedConfigServiceServer
}

func NewHandler() *Handler {
	cs := di.GetConfigService()
	return &Handler{
		service: cs,
	}
}

func (h *Handler) AddClient(ctx context.Context, req *pb.AddClientRequest) (*pb.AddClientResponse, error) {
	uuidID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, mapError.MapError(err)
	}
	file, err := h.service.CreateConfig(uuidID)
	if err != nil {
		return nil, mapError.MapError(err)
	}
	return &pb.AddClientResponse{
		Message: "Good",
		File:    file,
	}, nil

}

func (h *Handler) DeleteClient(ctx context.Context, req *pb.DeleteClientRequest) (*pb.DeleteClientResponse, error) {
	uuidID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, mapError.MapError(err)
	}
	if err = h.service.DeleteConfig(uuidID); err != nil {
		return nil, mapError.MapError(err)
	}
	return &pb.DeleteClientResponse{
		Message: "Good",
	}, nil
}

func (h *Handler) DeleteClients(ctx context.Context, req *pb.DeleteClientsRequest) (*pb.DeleteClientsResponse, error) {
	uuidIDs := make([]uuid.UUID, len(req.Ids))
	for i, idStr := range req.Ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			fmt.Printf("Invalid UUID at index %d: %v\n", i, err)
			continue
		}
		uuidIDs[i] = id
	}

	if err := h.service.DeleteConfigs(uuidIDs); err != nil {
		return nil, mapError.MapError(err)
	}
	return &pb.DeleteClientsResponse{
		Message: "Good",
	}, nil
}

func (h *Handler) StartServer(ctx context.Context, req *emptypb.Empty) (*pb.StartServerResponse, error) {
	if err := h.service.StartServerConfig(); err != nil {
		return nil, mapError.MapError(err)
	}
	return &pb.StartServerResponse{
		Message: "Good",
	}, nil
}
