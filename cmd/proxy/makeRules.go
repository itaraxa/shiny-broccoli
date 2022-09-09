package main

import (
	"fmt"
	"strings"

	"github.com/itaraxa/shiny-broccoli/internal/diagConfig"
	"github.com/itaraxa/shiny-broccoli/internal/proxyRules"
	"github.com/sirupsen/logrus"
	"github.com/slayercat/GoSNMPServer"
	"github.com/urfave/cli/v2"
)

/* Generate proxy-configuration from Async configuration
 */
func makeRules(c *cli.Context) error {
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

	conf := diagConfig.NewDiagConf()
	if err := conf.LoadXML(c.String("configXML")); err != nil {
		logger.Fatalf("Fatal error loading xml file %s: %v\n", c.String("configXML"), err)
	}
	logger.Infof("File %s loaded\n", c.String("configXML"))

	logger.Debugf("Data from %s:\n%s\n", c.String("configXML"), conf.String())

	pr := proxyRules.NewProxyRules()

	if err := pr.DumpProxyRulesJSON(c.String("rules")); err != nil {
		logger.Fatalf("Fatal error dumping proxy rules: %v", err)
	}
	logger.Infof("File %s created\n", c.String("rules"))

	if c.Bool("show") {
		fmt.Printf("\n%s\n", pr.String())
	}

	return nil
}
