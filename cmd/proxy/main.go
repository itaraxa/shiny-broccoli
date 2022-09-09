/* proxy: main.go
 */
package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func makeApp() *cli.App {
	return &cli.App{
		Name:        "gosnmpproxy",
		Description: "Proxy server for converting SNMP query from version 2c to version 3",
		Commands: []*cli.Command{
			{
				Name:    "startProxy",
				Usage:   "Start SNMP proxy",
				Aliases: []string{"s"},
				Action:  startProxy,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "logLevel", Value: "info", Usage: "Logging level: fatal/error/info/debug/trace"},
					&cli.StringFlag{Name: "config", Value: "proxyConfig.json", Usage: "Proxy configuration file"},
					&cli.StringFlag{Name: "rules", Value: "proxyRules.json", Usage: "Proxy rules file"},
					&cli.StringFlag{Name: "version", Value: "2", Usage: "Version of proxy server function"},
				},
			},
			{
				Name:    "makeConfig",
				Usage:   "Make skeleton of general configuration file",
				Aliases: []string{"mc"},
				Action:  makeConfig,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "logLevel", Value: "info", Usage: "Logging level: fatal/error/info/debug/trace"},
					&cli.StringFlag{Name: "fileName", Value: "proxyConfig.json.skeleton", Usage: "Name for template configuration"},
				},
			},
			{
				Name:    "makeRules",
				Usage:   "Make Proxy rules template file from ASync2 configuration",
				Aliases: []string{"mr"},
				Action:  makeRules,
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
