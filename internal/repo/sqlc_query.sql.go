// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: sqlc_query.sql

package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

const createDahuaEvent = `-- name: CreateDahuaEvent :one
INSERT INTO dahua_events (
  device_id,
  code,
  action,
  ` + "`" + `index` + "`" + `,
  data,
  created_at
) VALUES (
  ?, ?, ?, ?, ?, ?
) RETURNING id
`

type CreateDahuaEventParams struct {
	DeviceID  int64
	Code      string
	Action    string
	Index     int64
	Data      json.RawMessage
	CreatedAt types.Time
}

func (q *Queries) CreateDahuaEvent(ctx context.Context, arg CreateDahuaEventParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createDahuaEvent,
		arg.DeviceID,
		arg.Code,
		arg.Action,
		arg.Index,
		arg.Data,
		arg.CreatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createDahuaEventRule = `-- name: CreateDahuaEventRule :exec
INSERT INTO dahua_event_rules(
  code,
  ignore_db,
  ignore_live,
  ignore_mqtt
) VALUES(
  ?,
  ?,
  ?,
  ?
)
`

type CreateDahuaEventRuleParams struct {
	Code       string
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
}

func (q *Queries) CreateDahuaEventRule(ctx context.Context, arg CreateDahuaEventRuleParams) error {
	_, err := q.db.ExecContext(ctx, createDahuaEventRule,
		arg.Code,
		arg.IgnoreDb,
		arg.IgnoreLive,
		arg.IgnoreMqtt,
	)
	return err
}

const createDahuaEventWorkerState = `-- name: CreateDahuaEventWorkerState :exec
INSERT INTO dahua_event_worker_states(
  device_id,
  state,
  error,
  created_at
) VALUES(
  ?,
  ?,
  ?,
  ?
)
`

type CreateDahuaEventWorkerStateParams struct {
	DeviceID  int64
	State     models.DahuaEventWorkerState
	Error     sql.NullString
	CreatedAt types.Time
}

func (q *Queries) CreateDahuaEventWorkerState(ctx context.Context, arg CreateDahuaEventWorkerStateParams) error {
	_, err := q.db.ExecContext(ctx, createDahuaEventWorkerState,
		arg.DeviceID,
		arg.State,
		arg.Error,
		arg.CreatedAt,
	)
	return err
}

const createDahuaFile = `-- name: CreateDahuaFile :one
INSERT INTO dahua_files (
  device_id,
  channel,
  start_time,
  end_time,
  length,
  type,
  file_path,
  duration,
  disk,
  video_stream,
  flags,
  events,
  cluster,
  partition,
  pic_index,
  repeat,
  work_dir,
  work_dir_sn,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? 
) RETURNING id
`

type CreateDahuaFileParams struct {
	DeviceID    int64
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

func (q *Queries) CreateDahuaFile(ctx context.Context, arg CreateDahuaFileParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createDahuaFile,
		arg.DeviceID,
		arg.Channel,
		arg.StartTime,
		arg.EndTime,
		arg.Length,
		arg.Type,
		arg.FilePath,
		arg.Duration,
		arg.Disk,
		arg.VideoStream,
		arg.Flags,
		arg.Events,
		arg.Cluster,
		arg.Partition,
		arg.PicIndex,
		arg.Repeat,
		arg.WorkDir,
		arg.WorkDirSn,
		arg.UpdatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createDahuaFileScanLock = `-- name: CreateDahuaFileScanLock :one
INSERT INTO dahua_file_scan_locks (
  device_id, touched_at
) VALUES (
  ?, ?
) RETURNING device_id, touched_at
`

type CreateDahuaFileScanLockParams struct {
	DeviceID  int64
	TouchedAt types.Time
}

func (q *Queries) CreateDahuaFileScanLock(ctx context.Context, arg CreateDahuaFileScanLockParams) (DahuaFileScanLock, error) {
	row := q.db.QueryRowContext(ctx, createDahuaFileScanLock, arg.DeviceID, arg.TouchedAt)
	var i DahuaFileScanLock
	err := row.Scan(&i.DeviceID, &i.TouchedAt)
	return i, err
}

const deleteDahuaDevice = `-- name: DeleteDahuaDevice :exec
DELETE FROM dahua_devices WHERE id = ?
`

func (q *Queries) DeleteDahuaDevice(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteDahuaDevice, id)
	return err
}

const deleteDahuaEventRule = `-- name: DeleteDahuaEventRule :exec
DELETE FROM dahua_event_rules WHERE id = ?
`

func (q *Queries) DeleteDahuaEventRule(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteDahuaEventRule, id)
	return err
}

const deleteDahuaFile = `-- name: DeleteDahuaFile :exec
DELETE FROM dahua_files
WHERE
  updated_at < ?1 AND
  device_id = ?2 AND
  start_time <= ?3 AND
  ?4 < start_time
`

type DeleteDahuaFileParams struct {
	UpdatedAt types.Time
	DeviceID  int64
	End       types.Time
	Start     types.Time
}

func (q *Queries) DeleteDahuaFile(ctx context.Context, arg DeleteDahuaFileParams) error {
	_, err := q.db.ExecContext(ctx, deleteDahuaFile,
		arg.UpdatedAt,
		arg.DeviceID,
		arg.End,
		arg.Start,
	)
	return err
}

const deleteDahuaFileScanLock = `-- name: DeleteDahuaFileScanLock :exec
DELETE FROM dahua_file_scan_locks WHERE device_id = ?
`

func (q *Queries) DeleteDahuaFileScanLock(ctx context.Context, deviceID int64) error {
	_, err := q.db.ExecContext(ctx, deleteDahuaFileScanLock, deviceID)
	return err
}

const deleteDahuaFileScanLockByAge = `-- name: DeleteDahuaFileScanLockByAge :exec
DELETE FROM dahua_file_scan_locks WHERE touched_at < ?
`

func (q *Queries) DeleteDahuaFileScanLockByAge(ctx context.Context, touchedAt types.Time) error {
	_, err := q.db.ExecContext(ctx, deleteDahuaFileScanLockByAge, touchedAt)
	return err
}

const getDahuaDevice = `-- name: GetDahuaDevice :one
SELECT dahua_devices.id, dahua_devices.name, dahua_devices.address, dahua_devices.username, dahua_devices.password, dahua_devices.location, dahua_devices.created_at, dahua_devices.updated_at, coalesce(seed, id) FROM dahua_devices 
LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
WHERE id = ? LIMIT 1
`

type GetDahuaDeviceRow struct {
	ID        int64
	Name      string
	Address   string
	Username  string
	Password  string
	Location  types.Location
	CreatedAt types.Time
	UpdatedAt types.Time
	Seed      int64
}

func (q *Queries) GetDahuaDevice(ctx context.Context, id int64) (GetDahuaDeviceRow, error) {
	row := q.db.QueryRowContext(ctx, getDahuaDevice, id)
	var i GetDahuaDeviceRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Username,
		&i.Password,
		&i.Location,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Seed,
	)
	return i, err
}

