package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/theb0imanuu/salam/internal/config"
	"github.com/theb0imanuu/salam/internal/models"
	"github.com/theb0imanuu/salam/internal/reporter"
)

type Options struct {
	CPUOnly     bool
	MemoryOnly  bool
	DiskOnly    bool
	NetworkOnly bool
	JSONOutput  bool
	Interval    int
	Threshold   int
	WebhookURL  string
	Config      *config.Config
}

type Monitor struct {
	opts   Options
	cpu    *CPUMonitor
	memory *MemoryMonitor
	disk   *DiskMonitor
	net    *NetworkMonitor
}

func New(opts Options) *Monitor {
	return &Monitor{
		opts:   opts,
		cpu:    NewCPUMonitor(),
		memory: NewMemoryMonitor(),
		disk:   NewDiskMonitor(),
		net:    NewNetworkMonitor(),
	}
}

func (m *Monitor) RunOnce() error {
	spinner := reporter.NewSpinner("Collecting system health data...")
	spinner.Start()

	data, err := m.collect()
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to collect metrics: %w", err)
	}

	if m.opts.JSONOutput {
		return m.outputJSON(data)
	}

	m.outputPretty(data)

	if m.opts.WebhookURL != "" {
		if err := reporter.SendWebhook(m.opts.WebhookURL, data); err != nil {
			color.Yellow("⚠ Webhook failed: %v", err)
		}
	}

	if !data.Healthy {
		os.Exit(1)
	}
	return nil
}

func (m *Monitor) StartWatching() error {
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n🔍 Salam Watch Mode Started (interval: %ds, threshold: %d%%)\n\n",
		m.opts.Interval, m.opts.Threshold)

	// First run
	if err := m.RunOnce(); err != nil {
		return err
	}

	// Setup ticker
	ticker := time.NewTicker(time.Duration(m.opts.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Print("\033[H\033[2J") // Clear screen
		if err := m.RunOnce(); err != nil {
			color.Red("Error: %v", err)
		}
	}

	return nil
}

func (m *Monitor) collect() (*models.HealthData, error) {
	hostname, _ := os.Hostname()

	data := &models.HealthData{
		Timestamp: time.Now().Format(time.RFC3339),
		Hostname:  hostname,
		Platform:  fmt.Sprintf("%s %s", os.Getenv("GOOS"), os.Getenv("GOARCH")),
		Healthy:   true,
		Alerts:    []string{},
	}

	// Collect based on flags or all if none specified
	collectAll := !m.opts.CPUOnly && !m.opts.MemoryOnly && !m.opts.DiskOnly && !m.opts.NetworkOnly

	if collectAll || m.opts.CPUOnly {
		cpu, err := m.cpu.Collect()
		if err != nil {
			return nil, err
		}
		data.CPU = cpu
		if cpu.Usage > float64(m.opts.Config.Thresholds.CPU) {
			data.Alerts = append(data.Alerts,
				fmt.Sprintf("CPU usage %.1f%% exceeds threshold %d%%", cpu.Usage, m.opts.Config.Thresholds.CPU))
			data.Healthy = false
		}
	}

	if collectAll || m.opts.MemoryOnly {
		mem, err := m.memory.Collect()
		if err != nil {
			return nil, err
		}
		data.Memory = mem
		if mem.Usage > float64(m.opts.Config.Thresholds.Memory) {
			data.Alerts = append(data.Alerts,
				fmt.Sprintf("Memory usage %.1f%% exceeds threshold %d%%", mem.Usage, m.opts.Config.Thresholds.Memory))
			data.Healthy = false
		}
	}

	if collectAll || m.opts.DiskOnly {
		disks, err := m.disk.Collect()
		if err != nil {
			return nil, err
		}
		data.Disk = disks
		for _, d := range disks {
			if d.Usage > float64(m.opts.Config.Thresholds.Disk) {
				data.Alerts = append(data.Alerts,
					fmt.Sprintf("Disk %s usage %.1f%% exceeds threshold %d%%", d.Path, d.Usage, m.opts.Config.Thresholds.Disk))
				data.Healthy = false
			}
		}
	}

	if collectAll || m.opts.NetworkOnly {
		net, err := m.net.Collect()
		if err != nil {
			return nil, err
		}
		data.Network = net
	}

	return data, nil
}

func (m *Monitor) outputJSON(data *models.HealthData) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (m *Monitor) outputPretty(data *models.HealthData) {
	reporter.PrintHeader("SALAM HEALTH REPORT")
	reporter.PrintMeta(data.Hostname, data.Platform, data.Timestamp)

	if data.CPU != nil {
		reporter.PrintCPU(data.CPU, m.opts.Config.Thresholds.CPU)
	}
	if data.Memory != nil {
		reporter.PrintMemory(data.Memory, m.opts.Config.Thresholds.Memory)
	}
	if len(data.Disk) > 0 {
		reporter.PrintDisk(data.Disk, m.opts.Config.Thresholds.Disk)
	}
	if len(data.Network) > 0 {
		reporter.PrintNetwork(data.Network)
	}

	if len(data.Alerts) > 0 {
		reporter.PrintAlerts(data.Alerts)
	} else {
		reporter.PrintSuccess("All systems healthy")
	}
}
