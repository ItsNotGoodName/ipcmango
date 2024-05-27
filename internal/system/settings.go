package system

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/jmoiron/sqlx"
)

type Settings struct {
	ID                 int64
	Location           types.Location
	Latitude           float64
	Longitude          float64
	Sunrise_Offset     types.Duration
	Sunset_Offset      types.Duration
	Sync_Video_In_Mode bool
	Updated_At         types.Time
}

func NewSettings() Settings {
	return Settings{
		Location:       types.NewLocation(time.Local),
		Latitude:       0,
		Longitude:      0,
		Sunrise_Offset: types.NewDuration(0),
		Sunset_Offset:  types.NewDuration(0),
		Updated_At:     types.NewTime(time.Now()),
	}
}

func initializeSettings(ctx context.Context, tx sqlx.ExecerContext) error {
	settings := NewSettings()
	_, err := tx.ExecContext(ctx, `
		INSERT INTO settings (
			location,
			latitude,
			longitude,
			sunrise_offset,
			sunset_offset,
			sync_video_in_mode,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT DO NOTHING
	`,
		settings.Location,
		settings.Latitude,
		settings.Longitude,
		settings.Sunrise_Offset,
		settings.Sunrise_Offset,
		settings.Sync_Video_In_Mode,
		settings.Updated_At,
	)
	return err
}

func InitializeSettings(ctx context.Context, db *sqlx.DB) error {
	return initializeSettings(ctx, db)
}

func DefaultSettings(ctx context.Context, db *sqlx.DB) (Settings, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return Settings{}, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		DELETE FROM settings
	`)
	if err != nil {
		return Settings{}, err
	}

	if err := initializeSettings(ctx, tx); err != nil {
		return Settings{}, err
	}

	var settings Settings
	err = tx.Get(&settings, `
		SELECT * FROM settings
	`)
	if err != nil {
		return Settings{}, err
	}

	return settings, tx.Commit()
}

type UpdateSettingsArgs struct {
	Location        types.Location
	Latitude        float64
	Longitude       float64
	SunriseOffset   types.Duration
	SunsetOffset    types.Duration
	SyncVideoInMode bool
}

func UpdateSettings(ctx context.Context, db *sqlx.DB, args UpdateSettingsArgs) (Settings, error) {
	updatedAt := types.NewTime(time.Now())

	var settings Settings
	err := db.GetContext(ctx, &settings, `
		UPDATE settings SET
			location = ?,
			latitude = ?,
			longitude = ?,
			sunrise_offset = ?,
			sunset_offset = ?,
			sync_video_in_mode = ?,
			updated_at = ?
		RETURNING *;
	`,
		args.Location,
		args.Latitude,
		args.Longitude,
		args.SunriseOffset,
		args.SunsetOffset,
		args.SyncVideoInMode,
		updatedAt,
	)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
