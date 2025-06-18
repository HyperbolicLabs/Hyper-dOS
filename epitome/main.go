package main

import (
	"flag"

	"epitome.hyperbolic.xyz/cluster"
	"epitome.hyperbolic.xyz/config"
	"epitome.hyperbolic.xyz/helper"
	"epitome.hyperbolic.xyz/mode/jungle"
	"epitome.hyperbolic.xyz/mode/maintain"
	"epitome.hyperbolic.xyz/mode/monkey"
	"epitome.hyperbolic.xyz/mode/sh"
	"k8s.io/client-go/tools/clientcmd"

	// Added new mode
	env11 "github.com/caarlos0/env/v11"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	help := flag.Bool("help", false, "Show help")
	mode := flag.String("mode", "jungle", "Specify the mode to run epitome in (jungle | maintain | monkey | sh)")
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

	// set up logger
	logCfg.EncoderConfig.TimeKey = "timestamp"
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := logCfg.Build()
	logger = logger.With(zap.String("mode", *mode))
	logger.Info("launching epitome")

	if cfg.HYPERDOS_VERSION != nil {
		logger.Info("hyperdos version", zap.String("version", *cfg.HYPERDOS_VERSION))
	}

	var kubeconfigPath *string
	if cfg.KUBECONFIG != "" {
		kubeconfigPath = &cfg.KUBECONFIG
	}

	// check if no role is set
	noroles := config.JungleRole{}
	if cfg.Role == noroles && *mode != "sh" {
		// TODO ronin functionality
		logger.Fatal("no jungle role set")
	}

	switch *mode {
	case "jungle":
		clientset, dynamicClient, _ := cluster.MustConnect(kubeconfigPath)
		err := jungle.Run(
			cfg,
			logger,
			clientset,
			dynamicClient,
		)
		logger.Fatal("hyperweb runloop exited unexpectedly", zap.Error(err))

	case "maintain":
		clientset, dynamicClient, argoClient := cluster.MustConnect(kubeconfigPath)
		err := maintain.Run(
			cfg,
			logger,
			clientset,
			dynamicClient,
			argoClient,
		)
		logger.Fatal("maintain runloop exited unexpectedly", zap.Error(err))

	case "monkey":
		clientset, dynamicClient, _ := cluster.MustConnect(kubeconfigPath)
		err := monkey.Run(
			cfg,
			logger,
			clientset,
			dynamicClient,
		)
		logger.Fatal("monkey runloop exited unexpectedly", zap.Error(err))

	case "sh":
		if kubeconfigPath == nil {
			// use default kubeconfig
			kubeconfigPath = &clientcmd.RecommendedHomeFile
		}

		clientset, _, _, err := cluster.GenerateClientsets(kubeconfigPath)
		if err != nil {
			logger.Info("no cluster detected")
		}
		err = sh.Run(
			cfg,
			logger,
			clientset,
		)
		if err != nil {
			logger.Fatal("epitomesh exited with error", zap.Error(err))
		}

		// otherwise, exit 0
	default:
		logger.Fatal("unknown mode", zap.String("mode", *mode))
	}
}
