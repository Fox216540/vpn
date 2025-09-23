package interceptor

import "vpn/src/core/exception"

const (
	layer  = "Interceptor"
	domain = "Config"
)

type InterceptorUnauthenticated struct {
	*exception.Unauthenticated
}

func NewInterceptorUnauthenticated(msg string, err error) *InterceptorUnauthenticated {
	return &InterceptorUnauthenticated{
		Unauthenticated: exception.NewUnauthenticated(msg, domain, err),
	}
}

func (e *InterceptorUnauthenticated) Error() string {
	return e.Unauthenticated.Error()
}

type MissingMetadata struct {
	*InterceptorUnauthenticated
}

func NewMissingMetadata(err error) *MissingMetadata {
	return &MissingMetadata{
		InterceptorUnauthenticated: NewInterceptorUnauthenticated("Missing metadata", err),
	}
}

func (e *MissingMetadata) Error() string {
	return e.InterceptorUnauthenticated.Error()
}

type InvalidAuthHeader struct {
	*InterceptorUnauthenticated
}

func NewInvalidAuthHeader(err error) *InvalidAuthHeader {
	return &InvalidAuthHeader{
		InterceptorUnauthenticated: NewInterceptorUnauthenticated("Invalid authorization header", err),
	}
}

func (e *InvalidAuthHeader) Error() string {
	return e.InterceptorUnauthenticated.Error()
}

type InvalidToken struct {
	*InterceptorUnauthenticated
}

func NewInvalidToken(err error) *InvalidToken {
	return &InvalidToken{
		InterceptorUnauthenticated: NewInterceptorUnauthenticated("Invalid token", err),
	}
}

func (e *InvalidToken) Error() string {
	return e.InterceptorUnauthenticated.Error()
}
