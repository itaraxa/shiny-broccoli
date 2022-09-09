package internalStorage

import (
	"errors"
	"sync"
	"time"

	"github.com/itaraxa/shiny-broccoli/internal/models"
)

type internalStorage models.InternalStorage

func NewInternalStorage() *internalStorage {
	return new(internalStorage)
}

/* Инициализация структуры внутреннего хранилища
 */
func (is *internalStorage) Init(pr *models.ProxyRules) error {

	return nil
}

/* Запись данных во внутреннее хранилище
 */
func (is *internalStorage) Add(hostname string, data *map[string]interface{}, m *sync.Mutex) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("data add error")
			m.Unlock()
		}
	}()

	// Update data in internalStorage struct
	m.Lock()
	for oid, val := range *data {
		is.Hosts[hostname].OIDs[oid] = val
	}
	m.Unlock()

	return err
}

/* Получить данные из внутреннего хранилища
 */
func (is *internalStorage) Get(hostname string) (status string, lastUpdate time.Time, data map[string]interface{}, err error) {
	status = is.Hosts[hostname].Status
	lastUpdate = is.Hosts[hostname].LastUpdate
	data = is.Hosts[hostname].OIDs

	return status, lastUpdate, data, nil
}
