package config

import "vpn/src/app/globalError"

const domain = "Config"

type ConfigError struct {
	*globalError.AppServerError
}

func NewConfigError(msg string, err error) *ConfigError {
	return &ConfigError{
		AppServerError: globalError.NewAppServerError(msg, domain, err),
	}
}

func (e *ConfigError) Error() string {
	return e.AppServerError.Error()
}

type InvalidCreateConfig struct {
	*ConfigError
}

func NewInvalidCreateConfig(err error) *InvalidCreateConfig {
	return &InvalidCreateConfig{
		ConfigError: NewConfigError("Invalid create config", err),
	}
}

type InvalidReadFile struct {
	*ConfigError
}

func NewInvalidReadFile(err error) *InvalidReadFile {
	return &InvalidReadFile{
		ConfigError: NewConfigError("Invalid read file", err),
	}
}

func (e *InvalidReadFile) Error() string {
	return e.ConfigError.Error()
}

type InvalidRemoveFile struct {
	*ConfigError
}

func NewInvalidRemoveFile(err error) *InvalidRemoveFile {
	return &InvalidRemoveFile{
		ConfigError: NewConfigError("Invalid remove file", err),
	}
}

func (e *InvalidRemoveFile) Error() string {
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

type InvalidStartServer struct {
	*ConfigError
}

func NewInvalidStartServer(err error) *InvalidStartServer {
	return &InvalidStartServer{
		ConfigError: NewConfigError("Invalid start server", err),
	}
}

func (e *InvalidStartServer) Error() string {
	return e.ConfigError.Error()
}
