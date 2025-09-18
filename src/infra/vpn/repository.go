package vpn

import (
	"bufio"
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
	clientName := configID.String()
	cmd := exec.Command("./openvpn-install.sh")

	cmd.Env = append(os.Environ(),
		"MENU_OPTION=1",
		"CLIENT="+configID.String(),
		"PASS=1",
	)

	if err := cmd.Run(); err != nil {
		// TODO: return кастомную ошибку
		fmt.Println("❌ Ошибка запуска:", err)
		return "", err
	}

	vpnPath := filepath.Join(settings.Config.ClientConfigPath, clientName+".ovpn")
	return vpnPath, nil
}

func (r *Repository) openTempFile(filePath string) (*os.File, string, error) {
	tmpPath := filePath + ".tmp"
	f, err := os.Create(tmpPath)
	return f, tmpPath, err
}

func (r *Repository) writeLines(lines []string, filePath string) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, line := range lines {
		if _, err := outFile.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) scanFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "R") {
			lines = append(lines, line)
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func (r *Repository) cleanDeletedIDs() error {
	filePath := settings.Config.FilePath
	lines, err := r.scanFile(filePath)
	if err != nil {
		return err
	}
	return r.writeLines(lines, filePath)
}

func (r *Repository) deleteRSA(configID uuid.UUID) error {
	dir := "/etc/openvpn/easy-rsa"

	// безопасный вызов easyrsa
	cmd := exec.Command(filepath.Join(dir, "easyrsa"), "--batch", "revoke", configID.String())
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command(filepath.Join(dir, "easyrsa"), "gen-crl")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	// удаляем ключ и crt
	files := []string{
		filepath.Join(dir, "pki/issued", configID.String()+".crt"),
		filepath.Join(dir, "pki/private", configID.String()+".key"),
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return err
		}
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
		cmd := exec.Command("curl", "-o", scriptPath, "https://raw.githubusercontent.com/angristan/openvpn-install/master/openvpn-install.sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("failed to download script: %w", err)
		}
	}

	// chmod +x
	cmd := exec.Command("chmod", "+x", scriptPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed chmod: %w", err)
	}

	cmd = exec.Command("sudo", "-E", "bash", "-c", scriptPath)

	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run script start: %w", err)
	}

	return nil
}
