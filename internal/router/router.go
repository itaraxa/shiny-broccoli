package router

import (
	"net"

	"github.com/itaraxa/shiny-broccoli/internal/models"
)

/* Change destination IP address in SNMP-query
 */
func RedirectQuery(in *models.Entity, dest net.IPAddr) (err error) {

	return
}

/* Change destination IP address in SNMP-answer
 */
func RedirectAnswer(in *models.Entity, dest net.IPAddr) (err error) {

	return
}
