package monitoring

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	totalMbps, err := h.service.GetTotalMbps()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetTotalMbpsResponse{
		Mbps: totalMbps,
	}, nil
}
func (h *Handler) GetSplitMbps(ctx context.Context, req *emptypb.Empty) (*pb.GetSplitMbpsResponse, error) {
	down, up, err := h.service.GetSplitMbps()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetSplitMbpsResponse{
		Down: down,
		Up:   up,
	}, nil
}

func (h *Handler) GetConnections(ctx context.Context, req *emptypb.Empty) (*pb.GetConnectionsResponse, error) {
	connections, err := h.service.GetConnections()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetConnectionsResponse{
		Connections: int32(connections),
	}, nil
}

func (h *Handler) GetCPUPercent(ctx context.Context, req *emptypb.Empty) (*pb.GetCPUPercentResponse, error) {
	cpuPercent, err := h.service.GetCPUPercent()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetCPUPercentResponse{
		Percent: cpuPercent,
	}, nil
}

func (h *Handler) GetMemoryPercent(ctx context.Context, req *emptypb.Empty) (*pb.GetMemoryPercentResponse, error) {
	memoryPercent, err := h.service.GetMemoryPercent()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetMemoryPercentResponse{
		Percent: memoryPercent,
	}, nil
}
