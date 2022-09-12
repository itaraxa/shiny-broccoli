package proxyRules

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itaraxa/shiny-broccoli/internal/models"
)

/* proxyRules struct constructor
 */
func NewProxyRules() *models.ProxyRules {
	return new(models.ProxyRules)
}

/* Save proxy rules structure to JSON file
 */
func DumpProxyRulesJSON(pr *models.ProxyRules, fileName string) error {
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
func LoadProxyRules(pr *models.ProxyRules, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", fileName, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(pr); err != nil {
		return err
	}

	return nil
}
