package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itaraxa/shiny-broccoli/internal/models"
)

/* Чтение перечня устройств и параметров доступа по SNMP из конфигурационного файла в JSON формате.
 */
func LoadConfigFromJSON(fileName string) (*models.TSs, error) {
	TSs := new(models.TSs)

	return TSs, nil
}

func GenerateSkeletonConfigJSON(fileName string) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error opening file for writing congfig skelton: %v", err)
	}

	myTSs := new(models.TSs)
	myTS := new(models.TS)
	myTSs.ListOfTS = append(myTSs.ListOfTS, *myTS)
	jsonData, err := json.MarshalIndent(myTSs, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshall struct to json")
	}
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error write data to file: %v", err)
	}

	return
}
