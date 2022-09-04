package main

import (
	"strings"

	"github.com/itaraxa/shiny-broccoli/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/slayercat/GoSNMPServer"
	"github.com/urfave/cli/v2"
)

func generateConfig(c *cli.Context) (err error) {
	// Create and setup logger
	logger := GoSNMPServer.NewDefaultLogger()
	switch strings.ToLower(c.String("logLevel")) {
	case "fatal":
		logger.(*GoSNMPServer.DefaultLogger).Level = logrus.FatalLevel
	case "error":
		logger.(*GoSNMPServer.DefaultLogger).Level = logrus.ErrorLevel
	case "info":
		logger.(*GoSNMPServer.DefaultLogger).Level = logrus.InfoLevel
	case "debug":
		logger.(*GoSNMPServer.DefaultLogger).Level = logrus.DebugLevel
	case "trace":
		logger.(*GoSNMPServer.DefaultLogger).Level = logrus.TraceLevel
	}

	logger.Info("start programm")
	logger.Info("generate skeleton configuration file to ", c.String("fileName"))
	err = config.GenerateSkeletonConfigJSON(c.String("fileName"))
	if err != nil {
		logger.Fatalf("error creating skeleton config file: %v", err)
	}
	logger.Info("skeleton configuration file was saved to ", c.String("fileName"))
	logger.Info("exit programm")

	return
}
