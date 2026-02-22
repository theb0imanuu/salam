package monitor

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/theb0imanuu/salam/internal/models"
)

type CPUMonitor struct{}

func NewCPUMonitor() *CPUMonitor {
	return &CPUMonitor{}
}

func (m *CPUMonitor) Collect() (*models.CPUInfo, error) {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	cores, _ := cpu.Counts(true)

	loadStat, err := load.Avg()
	loadAvg := 0.0
	if err == nil {
		loadAvg = loadStat.Load1
	}

	usage := 0.0
	if len(percent) > 0 {
		usage = percent[0]
	}

	return &models.CPUInfo{
		Usage:       usage,
		Cores:       cores,
		LoadAverage: loadAvg,
	}, nil
}
