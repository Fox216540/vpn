package globalError

import "vpn/src/core/exception"

const layer = "Infra"

type InfraGlobalError struct {
	*exception.ServerError
}

func NewInfraGlobalError(msg, domain string, err error) *InfraGlobalError {
	return &InfraGlobalError{
		ServerError: exception.NewServerError(msg, domain, layer, err),
	}
}

func (e *InfraGlobalError) Error() string {
	return e.ServerError.Error()
}
