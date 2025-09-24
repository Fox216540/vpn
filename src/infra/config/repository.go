package config

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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
		return "", NewInvalidCreateConfig(err)
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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return NewInvalidDeleteConfig(err)
	}

	return nil
}

func (r *Repository) DeleteIDs(ids []uuid.UUID) error {
	rsaPath := "/etc/openvpn/server/easy-rsa/"
	pkiDir := rsaPath + "pki"
	workDir := fmt.Sprintf("/tmp/ovpn_revoke_%d", os.Getpid())

	clientStrings, clientSet := r.convertUUIDs(ids)

	if err := os.MkdirAll(workDir, 0755); err != nil {
		return fmt.Errorf("ошибка создания рабочей директории: %w", err)
	}
	defer os.RemoveAll(workDir)

	fmt.Printf("Запуск отзыва %d клиентов...\n", len(clientStrings))

	// Параллельный отзыв клиентов
	successCount, err := r.revokeClients(clientStrings, pkiDir)
	if err != nil {
		return fmt.Errorf("ошибки отзыва: %w", err)
	}

	// Финальные операции
	if err := r.finalizeCRL(clientSet, rsaPath); err != nil {
		return err
	}

	if err := r.disconnectClients(clientStrings); err != nil {
		return fmt.Errorf("ошибки отключения клиентов: %w", err)
	}

	fmt.Printf("✅ Успешно отозвано: %d клиентов\n", successCount)
	return nil
}

// Конвертация UUID → строки и set
func (r *Repository) convertUUIDs(ids []uuid.UUID) ([]string, map[string]struct{}) {
	clientStrings := make([]string, len(ids))
	clientSet := make(map[string]struct{}, len(ids))
	for i, id := range ids {
		s := id.String()
		clientStrings[i] = s
		clientSet[s] = struct{}{}
	}
	return clientStrings, clientSet
}

// Параллельный отзыв клиентов с использованием errgroup
func (r *Repository) revokeClients(clients []string, pkiDir string) (int, error) {
	clientCount := len(clients)
	if clientCount == 0 {
		return 0, nil
	}

	parallel := r.min(r.max(clientCount/3, 4), 16)
	sem := make(chan struct{}, parallel)

	g, ctx := errgroup.WithContext(context.Background())

	var successCount int32
	var mu sync.Mutex
	var errors []string

	for _, client := range clients {
		client := client

		// Ожидание свободного слота с проверкой контекста
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			return int(atomic.LoadInt32(&successCount)), ctx.Err()
		}

		g.Go(func() error {
			defer func() { <-sem }()

			if err := r.revokeClient(client, pkiDir); err != nil {
				mu.Lock()
				errors = append(errors, err.Error())
				mu.Unlock()
				return nil // Не возвращаем ошибку, чтобы дождаться всех горутин
			}

			atomic.AddInt32(&successCount, 1)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return int(atomic.LoadInt32(&successCount)), err
	}

	if len(errors) > 0 {
		return int(atomic.LoadInt32(&successCount)), fmt.Errorf(strings.Join(errors, "; "))
	}

	return int(atomic.LoadInt32(&successCount)), nil
}

func (r *Repository) disconnectClients(clients []string) error {
	clientCount := len(clients)
	if clientCount == 0 {
		return nil
	}

	parallel := r.min(r.max(clientCount/3, 4), 16)
	sem := make(chan struct{}, parallel)

	g, ctx := errgroup.WithContext(context.Background())
	var errors []string
	var mu sync.Mutex

	for _, client := range clients {
		client := client

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			return ctx.Err()
		}

		g.Go(func() error {
			defer func() { <-sem }()

			if err := r.disconnectClient(client); err != nil {
				mu.Lock()
				errors = append(errors, err.Error())
				mu.Unlock()
				return nil
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}

func (r *Repository) disconnectClient(client string) error {
	cmd := exec.Command("sh", "-c",
		fmt.Sprintf("echo 'kill %s\nexit' | nc -w 2 127.0.0.1 7505", client))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка отключения клиента %s: %w", client, err)
	}

	fmt.Printf("✓ Клиент %s отключен\n", client)
	return nil
}

// Упрощенные финальные операции
func (r *Repository) finalizeCRL(clientSet map[string]struct{}, rsaPath string) error {
	pkiDir := rsaPath + "/pki"

	if err := r.updateIndexFile(clientSet, pkiDir); err != nil {
		return err
	}

	if err := r.generateNewCRL(rsaPath); err != nil {
		return err
	}

	return nil
}

// Обновление index.txt - удаление записей отозванных клиентов
func (r *Repository) updateIndexFile(clientSet map[string]struct{}, pkiDir string) error {
	indexPath := filepath.Join(pkiDir, "index.txt")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения index.txt: %w", err)
	}

	if err := r.backupFile(indexPath, data); err != nil {
		return err
	}

	filteredLines := r.markRevokedByClient(data, clientSet)

	if err := r.writeFileWithBackup(indexPath, filteredLines, data); err != nil {
		return err
	}

	fmt.Println("✅ index.txt успешно обновлен")
	return nil
}

// Создание бэкапа файла
func (r *Repository) backupFile(path string, data []byte) error {
	backupPath := path + ".backup"
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("ошибка создания бэкапа: %w", err)
	}
	return nil
}

// Фильтрация строк index.txt по clientSet
func (r *Repository) markRevokedByClient(data []byte, clientSet map[string]struct{}) []string {
	lines := strings.Split(string(data), "\n")
	var result []string
	now := time.Now().UTC().Format("060102150405Z") // дата отзыва в формате YYMMDDHHMMSSZ

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			result = append(result, line)
			continue
		}

		cn := r.extractCN(line)
		if cn == "" || !r.contains(clientSet, cn) {
			// Если CN не найден или клиент не в списке — оставляем строку как есть
			result = append(result, line)
			continue
		}

		// CN есть и клиент в списке → ревокация
		serial := r.extractSerial(line)
		newLine := fmt.Sprintf("R\t%s\t%s\t%s", now, serial, cn)
		result = append(result, newLine)
	}

	return result
}

