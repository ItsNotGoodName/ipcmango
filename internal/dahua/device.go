package dahua

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CreateDeviceArgs struct {
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

func CreateDevice(ctx context.Context, db *sqlx.DB, args CreateDeviceArgs) (DahuaDevice, error) {
	uuid := uuid.NewString()
	createdAt := types.NewTime(time.Now())
	updatedAt := types.NewTime(time.Now())

	var device DahuaDevice
	err := db.GetContext(ctx, &device, `
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
			location, 
			features,
			email, 
			created_at, 
			updated_at, 
			latitude, 
			longitude,
			sunrise_offset,
			sunset_offset,
			sync_video_in_mode
		) 
		VALUES ((SELECT value FROM generate_series WHERE value NOT IN (SELECT seed from dahua_devices) LIMIT 1), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
		RETURNING *;
	`,
		uuid,
		args.Name,
		args.IP,
		args.Username,
		args.Password,
		args.Location,
		args.Features,
		args.Email,
		createdAt,
		updatedAt,
		args.Latitude,
		args.Longitude,
		args.SunriseOffset,
		args.SunsetOffset,
		args.SyncVideoInMode,
	)
	if err != nil {
		return device, err
	}

	bus.Publish(bus.DeviceCreated{
		DeviceKey: device.Key,
	})

	return device, nil
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
	SyncVideoInMode bool
}

func UpdateDevice(ctx context.Context, db *sqlx.DB, args UpdateDeviceArgs) (DahuaDevice, error) {
	updatedAt := types.NewTime(time.Now())

	var device DahuaDevice
	err := db.GetContext(ctx, &device, `
		UPDATE dahua_devices SET
			name = ?,
			ip = ?,
			username = ?,
			password = coalesce(?, password),
			location = ?,
			features =  ?,
			email = ?,
			latitude = ?,
			longitude = ?,
			sunrise_offset = ?,
			sunset_offset = ?,
			sync_video_in_mode = ?,
			updated_at = ?
		WHERE uuid = ?
		RETURNING *;
	`,
		args.Name,
		args.IP,
		args.Username,
		args.Password,
		args.Location,
		args.Email,
		args.Features,
		args.Latitude,
		args.Longitude,
		args.SunriseOffset,
		args.SunsetOffset,
		args.SyncVideoInMode,
		updatedAt,
		args.UUID,
	)
	if err != nil {
		return device, err
	}

	bus.Publish(bus.DeviceUpdated{
		DeviceKey: device.Key,
	})

	return device, nil
}

func DeleteDevice(ctx context.Context, db *sqlx.DB, uuid string) error {
	var key core.Key
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
