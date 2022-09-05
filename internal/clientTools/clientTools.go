package clienttools

import (
	g "github.com/gosnmp/gosnmp"
)

/* Get data by OID
 */
func GetData(params *g.GoSNMP, OIDs []string) (result []struct {
	OID   string
	Value interface{}
}, err error) {

	if err = params.Connect(); err != nil {
		return
	}
	defer params.Conn.Close()

	var a, b int
	a = 0
	if len(OIDs) > params.MaxOids {
		b = params.MaxOids
	} else {
		b = len(OIDs)
	}

	for {
		answer, err2 := params.Get(OIDs[a:b])
		if err != nil {
			return result, err2
		}

		for _, variable := range answer.Variables {
			switch variable.Type {
			// Ответ строкой
			case g.OctetString:
				result = append(result, struct {
					OID   string
					Value interface{}
				}{OID: variable.Name,
					Value: string(variable.Value.([]byte))})
			// Ответ числом
			default:
				result = append(result, struct {
					OID   string
					Value interface{}
				}{OID: variable.Name,
					Value: g.ToBigInt(variable.Value)})
			}
		}

		// Сдвигаем границы списка OIDов
		a = b
		if b+params.MaxOids > len(OIDs) {
			b = len(OIDs)
		} else {
			b = b + params.MaxOids
		}
		// Выход когда прошли все OID
		if a == len(OIDs) {
			break
		}
	}

	return result, nil
}