const getDahuaEventData = `-- name: GetDahuaEventData :one
SELECT data FROM dahua_events WHERE id = ?
`

func (q *Queries) GetDahuaEventData(ctx context.Context, id int64) (json.RawMessage, error) {
	row := q.db.QueryRowContext(ctx, getDahuaEventData, id)
	var data json.RawMessage
	err := row.Scan(&data)
	return data, err
}

const getDahuaEventRule = `-- name: GetDahuaEventRule :one
SELECT id, code, ignore_db, ignore_live, ignore_mqtt FROM dahua_event_rules
WHERE id = ?
`

func (q *Queries) GetDahuaEventRule(ctx context.Context, id int64) (DahuaEventRule, error) {
	row := q.db.QueryRowContext(ctx, getDahuaEventRule, id)
	var i DahuaEventRule
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.IgnoreDb,
		&i.IgnoreLive,
		&i.IgnoreMqtt,
	)
	return i, err
}

const getDahuaFileByFilePath = `-- name: GetDahuaFileByFilePath :one
SELECT id, device_id, channel, start_time, end_time, length, type, file_path, duration, disk, video_stream, flags, events, cluster, "partition", pic_index, repeat, work_dir, work_dir_sn, updated_at
FROM dahua_files
WHERE device_id = ? and file_path = ?
`

