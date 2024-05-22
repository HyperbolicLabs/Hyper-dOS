package main

import (
	"flag"
	"os"

	"epitome.hyperbolic.xyz/cluster"
	"epitome.hyperbolic.xyz/helper"
	"epitome.hyperbolic.xyz/hyperweb"
	"github.com/sirupsen/logrus"
)

const VERSION = "v1-alpha"

func main() {
	var token = os.Getenv("HYPERBOLIC_TOKEN")
	var gatewayUrl = os.Getenv("HYPERBOLIC_GATEWAY_URL")
	if token == "" {
		logrus.Fatalf("token not set")
	}
	if gatewayUrl == "" {
		logrus.Fatalf("gatewayUrl not set")
	}

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
	go hyperweb.Run(clientset, gatewayUrl, token)

	for {
		// do nothing
		<-make(chan struct{})
		logrus.Infof("exiting") // this will never happen
		return
	}
}
