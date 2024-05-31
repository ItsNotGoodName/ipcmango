package bus

import (
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
)

type DeviceCreated struct {
	DeviceKey core.Key
}

type DeviceUpdated struct {
	DeviceKey core.Key
}

type DeviceDeleted struct {
	DeviceKey core.Key
}

type EventCreated struct {
	EventID    string
	DeviceKey  core.Key
	IgnoreDB   bool
	IgnoreMQTT bool
	IgnoreLive bool
	Event      dahuacgi.Event
	CreatedAt  time.Time
}

type EmailCreated struct {
	DeviceKey  core.Key
	MessageKey core.Key
}

type FileScanProgressed struct {
	DeviceKey core.Key
	Progress  float64
}

type FileScanFinished struct {
	DeviceKey    core.Key
	CreatedCount int64
	UpdatedCount int64
	DeletedCount int64
}
