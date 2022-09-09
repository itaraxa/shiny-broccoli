package models

import (
	g "github.com/gosnmp/gosnmp"
)

/* Структура для хранения параметров SNMP запроса
 */
type Entity struct {
	Params *g.GoSNMP
	// Ctx    context.Context
	// Logger log.Logger
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
