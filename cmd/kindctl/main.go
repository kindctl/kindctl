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
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kindctl",
		Short: "kindctl is a CLI tool to manage local Kubernetes clusters using Kind",
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "kindctl.yaml", "Path to configuration file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Kind cluster and create a default config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewLogger(logLevel)
			return cluster.Initialize(log, configFile)
		},
	}

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

	rootCmd.AddCommand(initCmd, updateCmd, destroyCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
