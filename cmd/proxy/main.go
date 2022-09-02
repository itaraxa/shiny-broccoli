/* proxy: main.go
 */
package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/GoSNMPServer/mibImps"
	"github.com/urfave/cli/v2"
)

func makeApp() *cli.App {
	return &cli.App{
		Name:        "gosnmpproxy",
		Description: "Proxy server for converting SNMP query from version 2c to version 3",
		Commands: []*cli.Command{
			{
				Name:    "start",
				Usage:   "Start SNMP proxy-server",
				Aliases: []string{"run"},
				Action:  startProxy,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "logLevel", Value: "info"},
					&cli.StringFlag{Name: "config", Value: "proxy.json"},
				},
			},
		},
	}
}

func startProxy(c *cli.Context) error {
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

	// Setup SNMP master: listen community "public" with default OIDs
	master := GoSNMPServer.MasterAgent{
		Logger: logger,
		SecurityConfig: GoSNMPServer.SecurityConfig{
			NoSecurity: true, // disable authorisation
		},
		SubAgents: []*GoSNMPServer.SubAgent{
			{
				// Setup listening community
				CommunityIDs: []string{"public"},
				// Setup list of OIDs with middleware
				// TODO: add midleware for resending query
				OIDs: mibImps.All(),
			},
		},
	}

	// Start SNMP server with master
	server := GoSNMPServer.NewSNMPServer(master)
	err := server.ListenUDP("udp", "127.0.0.1:6161")
	if err != nil {
		logger.Errorf("Error in listen: %+v", err)
	}
	server.ServeForever()

	return nil
}

func main() {
	app := makeApp()
	app.Run(os.Args)
}
