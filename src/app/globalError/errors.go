package globalError

import "vpn/src/core/exception"

const layer = "App"

type AppServerError struct {
	*exception.ServerError
}

func NewAppServerError(msg, domain string, err error) *AppServerError {
	return &AppServerError{
		ServerError: exception.NewServerError(msg, domain, layer, err),
	}
}

func (e *AppServerError) Error() string {
	return e.ServerError.Error()
}