func (r *Repository) extractSerial(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return "" // нет serial
	}
	return fields[2]
}

// Извлечение CN из строки index.txt
func (r *Repository) extractCN(line string) string {
	if idx := strings.Index(line, "/CN="); idx != -1 {
		cnPart := line[idx+4:]
		return strings.Fields(cnPart)[0]
	}
	return ""
}

// Проверка наличия в set
func (r *Repository) contains(clientSet map[string]struct{}, cn string) bool {
	_, ok := clientSet[cn]
	return ok
}

// Запись с восстановлением из бэкапа при ошибке
func (r *Repository) writeFileWithBackup(path string, lines []string, backup []byte) error {
	newContent := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		os.WriteFile(path, backup, 0644)
		return fmt.Errorf("ошибка записи index.txt: %w", err)
	}
	return nil
}

// Генерация нового CRL
func (r *Repository) generateNewCRL(rsaPath string) error {
	// Генерируем CRL через easyrsa
	cmd := exec.Command("./easyrsa", "--batch", "--days=3650", "gen-crl")
	cmd.Dir = rsaPath

	pkiDir := rsaPath + "/pki/"

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ошибка выполнения easyrsa: %w, output: %s", err, string(output))
	}

	if err := exec.Command("rm", "-f", "/etc/openvpn/server/crl.pem").Run(); err != nil {
		return fmt.Errorf("ошибка удаления старого crl.pem: %w", err)
	}

	cmd = exec.Command("cp", fmt.Sprintf("%s/crl.pem", pkiDir), "/etc/openvpn/server/crl.pem")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ошибка копирования crl.pem: %w, output: %s", err, string(output))
	}

	if err := exec.Command("chown", "nobody:nogroup", "/etc/openvpn/server/crl.pem").Run(); err != nil {
		return fmt.Errorf("error chown")
	}

	fmt.Println("✅ CRL успешно обновлен")
	return nil
}

// Упрощенная функция отзыва
func (r *Repository) revokeClient(client string, pkiDir string) error {
	// Удаление файлов клиента
	files := []string{
		filepath.Join(pkiDir, "private", client+".key"),
		filepath.Join(pkiDir, "reqs", client+".req"),
		filepath.Join(pkiDir, "issued", client+".crt"),
		filepath.Join(pkiDir, "inline", client+".inline"),
		filepath.Join(pkiDir, "inline/private", client+".inline"),
	}

	var removed bool
	for _, file := range files {
		if err := os.Remove(file); err == nil {
			removed = true
		}
	}

	if !removed {
		return fmt.Errorf("файлы клиента не найдены")
	}

	fmt.Printf("✓ %s отозван\n", client)
	return nil
}

// Вспомогательные функции
func (r *Repository) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *Repository) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (r *Repository) CreateServer() error {
	config := settings.Config
	script := config.ScriptName
	path := config.HomePath
	scriptPath := path + script
	if _, err := os.Stat(script); os.IsNotExist(err) {
		// скачать скрипт
		cmd := exec.Command("curl", "-o", scriptPath, "https://raw.githubusercontent.com/Fox216540/openvpn-installer/main/openvpn-install.sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err = cmd.Run(); err != nil {
			fmt.Println(err)
			return NewInvalidDownloadScript(err)
		}

	}

	// chmod +x
	cmd := exec.Command("chmod", "+x", scriptPath)

	if err := cmd.Run(); err != nil {
		return NewInvalidChmodScript(err)
	}

	cmd = exec.Command("sudo", "-E", "bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	fmt.Println(os.Environ())
	fmt.Println(cmd.Env)

	if err := cmd.Run(); err != nil {
		return NewInvalidCreateServer(err)
	}

	return nil
}
