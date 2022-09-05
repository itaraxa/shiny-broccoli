package models

import (
	"time"

	g "github.com/gosnmp/gosnmp"
)

/* Структура для хранения параметров SNMP запроса
 */
type Entity struct {
	Params *g.GoSNMP
	// Ctx    context.Context
	// Logger log.Logger
}

/* Общие настройки
 */
type GlobalConfig struct {
	LogFile     string               `json:"File for logs"`
	LogLevel    string               `json:"Logging level (fatal/error/info/debug/trace)"`
	NProcs      int                  `json:"Number of parallel processes"`
	DiagXMLfile string               `json:"Path to diag.xml file"`
	SNMPv3      SNMPv3GlobalSettings `json:"Common SNMPv3 settings"`
}

/* Общие настройки SNMPv3
 */
type SNMPv3GlobalSettings struct {
	Level      string
	UserName   string
	Context    string
	AuthMethod string
	PrivMethod string
	AuthPass   string
	PrivPass   string
}

/*
Стурктура для хранения конфигурации технических средств/серверов SNMP
*/
type Nodes struct {
	Stat struct {
		Nodes  uint
		Points uint
	}
	ListOfTS []Node
}

/* Структура для хранения информации об узле
 */
type Node struct {
	NodeName        string `json:"KKS"`
	NodeIPMain      string `json:"IP address of TS"`
	NodePortMain    string
	NodeIPStandby   string
	NodePortStandby string
	Community       string
	OIDs            []struct {
		OID  string
		Type string
		Id   string
	} `json:"List of OIDs"`
	SNMPVersion string `json:"Node SNMP version"`
	Period      int
}

type ProxyRules struct {
	Stat  struct{}
	Nodes []struct {
		NodeId string
		Local  struct {
			IP   string
			Port int
			SNMP struct {
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
			IP   string
			Port int
			SNMP struct {
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
	}
}

/* Структура для хранения результатов опроса устройств по SNMPv3
 */
type InternalStorage struct {
	Hosts map[string]struct {
		Status     string
		LastUpdate time.Time
		OIDs       []struct {
			OID   string
			Value interface{}
		}
	}
}
