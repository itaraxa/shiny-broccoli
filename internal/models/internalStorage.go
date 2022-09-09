package models

import (
	"time"
)

/* Структура для хранения результатов опроса устройств по SNMPv3
 */
type InternalStorage struct {
	Hosts map[string]struct {
		Status     string
		LastUpdate time.Time
		OIDs       map[string]interface{}
	}
}
