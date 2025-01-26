package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/EnSync-engine/CLI/app/api"
	"github.com/EnSync-engine/CLI/app/config"
)

var (
	cfgFile string
	debug   bool
)

// Execute is the entry point for the CLI application.
func Execute() error {
	rootCmd := setupRootCommand()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize the API client
	client := api.NewClient(
		cfg.BaseURL,
		api.WithLogger(zap.L()),
		api.WithRateLimit(10, 20),
	)

	// Register subcommands
	rootCmd.AddCommand(
		newEventCmd(client),
		newAccessKeyCmd(client),
		newVersionCmd(),
	)

	// Execute the root command
	return rootCmd.Execute()
}

// setupRootCommand configures and returns the root Cobra command.
func setupRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ensync",
		Short: "EnSync CLI tool",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Load configuration
			cfg, err := config.Load()
			if err != nil {
				zap.L().Fatal("Failed to load configuration", zap.Error(err))
			}
			// Set up logging level based on debug flag or configuration
			logLevel := determineLogLevel(cfg)
			logger := newLogger(logLevel)
			zap.ReplaceGlobals(logger)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ensync.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")

	return rootCmd
}

// determineLogLevel determines the appropriate log level based on the debug flag and configuration.
func determineLogLevel(cfg *config.Config) zapcore.Level {
	if debug || cfg.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

// newLogger creates and returns a new Zap logger with the specified log level.
func newLogger(level zapcore.Level) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "level"
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.MessageKey = "message"

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		level,
	)

	return zap.New(core)
}
