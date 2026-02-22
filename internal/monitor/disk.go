package monitor

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/theb0imanuu/salam/internal/models"
)

type DiskMonitor struct{}

func NewDiskMonitor() *DiskMonitor {
	return &DiskMonitor{}
}

func (m *DiskMonitor) Collect() ([]models.DiskInfo, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []models.DiskInfo
	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}

		disks = append(disks, models.DiskInfo{
			Path:   part.Mountpoint,
			Total:  models.FormatBytes(usage.Total),
			Used:   models.FormatBytes(usage.Used),
			Free:   models.FormatBytes(usage.Free),
			Usage:  usage.UsedPercent,
			FSType: part.Fstype,
		})
	}

	return disks, nil
}
