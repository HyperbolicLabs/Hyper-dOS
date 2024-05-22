package main

import (
	"flag"

	"epitome.hyperbolic.xyz/cluster"
	"epitome.hyperbolic.xyz/helper"
	"epitome.hyperbolic.xyz/hyperweb"
	"github.com/sirupsen/logrus"
)

const VERSION = "v1-alpha"

func main() {
	var help = flag.Bool("help", false, "Show help")
	var loglevel = flag.String("loglevel", "info", "debug, info, error")
	flag.Parse()

	helper.SetLogLevel(*loglevel)
	logrus.Infof("version: %v", VERSION)

	if *help {
		flag.PrintDefaults()
		return
	}

	logrus.Infof("connecting to in-cluster kube api-server")
	clientset := cluster.MustConnect()
	hyperweb.Initialize(clientset)
}
