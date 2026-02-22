package monitor

import (
	"github.com/shirou/gopsutil/v3/net"
	"github.com/theb0imanuu/salam/internal/models"
)

type NetworkMonitor struct{}

func NewNetworkMonitor() *NetworkMonitor {
	return &NetworkMonitor{}
}

func (m *NetworkMonitor) Collect() ([]models.NetInfo, error) {
	io, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	var nets []models.NetInfo
	for _, iface := range io {
		// Skip loopback
		if iface.Name == "lo" || iface.Name == "Loopback Pseudo-Interface 1" {
			continue
		}
		nets = append(nets, models.NetInfo{
			Name:        iface.Name,
			BytesSent:   iface.BytesSent,
			BytesRecv:   iface.BytesRecv,
			PacketsSent: iface.PacketsSent,
			PacketsRecv: iface.PacketsRecv,
		})
	}

	return nets, nil
}
