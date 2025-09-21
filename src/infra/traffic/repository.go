package traffic

import (
	"log"
	"strconv"
	"sync"
	"time"
	"vpn/src/core/settings"
	"vpn/src/domain/traffic"

	"github.com/shirou/gopsutil/v3/net"
)

type Repository struct {
	data traffic.Traffic
	mu   sync.RWMutex
}

func NewRepository() *Repository {
	r := &Repository{}
	go r.startWriter()
	return r
}

func (r *Repository) startWriter() {
	var prevRx, prevTx uint64

	trafficIntervalStr := settings.Config.TrafficInterval

	trafficInterval, err := strconv.Atoi(trafficIntervalStr)
	if err != nil {
		log.Fatalf("Error converting traffic interval: %v", err)
	}

	for {
		counters, err := net.IOCounters(true)
		if err != nil {
		}

		for _, c := range counters {
			if c.Name == "tun0" {
				rxDiff := c.BytesRecv - prevRx
				txDiff := c.BytesSent - prevTx

				prevRx = c.BytesRecv
				prevTx = c.BytesSent

				down := float64(rxDiff*8) / 1_000_000 // байты -> биты -> мегабиты
				up := float64(txDiff*8) / 1_000_000

				r.mu.Lock()
				r.data.Down = down
				r.data.Up = up
				r.mu.Unlock()
			}
		}
		time.Sleep(time.Duration(trafficInterval) * time.Second)
	}
}

func (r *Repository) GetTraffic() (down, up float64, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data.Down, r.data.Up, nil
}
