package models

/* Proxy rules configuration
 */
type ProxyRules struct {
	Stat  struct{}
	Nodes []struct {
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
	}
}
