package oids

import (
	"strings"

	"github.com/itaraxa/shiny-broccoli/internal/models"
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

func NewOIDs(b models.Node) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{}

	for _, item := range b.OIDs {
		t := new(GoSNMPServer.PDUValueControlItem)
		t.OID = item.OID
		switch strings.ToLower(item.Type) {
		case "int":
			t.Type = gosnmp.Integer
		case "ana":
			t.Type = gosnmp.Integer
		default:
			t.Type = gosnmp.Integer
		}
		t.Document = item.Id

		toRet = append(toRet, t)
	}

	return toRet
}
