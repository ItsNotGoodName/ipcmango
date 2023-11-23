// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

const createDahuaEvent = `-- name: CreateDahuaEvent :one
INSERT INTO dahua_events (
  camera_id,
  content_type,
  content_length,
  code,
  action,
  ` + "`" + `index` + "`" + `,
  data,
  created_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING id
`

type CreateDahuaEventParams struct {
	CameraID      int64
	ContentType   string
	ContentLength int64
	Code          string
	Action        string
	Index         int64
	Data          json.RawMessage
	CreatedAt     time.Time
}

func (q *Queries) CreateDahuaEvent(ctx context.Context, arg CreateDahuaEventParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createDahuaEvent,
		arg.CameraID,
		arg.ContentType,
		arg.ContentLength,
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

const createDahuaFile = `-- name: CreateDahuaFile :one
INSERT INTO dahua_files (
  camera_id,
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
	CameraID    int64
	Channel     int64
	StartTime   time.Time
	EndTime     time.Time
	Length      int64
	Type        string
	FilePath    string
	Duration    int64
	Disk        int64
	VideoStream string
	Flags       models.StringSlice
	Events      models.StringSlice
	Cluster     int64
	Partition   int64
	PicIndex    int64
	Repeat      int64
	WorkDir     string
	WorkDirSn   int64
	UpdatedAt   time.Time
}

func (q *Queries) CreateDahuaFile(ctx context.Context, arg CreateDahuaFileParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createDahuaFile,
		arg.CameraID,
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
  camera_id, created_at
) VALUES (
  ?, ?
) RETURNING camera_id, created_at
`

type CreateDahuaFileScanLockParams struct {
	CameraID  int64
	CreatedAt time.Time
}

func (q *Queries) CreateDahuaFileScanLock(ctx context.Context, arg CreateDahuaFileScanLockParams) (DahuaFileScanLock, error) {
	row := q.db.QueryRowContext(ctx, createDahuaFileScanLock, arg.CameraID, arg.CreatedAt)
	var i DahuaFileScanLock
	err := row.Scan(&i.CameraID, &i.CreatedAt)
	return i, err
}

const deleteDahuaCamera = `-- name: DeleteDahuaCamera :exec
DELETE FROM dahua_cameras WHERE id = ?
`

func (q *Queries) DeleteDahuaCamera(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteDahuaCamera, id)
	return err
}

const deleteDahuaFile = `-- name: DeleteDahuaFile :exec
DELETE FROM dahua_files
WHERE
  updated_at < ?1 AND
  camera_id = ?2 AND
  start_time <= ?3 AND
  ?4 < start_time
`

type DeleteDahuaFileParams struct {
	UpdatedAt time.Time
	CameraID  int64
	End       time.Time
	Start     time.Time
}

func (q *Queries) DeleteDahuaFile(ctx context.Context, arg DeleteDahuaFileParams) error {
	_, err := q.db.ExecContext(ctx, deleteDahuaFile,
		arg.UpdatedAt,
		arg.CameraID,
		arg.End,
		arg.Start,
	)
	return err
}

const deleteDahuaFileScanLock = `-- name: DeleteDahuaFileScanLock :exec
DELETE FROM dahua_file_scan_locks WHERE camera_id = ?
`

func (q *Queries) DeleteDahuaFileScanLock(ctx context.Context, cameraID int64) error {
	_, err := q.db.ExecContext(ctx, deleteDahuaFileScanLock, cameraID)
	return err
}

const getDahuaCamera = `-- name: GetDahuaCamera :one
SELECT id, name, address, username, password, location, created_at, coalesce(seed, id) FROM dahua_cameras 
LEFT JOIN dahua_seeds ON dahua_seeds.camera_id = dahua_cameras.id
WHERE id = ? LIMIT 1
`

type GetDahuaCameraRow struct {
	ID        int64
	Name      string
	Address   string
	Username  string
	Password  string
	Location  models.Location
	CreatedAt time.Time
	Seed      int64
}

func (q *Queries) GetDahuaCamera(ctx context.Context, id int64) (GetDahuaCameraRow, error) {
	row := q.db.QueryRowContext(ctx, getDahuaCamera, id)
	var i GetDahuaCameraRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Username,
		&i.Password,
		&i.Location,
		&i.CreatedAt,
		&i.Seed,
	)
	return i, err
}

const getDahuaFileCursor = `-- name: GetDahuaFileCursor :one
SELECT camera_id, quick_cursor, full_cursor, full_epoch, full_complete FROM dahua_file_cursors 
WHERE camera_id = ?
`

