// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package repo

import (
	"database/sql"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

type DahuaCamera struct {
	ID        int64
	Name      string
	Address   string
	Username  string
	Password  string
	Location  types.Location
	CreatedAt types.Time
	UpdatedAt types.Time
}

type DahuaEvent struct {
	ID        int64
	CameraID  int64
	Code      string
	Action    string
	Index     int64
	Data      json.RawMessage
	CreatedAt types.Time
}

type DahuaEventDefaultRule struct {
	ID         int64
	Code       string
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
}

type DahuaEventRule struct {
	CameraID   int64
	Code       string
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
}

type DahuaFile struct {
	ID          int64
	CameraID    int64
	Channel     int64
	StartTime   types.Time
	EndTime     types.Time
	Length      int64
	Type        string
	FilePath    string
	Duration    int64
	Disk        int64
	VideoStream string
	Flags       types.StringSlice
	Events      types.StringSlice
	Cluster     int64
	Partition   int64
	PicIndex    int64
	Repeat      int64
	WorkDir     string
	WorkDirSn   int64
	UpdatedAt   types.Time
}

type DahuaFileCursor struct {
	CameraID     int64
	QuickCursor  types.Time
	FullCursor   types.Time
	FullEpoch    types.Time
	FullComplete bool
}

type DahuaFileScanLock struct {
	CameraID  int64
	CreatedAt types.Time
}

type DahuaSeed struct {
	Seed     int64
	CameraID sql.NullInt64
}

type Setting struct {
	SiteName        string
	DefaultLocation types.Location
}
