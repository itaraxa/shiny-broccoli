package models

import "encoding/xml"

/* Структура опиcывающая конфигурационный файл Async (/etc/diag/diag.xml)
 */
type DiagConf struct {
	XMLName        xml.Name `xml:"doc"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Logging        struct {
		Text  string `xml:",chardata"`
		Level string `xml:"level,attr"`
	} `xml:"logging"`
	Trapsrv struct {
		Text string `xml:",chardata"`
		Port string `xml:"port,attr"`
	} `xml:"trapsrv"`
	Nodes struct {
		Text string `xml:",chardata"`
		Node []struct {
			Text   string `xml:",chardata"`
			ID     string `xml:"id,attr"`
			Period string `xml:"period,attr"`
			Name   string `xml:"name,attr"`
			Snmp   struct {
				Text        string `xml:",chardata"`
				Version     string `xml:"version,attr"`
				Mainlink    string `xml:"mainlink,attr"`
				Community   string `xml:"community,attr"`
				Standbylink string `xml:"standbylink,attr"`
				Params      struct {
					Text  string `xml:",chardata"`
					Param []struct {
						Text     string `xml:",chardata"`
						Ref      string `xml:"ref,attr"`
						Namefull string `xml:"namefull,attr"`
						Type     string `xml:"type,attr"`
						Name     string `xml:"name,attr"`
						ID       string `xml:"id,attr"`
					} `xml:"param"`
				} `xml:"params"`
			} `xml:"snmp"`
			Ntp struct {
				Text        string `xml:",chardata"`
				Standbylink string `xml:"standbylink,attr"`
				Mainlink    string `xml:"mainlink,attr"`
				Params      struct {
					Text  string `xml:",chardata"`
					Param struct {
						Text     string `xml:",chardata"`
						Ref      string `xml:"ref,attr"`
						Namefull string `xml:"namefull,attr"`
						Type     string `xml:"type,attr"`
						Name     string `xml:"name,attr"`
						ID       string `xml:"id,attr"`
					} `xml:"param"`
				} `xml:"params"`
			} `xml:"ntp"`
		} `xml:"node"`
	} `xml:"nodes"`
	Dts struct {
		Text    string `xml:",chardata"`
		Port    string `xml:"port,attr"`
		Channel struct {
			Text  string `xml:",chardata"`
			ID    string `xml:"id,attr"`
			Value []struct {
				Text  string `xml:",chardata"`
				Node  string `xml:"node,attr"`
				ID    string `xml:"id,attr"`
				Param string `xml:"param,attr"`
			} `xml:"value"`
		} `xml:"channel"`
	} `xml:"dts"`
}
