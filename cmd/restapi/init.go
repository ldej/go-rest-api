package main

import (
	"fmt"
	"os"
	"strings"

	config "github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// Load config defaults by importing package
	_ "github.com/ldej/go-rest-example/conf"
)

var (
	logger *zap.SugaredLogger
)

func init() {
	initConfig()
	initLogger()
}

func initConfig() {
	// Sets up the config file, environment etc

	// If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	config.SetTypeByDefaultValue(true)
	// Automatically use environment variables where available
	config.AutomaticEnv()
	// Environment variables use underscores instead of periods
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initLogger() {
	logConfig := zap.NewProductionConfig()

	// Log Level
	var logLevel zapcore.Level
	if err := logLevel.Set(config.GetString("logger.level")); err != nil {
		zap.S().Fatalw("Could not determine logger.level", "error", err)
	}
	logConfig.Level.SetLevel(logLevel)

	// Settings
	logConfig.Encoding = config.GetString("logger.encoding")
	logConfig.Development = config.GetBool("logger.dev_mode")

	// Enable Color
	if config.GetBool("logger.color") {
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Use sane timestamp when logging to console
	if logConfig.Encoding == "console" {
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// JSON Fields
	logConfig.EncoderConfig.MessageKey = "msg"
	logConfig.EncoderConfig.LevelKey = "level"
	logConfig.EncoderConfig.CallerKey = "caller"

	// Build the logger
	globalLogger, err := logConfig.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "whelp: %v\n", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(globalLogger)
	logger = globalLogger.Sugar().With("package", "cmd")
}
