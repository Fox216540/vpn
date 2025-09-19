package config

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	cmd := exec.Command("sudo", "bash", scriptPath)
	cmd.Env = append(os.Environ(),
		"CLIENT_NAME="+configID.String(),
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error create config")
	}

	pathFile := path + configID.String() + ".ovpn"

	return pathFile, nil
}

func (r *Repository) writeLines(lines []string, filePath string) error {
	config := settings.Config
	fileName := config.FileName
	tmpFile, err := os.CreateTemp("/tmp", fileName+"-*.txt")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	for _, line := range lines {
		if _, err = tmpFile.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("58", err)
		}
	}
	tmpFile.Close()

	cmd := exec.Command("sudo", "cp", tmpPath, filePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("sudo copy failed: %w", err)
	}

	return nil
}

func (r *Repository) scanFile(filePath string) ([]string, error) {
	cmd := exec.Command("sudo", "cat", filePath)
	output, err := cmd.Output()
	fmt.Println(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to read file with sudo: %w", err)
	}
	var lines []string

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "R") {
			lines = append(lines, line)
		}
	}
	fmt.Println(lines)
	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("99", err)
	}
	return lines, nil
}

func (r *Repository) cleanDeletedIDs() error {
	filePath := settings.Config.FilePath
	lines, err := r.scanFile(filePath)
	if err != nil {
		return fmt.Errorf("108", err)
	}
	return r.writeLines(lines, filePath)
}

func (r *Repository) recallCertificate(dir, path string, ID uuid.UUID) error {
	cmd := exec.Command("sudo", path, "--batch", "revoke", ID.String())
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("revoke failed: %w", err)
	}

	return nil
}

func (r *Repository) generateCertificate(dir, path string) error {
	crlCmd := exec.Command("sudo", path, "--batch", "--days=3650", "gen-crl")
	crlCmd.Dir = dir
	crlCmd.Stdout = os.Stdout
	crlCmd.Stderr = os.Stderr

	if err := crlCmd.Run(); err != nil {
		return fmt.Errorf("CRL generation failed: %w", err)
	}
	return nil
}

func (r *Repository) copyCRL(dir string) error {
	copyCmd := exec.Command("sudo", "cp",
		filepath.Join(dir, "pki", "crl.pem"),
		"/etc/openvpn/server/crl.pem")

	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("CRL copy failed: %w", err)
	}
	return nil
}

func (r *Repository) removeCertFiles(dir string, ID uuid.UUID) error {
	files := []string{
		filepath.Join(dir, "pki", "issued", ID.String()+".crt"),
		filepath.Join(dir, "pki", "private", ID.String()+".key"),
		filepath.Join(dir, "pki", "reqs", ID.String()+".req"),
	}

	for _, f := range files {
		err := os.Remove(f) // игнорируем ошибки, можно логировать
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) deleteRSA(configID uuid.UUID) error {
	dir := "/etc/openvpn/server/easy-rsa"

	// 1. Проверяем существование easyrsa
	easyrsaPath := filepath.Join(dir, "easyrsa")

	//2. Отзываем сертификат
	if err := r.recallCertificate(dir, easyrsaPath, configID); err != nil {
		return err
	}

	//3. Генерируем CRL (Certificate Revocation List)
	if err := r.generateCertificate(dir, easyrsaPath); err != nil {
		return err
	}

	//4. Копируем CRL в директорию OpenVPN (требует sudo)
	if err := r.copyCRL(dir); err != nil {
		return err
	}

	// 5. Удаляем файлы сертификата и ключа
	if err := r.removeCertFiles(dir, configID); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(configID uuid.UUID) error {
	err := r.deleteRSA(configID)

	if err != nil {
		fmt.Println(err)
	}

	err = r.cleanDeletedIDs()

	if err != nil {
		fmt.Println(err)
	}

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
			return fmt.Errorf("failed to download script: %w", err)
		}
	}

	// chmod +x
	cmd := exec.Command("chmod", "+x", scriptPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed chmod: %w", err)
	}

	cmd = exec.Command("sudo", "bash", scriptPath)
	cmd.Env = os.Environ()
	fmt.Println(cmd.Env)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed create server: %w", err)
	}

	return nil
}
