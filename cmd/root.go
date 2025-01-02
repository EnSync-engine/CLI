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

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "ensync",
		Short: "EnSync CLI tool",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				zap.L().Fatal("Failed to load configuration", zap.Error(err))
			}

			var logLevel zapcore.Level
			if debug || cfg.Debug {
				logLevel = zapcore.DebugLevel
			} else {
				logLevel = zapcore.InfoLevel
			}

			logger := newLogger(logLevel)
			zap.ReplaceGlobals(logger)
		},
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	client := api.NewClient(
		cfg.BaseURL,
		cfg.APIKey,
		api.WithLogger(zap.L()),
		api.WithRateLimit(10, 20),
	)

	rootCmd.AddCommand(
		newEventCmd(client),
		newAccessKeyCmd(client),
		newVersionCmd(),
	)

	return rootCmd.Execute()
}

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
