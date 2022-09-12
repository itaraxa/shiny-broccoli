package internalStorage

import (
	"errors"
	"sync"
	"time"

	"github.com/itaraxa/shiny-broccoli/internal/models"
)

// type internalStorage models.InternalStorage

func NewInternalStorage() *models.InternalStorage {
	return new(models.InternalStorage)
}

/* Инициализация структуры внутреннего хранилища
 */
func Init(is *models.InternalStorage, pr *models.ProxyRules, dc *models.DiagConfig) (err error) {
	for _, rule := range pr.Nodes {
		is.Hosts[rule.NodeId] = struct {
			Status     string
			LastUpdate time.Time
			OIDs       map[string]interface{}
		}{Status: "Not answered",
			LastUpdate: time.Now()}
		// for _, OID := range dc.Nodes.Node
	}

	return
}

/* Запись данных во внутреннее хранилище
 */
func Add(hostname string, data *map[string]interface{}, m *sync.Mutex, is *models.InternalStorage) (err error) {
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
func Get(is *models.InternalStorage, hostname string) (status string, lastUpdate time.Time, data map[string]interface{}, err error) {
	status = is.Hosts[hostname].Status
	lastUpdate = is.Hosts[hostname].LastUpdate
	data = is.Hosts[hostname].OIDs

	return status, lastUpdate, data, nil
}