func (q *Queries) GetDahuaFileCursor(ctx context.Context, cameraID int64) (DahuaFileCursor, error) {
	row := q.db.QueryRowContext(ctx, getDahuaFileCursor, cameraID)
	var i DahuaFileCursor
	err := row.Scan(
		&i.CameraID,
		&i.QuickCursor,
		&i.FullCursor,
		&i.FullEpoch,
		&i.FullComplete,
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

const listDahuaCamera = `-- name: ListDahuaCamera :many
SELECT id, name, address, username, password, location, created_at, coalesce(seed, id) FROM dahua_cameras 
LEFT JOIN dahua_seeds ON dahua_seeds.camera_id = dahua_cameras.id
`

type ListDahuaCameraRow struct {
	ID        int64
	Name      string
	Address   string
	Username  string
	Password  string
	Location  models.Location
	CreatedAt time.Time
	Seed      int64
}

func (q *Queries) ListDahuaCamera(ctx context.Context) ([]ListDahuaCameraRow, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaCamera)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDahuaCameraRow
	for rows.Next() {
		var i ListDahuaCameraRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Username,
			&i.Password,
			&i.Location,
			&i.CreatedAt,
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

const listDahuaEvent = `-- name: ListDahuaEvent :many
SELECT id, camera_id, content_type, content_length, code, "action", ` + "`" + `index` + "`" + `, data, created_at FROM dahua_events
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

type ListDahuaEventParams struct {
	Limit  int64
	Offset int64
}

func (q *Queries) ListDahuaEvent(ctx context.Context, arg ListDahuaEventParams) ([]DahuaEvent, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaEvent, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaEvent
	for rows.Next() {
		var i DahuaEvent
		if err := rows.Scan(
			&i.ID,
			&i.CameraID,
			&i.ContentType,
			&i.ContentLength,
			&i.Code,
			&i.Action,
			&i.Index,
			&i.Data,
			&i.CreatedAt,
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
  c.camera_id, c.quick_cursor, c.full_cursor, c.full_epoch, c.full_complete,
  count(f.camera_id) as files
FROM dahua_file_cursors AS c
LEFT JOIN dahua_files as f ON f.camera_id = c.camera_id
GROUP BY c.camera_id
`

type ListDahuaFileCursorRow struct {
	CameraID     int64
	QuickCursor  time.Time
	FullCursor   time.Time
	FullEpoch    time.Time
	FullComplete bool
	Files        int64
}

func (q *Queries) ListDahuaFileCursor(ctx context.Context) ([]ListDahuaFileCursorRow, error) {
	rows, err := q.db.QueryContext(ctx, listDahuaFileCursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDahuaFileCursorRow
	for rows.Next() {
		var i ListDahuaFileCursorRow
		if err := rows.Scan(
			&i.CameraID,
			&i.QuickCursor,
			&i.FullCursor,
			&i.FullEpoch,
			&i.FullComplete,
			&i.Files,
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

const updateDahuaCamera = `-- name: UpdateDahuaCamera :one
UPDATE dahua_cameras 
SET name = ?, address = ?, username = ?, password = ?, location = ?
WHERE id = ?
RETURNING id
`

type UpdateDahuaCameraParams struct {
	Name     string
	Address  string
	Username string
	Password string
	Location models.Location
	ID       int64
}

func (q *Queries) UpdateDahuaCamera(ctx context.Context, arg UpdateDahuaCameraParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateDahuaCamera,
		arg.Name,
		arg.Address,
		arg.Username,
		arg.Password,
		arg.Location,
		arg.ID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
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
WHERE camera_id = ? AND file_path = ?
RETURNING id
`

type UpdateDahuaFileParams struct {
	Channel     int64
	StartTime   time.Time
	EndTime     time.Time
	Length      int64
	Type        string
	Duration    int64
	Disk        int64
	VideoStream string
	Flags       models.StringSlice
	Events      models.StringSlice
	Cluster     int64
	Partition   int64
	PicIndex    int64
	Repeat      int64
	WorkDir     string
	WorkDirSn   int64
	UpdatedAt   time.Time
	CameraID    int64
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
		arg.CameraID,
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
  full_epoch = ?
WHERE camera_id = ?
RETURNING camera_id, quick_cursor, full_cursor, full_epoch, full_complete
`

type UpdateDahuaFileCursorParams struct {
	QuickCursor time.Time
	FullCursor  time.Time
	FullEpoch   time.Time
	CameraID    int64
}

func (q *Queries) UpdateDahuaFileCursor(ctx context.Context, arg UpdateDahuaFileCursorParams) (DahuaFileCursor, error) {
	row := q.db.QueryRowContext(ctx, updateDahuaFileCursor,
		arg.QuickCursor,
		arg.FullCursor,
		arg.FullEpoch,
		arg.CameraID,
	)
	var i DahuaFileCursor
	err := row.Scan(
		&i.CameraID,
		&i.QuickCursor,
		&i.FullCursor,
		&i.FullEpoch,
		&i.FullComplete,
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
	DefaultLocation models.Location
	SiteName        sql.NullString
}

func (q *Queries) UpdateSettings(ctx context.Context, arg UpdateSettingsParams) (Setting, error) {
	row := q.db.QueryRowContext(ctx, updateSettings, arg.DefaultLocation, arg.SiteName)
	var i Setting
	err := row.Scan(&i.SiteName, &i.DefaultLocation)
	return i, err
}

const createDahuaCamera = `-- name: createDahuaCamera :one
INSERT INTO dahua_cameras (
  name, address, username, password, location, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
) RETURNING id
`

type createDahuaCameraParams struct {
	Name      string
	Address   string
	Username  string
	Password  string
	Location  models.Location
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) createDahuaCamera(ctx context.Context, arg createDahuaCameraParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createDahuaCamera,
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
  camera_id,
  quick_cursor,
  full_cursor,
  full_epoch
) VALUES (
  ?, ?, ?, ?
)
`

type createDahuaFileCursorParams struct {
	CameraID    int64
	QuickCursor time.Time
	FullCursor  time.Time
	FullEpoch   time.Time
}

func (q *Queries) createDahuaFileCursor(ctx context.Context, arg createDahuaFileCursorParams) error {
	_, err := q.db.ExecContext(ctx, createDahuaFileCursor,
		arg.CameraID,
		arg.QuickCursor,
		arg.FullCursor,
		arg.FullEpoch,
	)
	return err
}

const setDahuaSeed = `-- name: setDahuaSeed :exec
UPDATE dahua_seeds 
SET camera_id = ?1
WHERE seed = (SELECT seed FROM dahua_seeds WHERE camera_id = ?1 OR camera_id IS NULL ORDER BY camera_id asc LIMIT 1)
`

func (q *Queries) setDahuaSeed(ctx context.Context, cameraID sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, setDahuaSeed, cameraID)
	return err
}
