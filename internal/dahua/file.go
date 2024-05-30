package dahua

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

var FileScanEpoch time.Time = core.Must2(time.ParseInLocation(time.DateTime, "2009-12-31 00:00:00", time.UTC))

const FileScanMaxPeriod = 30 * 24 * time.Hour

func NewFileScanRange(start, end time.Time, period time.Duration, ascending bool) *FileScanRange {
	if period <= 0 {
		panic("period is too short")
	}
	if start.After(end) {
		panic("invalid time range")
	}

	cursor := end.Add(period)
	if ascending {
		cursor = start.Add(-period)
	}

	return &FileScanRange{
		start:     start,
		end:       end,
		period:    period,
		ascending: ascending,
		cursor:    cursor,
	}
}

type FileScanRange struct {
	start     time.Time
	end       time.Time
	period    time.Duration
	ascending bool

	cursor time.Time
}

func (r *FileScanRange) Cursor() time.Time {
	return r.cursor
}

func (r *FileScanRange) Percent() float64 {
	if r.ascending {
		if r.cursor.Equal(r.end) {
			return 100.0
		}
		return (r.cursor.Sub(r.start).Hours() / r.end.Sub(r.start).Hours()) * 100
	} else {
		if r.cursor.Equal(r.start) {
			return 100.0
		}
		return (r.end.Sub(r.cursor).Hours() / r.end.Sub(r.start).Hours()) * 100
	}
}

func (r *FileScanRange) Range() (time.Time, time.Time) {
	if r.ascending {
		end := r.cursor.Add(r.period)
		if end.After(r.end) {
			end = r.end
		}
		return r.cursor, end
	} else {
		start := r.cursor.Add(-r.period)
		if start.Before(r.start) {
			start = r.start
		}
		return start, r.cursor
	}
}

func (r *FileScanRange) Next() bool {
	var cursor time.Time
	if r.ascending {
		cursor = r.cursor.Add(r.period)
		if cursor.After(r.end) {
			r.cursor = r.end
			return false
		}
	} else {
		cursor = r.cursor.Add(-r.period)
		if cursor.Before(r.start) {
			r.cursor = r.start
			return false
		}
	}

	r.cursor = cursor

	return true
}

func NewCondition(ctx context.Context, scanRange *FileScanRange, location *time.Location) mediafilefind.Condition {
	start, end := scanRange.Range()
	startTs, endTs := dahuarpc.NewTimestamp(start, location), dahuarpc.NewTimestamp(end, location)
	condition := mediafilefind.NewCondtion(startTs, endTs)
	if scanRange.ascending {
		condition.Order = mediafilefind.ConditionOrderAscent
	} else {
		condition.Order = mediafilefind.ConditionOrderDescent
	}
	return condition
}

type FileScanResult struct {
	CreatedCount int64 `json:"created_count"`
	UpdatedCount int64 `json:"updated_count"`
	DeletedCount int64 `json:"deleted_count"`
}

func FileScan(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64, start, end time.Time) (FileScanResult, error) {
	var data struct {
		core.Key
		Location types.Location
		Seed     int64
		Name     string
	}
	err := db.GetContext(ctx, &data, `
		SELECT d.uuid, d.id, coalesce(d.location, s.location) AS location, d.seed, d.name
		FROM dahua_devices AS d, settings as s
		WHERE d.id = ?
	`, deviceID)
	if err != nil {
		return FileScanResult{}, err
	}

	updatedAt := types.NewTime(time.Now())

	scanRange := NewFileScanRange(start, end, FileScanMaxPeriod, true)

	var createdCount int64
	var updatedCount int64
	var stream *mediafilefind.Stream
	defer func() {
		if stream != nil {
			stream.Close()
		}
	}()

	// For each time range
	for scanRange.Next() {
		progress := scanRange.Percent()
		bus.Publish(bus.FileScanProgress{
			DeviceKey: data.Key,
			Progress:  progress,
		})
		condition := NewCondition(ctx, scanRange, data.Location.Location)

		// For picture and video conditions
		for _, kind := range []string{"picture", "video"} {
			switch kind {
			case "picture":
				stream, err = mediafilefind.OpenStream(ctx, conn, condition.Picture())
				if err != nil {
					return FileScanResult{}, err
				}
			case "video":
				stream, err = mediafilefind.OpenStream(ctx, conn, condition.Video())
				if err != nil {
					return FileScanResult{}, err
				}
			default:
				panic(fmt.Sprintf("invalid kind %s", kind))
			}

			// Until stream is empty
			for {
				files, next, err := stream.Next(ctx)
				if err != nil {
					return FileScanResult{}, err
				}
				if !next {
					break
				}

				// For each file
				for _, v := range files {
					result, err := db.ExecContext(ctx, `
						UPDATE dahua_files SET updated_at = ? WHERE device_id = ? AND file_path = ?
					`, updatedAt, deviceID, v.FilePath)
					if err != nil {
						return FileScanResult{}, err
					}
					count, err := result.RowsAffected()
					if err != nil {
						return FileScanResult{}, err
					}
					if count != 0 {
						updatedCount++
						continue
					}

					startTime, endTime, err := v.UniqueTime(int(data.Seed), data.Location.Location)
					if err != nil {
						slog.Error("Failed to get unique time for file", "error", err, "device", data.Name)
						continue
					}

					fileID := ulid.MustNew(ulid.Timestamp(startTime), ulid.DefaultEntropy()).String()
					storage := StorageFromFilePath(v.FilePath)

					_, err = db.ExecContext(ctx, `
						INSERT INTO dahua_files (
							id,
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
							storage,
							updated_at
						)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
						ON CONFLICT (start_time) DO UPDATE SET id = id RETURNING id -- Assume that files with the same time are really rare
					`,
						fileID,
						deviceID,
						v.Channel,
						types.NewTime(startTime),
						types.NewTime(endTime),
						v.Length,
						v.Type,
						v.FilePath,
						v.Duration,
						v.Disk,
						v.VideoStream,
						types.NewSlice(v.Flags),
						types.NewSlice(v.Events),
						v.Cluster,
						v.Partition,
						v.PicIndex,
						v.Repeat,
						v.WorkDir,
						v.WorkDirSN,
						storage,
						updatedAt,
					)
					if err != nil {
						return FileScanResult{}, err
					}

					createdCount++
				}
			}

			stream.Close()
		}
	}
	progress := scanRange.Percent()
	bus.Publish(bus.FileScanProgress{
		DeviceKey: data.Key,
		Progress:  progress,
	})

	// Delete stale files
	result, err := db.ExecContext(ctx, `
		DELETE FROM dahua_files
		WHERE
			device_id = ?
			AND ? < start_time
			AND start_time <= ?
			AND updated_at < ?
	`, deviceID, types.NewTime(start), types.NewTime(end), updatedAt)
	if err != nil {
		return FileScanResult{}, err
	}

	deletedCount, err := result.RowsAffected()
	if err != nil {
		return FileScanResult{}, err
	}

	return FileScanResult{
		CreatedCount: createdCount,
		UpdatedCount: updatedCount,
		DeletedCount: deletedCount,
	}, err
}
