package traffic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"vpn/src/api/traffic/di"
	pb "vpn/src/api/traffic/proto"
	"vpn/src/app/traffic"
)

type Handler struct {
	service traffic.UseCase
	pb.UnimplementedTrafficServiceServer
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
