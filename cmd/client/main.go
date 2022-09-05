/* client: main.go
 */
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func makeApp() *cli.App {
	return &cli.App{
		Name:        "SNMP client",
		Description: "SNMP client for testing SNMP proxy",
		Commands: []*cli.Command{
			{
				Name:    "start",
				Usage:   "start SNMP cient",
				Aliases: []string{"s"},
				Action:  startClient,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "logLevel", Value: "info"},
					&cli.StringFlag{Name: "config", Value: "client.json"},
				},
			},
		},
	}
}

func main() {
	app := makeApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Critical error. Cannot start application: %v\n", err)
	}
}

func startClient(c *cli.Context) (err error) {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{DisableLevelTruncation: true, FullTimestamp: true})
	logger.SetOutput(os.Stdout)
	switch strings.ToLower(c.String("logLevel")) {
	case "fatal":
		logger.Level = log.FatalLevel
	case "error":
		logger.Level = log.ErrorLevel
	case "info":
		logger.Level = log.InfoLevel
	case "debug":
		logger.Level = log.DebugLevel
	case "trace":
		logger.Level = log.TraceLevel
	}

	logger.WithFields(log.Fields{"logLevel": c.String("logLevel"),
		"config": c.String("config")}).Info("Started configuration")

	OIDs := []string{"1.3.6.1.2.1.2.2.1.1.0",
		"1.3.6.1.2.1.2.2.1.1.1",
		"1.3.6.1.2.1.2.2.1.1.2",
		"1.3.6.1.2.1.2.2.1.1.3",
		"1.3.6.1.2.1.2.2.1.2.0",
		"1.3.6.1.2.1.2.2.1.2.1",
		"1.3.6.1.2.1.2.2.1.2.2",
		"1.3.6.1.2.1.2.2.1.2.3",
		"1.3.6.1.2.1.2.2.1.3.0",
		"1.3.6.1.2.1.2.2.1.3.1",
		"1.3.6.1.2.1.2.2.1.3.2",
		"1.3.6.1.2.1.2.2.1.3.3",
	}

	params := &g.GoSNMP{
		Target:  "127.0.0.1",
		Port:    1161,
		Version: g.Version3,
		// Community:     "public",
		ContextName:   "public",
		SecurityModel: g.UserSecurityModel,
		MsgFlags:      g.AuthPriv,
		Timeout:       time.Duration(5) * time.Second,
		Retries:       3,
		SecurityParameters: &g.UsmSecurityParameters{
			UserName:                 "testuser",
			AuthenticationProtocol:   g.MD5,
			AuthenticationPassphrase: "testauth",
			PrivacyProtocol:          g.DES,
			PrivacyPassphrase:        "testpriv",
		},
		MaxOids: 16,
		Logger:  g.NewLogger(logger),
	}

	logger.WithFields(log.Fields{"Version": params.Version.String(),
		"Community": params.Community,
		"MsgFlags":  params.MsgFlags,
		"UserName":  params.SecurityParameters.Description,
	}).Info("SNMP connection parametres")

	err = params.Connect()
	if err != nil {
		logger.WithFields(log.Fields{"error": fmt.Sprintf("%v", err)}).Fatal("connection error")
	}
	defer params.Conn.Close()

	result, err := params.Get(OIDs)
	if err != nil {
		logger.WithFields(log.Fields{"error": fmt.Sprintf("%v", err)}).Fatal("get OID error")
	}

	fmt.Printf("Len result.Variables = %d", len(result.Variables))
	for _, variable := range result.Variables {
		switch variable.Type {
		case g.OctetString:
			logger.WithFields(log.Fields{"OID": variable.Name, "Value": string(variable.Value.([]byte))}).Info("get data")
		default:
			logger.WithFields(log.Fields{"OID": variable.Name, "Value": g.ToBigInt(variable.Value)}).Info("get data")
		}
	}
	return
}
