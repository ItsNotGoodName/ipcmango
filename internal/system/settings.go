package system

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/jmoiron/sqlx"
)

type Settings struct {
	Location           types.Location
	Latitude           float64
	Longitude          float64
	Sunrise_Offset     types.Duration
	Sunset_Offset      types.Duration
	Sync_Video_In_Mode bool
}

const (
	KeyLocation        = "location"
	KeyLatitude        = "latitude"
	KeyLongitude       = "longitude"
	KeySunriseOffset   = "sunrise_offset"
	KeySunsetOffset    = "sunset_offset"
	KeySyncVideoInMode = "sync_video_in_mode"
)

func NewSettings() Settings {
	return Settings{
		Location:       types.NewLocation(time.Local),
		Latitude:       0,
		Longitude:      0,
		Sunrise_Offset: types.NewDuration(0),
		Sunset_Offset:  types.NewDuration(0),
	}
}

func DefaultSettings(ctx context.Context, db *sqlx.DB) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		DELETE FROM settings
	`)
	if err != nil {
		return err
	}

	if err := initializeSettings(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

func InitializeSettings(ctx context.Context, db *sqlx.DB) error {
	return initializeSettings(ctx, db)
}

func initializeSettings(ctx context.Context, tx sqlx.ExecerContext) error {
	settings := NewSettings()
	updatedAt := types.NewTime(time.Now())
	_, err := tx.ExecContext(ctx, `
		INSERT INTO settings (
			key,
			value,
			updated_at
		) VALUES 
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?)
		ON CONFLICT DO NOTHING
	`,
		KeyLocation, settings.Location, updatedAt,
		KeyLatitude, settings.Latitude, updatedAt,
		KeyLongitude, settings.Longitude, updatedAt,
		KeySunriseOffset, settings.Sunrise_Offset, updatedAt,
		KeySunsetOffset, settings.Sunset_Offset, updatedAt,
		KeySyncVideoInMode, settings.Sync_Video_In_Mode, updatedAt,
	)
	return err
}

type UpdateSettingsArgs struct {
	Location        types.Location
	Latitude        float64
	Longitude       float64
	SunriseOffset   types.Duration
	SunsetOffset    types.Duration
	SyncVideoInMode bool
}

func UpdateSettings(ctx context.Context, db *sqlx.DB, settings UpdateSettingsArgs) error {
	updatedAt := types.NewTime(time.Now())
	_, err := db.ExecContext(ctx, `
		REPLACE INTO settings (
			key,
			value,
			updated_at
		) VALUES 
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?)
	`,
		KeyLocation, settings.Location, updatedAt,
		KeyLatitude, settings.Latitude, updatedAt,
		KeyLongitude, settings.Longitude, updatedAt,
		KeySunriseOffset, settings.SunriseOffset, updatedAt,
		KeySunsetOffset, settings.SunsetOffset, updatedAt,
		KeySyncVideoInMode, settings.SyncVideoInMode, updatedAt,
	)
	return err
}

func GetSettings(ctx context.Context, db *sqlx.DB) (Settings, error) {
	var settings Settings
	err := db.GetContext(ctx, &settings, `
		SELECT 
			(SELECT value FROM settings WHERE key = ?) AS location,
			(SELECT value FROM settings WHERE key = ?) AS latitude,
			(SELECT value FROM settings WHERE key = ?) AS longitude,
			(SELECT value FROM settings WHERE key = ?) AS sunrise_offset,
			(SELECT value FROM settings WHERE key = ?) AS sunset_offset,
			(SELECT value FROM settings WHERE key = ?) AS sync_video_in_mode
	`,
		KeyLocation,
		KeyLatitude,
		KeyLongitude,
		KeySunriseOffset,
		KeySunsetOffset,
		KeySyncVideoInMode,
	)
	if err != nil {
		return Settings{}, err
	}

	return settings, nil
}
