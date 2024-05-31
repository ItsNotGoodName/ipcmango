package bus

import (
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
)

type DeviceCreated struct {
	DeviceKey types.Key
}

type DeviceUpdated struct {
	DeviceKey types.Key
}

type DeviceDeleted struct {
	DeviceKey types.Key
}

type EventCreated struct {
	EventID    string
	DeviceKey  types.Key
	IgnoreDB   bool
	IgnoreMQTT bool
	IgnoreLive bool
	Event      dahuacgi.Event
	CreatedAt  time.Time
}

type EmailCreated struct {
	DeviceKey  types.Key
	MessageKey types.Key
}

type FileScanProgressed struct {
	DeviceKey types.Key
	Progress  float64
}

type FileScanFinished struct {
	DeviceKey    types.Key
	CreatedCount int64
	UpdatedCount int64
	DeletedCount int64
}
