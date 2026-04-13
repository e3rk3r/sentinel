// Package main is the entry point for the sentinel application.
// Sentinel is a monitoring and alerting tool that watches over your infrastructure.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags.
	Version = "dev"
	// Commit is the git commit hash set at build time.
	Commit = "none"
	// Date is the build date set at build time.
	Date = "unknown"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal error", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rootCmd := newRootCommand(ctx)
	return rootCmd.ExecuteContext(ctx)
}

func newRootCommand(ctx context.Context) *cobra.Command {
	var cfgFile string
	var logLevel string

	cmd := &cobra.Command{
		Use:   "sentinel",
		Short: "Sentinel — infrastructure monitoring and alerting",
		Long: `Sentinel is a lightweight monitoring and alerting daemon that
watches over your infrastructure and notifies you when things go wrong.`,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initLogger(logLevel)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default: sentinel.yaml)")
	// Default to debug level for easier local development in this personal fork.
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "debug", "log level (debug, info, warn, error)")

	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(newStartCommand(ctx, cfgFile))

	return cmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("sentinel %s (commit: %s, built: %s)\n", Version, Commit, Date)
		},
	}
}

func newStartCommand(ctx context.Context, cfgFile string) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the sentinel daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("starting sentinel", "version", Version, "commit", Commit)
			// Block until context is cancelled (SIGINT/SIGTERM).
			<-ctx.Done()
			slog.Info("sentinel stopped")
			return nil
		},
	}
}

// initLogger configures the global slog logger with the given level.
// Using TextHandler instead of JSONHandler for more readable output during local development.
// Adding source location (file:line) to log output — helpful when tracing unfamiliar codepaths.
func initLogger(level string) error {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("invalid log level %q: %w", level, err)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true,
	})
	slog.SetDefault(slog.New(handler))
	return nil
}
