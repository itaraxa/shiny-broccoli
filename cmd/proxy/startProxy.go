package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/itaraxa/shiny-broccoli/internal/diagConfig"
	"github.com/itaraxa/shiny-broccoli/internal/globalConfig"
	"github.com/itaraxa/shiny-broccoli/internal/internalStorage"
	"github.com/itaraxa/shiny-broccoli/internal/proxyRules"
	"github.com/sirupsen/logrus"
	"github.com/slayercat/GoSNMPServer"
	"github.com/urfave/cli/v2"
)

func startProxy(c *cli.Context) error {
	switch c.String("version") {
	case "1":
		if err := startProxyV1(c); err != nil {
			return err
		}
	case "2":
		if err := startProxyV2(c); err != nil {
			return err
		}
	default:
		if err := startProxyV1(c); err != nil {
			return err
		}
	}
	return nil
}

func startProxyV2(c *cli.Context) error {
	// Настроить логирование
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

	// Graceful shutdown
	defer func() {
		if r := recover(); r != nil {
			logger.Infof("SNMPProxy stoped. Reason: %v", r)
		}
	}()

	// Catch Ctrl+C
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		logger.Info("SNMPProxy stoped. Reason: Keyboard interrupt")
		os.Exit(1)
	}()

	// Прочитать глобальную конфигурацию
	gc, err := globalConfig.LoadConfigFromJSON(c.String("config"))
	if err != nil {
		logger.Fatalf("Error loading global configuration from %s: %v", c.String("config"), err)
	}
	logger.Infof("Global config: NProc = %d", gc.NProcs)

	// Прочитать правила proxy
	pr := proxyRules.NewProxyRules()
	if err = proxyRules.LoadProxyRules(pr, c.String("rules")); err != nil {
		log.Fatalf("Error loading Proxy rules from %s: %v", c.String("rules"), err)
	}
	logger.Infof("Proxy rules loaded from %s\n%s", c.String("rules"), pr.String())

	// Прочитать конфигурацию async для получения перечня OID
	dc := diagConfig.NewDiagConf()

	// Создать струтуры для хранения данных
	is := internalStorage.NewInternalStorage()
	if err = internalStorage.Init(is, pr, dc); err != nil {
		log.Fatalf("Error initialize internal storage: %v", err)
	}
	logger.Infof("Internal storage initialized")

	// Запустить горутины SNMPv3 clients

	// Run goroutines SNMPv3 servers

	logger.Info("SNMPProxy stoped. Reason: end programm")
	return nil
}

/* Proxy-сервер: Версия 1
 */
func startProxyV1(c *cli.Context) error {
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

	// myConfig, err := config.LoadConfigFromJSON(c.String("config"))
	// if err != nil {
	// 	logger.Fatalf("Error loading config: %v", err)
	// }
	// logger.Infof("Configuration file %s loaded", c.String("config"))

	// if err != nil {
	// 	logger.Fatalf("Critical errorr convert node info from diag xml: %v", err)
	// }
	// logger.Infof("Node list is ready. Tottaly %d nodes", myNodes.Stat.Nodes)

	// var wg sync.WaitGroup

	// for _, client := range myNodes.ListOfTS {
	// 	wg.Add(1)
	// 	go func(b models.Node) {
	// 		// Setup SNMP master: listen community "public" with default OIDs
	// 		logger.Infof("Start goroutine for %s serving", b.NodeName)
	// 		defer wg.Done()
	// 		master := GoSNMPServer.MasterAgent{
	// 			Logger: logger,
	// 			SecurityConfig: GoSNMPServer.SecurityConfig{
	// 				// NoSecurity: false, // disable authorisation
	// 				AuthoritativeEngineBoots: 1,
	// 				Users: []gosnmp.UsmSecurityParameters{
	// 					{
	// 						UserName:                 "test1",
	// 						AuthenticationProtocol:   gosnmp.MD5,
	// 						PrivacyProtocol:          gosnmp.DES,
	// 						AuthenticationPassphrase: "test1test",
	// 						PrivacyPassphrase:        "test1test",
	// 					},
	// 					{
	// 						UserName:                 "test2",
	// 						AuthenticationProtocol:   gosnmp.MD5,
	// 						PrivacyProtocol:          gosnmp.DES,
	// 						AuthenticationPassphrase: "test2test",
	// 						PrivacyPassphrase:        "test2test",
	// 					},
	// 				},
	// 			},
	// 			SubAgents: []*GoSNMPServer.SubAgent{
	// 				{
	// 					// Setup listening community
	// 					CommunityIDs: []string{"public"},
	// 					// Setup list of OIDs with middleware
	// 					// TODO: add midleware for resending query
	// 					OIDs: mibImps.All(),
	// 				},
	// 			},
	// 		}

	// 		// Start SNMP server with master
	// 		server := GoSNMPServer.NewSNMPServer(master)
	// 		// Привязка к привеллигированному порту (<1000) требует прав root
	// 		err = server.ListenUDP("udp", fmt.Sprintf("%s:%s", "127.0.0.1", "161"))
	// 		if err != nil {
	// 			logger.Errorf("Error in listen: %+v", err)
	// 		}
	// 		server.ServeForever()
	// 	}(client)
	// }

	// wg.Wait()

	return nil
}
