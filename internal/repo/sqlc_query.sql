-- name: createDahuaCamera :one
INSERT INTO dahua_cameras (
  name, address, username, password, location, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
) RETURNING id;

-- name: dahuaCameraExists :one
SELECT COUNT(id) FROM dahua_cameras WHERE id = ?;

-- name: UpdateDahuaCamera :one
UPDATE dahua_cameras 
SET name = ?, address = ?, username = ?, password = ?, location = ?, updated_at = ?
WHERE id = ?
RETURNING id;

-- name: GetDahuaCamera :one
SELECT dahua_cameras.*, coalesce(seed, id) FROM dahua_cameras 
LEFT JOIN dahua_seeds ON dahua_seeds.camera_id = dahua_cameras.id
WHERE id = ? LIMIT 1;

-- name: ListDahuaCamera :many
SELECT dahua_cameras.*, coalesce(seed, id) FROM dahua_cameras 
LEFT JOIN dahua_seeds ON dahua_seeds.camera_id = dahua_cameras.id;

-- name: ListDahuaCameraByIDs :many
SELECT dahua_cameras.*, coalesce(seed, id) FROM dahua_cameras 
LEFT JOIN dahua_seeds ON dahua_seeds.camera_id = dahua_cameras.id
WHERE id IN (sqlc.slice('ids'));

-- name: DeleteDahuaCamera :exec
DELETE FROM dahua_cameras WHERE id = ?;

-- name: GetSettings :one
SELECT * FROM settings
LIMIT 1;

-- name: UpdateSettings :one
UPDATE settings
SET
  default_location = coalesce(sqlc.narg('default_location'), default_location),
  site_name = coalesce(sqlc.narg('site_name'), site_name)
WHERE 1 = 1
RETURNING *;

-- name: allocateDahuaSeed :exec
UPDATE dahua_seeds 
SET camera_id = ?1
WHERE seed = (SELECT seed FROM dahua_seeds WHERE camera_id = ?1 OR camera_id IS NULL ORDER BY camera_id asc LIMIT 1);

-- name: CreateDahuaFileScanLock :one
INSERT INTO dahua_file_scan_locks (
  camera_id, touched_at
) VALUES (
  ?, ?
) RETURNING *;

-- name: DeleteDahuaFileScanLock :exec
DELETE FROM dahua_file_scan_locks WHERE camera_id = ?;

-- name: DeleteDahuaFileScanLockByAge :exec
DELETE FROM dahua_file_scan_locks WHERE touched_at < ?;

-- name: TouchDahuaFileScanLock :exec
UPDATE dahua_file_scan_locks
SET touched_at = ?
WHERE camera_id = ?;

-- name: UpdateDahuaFileCursorPercent :one
UPDATE dahua_file_cursors 
SET
  percent = ?
WHERE camera_id = ?
RETURNING *;

-- name: ListDahuaFileCursor :many
SELECT 
  c.*,
  count(f.camera_id) AS files,
  coalesce(l.touched_at > ?, false) AS locked
FROM dahua_file_cursors AS c
LEFT JOIN dahua_files AS f ON f.camera_id = c.camera_id
LEFT JOIN dahua_file_scan_locks AS l ON l.camera_id = c.camera_id
GROUP BY c.camera_id;

-- name: UpdateDahuaFileCursor :one
UPDATE dahua_file_cursors
SET 
  quick_cursor = ?,
  full_cursor = ?,
  full_epoch = ?,
  percent = ?
WHERE camera_id = ?
RETURNING *;

-- name: createDahuaFileCursor :exec
INSERT INTO dahua_file_cursors (
  camera_id,
  quick_cursor,
  full_cursor,
  full_epoch,
  percent
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: ListDahuaFileTypes :many
SELECT DISTINCT type
FROM dahua_files;

-- name: CreateDahuaFile :one
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
) RETURNING id;

-- name: GetDahuaFileByFilePath :one
SELECT *
FROM dahua_files
WHERE camera_id = ? and file_path = ?;

-- name: UpdateDahuaFile :one
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
RETURNING id;

-- name: DeleteDahuaFile :exec
DELETE FROM dahua_files
WHERE
  updated_at < sqlc.arg('updated_at') AND
  camera_id = sqlc.arg('camera_id') AND
  start_time <= sqlc.arg('end') AND
  sqlc.arg('start') < start_time;

-- name: CreateDahuaEvent :one
INSERT INTO dahua_events (
  camera_id,
  code,
  action,
  `index`,
  data,
  created_at
) VALUES (
  ?, ?, ?, ?, ?, ?
) RETURNING id;

-- name: ListDahuaEventCodes :many
SELECT DISTINCT code FROM dahua_events;

-- name: ListDahuaEventActions :many
SELECT DISTINCT action FROM dahua_events;

-- name: GetDahuaEventData :one
SELECT data FROM dahua_events WHERE id = ?;

-- name: getDahuaEventRuleByEvent :many
SELECT 
  ignore_db,
  ignore_live,
  ignore_mqtt,
  code
FROM dahua_event_camera_rules 
WHERE camera_id = sqlc.arg('camera_id') AND (dahua_event_camera_rules.code = sqlc.arg('code') OR dahua_event_camera_rules.code = '')
UNION ALL
SELECT 
  ignore_db,
  ignore_live,
  ignore_mqtt,
  code
FROM dahua_event_rules
WHERE dahua_event_rules.code = sqlc.arg('code') OR dahua_event_rules.code = ''
ORDER BY code DESC;

-- name: GetDahuaEventRule :one
SELECT * FROM dahua_event_rules
WHERE id = ?;

-- name: ListDahuaEventRule :many
SELECT * FROM dahua_event_rules;

-- name: UpdateDahuaEventRule :exec
UPDATE dahua_event_rules 
SET 
  code = ?,
  ignore_db = ?,
  ignore_live = ?,
  ignore_mqtt = ?
WHERE id = ?;

-- name: CreateDahuaEventRule :exec
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
);

-- name: DeleteDahuaEventRule :exec
DELETE FROM dahua_event_rules WHERE id = ?;

-- name: CreateDahuaEventWorkerState :exec
INSERT INTO dahua_event_worker_states(
  camera_id,
  state,
  error,
  created_at
) VALUES(
  ?,
  ?,
  ?,
  ?
);

-- name: ListDahuaEventWorkerState :many
SELECT *,max(created_at) FROM dahua_event_worker_states GROUP BY camera_id;
