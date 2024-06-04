package dahua

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

// FileScanEpoch is the oldest a file can be.
var FileScanEpoch time.Time = core.Must2(time.ParseInLocation(time.DateTime, "2009-12-31 00:00:00", time.UTC))

const (
	// fileScanVolatilePeriod is latest time period that the device could still be writing files to disk.
	fileScanVolatilePeriod = 8 * time.Hour
	// fileScanMaxPeriod is the maximum time period a device can handle when scanning files before they give weird results.
	fileScanMaxPeriod = 30 * 24 * time.Hour
)

func NewFileCursorQuick(end, now time.Time) time.Time {
	return core.Oldest(end, now.Add(-fileScanVolatilePeriod))
}

func NewFileCursorFull(start, epoch time.Time) time.Time {
	return core.Newest(start, epoch)
}

func ResetFileScanCursor(ctx context.Context, db sqlx.ExecerContext, deviceID int64) error {
	now := time.Now()
	quickCursor := types.NewTime(now.Add(-fileScanVolatilePeriod))
	fullCursor := types.NewTime(now)
	fullEpoch := types.NewTime(FileScanEpoch)
	updatedAt := types.NewTime(now)

	_, err := db.ExecContext(ctx, `
		INSERT OR REPLACE INTO dahua_file_cursors (
			device_id,
			quick_cursor,
			full_cursor,
			full_epoch,
			updated_at
		) VALUES (?, ?, ?, ?, ?)
	`,
		deviceID,
		quickCursor,
		fullCursor,
		fullEpoch,
		updatedAt,
	)
	return err
}

func NewScanRange(start, end time.Time) (ScanRange, ScanCursor, bool) {
	scanRange := ScanRange{
		Start: start,
		End:   end,
	}
	scanCursor, ok := scanRange.NextScanCursor(start)
	return scanRange, scanCursor, ok
}

type ScanRange struct {
	Start time.Time
	End   time.Time
}

type ScanCursor struct {
	ScanStart time.Time
	ScanEnd   time.Time
	Cursor    time.Time
}

func (r ScanRange) NextScanCursor(cursor time.Time) (ScanCursor, bool) {
	if r.Start.Equal(r.End) || cursor.Equal(r.End) {
		return ScanCursor{}, false
	}

	if r.Start.Before(r.End) {
		nextCursor := core.Oldest(cursor.Add(fileScanMaxPeriod), r.End)
		return ScanCursor{
			ScanStart: cursor,
			ScanEnd:   nextCursor,
			Cursor:    nextCursor,
		}, true
	} else {
		nextCursor := core.Newest(cursor.Add(-fileScanMaxPeriod), r.End)
		return ScanCursor{
			ScanStart: cursor,
			ScanEnd:   nextCursor,
			Cursor:    nextCursor,
		}, true
	}
}

func (r ScanRange) Percent(cursor time.Time) float64 {
	top := cursor.Sub(r.Start).Abs().Hours()
	bottom := r.End.Sub(r.Start).Abs().Hours()
	percent := (top / bottom) * 100
	return percent
}

func NewCondition(ctx context.Context, scanStart, scanEnd time.Time, location *time.Location) mediafilefind.Condition {
	start, end, order := scanStart, scanEnd, mediafilefind.ConditionOrderAscent
	if scanStart.After(end) {
		start, end, order = scanEnd, scanStart, mediafilefind.ConditionOrderDescent
	}

	startTs, endTs := dahuarpc.NewTimestamp(start, location), dahuarpc.NewTimestamp(end, location)
	condition := mediafilefind.NewCondtion(startTs, endTs, order)
	return condition
}

type ScanResult struct {
	CreatedCount int64
	UpdatedCount int64
	DeletedCount int64
}

// Scan files on a device by a time range.
func Scan(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64, start, end time.Time) (ScanResult, error) {
	// Get device context
	var data struct {
		types.Key
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
		return ScanResult{}, err
	}

	updatedAt := types.NewTime(time.Now())

	dahuaStart, dahuaEnd, dahuaOrder := start, end, mediafilefind.ConditionOrderAscent
	if dahuaEnd.Before(dahuaStart) {
		dahuaStart, dahuaEnd, dahuaOrder = end, start, mediafilefind.ConditionOrderDescent
	}

	condition := mediafilefind.NewCondtion(
		dahuarpc.NewTimestamp(dahuaStart, data.Location.Location),
		dahuarpc.NewTimestamp(dahuaEnd, data.Location.Location),
		dahuaOrder,
	)

	var createdCount int64
	var updatedCount int64
	var stream *mediafilefind.Stream
	defer func() {
		if stream != nil {
			stream.Close()
		}
	}()

	// For picture and video conditions
	for _, kind := range []string{"picture", "video"} {
		switch kind {
		case "picture":
			stream, err = mediafilefind.OpenStream(ctx, conn, condition.Picture())
			if err != nil {
				return ScanResult{}, err
			}
		case "video":
			stream, err = mediafilefind.OpenStream(ctx, conn, condition.Video())
			if err != nil {
				return ScanResult{}, err
			}
		default:
			panic(fmt.Sprintf("invalid kind %s", kind))
		}

		// Until stream is empty
		for {
			files, next, err := stream.Next(ctx)
			if err != nil {
				return ScanResult{}, err
			}
			if !next {
				break
			}

			// For each file
			for _, v := range files {
				startTime, endTime, err := v.UniqueTime(int(data.Seed), data.Location.Location)
				if err != nil {
					slog.Error("Failed to get unique time for file", "error", err, "device", data.Name)
					continue
				}
				storage := StorageFromFilePath(v.FilePath)
				events := v.CleanEvents()

				result, err := db.ExecContext(ctx, `
						UPDATE dahua_files SET 
							channel = ?,
							start_time = ?,
							end_time = ?,
							length = ?,
							type = ?,
							file_path = ?,
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
							storage = ?,
							updated_at = ?
						WHERE device_id = ? AND file_path = ?
					`,
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
					types.NewSlice(events),
					v.Cluster,
					v.Partition,
					v.PicIndex,
					v.Repeat,
					v.WorkDir,
					v.WorkDirSN,
					storage,
					updatedAt,
					deviceID,
					v.FilePath,
				)
				if err != nil {
					return ScanResult{}, err
				}
				count, err := result.RowsAffected()
				if err != nil {
					return ScanResult{}, err
				}
				if count != 0 {
					updatedCount++
					continue
				}

				fileID := ulid.MustNew(ulid.Timestamp(startTime), ulid.DefaultEntropy()).String()

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
					types.NewSlice(events),
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
					return ScanResult{}, err
				}

				createdCount++
			}
		}

		stream.Close()
	}

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
		return ScanResult{}, err
	}

	deletedCount, err := result.RowsAffected()
	if err != nil {
		return ScanResult{}, err
	}

	return ScanResult{
		CreatedCount: createdCount,
		UpdatedCount: updatedCount,
		DeletedCount: deletedCount,
	}, err
}
