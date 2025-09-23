package monitoring

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"vpn/src/api/monitoring/di"
	pb "vpn/src/api/monitoring/proto"
	"vpn/src/app/monitoring"
)

type Handler struct {
	service monitoring.UseCase
	pb.UnimplementedMonitoringServiceServer
}

func NewHandler() *Handler {
	ts := di.GetTrafficService()
	return &Handler{
		service: ts,
	}
}

func (h *Handler) GetTotalMbps(ctx context.Context, req *emptypb.Empty) (*pb.GetTotalMbpsResponse, error) {
	totalMbps := h.service.GetTotalMbps()
	return &pb.GetTotalMbpsResponse{
		Mbps: totalMbps,
	}, nil
}
func (h *Handler) GetSplitMbps(ctx context.Context, req *emptypb.Empty) (*pb.GetSplitMbpsResponse, error) {
	down, up := h.service.GetSplitMbps()
	return &pb.GetSplitMbpsResponse{
		Down: down,
		Up:   up,
	}, nil
}

func (h *Handler) GetConnections(ctx context.Context, req *emptypb.Empty) (*pb.GetConnectionsResponse, error) {
	connections := h.service.GetConnections()
	return &pb.GetConnectionsResponse{
		Connections: int32(connections),
	}, nil
}

func (h *Handler) GetCPUPercent(ctx context.Context, req *emptypb.Empty) (*pb.GetCPUPercentResponse, error) {
	cpuPercent := h.service.GetCPUPercent()
	return &pb.GetCPUPercentResponse{
		Percent: cpuPercent,
	}, nil
}

func (h *Handler) GetMemoryPercent(ctx context.Context, req *emptypb.Empty) (*pb.GetMemoryPercentResponse, error) {
	memoryPercent := h.service.GetMemoryPercent()
	return &pb.GetMemoryPercentResponse{
		Percent: memoryPercent,
	}, nil
}
