package main

import (
	"fmt"
	"kindctl/internal/config"
	"os"

	"github.com/spf13/cobra"
	"kindctl/internal/cluster"
	"kindctl/internal/logger"
	"kindctl/internal/tools"
)

var (
	configFile  string
	logLevel    string
	version     = "dev"
	showVersion bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kindctl",
		Short: "kindctl is a CLI tool to manage local Kubernetes clusters using Kind",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("kindctl version %s\n", version)
				os.Exit(0)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "kindctl.yaml", "Path to configuration file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Print the version of kindctl")

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

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of kindctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kindctl version: %s\n", version)
		},
	}

	rootCmd.AddCommand(initCmd, updateCmd, destroyCmd, versionCmd)
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("kindctl version: %s\n", version)
			os.Exit(0)
		}
	}
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
