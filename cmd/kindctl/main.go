package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"kindctl/internal/cluster"
	"kindctl/internal/config"
	"kindctl/internal/logger"
	"kindctl/internal/tools"
)

var (
	configFile string
	logLevel   string
	version    = "dev"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kindctl",
		Short: "kindctl is a CLI tool to manage local Kubernetes clusters using Kind",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Global --version flag support
			for _, arg := range args {
				if arg == "--version" || arg == "-v" {
					fmt.Println("kindctl version", version)
					os.Exit(0)
				}
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "kindctl.yaml", "Path to configuration file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")

	// Init command
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Kind cluster and create a default config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewLogger(logLevel)
			return cluster.Initialize(log, configFile)
		},
	}

	// Update command
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update the Kind cluster with tools specified in the config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewLogger(logLevel)
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			return tools.UpdateCluster(log, cfg)
		},
	}

	// Destroy command
	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Delete the Kind cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewLogger(logLevel)
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			return cluster.Destroy(log, cfg.Cluster.Name)
		},
	}

	// Version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of kindctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("kindctl version", version)
		},
	}

	// Add all commands
	rootCmd.AddCommand(initCmd, updateCmd, destroyCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
