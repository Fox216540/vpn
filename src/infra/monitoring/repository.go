package monitoring

import (
	"bufio"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"vpn/src/core/settings"
	"vpn/src/domain/monitoring"

	"github.com/shirou/gopsutil/v3/net"
)

type Repository struct {
	data monitoring.Monitoring
	mu   sync.RWMutex
}

func NewRepository() *Repository {
	r := &Repository{}
	go r.startWriterTraffic()
	go r.startWriterCPUPercent()
	go r.startWriterMemoryPercent()
	go r.startWriterGetActiveConnections()
	return r
}

func (r *Repository) searchTun(nameTun string) net.IOCountersStat {
	counters, _ := net.IOCounters(true)
	for _, c := range counters {
		if c.Name == nameTun {
			return c
		}
	}
	return net.IOCountersStat{}
}

func (r *Repository) startWriterTraffic() {
	var prevRx, prevTx uint64

	trafficIntervalStr := settings.Config.TrafficInterval

	trafficInterval, _ := strconv.Atoi(trafficIntervalStr)

	for {
		c := r.searchTun("tun0")

		rxDiff := c.BytesRecv - prevRx
		txDiff := c.BytesSent - prevTx

		prevRx = c.BytesRecv
		prevTx = c.BytesSent

		down := float64(rxDiff*8) / 1_000_000 // байты -> биты -> мегабиты
		up := float64(txDiff*8) / 1_000_000

		r.mu.Lock()
		r.data.MbpsDown = down
		r.data.MbpsUp = up
		r.mu.Unlock()
		time.Sleep(time.Duration(trafficInterval) * time.Second)
	}
}

func (r *Repository) GetTrafficUsage() (down, up float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data.MbpsDown, r.data.MbpsUp
}

func (r *Repository) startWriterCPUPercent() {
	intervalStr := settings.Config.CPUInterval
	interval, _ := strconv.Atoi(intervalStr)

	for {
		percent, _ := cpu.Percent(time.Duration(interval)*time.Second, false) // среднее за 1 секунду
		r.mu.Lock()
		r.data.CPUPercentUsage = percent[0]
		r.mu.Unlock()
	}
}

func (r *Repository) GetCPUPercentUsage() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data.CPUPercentUsage
}

func (r *Repository) startWriterMemoryPercent() {
	intervalStr := settings.Config.MemoryInterval

	interval, _ := strconv.Atoi(intervalStr)

	for {
		v, _ := mem.VirtualMemory()
		r.mu.Lock()
		r.data.MemoryPercentUsage = v.UsedPercent
		r.mu.Unlock()

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (r *Repository) GetMemoryPercentUsage() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data.MemoryPercentUsage
}

func (r *Repository) readFile(path string) int {
	file, _ := os.Open(path)
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "CLIENT_LIST") {
			count++
		}
	}
	return count
}

func (r *Repository) startWriterGetActiveConnections() {
	intervalStr := settings.Config.ConnectionsInterval

	interval, _ := strconv.Atoi(intervalStr)

	for {
		count := r.readFile("/run/openvpn-server/status-server.log")
		r.mu.Lock()
		r.data.ActiveConnections = count
		r.mu.Unlock()
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (r *Repository) GetActiveConnections() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data.ActiveConnections
}
