package monitor

import (
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/theb0imanuu/salam/internal/models"
)

type MemoryMonitor struct{}

func NewMemoryMonitor() *MemoryMonitor {
	return &MemoryMonitor{}
}

func (m *MemoryMonitor) Collect() (*models.MemoryInfo, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &models.MemoryInfo{
		Total: models.FormatBytes(vm.Total),
		Used:  models.FormatBytes(vm.Used),
		Free:  models.FormatBytes(vm.Free),
		Usage: vm.UsedPercent,
	}, nil
}
