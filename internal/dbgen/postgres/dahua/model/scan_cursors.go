//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type ScanCursors struct {
	CameraID     int32
	QuickCursor  time.Time
	FullCursor   time.Time
	FullEpoch    time.Time
	FullEpochEnd time.Time
	FullComplete bool
}
