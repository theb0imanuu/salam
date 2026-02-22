package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/theb0imanuu/salam/internal/config"
	"github.com/theb0imanuu/salam/internal/monitor"
)

var (
	version = "1.0.0"
	cfgFile string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "salam",
		Short: "A comprehensive server health monitoring CLI",
		Long: `
╔══════════════════════════════════════════╗
║   SALAM - Server Health Monitor          ║
║   Peace of mind for your infrastructure  ║
╚══════════════════════════════════════════╝`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(cmd, args)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.salam.yaml)")

	// Check command
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Run a one-time health check",
		RunE:  runCheck,
	}
	checkCmd.Flags().Bool("cpu", false, "Check CPU only")
	checkCmd.Flags().Bool("memory", false, "Check memory only")
	checkCmd.Flags().Bool("disk", false, "Check disk only")
	checkCmd.Flags().Bool("network", false, "Check network only")
	checkCmd.Flags().Bool("json", false, "Output as JSON")
	checkCmd.Flags().String("webhook", "", "Send results to webhook URL")

	// Add flags to rootCmd as well to support default execution
	rootCmd.Flags().AddFlagSet(checkCmd.Flags())

	// Watch command
	watchCmd := &cobra.Command{
		Use:   "watch",
		Short: "Continuously monitor server health",
		RunE:  runWatch,
	}
	watchCmd.Flags().Int("interval", 30, "Check interval in seconds")
	watchCmd.Flags().Int("threshold", 80, "Alert threshold percentage")
	watchCmd.Flags().String("webhook", "", "Send alerts to webhook URL")

	// Dashboard command
	dashboardCmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Launch interactive dashboard",
		RunE:  runDashboard,
	}

	// Config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Generate default configuration file",
		RunE:  runConfig,
	}

	rootCmd.AddCommand(checkCmd, watchCmd, dashboardCmd, configCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCheck(cmd *cobra.Command, args []string) error {
	cfg := config.Load(cfgFile)

	opts := monitor.Options{
		CPUOnly:     mustGetBool(cmd, "cpu"),
		MemoryOnly:  mustGetBool(cmd, "memory"),
		DiskOnly:    mustGetBool(cmd, "disk"),
		NetworkOnly: mustGetBool(cmd, "network"),
		JSONOutput:  mustGetBool(cmd, "json"),
		WebhookURL:  mustGetString(cmd, "webhook"),
		Config:      cfg,
	}

	m := monitor.New(opts)
	return m.RunOnce()
}

func runWatch(cmd *cobra.Command, args []string) error {
	cfg := config.Load(cfgFile)

	opts := monitor.Options{
		Interval:   mustGetInt(cmd, "interval"),
		Threshold:  mustGetInt(cmd, "threshold"),
		WebhookURL: mustGetString(cmd, "webhook"),
		Config:     cfg,
	}

	m := monitor.New(opts)
	return m.StartWatching()
}

func runDashboard(cmd *cobra.Command, args []string) error {
	fmt.Println("Dashboard mode - coming in v2.0!")
	fmt.Println("For now, use: salam watch")
	return nil
}

func runConfig(cmd *cobra.Command, args []string) error {
	return config.Generate()
}

func mustGetBool(cmd *cobra.Command, name string) bool {
	v, _ := cmd.Flags().GetBool(name)
	return v
}

func mustGetString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

func mustGetInt(cmd *cobra.Command, name string) int {
	v, _ := cmd.Flags().GetInt(name)
	return v
}
