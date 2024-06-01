package dahua

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/jmoiron/sqlx"
)

type Device struct {
	types.Key
	types.Timestamp
	Name               string
	IP                 string
	Username           string
	Password           string
	Location           sql.Null[types.Location]
	Features           types.Slice[Feature]
	Email              sql.NullString
	Seed               int64
	Latitude           sql.Null[float64]
	Longitude          sql.Null[float64]
	Sunrise_Offset     sql.Null[types.Duration]
	Sunset_Offset      sql.Null[types.Duration]
	Sync_Video_In_Mode sql.Null[bool]
}

type CreateDeviceArgs struct {
	UUID            string
	Name            string
	IP              string
	Username        string
	Password        string
	Location        *types.Location
	Features        types.Slice[Feature]
	Email           sql.Null[string]
	Latitude        sql.Null[float64]
	Longitude       sql.Null[float64]
	SunriseOffset   *types.Duration
	SunsetOffset    *types.Duration
	SyncVideoInMode sql.Null[bool]
}

func CreateDevice(ctx context.Context, db *sqlx.DB, args CreateDeviceArgs) (types.Key, error) {
	deviceKey, err := createDevice(ctx, db, args)
	if err != nil {
		return deviceKey, err
	}

	bus.Publish(bus.DeviceCreated{
		DeviceKey: deviceKey,
	})

	return deviceKey, nil
}

func PutDevices(ctx context.Context, db *sqlx.DB, args []CreateDeviceArgs) ([]types.Key, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var deletedKeys []types.Key
	err = db.SelectContext(ctx, &deletedKeys, `
		DELETE FROM dahua_devices RETURNING id, uuid
	`)
	if err != nil {
		return nil, err
	}

	var keys []types.Key
	for _, arg := range args {
		key, err := createDevice(ctx, tx, arg)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	for _, deletedKey := range deletedKeys {
		bus.Publish(bus.DeviceDeleted{
			DeviceKey: deletedKey,
		})
	}

	for _, key := range keys {
		bus.Publish(bus.DeviceCreated{
			DeviceKey: key,
		})
	}

	return keys, nil
}

func createDevice(ctx context.Context, db sqlx.ExtContext, args CreateDeviceArgs) (types.Key, error) {
	createdAt := types.NewTime(time.Now())
	updatedAt := types.NewTime(time.Now())

	var key types.Key
	err := sqlx.GetContext(ctx, db, &key, `
		WITH RECURSIVE generate_series(value) AS (
			SELECT 1
			UNION ALL
			SELECT value+1 FROM generate_series WHERE value+1<=999
		)
		INSERT INTO dahua_devices (
			seed, 
			uuid, 
			name, 
			ip, 
			username, 
			password, 
			email, 
			features,
			location, 
			latitude, 
			longitude,
			sunrise_offset,
			sunset_offset,
			sync_video_in_mode,
			created_at,
			updated_at
		) 
		VALUES ((SELECT value FROM generate_series WHERE value NOT IN (SELECT seed from dahua_devices) LIMIT 1), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
		RETURNING id, uuid;
	`,
		args.UUID,
		args.Name,
		args.IP,
		args.Username,
		args.Password,
		args.Email,
		args.Features,
		args.Location,
		args.Latitude,
		args.Longitude,
		args.SunriseOffset,
		args.SunsetOffset,
		args.SyncVideoInMode,
		createdAt,
		updatedAt,
	)
	if err != nil {
		return types.Key{}, err
	}

	if err := ResetFileScanCursor(ctx, db, key.ID); err != nil {
		return types.Key{}, err
	}

	return key, nil
}

type UpdateDeviceArgs struct {
	UUID            string
	Name            string
	IP              string
	Username        string
	Password        sql.Null[string]
	Location        *types.Location
	Features        types.Slice[Feature]
	Email           *string
	Latitude        sql.Null[float64]
	Longitude       sql.Null[float64]
	SunriseOffset   *types.Duration
	SunsetOffset    *types.Duration
	SyncVideoInMode sql.Null[bool]
}

func UpdateDevice(ctx context.Context, db *sqlx.DB, args UpdateDeviceArgs) (types.Key, error) {
	updatedAt := types.NewTime(time.Now())

	var key types.Key
	err := db.GetContext(ctx, &key, `
		UPDATE dahua_devices SET
			name = ?,
			ip = ?,
			username = ?,
			password = coalesce(?, password),
			email = ?,
			features =  ?,
			location = ?,
			latitude = ?,
			longitude = ?,
			sunrise_offset = ?,
			sunset_offset = ?,
			sync_video_in_mode = ?,
			updated_at = ?
		WHERE uuid = ?
		RETURNING id, uuid;
	`,
		args.Name,
		args.IP,
		args.Username,
		args.Password,
		args.Email,
		args.Features,
		args.Location,
		args.Latitude,
		args.Longitude,
		args.SunriseOffset,
		args.SunsetOffset,
		args.SyncVideoInMode,
		updatedAt,
		args.UUID,
	)
	if err != nil {
		return key, err
	}

	bus.Publish(bus.DeviceUpdated{
		DeviceKey: key,
	})

	return key, nil
}

func DeleteDevice(ctx context.Context, db *sqlx.DB, uuid string) error {
	var key types.Key
	err := db.GetContext(ctx, &key, `
		DELETE FROM dahua_devices WHERE uuid = ? RETURNING uuid, id
	`, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	bus.Publish(bus.DeviceDeleted{
		DeviceKey: key,
	})

	return nil
}

type DevicePosition struct {
	Location       types.Location
	Latitude       float64
	Longitude      float64
	Sunrise_Offset types.Duration
	Sunset_Offset  types.Duration
}

func GetDevicePosition(ctx context.Context, db *sqlx.DB, id int64) (DevicePosition, error) {
	var position DevicePosition
	err := db.GetContext(ctx, &position, `
		SELECT 
			coalesce(d.location, s.location) AS location,
			coalesce(d.latitude, s.latitude) AS latitude,
			coalesce(d.longitude, s.longitude) AS longitude,
			coalesce(d.sunrise_offset, s.sunrise_offset) AS sunrise_offset,
			coalesce(d.sunset_offset, s.sunset_offset) AS sunset_offset
		FROM dahua_devices AS d, settings as s
		WHERE d.id = ?
	`, id)
	return position, err
}
