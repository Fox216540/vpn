package vpn

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"log"
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
	config := settings.Config
	script := config.ScriptName
	path := config.HomePath
	scriptPath := path + script
	cmd := exec.Command("sudo", "-E", "bash", "-c", scriptPath)

	cmd.Env = append(os.Environ(),
		"MENU_OPTION=1",
		"CLIENT="+configID.String(),
		"PASS=1",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// TODO: return кастомную ошибку
		fmt.Println("❌ Ошибка запуска:", err)
		return "", err
	}

	vpnPath := filepath.Join(settings.Config.ClientConfigPath, clientName+".ovpn")
	return vpnPath, nil
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
	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("93", err)
	}
	return lines, nil
}

func (r *Repository) cleanDeletedIDs() error {
	filePath := settings.Config.FilePath
	lines, err := r.scanFile(filePath)
	if err != nil {
		return fmt.Errorf("102", err)
	}
	return r.writeLines(lines, filePath)
}

func (r *Repository) deleteRSA(configID uuid.UUID) error {
	dir := "/etc/openvpn/easy-rsa"

	// 1. Проверяем существование easyrsa
	easyrsaPath := filepath.Join(dir, "easyrsa")
	if _, err := os.Stat(easyrsaPath); os.IsNotExist(err) {
		return fmt.Errorf("easyrsa not found in %s", dir)
	}

	// 2. Отзываем сертификат
	cmd := exec.Command("sudo", easyrsaPath, "--batch", "revoke", configID.String())
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("revoke failed: %w", err)
	}

	// 3. Генерируем CRL (Certificate Revocation List)
	//crlCmd := exec.Command(easyrsaPath, "--batch", "gen-crl")
	//crlCmd.Dir = dir
	//crlCmd.Stdout = os.Stdout
	//crlCmd.Stderr = os.Stderr
	//
	//if err := crlCmd.Run(); err != nil {
	//	return fmt.Errorf("CRL generation failed: %w", err)
	//}

	// 4. Копируем CRL в директорию OpenVPN (требует sudo)
	//copyCmd := exec.Command("sudo", "cp",
	//	filepath.Join(dir, "pki", "crl.pem"),
	//	"/etc/openvpn/server/crl.pem")
	//
	//if err := copyCmd.Run(); err != nil {
	//	return fmt.Errorf("CRL copy failed: %w", err)
	//}

	// 5. Удаляем файлы сертификата и ключа
	filesToRemove := []string{
		filepath.Join(dir, "pki", "issued", configID.String()+".crt"),
		filepath.Join(dir, "pki", "private", configID.String()+".key"),
		filepath.Join(dir, "pki", "reqs", configID.String()+".req"), // обычно есть и этот файл
	}

	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: could not remove %s: %v", file, err)
			// Не прерываем выполнение, если файл не найден
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