type GetDahuaFileByFilePathParams struct {
	DeviceID int64
	FilePath string
}

func (q *Queries) GetDahuaFileByFilePath(ctx context.Context, arg GetDahuaFileByFilePathParams) (DahuaFile, error) {
	row := q.db.QueryRowContext(ctx, getDahuaFileByFilePath, arg.DeviceID, arg.FilePath)
	var i DahuaFile
	err := row.Scan(
		&i.ID,
		&i.DeviceID,
		&i.Channel,
		&i.StartTime,
		&i.EndTime,
		&i.Length,
		&i.Type,
		&i.FilePath,
		&i.Duration,
		&i.Disk,
		&i.VideoStream,
		&i.Flags,
		&i.Events,
		&i.Cluster,
		&i.Partition,
		&i.PicIndex,
		&i.Repeat,
		&i.WorkDir,
		&i.WorkDirSn,
		&i.UpdatedAt,
	)
	return i, err
}

const getSettings = `-- name: GetSettings :one
SELECT site_name, default_location FROM settings
LIMIT 1
`

func (q *Queries) GetSettings(ctx context.Context) (Setting, error) {
	row := q.db.QueryRowContext(ctx, getSettings)
	var i Setting
	err := row.Scan(&i.SiteName, &i.DefaultLocation)
	return i, err
}

const listDahuaDevice = `-- name: ListDahuaDevice :many
SELECT dahua_devices.id, dahua_devices.name, dahua_devices.address, dahua_devices.username, dahua_devices.password, dahua_devices.location, dahua_devices.created_at, dahua_devices.updated_at, coalesce(seed, id) FROM dahua_devices 
LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
`

type ListDahuaDeviceRow struct {
	ID        int64
	Name      string
	Address   string
	Username  string
	Password  string
	Location  types.Location
	CreatedAt types.Time
	UpdatedAt types.Time
	Seed      int64
}

