package config

import (
	"encoding/xml"
	"fmt"
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

func NewDiagConf() *DiagConf {
	return new(DiagConf)
}

/* Load config from <config.xml> file
 */
func (dc *DiagConf) LoadXML(fileName string) error {
	xmlFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	data, err := io.ReadAll(xmlFile)
	if err != nil {
		return err
	}

	if err = xml.Unmarshal(data, dc); err != nil {
		return err
	}

	return nil
}

/*
	Create PROXY-config

TO-DO: Add encoding from CP-1251, UTF-8
*/
func (dc *DiagConf) DumpXML(fileName string) error {
	xmlFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	data, err := xml.MarshalIndent(dc, "", "\t")
	if err != nil {
		return err
	}

	if _, err = xmlFile.Write(data); err != nil {
		return err
	}

	return nil
}

func (dc *DiagConf) String() string {
	var res string
	res = fmt.Sprintf("Logging: %s\n", dc.Logging.Level)
	res += fmt.Sprintf("Trap Server: %s\n", dc.Trapsrv.Port)
	for _, node := range dc.Nodes.Node {
		res += fmt.Sprintf("\tid=%s\tIP1=%s\tIP2=%s\n", node.ID, node.Snmp.Mainlink, node.Snmp.Standbylink)
		for _, OID := range node.Snmp.Params.Param {
			res += fmt.Sprintf("\t\tRef=%s\tName=%s\tID=%s\n", OID.Ref, OID.Name, OID.ID)
		}
	}

	return res
}

/* Fill proxy rules from config.xml
 */
func (dc *DiagConf) NewProxyRules() (*models.ProxyRules, error) {
	pr := new(models.ProxyRules)

	for _, node := range dc.Nodes.Node {
		t := new(struct {
			NodeId string
			Local  struct {
				IP1   string
				Port1 int
				IP2   string
				Port2 int
				SNMP  struct {
					Version string
					// For SNMPv1 and SNMPv2c
					Community string
					// For SNMPv3
					Level      string
					Context    string
					AuthMethod string
					AuthPass   string
					PrivMethod string
					PrivPass   string
				}
			}
			Remote struct {
				IP1   string
				Port1 int
				IP2   string
				Port2 int
				SNMP  struct {
					Version string
					// For SNMPv1 and SNMPv2c
					Community string
					// For SNMPv3
					Level      string
					Context    string
					AuthMethod string
					AuthPass   string
					PrivMethod string
					PrivPass   string
				}
			}
		})

		t.NodeId = node.ID
		if strings.Contains(node.Snmp.Mainlink, ":") {
			t.Local.IP1 = strings.Split(node.Snmp.Mainlink, ":")[0]
			t.Local.Port1, _ = strconv.Atoi(strings.Split(node.Snmp.Mainlink, ":")[1])
		} else {
			t.Local.IP1 = node.Snmp.Mainlink
			t.Local.Port1 = 161
		}

		if node.Snmp.Standbylink != "" {
			if strings.Contains(node.Snmp.Standbylink, ":") {
				t.Local.IP2 = strings.Split(node.Snmp.Standbylink, ":")[0]
				t.Local.Port2, _ = strconv.Atoi(strings.Split(node.Snmp.Standbylink, ":")[1])
			} else {
				t.Local.IP2 = node.Snmp.Standbylink
				t.Local.Port2 = 161
			}
		}

		// if strings.Contains(node.Snmp.Standbylink, ":") {
		// 	t.Local.IP = strings.Split(node.Snmp.Standbylink, ":")[0]
		// 	t.Local.Port, _ = strconv.Atoi(strings.Split(node.Snmp.Standbylink, ":")[1])
		// } else {
		// 	t.Local.IP = node.Snmp.Standbylink
		// 	t.Local.Port = 161
		// }

		t.Local.SNMP.Version = node.Snmp.Version
		t.Remote.SNMP.Version = "3"

		t.Local.SNMP.Community = node.Snmp.Community
		t.Remote.SNMP.Context = node.Snmp.Community

		t.Remote.SNMP.Level = "AuthPriv"
		t.Remote.SNMP.AuthMethod = "MD5"
		t.Remote.SNMP.AuthPass = "SNMPv3AuthPass"
		t.Remote.SNMP.PrivMethod = "DES"
		t.Remote.SNMP.PrivPass = "SNMPv3PrivPass"

		pr.Nodes = append(pr.Nodes, *t)
	}

	return pr, nil
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
