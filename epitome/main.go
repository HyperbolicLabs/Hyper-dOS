package main

import (
	"flag"
	"time"

	"epitome.hyperbolic.xyz/cluster"
	"epitome.hyperbolic.xyz/config"
	"epitome.hyperbolic.xyz/helper"
	"epitome.hyperbolic.xyz/hyperweb"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const VERSION = "v1-alpha"

func main() {
	help := flag.Bool("help", false, "Show help")
	mode := flag.String("mode", "default", "Specify the mode to run.")
	flag.Parse()

	var cfg config.Config
	env.Parse(&cfg)
	helper.SetLogLevel(cfg.LOG_LEVEL)

	logCfg := zap.NewProductionConfig()
	switch cfg.LOG_LEVEL {
	case "debug":
		logCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		logCfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		logCfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		logCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	if *help {
		flag.PrintDefaults()
		return
	}

	logCfg.EncoderConfig.TimeKey = "timestamp"
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := logCfg.Build()

	logger = logger.With(zap.String("mode", *mode))
	logger.Info("launching epitome")

	clientset, dynamicClient := cluster.MustConnect(cfg.KUBECONFIG)
	switch *mode {
	case "default":
		hyperweb.RunLoop(
			cfg,
			logger,
			clientset,
			dynamicClient,
			60*time.Second,
		)
		logger.Fatal("hyperweb runloop exited unexpectedly")
	default:
		logger.Fatal("unknown mode", zap.String("mode", *mode))
	}
}
