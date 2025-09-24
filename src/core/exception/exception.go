package exception

import "fmt"

type Exception struct {
	Domain  string
	Layer   string
	Message string
	Err     error
}

func (e Exception) Error() string {
	return fmt.Sprintf("Domain: %s\nLayer: %s\nMessage: %s\nCause: %v",
		e.Domain, e.Layer, e.Message, e.Err)
}

func (e Exception) Unwrap() error {
	return e.Err
}

func NewException(msg, domain, layer string, err error) *Exception {
	return &Exception{
		Message: msg,
		Err:     err,
		Layer:   layer,
		Domain:  domain,
	}
}

const (
	domainLayer = "Domain"
)

type DomainError struct {
	*Exception
}

func (e DomainError) Error() string {
	return e.Exception.Error()
}

func (e DomainError) Unwrap() error { return e.Exception }

func NewDomainException(msg, domain string, err error) *DomainError {
	return &DomainError{
		Exception: NewException(msg, domain, domainLayer, err),
	}
}

type Unauthenticated struct {
	*DomainError
}

func (e Unauthenticated) Error() string {
	return e.DomainError.Error()
}

func (e Unauthenticated) Unwrap() error { return e.DomainError }

func NewUnauthenticated(msg, domain string, err error) *Unauthenticated {
	return &Unauthenticated{
		DomainError: NewDomainException(msg, domain, err),
	}
}

type BadRequestError struct {
	*DomainError
}

func (e BadRequestError) Error() string {
	return e.DomainError.Error()
}

func (e BadRequestError) Unwrap() error { return e.DomainError }

func NewBadRequestError(msg, domain string, err error) *BadRequestError {
	return &BadRequestError{
		DomainError: NewDomainException(msg, domain, err),
	}
}

type ServerError struct {
	*Exception
}

func (e ServerError) Error() string {
	return e.Exception.Error()
}

func (e ServerError) Unwrap() error { return e.Exception }

func NewServerError(msg, domain, layer string, err error) *ServerError {
	return &ServerError{
		Exception: NewException(msg, domain, layer, err),
	}
}
