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
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger := initLogger(cfg)
	zap.ReplaceGlobals(logger)

	client := api.NewClient(
		cfg.BaseURL,
		api.WithLogger(logger),
		api.WithRateLimit(10, 20),
	)

	rootCmd := newRootCmd()
	rootCmd.AddCommand(
		newEventCmd(client),
		newAccessKeyCmd(client),
		newWorkspaceCmd(client),
		newVersionCmd(),
	)

	return rootCmd.Execute()
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ensync",
		Short: "EnSync CLI tool for managing events and access keys",
		Long: `EnSync CLI provides commands for managing events and access keys
in the EnSync real-time messaging system.

Use --access-key flag with subcommands to authenticate API requests.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ensync/config.yaml)")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")

	return cmd
}

func initLogger(cfg *config.Config) *zap.Logger {
	level := zapcore.InfoLevel
	if debug || cfg.Debug {
		level = zapcore.DebugLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "level"
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.MessageKey = "message"

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)

	return zap.New(core)
}
