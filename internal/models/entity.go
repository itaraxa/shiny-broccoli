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
