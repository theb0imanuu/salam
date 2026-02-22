package reporter

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/theb0imanuu/salam/internal/models"
)

func NewSpinner(msg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + msg
	return s
}

func PrintHeader(title string) {
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n╔══════════════════════════════════════════╗\n")
	cyan.Printf("║  %-40s║\n", title)
	cyan.Printf("╚══════════════════════════════════════════╝\n\n")
}

func PrintMeta(hostname, platform, timestamp string) {
	gray := color.New(color.FgHiBlack)
	gray.Printf("Host: %s | %s\n", hostname, platform)
	gray.Printf("Time: %s\n\n", timestamp)
}

func PrintCPU(info *models.CPUInfo, threshold int) {
	title := color.New(color.Bold).Sprint("CPU Usage:")

	usageColor := color.GreenString
	if info.Usage > float64(threshold) {
		usageColor = color.RedString
	}

	fmt.Printf("%s\n", title)
	fmt.Printf("  %s across %d cores\n", usageColor("%.1f%%", info.Usage), info.Cores)
	fmt.Printf("  Load Average: %.2f\n\n", info.LoadAverage)
}

func PrintMemory(info *models.MemoryInfo, threshold int) {
	title := color.New(color.Bold).Sprint("Memory:")

	usageColor := color.GreenString
	if info.Usage > float64(threshold) {
		usageColor = color.RedString
	}

	fmt.Printf("%s\n", title)
	fmt.Printf("  Total: %s\n", info.Total)
	fmt.Printf("  Used:  %s %s\n", info.Used, usageColor("(%.1f%%)", info.Usage))
	fmt.Printf("  Free:  %s\n\n", info.Free)
}

func PrintDisk(disks []models.DiskInfo, threshold int) {
	title := color.New(color.Bold).Sprint("Disk Usage:")
	fmt.Printf("%s\n", title)

	for _, d := range disks {
		usageColor := color.GreenString
		if d.Usage > float64(threshold) {
			usageColor = color.RedString
		}
		fmt.Printf("  %-12s %8s / %-8s %s\n",
			d.Path, d.Used, d.Total, usageColor("(%.1f%%)", d.Usage))
	}
	fmt.Println()
}

func PrintNetwork(interfaces []models.NetInfo) {
	title := color.New(color.Bold).Sprint("Network:")
	fmt.Printf("%s\n", title)

	for _, net := range interfaces {
		fmt.Printf("  %s:\n", net.Name)
		fmt.Printf("    ↓ %s | ↑ %s\n",
			models.FormatBytes(net.BytesRecv), models.FormatBytes(net.BytesSent))
	}
	fmt.Println()
}

func PrintAlerts(alerts []string) {
	red := color.New(color.FgRed, color.Bold)
	red.Println("⚠ ALERTS DETECTED")
	for _, alert := range alerts {
		red.Printf("  • %s\n", alert)
	}
	fmt.Println()
}

func PrintSuccess(msg string) {
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("✓ %s\n\n", msg)
}
