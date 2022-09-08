/* proxy: main.go
 */
package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/itaraxa/shiny-broccoli/internal/config"
	"github.com/itaraxa/shiny-broccoli/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/GoSNMPServer/mibImps"
	"github.com/slayercat/gosnmp"
	"github.com/urfave/cli/v2"
)

func makeApp() *cli.App {
	return &cli.App{
		Name:        "gosnmpproxy",
		Description: "Proxy server for converting SNMP query from version 2c to version 3",
		Commands: []*cli.Command{
			{
				Name:    "start",
				Usage:   "Start SNMP proxy",
				Aliases: []string{"s"},
				Action:  startProxy,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "logLevel", Value: "info", Usage: "Logging level: fatal/error/info/debug/trace"},
					&cli.StringFlag{Name: "config", Value: "proxyConfig.json", Usage: "Proxy configuration file"},
					&cli.StringFlag{Name: "rules", Value: "proxyRules.json", Usage: "Proxy rules file"},
				},
			},
			{
				Name:    "generate",
				Usage:   "Generate skeleton for Proxy configuration file",
				Aliases: []string{"g"},
				Action:  generateConfig,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "logLevel", Value: "info", Usage: "Logging level: fatal/error/info/debug/trace"},
					&cli.StringFlag{Name: "fileName", Value: "proxyConfig.json.skeleton", Usage: "Name for template configuration"},
				},
			},
			{
				Name:    "makeRules",
				Usage:   "Make Proxy rules template file from ASync2 configuration",
				Aliases: []string{"mr"},
				Action:  makeConfig,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "logLevel", Value: "info", Usage: "Logging level: fatal/error/info/debug/trace"},
					&cli.StringFlag{Name: "configXML", Value: "config.xml", Usage: "Configuration file for Async2"},
					&cli.StringFlag{Name: "rules", Value: "proxyRules.json", Usage: "Proxy rules file"},
					&cli.BoolFlag{Name: "show", Value: false, Usage: "Print result rules to screen"},
				},
			},
		},
	}
}

func main() {
	app := makeApp()
	app.Run(os.Args)
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

	myConfig, err := config.LoadConfigFromJSON(c.String("config"))
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	logger.Infof("Configuration file %s loaded", c.String("config"))

	diagConf, err := config.LoadDiagXML(myConfig.DiagXMLfile)
	if err != nil {
		logger.Fatalf("Critical error read diag configuration file %s:%v", myConfig.DiagXMLfile, err)
	}
	logger.Infof("Diag-xml configuration file %s loaded", myConfig.DiagXMLfile)

	myNodes, err := config.NewNodes(logger, diagConf)
	if err != nil {
		logger.Fatalf("Critical errorr convert node info from diag xml: %v", err)
	}
	logger.Infof("Node list is ready. Tottaly %d nodes", myNodes.Stat.Nodes)

	logger.Infof("Logging settings from file: Filename=%s logLevel=%s", myConfig.LogFile, myConfig.LogLevel)
	var wg sync.WaitGroup

	for _, client := range myNodes.ListOfTS {
		wg.Add(1)
		go func(b models.Node) {
			// Setup SNMP master: listen community "public" with default OIDs
			logger.Infof("Start goroutine for %s serving", b.NodeName)
			defer wg.Done()
			master := GoSNMPServer.MasterAgent{
				Logger: logger,
				SecurityConfig: GoSNMPServer.SecurityConfig{
					// NoSecurity: false, // disable authorisation
					AuthoritativeEngineBoots: 1,
					Users: []gosnmp.UsmSecurityParameters{
						{
							UserName:                 "test1",
							AuthenticationProtocol:   gosnmp.MD5,
							PrivacyProtocol:          gosnmp.DES,
							AuthenticationPassphrase: "test1test",
							PrivacyPassphrase:        "test1test",
						},
						{
							UserName:                 "test2",
							AuthenticationProtocol:   gosnmp.MD5,
							PrivacyProtocol:          gosnmp.DES,
							AuthenticationPassphrase: "test2test",
							PrivacyPassphrase:        "test2test",
						},
					},
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
			// Привязка к привеллигированному порту (<1000) требует прав root
			err = server.ListenUDP("udp", fmt.Sprintf("%s:%s", "127.0.0.1", "161"))
			if err != nil {
				logger.Errorf("Error in listen: %+v", err)
			}
			server.ServeForever()
		}(client)
	}

	wg.Wait()

	return nil
}

/* Generate proxy-configuration from Async configuration
 */
func makeConfig(c *cli.Context) error {
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

	conf := config.NewDiagConf()
	if err := conf.LoadXML(c.String("configXML")); err != nil {
		logger.Fatalf("Fatal error loading xml file %s: %v\n", c.String("configXML"), err)
	}
	logger.Infof("File %s loaded\n", c.String("configXML"))

	logger.Debugf("Data from %s:\n%s\n", c.String("configXML"), conf.String())

	pr, err := conf.NewProxyRules()
	if err != nil {
		logger.Fatalf("Fatal error creating proxy rules: %v", err)
	}
	logger.Info("Proxy rules created\n")
	if err = pr.DumpProxyRulesJSON(c.String("rules")); err != nil {
		logger.Fatalf("Fatal error dumping proxy rules: %v", err)
	}
	logger.Infof("File %s created\n", c.String("rules"))

	if c.Bool("show") {
		fmt.Printf("\n%s\n", pr.String())
	}

	return nil
}
