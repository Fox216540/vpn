package config

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"vpn/src/core/settings"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(configID uuid.UUID) (string, error) {
	config := settings.Config
	script := config.ScriptName
	path := config.HomePath
	scriptPath := path + script

	cmd := exec.Command("sudo", "-E", "bash", scriptPath)
	cmd.Env = append(os.Environ(),
		"CLIENT_NAME="+configID.String(),
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error create config")
	}

	pathFile := path + configID.String() + ".ovpn"

	return pathFile, nil
}

func (r *Repository) Delete(configID uuid.UUID) error {
	config := settings.Config
	script := config.ScriptName
	path := config.HomePath
	scriptPath := path + script

	cmd := exec.Command("sudo", "-E", "bash", scriptPath)
	cmd.Env = append(os.Environ(),
		"OPTION=2",
		"CLIENT_NAME="+configID.String(),
	)

	return nil
}

func (r *Repository) CreateServer() error {
	config := settings.Config
	script := config.ScriptName
	path := config.HomePath
	scriptPath := path + script
	if _, err := os.Stat(script); os.IsNotExist(err) {
		// скачать скрипт
		cmd := exec.Command("curl", "-o", scriptPath, "https://raw.githubusercontent.com/Fox216540/openvpn-installer/main/openvpn-install.sh")
		if err = cmd.Run(); err != nil {
			fmt.Println("Error creating server")
			return fmt.Errorf("failed to download script: %w", err)
		}
	}

	// chmod +x
	cmd := exec.Command("chmod", "+x", scriptPath)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error creating server 71")
		return fmt.Errorf("failed chmod: %w", err)
	}

	cmd = exec.Command("sudo", "-E", "bash", scriptPath)
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		fmt.Println("Error creating server 79")
		return fmt.Errorf("failed create server: %w", err)
	}

	return nil
}
