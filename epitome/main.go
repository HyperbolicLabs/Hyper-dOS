package main

import (
	"flag"

	"epitome.hyperbolic.xyz/cluster"
	"epitome.hyperbolic.xyz/config"
	"epitome.hyperbolic.xyz/helper"
	"epitome.hyperbolic.xyz/hyperweb"
	"epitome.hyperbolic.xyz/mode/maintain"
	"epitome.hyperbolic.xyz/mode/monkey"
	env11 "github.com/caarlos0/env/v11"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const VERSION = "v1-alpha"

func main() {
	help := flag.Bool("help", false, "Show help")
	mode := flag.String("mode", "default", "Specify the mode to run.")
	flag.Parse()

	var cfg config.Config
	env11.Parse(&cfg)
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
		err := hyperweb.RunLoop(
			cfg,
			logger,
			clientset,
			dynamicClient,
		)
		logger.Fatal("hyperweb runloop exited unexpectedly", zap.Error(err))
	case "maintain":
		err := maintain.Run(logger, &cfg, &clientset)
		logger.Fatal("maintain runloop exited unexpectedly", zap.Error(err))
	case "monkey":
		err := monkey.Run(cfg, logger)
		logger.Fatal("monkey runloop exited unexpectedly", zap.Error(err))
	default:
		logger.Fatal("unknown mode", zap.String("mode", *mode))
	}
}
