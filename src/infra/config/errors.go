package config

import (
	"vpn/src/infra/globalError"
)

const domain = "Config"

type ConfigError struct {
	*globalError.InfraGlobalError
}

func NewConfigError(msg string, err error) *ConfigError {
	return &ConfigError{
		InfraGlobalError: globalError.NewInfraGlobalError(msg, domain, err),
	}
}

func (e *ConfigError) Error() string {
	return e.InfraGlobalError.Error()
}

type InvalidCreateConfig struct {
	*ConfigError
}

func NewInvalidCreateConfig(err error) *InvalidCreateConfig {
	return &InvalidCreateConfig{
		ConfigError: NewConfigError("Invalid create config", err),
	}
}

func (e *InvalidCreateConfig) Error() string {
	return e.ConfigError.Error()
}

type InvalidDeleteConfig struct {
	*ConfigError
}

func NewInvalidDeleteConfig(err error) *InvalidDeleteConfig {
	return &InvalidDeleteConfig{
		ConfigError: NewConfigError("Invalid delete config", err),
	}
}

func (e *InvalidDeleteConfig) Error() string {
	return e.ConfigError.Error()
}

type InvalidDownloadScript struct {
	*ConfigError
}

func NewInvalidDownloadScript(err error) *InvalidDownloadScript {
	return &InvalidDownloadScript{
		ConfigError: NewConfigError("Invalid download script", err),
	}
}

func (e *InvalidDownloadScript) Error() string {
	return e.ConfigError.Error()
}

type InvalidChmodScript struct {
	*ConfigError
}

func NewInvalidChmodScript(err error) *InvalidChmodScript {
	return &InvalidChmodScript{
		ConfigError: NewConfigError("Invalid chmod script", err),
	}
}

func (e *InvalidChmodScript) Error() string {
	return e.ConfigError.Error()
}

type InvalidCreateServer struct {
	*ConfigError
}

func NewInvalidCreateServer(err error) *InvalidCreateServer {
	return &InvalidCreateServer{
		ConfigError: NewConfigError("Invalid create server", err),
	}
}

func (e *InvalidCreateServer) Error() string {
	return e.ConfigError.Error()
}

type InvalidMkdir struct {
	*ConfigError
}

func NewInvalidMkdir(err error) *InvalidMkdir {
	return &InvalidMkdir{
		ConfigError: NewConfigError("Invalid mkdir", err),
	}
}

func (e *InvalidMkdir) Error() string {
	return e.ConfigError.Error()
}

type InvalidRevokeClients struct {
	*ConfigError
}

func NewInvalidRevokeClients(err error) *InvalidRevokeClients {
	return &InvalidRevokeClients{
		ConfigError: NewConfigError("Invalid revoke clients", err),
	}
}

func (e *InvalidRevokeClients) Error() string {
	return e.ConfigError.Error()
}

type InvalidFinalizeCRL struct {
	*ConfigError
}

func NewInvalidFinalizeCRL(err error) *InvalidFinalizeCRL {
	return &InvalidFinalizeCRL{
		ConfigError: NewConfigError("Invalid finalize CRL", err),
	}
}

func (e *InvalidFinalizeCRL) Error() string {
	return e.ConfigError.Error()
}

type InvalidDisconnectClients struct {
	*ConfigError
}

func NewInvalidDisconnectClients(err error) *InvalidDisconnectClients {
	return &InvalidDisconnectClients{
		ConfigError: NewConfigError("Invalid disconnect clients", err),
	}
}

func (e *InvalidDisconnectClients) Error() string {
	return e.ConfigError.Error()
}
