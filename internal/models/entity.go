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

/* Общие настройки
 */
type GlobalConfig struct {
	LogFile  string `json:"File for logs"`
	LogLevel string `json:"Logging level (fatal/error/info/debug/trace)"`
	NProcs   int    `json:"Number of parallel processes"`
	TSs      `json:"List of TS"`
}

/*
Стурктура для хранения конфигурации технических средств
*/
type TSs struct {
	ListOfTS []TS
}

type TS struct {
	Name        string   `json:"KKS"`
	ListenPort  string   `json:"Port for getting SNMPv2c query"`
	TargetIP    string   `json:"IP address of TS"`
	TargetPort  string   `json:"SNMP port of TS"`
	OIDs        []string `json:"List of OIDs"`
	SNMPVersion string   `json:"TS SNMP version"`
	V2c         struct {
		Community string
	} `json:"Input SNMP parametres"`
	V3 struct {
		AuthLevel        string `json:"Auth level SNMPv3 (noAuthNoPriv/authNoPriv/authPriv)"`
		AuthName         string `json:"Auth name for noAuthNoPriv level"`
		AuthString       string `json:"MD5 or SHA auth string"`
		AuthMethod       string `json:"Hashing method for AuthString (MD5/SHA)"`
		EncryptionMethod string `json:"Data encryption method (DES). Only for authPriv level"`
	} `json:"Output SNMP parametres"`
}
