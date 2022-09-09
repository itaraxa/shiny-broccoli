package proxyRules

import (
	"encoding/json"
	"os"

	"github.com/itaraxa/shiny-broccoli/internal/models"
)

type proxyRules models.ProxyRules

/* proxyRules struct constructor
 */
func NewProxyRules() *proxyRules {
	return new(proxyRules)
}

/* Save proxy rules structure to JSON file
 */
func (pr *proxyRules) DumpProxyRulesJSON(fileName string) error {
	jsonFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(pr, "", "\t")
	if err != nil {
		return err
	}

	if _, err = jsonFile.Write(jsonData); err != nil {
		return err
	}

	return nil
}

/* Load ProxyRules from json file
 */
func (pr *proxyRules) LoadProxyRules(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(pr); err != nil {
		return err
	}

	return nil
}

func (pr *proxyRules) String() string {
	res, _ := json.MarshalIndent(pr, "", "\t")
	return string(res) + "\n"
}
