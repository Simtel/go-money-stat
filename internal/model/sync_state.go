package model

// SyncState хранит состояние последней синхронизации
type SyncState struct {
	ID              string `gorm:"primaryKey"`
	LastSyncedAt    int64  // Unix timestamp последней успешной синхронизации
	ServerTimestamp int64  // Timestamp с сервера Zenmoney
	UpdatedAt       int64  // Время обновления записи в БД
}

// TableName указывает имя таблицы
func (SyncState) TableName() string {
	return "sync_state"
}
