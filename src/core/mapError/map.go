package mapError

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"vpn/src/core/exception"
	"vpn/src/core/logger"
)

func MapError(e error) error {
	if e != nil {
		logger.Log.Error(e.Error()) // логируем текст ошибки
	}
	var se *exception.ServerError
	var be *exception.BadRequestError
	var ua *exception.Unauthenticated

	switch {
	case errors.As(e, &se):
		return status.Error(codes.Internal, "Internal server error")
	case errors.As(e, &be):
		return status.Error(codes.InvalidArgument, "Bad request")
	case errors.As(e, &ua):
		return status.Error(codes.Unauthenticated, "Unauthenticated")
	default:
		return nil
	}
}
