package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itaraxa/shiny-broccoli/internal/models"
)

/* Чтение перечня устройств и параметров доступа по SNMP из конфигурационного файла в JSON формате.
 */
func LoadConfigFromJSON(fileName string) (*models.GlobalConfig, error) {
	myConf := new(models.GlobalConfig)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening configuration file %s: %v", fileName, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(myConf)
	if err != nil {
		return nil, fmt.Errorf("error reading json config from file %s: %v", fileName, err)
	}

	return myConf, nil
}

/* Создать файл с пустой конфигурацией
 */
func GenerateSkeletonConfigJSON(fileName string) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error opening file for writing congfig skelton: %v", err)
	}

	myConf := new(models.GlobalConfig)
	myTS := new(models.TS)
	myConf.TSs.ListOfTS = append(myConf.TSs.ListOfTS, *myTS)
	jsonData, err := json.MarshalIndent(myConf, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshall struct to json")
	}
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error write data to file: %v", err)
	}

	return
}
