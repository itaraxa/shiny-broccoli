package convert

import (
	g "github.com/gosnmp/gosnmp"
	"github.com/itaraxa/shiny-broccoli/internal/models"
)

func V2cToV3(in *models.Entity) (out *models.Entity, err error) {
	out.Params = in.Params

	out.Params.Version = g.Version3

	return
}