func (q *Queries) ListDahuaDevice(ctx context.Context) ([]ListDahuaDeviceRow, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaDevice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDahuaDeviceRow
	for rows.Next() {
		var i ListDahuaDeviceRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Username,
			&i.Password,
			&i.Location,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Seed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDahuaDeviceByIDs = `-- name: ListDahuaDeviceByIDs :many
SELECT dahua_devices.id, dahua_devices.name, dahua_devices.address, dahua_devices.username, dahua_devices.password, dahua_devices.location, dahua_devices.created_at, dahua_devices.updated_at, coalesce(seed, id) FROM dahua_devices 
LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
WHERE id IN (/*SLICE:ids*/?)
`

type ListDahuaDeviceByIDsRow struct {
	ID        int64
	Name      string
	Address   string
	Username  string
	Password  string
	Location  types.Location
	CreatedAt types.Time
	UpdatedAt types.Time
	Seed      int64
}

func (q *Queries) ListDahuaDeviceByIDs(ctx context.Context, ids []int64) ([]ListDahuaDeviceByIDsRow, error) {
	query := listDahuaDeviceByIDs
	var queryParams []interface{}
	if len(ids) > 0 {
		for _, v := range ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDahuaDeviceByIDsRow
	for rows.Next() {
		var i ListDahuaDeviceByIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Username,
			&i.Password,
			&i.Location,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Seed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDahuaEventActions = `-- name: ListDahuaEventActions :many
SELECT DISTINCT action FROM dahua_events
`

func (q *Queries) ListDahuaEventActions(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaEventActions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var action string
		if err := rows.Scan(&action); err != nil {
			return nil, err
		}
		items = append(items, action)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDahuaEventCodes = `-- name: ListDahuaEventCodes :many
SELECT DISTINCT code FROM dahua_events
`

func (q *Queries) ListDahuaEventCodes(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaEventCodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		items = append(items, code)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDahuaEventRule = `-- name: ListDahuaEventRule :many
SELECT id, code, ignore_db, ignore_live, ignore_mqtt FROM dahua_event_rules
`

func (q *Queries) ListDahuaEventRule(ctx context.Context) ([]DahuaEventRule, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaEventRule)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaEventRule
	for rows.Next() {
		var i DahuaEventRule
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.IgnoreDb,
			&i.IgnoreLive,
			&i.IgnoreMqtt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDahuaEventWorkerState = `-- name: ListDahuaEventWorkerState :many
SELECT id, device_id, state, error, created_at,max(created_at) FROM dahua_event_worker_states GROUP BY device_id
`

type ListDahuaEventWorkerStateRow struct {
	ID        int64
	DeviceID  int64
	State     models.DahuaEventWorkerState
	Error     sql.NullString
	CreatedAt types.Time
	Max       interface{}
}

func (q *Queries) ListDahuaEventWorkerState(ctx context.Context) ([]ListDahuaEventWorkerStateRow, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaEventWorkerState)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDahuaEventWorkerStateRow
	for rows.Next() {
		var i ListDahuaEventWorkerStateRow
		if err := rows.Scan(
			&i.ID,
			&i.DeviceID,
			&i.State,
			&i.Error,
			&i.CreatedAt,
			&i.Max,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDahuaFileCursor = `-- name: ListDahuaFileCursor :many
SELECT 
  c.device_id, c.quick_cursor, c.full_cursor, c.full_epoch, c.full_complete, c.percent,
  count(f.device_id) AS files,
  coalesce(l.touched_at > ?, false) AS locked
FROM dahua_file_cursors AS c
LEFT JOIN dahua_files AS f ON f.device_id = c.device_id
LEFT JOIN dahua_file_scan_locks AS l ON l.device_id = c.device_id
GROUP BY c.device_id
`

type ListDahuaFileCursorRow struct {
	DeviceID     int64
	QuickCursor  types.Time
	FullCursor   types.Time
	FullEpoch    types.Time
	FullComplete bool
	Percent      float64
	Files        int64
	Locked       interface{}
}

func (q *Queries) ListDahuaFileCursor(ctx context.Context, touchedAt types.Time) ([]ListDahuaFileCursorRow, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaFileCursor, touchedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDahuaFileCursorRow
	for rows.Next() {
		var i ListDahuaFileCursorRow
		if err := rows.Scan(
			&i.DeviceID,
			&i.QuickCursor,
			&i.FullCursor,
			&i.FullEpoch,
			&i.FullComplete,
			&i.Percent,
			&i.Files,
			&i.Locked,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDahuaFileTypes = `-- name: ListDahuaFileTypes :many
SELECT DISTINCT type
FROM dahua_files
`

func (q *Queries) ListDahuaFileTypes(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaFileTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var type_ string
		if err := rows.Scan(&type_); err != nil {
			return nil, err
		}
		items = append(items, type_)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const touchDahuaFileScanLock = `-- name: TouchDahuaFileScanLock :exec
UPDATE dahua_file_scan_locks
SET touched_at = ?
WHERE device_id = ?
`

type TouchDahuaFileScanLockParams struct {
	TouchedAt types.Time
	DeviceID  int64
}

func (q *Queries) TouchDahuaFileScanLock(ctx context.Context, arg TouchDahuaFileScanLockParams) error {
	_, err := q.db.ExecContext(ctx, touchDahuaFileScanLock, arg.TouchedAt, arg.DeviceID)
	return err
}

const updateDahuaDevice = `-- name: UpdateDahuaDevice :one
UPDATE dahua_devices 
SET name = ?, address = ?, username = ?, password = ?, location = ?, updated_at = ?
WHERE id = ?
RETURNING id
`

type UpdateDahuaDeviceParams struct {
	Name      string
	Address   string
	Username  string
	Password  string
	Location  types.Location
	UpdatedAt types.Time
	ID        int64
}

func (q *Queries) UpdateDahuaDevice(ctx context.Context, arg UpdateDahuaDeviceParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateDahuaDevice,
		arg.Name,
		arg.Address,
		arg.Username,
		arg.Password,
		arg.Location,
		arg.UpdatedAt,
		arg.ID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateDahuaEventRule = `-- name: UpdateDahuaEventRule :exec
UPDATE dahua_event_rules 
SET 
  code = ?,
  ignore_db = ?,
  ignore_live = ?,
  ignore_mqtt = ?
WHERE id = ?
`

type UpdateDahuaEventRuleParams struct {
	Code       string
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
	ID         int64
}

func (q *Queries) UpdateDahuaEventRule(ctx context.Context, arg UpdateDahuaEventRuleParams) error {
	_, err := q.db.ExecContext(ctx, updateDahuaEventRule,
		arg.Code,
		arg.IgnoreDb,
		arg.IgnoreLive,
		arg.IgnoreMqtt,
		arg.ID,
	)
	return err
}

const updateDahuaFile = `-- name: UpdateDahuaFile :one
UPDATE dahua_files 
SET 
  channel = ?,
  start_time = ?,
  end_time = ?,
  length = ?,
  type = ?,
  duration = ?,
  disk = ?,
  video_stream = ?,
  flags = ?,
  events = ?,
  cluster = ?,
  partition = ?,
  pic_index = ?,
  repeat = ?,
  work_dir = ?,
  work_dir_sn = ?,
  updated_at = ?
WHERE device_id = ? AND file_path = ?
RETURNING id
`

type UpdateDahuaFileParams struct {
	Channel     int64
	StartTime   types.Time
	EndTime     types.Time
	Length      int64
	Type        string
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
	DeviceID    int64
	FilePath    string
}

func (q *Queries) UpdateDahuaFile(ctx context.Context, arg UpdateDahuaFileParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateDahuaFile,
		arg.Channel,
		arg.StartTime,
		arg.EndTime,
		arg.Length,
		arg.Type,
		arg.Duration,
		arg.Disk,
		arg.VideoStream,
		arg.Flags,
		arg.Events,
		arg.Cluster,
		arg.Partition,
		arg.PicIndex,
		arg.Repeat,
		arg.WorkDir,
		arg.WorkDirSn,
		arg.UpdatedAt,
		arg.DeviceID,
		arg.FilePath,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateDahuaFileCursor = `-- name: UpdateDahuaFileCursor :one
UPDATE dahua_file_cursors
SET 
  quick_cursor = ?,
  full_cursor = ?,
  full_epoch = ?,
  percent = ?
WHERE device_id = ?
RETURNING device_id, quick_cursor, full_cursor, full_epoch, full_complete, percent
`

type UpdateDahuaFileCursorParams struct {
	QuickCursor types.Time
	FullCursor  types.Time
	FullEpoch   types.Time
	Percent     float64
	DeviceID    int64
}

func (q *Queries) UpdateDahuaFileCursor(ctx context.Context, arg UpdateDahuaFileCursorParams) (DahuaFileCursor, error) {
	row := q.db.QueryRowContext(ctx, updateDahuaFileCursor,
		arg.QuickCursor,
		arg.FullCursor,
		arg.FullEpoch,
		arg.Percent,
		arg.DeviceID,
	)
	var i DahuaFileCursor
	err := row.Scan(
		&i.DeviceID,
		&i.QuickCursor,
		&i.FullCursor,
		&i.FullEpoch,
		&i.FullComplete,
		&i.Percent,
	)
	return i, err
}

const updateDahuaFileCursorPercent = `-- name: UpdateDahuaFileCursorPercent :one
UPDATE dahua_file_cursors 
SET
  percent = ?
WHERE device_id = ?
RETURNING device_id, quick_cursor, full_cursor, full_epoch, full_complete, percent
`

type UpdateDahuaFileCursorPercentParams struct {
	Percent  float64
	DeviceID int64
}

func (q *Queries) UpdateDahuaFileCursorPercent(ctx context.Context, arg UpdateDahuaFileCursorPercentParams) (DahuaFileCursor, error) {
	row := q.db.QueryRowContext(ctx, updateDahuaFileCursorPercent, arg.Percent, arg.DeviceID)
	var i DahuaFileCursor
	err := row.Scan(
		&i.DeviceID,
		&i.QuickCursor,
		&i.FullCursor,
		&i.FullEpoch,
		&i.FullComplete,
		&i.Percent,
	)
	return i, err
}

const updateSettings = `-- name: UpdateSettings :one
UPDATE settings
SET
  default_location = coalesce(?1, default_location),
  site_name = coalesce(?2, site_name)
WHERE 1 = 1
RETURNING site_name, default_location
`

type UpdateSettingsParams struct {
	DefaultLocation types.Location
	SiteName        sql.NullString
}

func (q *Queries) UpdateSettings(ctx context.Context, arg UpdateSettingsParams) (Setting, error) {
	row := q.db.QueryRowContext(ctx, updateSettings, arg.DefaultLocation, arg.SiteName)
	var i Setting
	err := row.Scan(&i.SiteName, &i.DefaultLocation)
	return i, err
}

const allocateDahuaSeed = `-- name: allocateDahuaSeed :exec
UPDATE dahua_seeds 
SET device_id = ?1
WHERE seed = (SELECT seed FROM dahua_seeds WHERE device_id = ?1 OR device_id IS NULL ORDER BY device_id asc LIMIT 1)
`

func (q *Queries) allocateDahuaSeed(ctx context.Context, deviceID sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, allocateDahuaSeed, deviceID)
	return err
}

const createDahuaDevice = `-- name: createDahuaDevice :one
INSERT INTO dahua_devices (
  name, address, username, password, location, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
) RETURNING id
`

type createDahuaDeviceParams struct {
	Name      string
	Address   string
	Username  string
	Password  string
	Location  types.Location
	CreatedAt types.Time
	UpdatedAt types.Time
}

func (q *Queries) createDahuaDevice(ctx context.Context, arg createDahuaDeviceParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createDahuaDevice,
		arg.Name,
		arg.Address,
		arg.Username,
		arg.Password,
		arg.Location,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createDahuaFileCursor = `-- name: createDahuaFileCursor :exec
INSERT INTO dahua_file_cursors (
  device_id,
  quick_cursor,
  full_cursor,
  full_epoch,
  percent
) VALUES (
  ?, ?, ?, ?, ?
)
`

type createDahuaFileCursorParams struct {
	DeviceID    int64
	QuickCursor types.Time
	FullCursor  types.Time
	FullEpoch   types.Time
	Percent     float64
}

func (q *Queries) createDahuaFileCursor(ctx context.Context, arg createDahuaFileCursorParams) error {
	_, err := q.db.ExecContext(ctx, createDahuaFileCursor,
		arg.DeviceID,
		arg.QuickCursor,
		arg.FullCursor,
		arg.FullEpoch,
		arg.Percent,
	)
	return err
}

const dahuaDeviceExists = `-- name: dahuaDeviceExists :one
SELECT COUNT(id) FROM dahua_devices WHERE id = ?
`

func (q *Queries) dahuaDeviceExists(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaDeviceExists, id)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getDahuaEventRuleByEvent = `-- name: getDahuaEventRuleByEvent :many
SELECT 
  ignore_db,
  ignore_live,
  ignore_mqtt,
  code
FROM dahua_event_device_rules 
WHERE device_id = ?1 AND (dahua_event_device_rules.code = ?2 OR dahua_event_device_rules.code = '')
UNION ALL
SELECT 
  ignore_db,
  ignore_live,
  ignore_mqtt,
  code
FROM dahua_event_rules
WHERE dahua_event_rules.code = ?2 OR dahua_event_rules.code = ''
ORDER BY code DESC
`

type getDahuaEventRuleByEventParams struct {
	DeviceID int64
	Code     string
}

type getDahuaEventRuleByEventRow struct {
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
	Code       string
}

func (q *Queries) getDahuaEventRuleByEvent(ctx context.Context, arg getDahuaEventRuleByEventParams) ([]getDahuaEventRuleByEventRow, error) {
	rows, err := q.db.QueryContext(ctx, getDahuaEventRuleByEvent, arg.DeviceID, arg.Code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []getDahuaEventRuleByEventRow
	for rows.Next() {
		var i getDahuaEventRuleByEventRow
		if err := rows.Scan(
			&i.IgnoreDb,
			&i.IgnoreLive,
			&i.IgnoreMqtt,
			&i.Code,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
