package config

import (
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/itaraxa/shiny-broccoli/internal/models"
	"github.com/slayercat/GoSNMPServer"
)

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

func LoadDiagXML(XMLfileName string) (d DiagConf, err error) {
	xmlFile, err := os.Open(XMLfileName)
	if err != nil {
		return
	}
	defer xmlFile.Close()

	data, err := io.ReadAll(xmlFile)
	if err != nil {
		return
	}

	err = xml.Unmarshal(data, &d)
	if err != nil {
		return
	}

	return d, nil
}

/* Трансляция данных config.xml -> Node struct
 */
func NewNodes(logger GoSNMPServer.ILogger, d DiagConf) (n models.Nodes, err error) {
	for _, item := range d.Nodes.Node {
		tn := new(models.Node)
		tn.NodeName = item.Name
		if strings.Contains(item.Snmp.Mainlink, ":") {
			tn.NodeIPMain = strings.Split(item.Snmp.Mainlink, ":")[0]
			tn.NodePortMain = strings.Split(item.Snmp.Mainlink, ":")[1]
		} else {
			tn.NodeIPMain = item.Snmp.Mainlink
			tn.NodePortMain = "161"
		}
		if strings.Contains(item.Snmp.Standbylink, ":") {
			tn.NodeIPStandby = strings.Split(item.Snmp.Standbylink, ":")[0]
			tn.NodePortStandby = strings.Split(item.Snmp.Standbylink, ":")[1]
		} else {
			tn.NodeIPStandby = item.Snmp.Standbylink
			tn.NodePortStandby = "161"
		}
		tn.Community = item.Snmp.Community
		tn.SNMPVersion = item.Snmp.Version
		tn.Period, err = strconv.Atoi(item.Period)
		if err != nil {
			logger.Errorf("For %s cannot parse string %s to int: %v", item.Name, item.Period, err)
		}
		tl := make([]struct {
			OID  string
			Type string
			Id   string
		}, 10)
		for _, oid := range item.Snmp.Params.Param {
			tl = append(tl, struct {
				OID  string
				Type string
				Id   string
			}{OID: oid.Ref, Type: oid.Type, Id: oid.ID})
		}
		tn.OIDs = append(tn.OIDs, tl...)
		logger.Infof("For %s added %d points", tn.NodeName, len(tn.OIDs))
		n.ListOfTS = append(n.ListOfTS, *tn)
	}

	return
}
