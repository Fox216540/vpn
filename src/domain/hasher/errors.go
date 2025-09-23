package hasher

import "vpn/src/core/exception"

const domain = "Hasher"

type DomainBadRequestError struct {
	*exception.BadRequestError
}

func NewDomainBadRequestError(msg string, err error) *DomainBadRequestError {
	return &DomainBadRequestError{
		BadRequestError: exception.NewBadRequestError(msg, domain, err),
	}
}

func (e *DomainBadRequestError) Error() string { return e.BadRequestError.Error() }

type WrongPasswordError struct {
	*DomainBadRequestError
}

func NewWrongPasswordError(err error) *WrongPasswordError {
	return &WrongPasswordError{
		DomainBadRequestError: NewDomainBadRequestError("Wrong password", err),
	}
}

func (e *WrongPasswordError) Error() string { return e.DomainBadRequestError.Error() }
