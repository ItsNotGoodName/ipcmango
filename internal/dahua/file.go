package dahua

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/jlaffaye/ftp"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func OpenFileFTP(ctx context.Context, db *sqlx.DB, filePath string) (io.ReadCloser, int64, error) {
	urL, err := url.Parse(filePath)
	if err != nil {
		return nil, 0, err
	}

	var dest StorageDestination
	err = db.GetContext(ctx, &dest, `
		SELECT * FROM dahua_storage_destinations
		WHERE server_address = ? AND storage = ?
	`, urL.Host, StorageFTP)
	if err != nil {
		return nil, 0, err
	}

	c, err := ftp.Dial(core.Address(dest.Server_Address, int(dest.Port)), ftp.DialWithContext(ctx))
	if err != nil {
		return nil, 0, err
	}

	err = c.Login(dest.Username, dest.Password)
	if err != nil {
		return nil, 0, err
	}

	username := "/" + dest.Username
	path, _ := strings.CutPrefix(urL.Path, username)

	contentLength, err := c.FileSize(path)
	if err != nil {
		c.Quit()
		return nil, 0, err
	}

	rd, err := c.Retr(path)
	if err != nil {
		c.Quit()
		return nil, 0, err
	}

	return core.MultiReadCloser{
		Reader:  rd,
		Closers: []func() error{rd.Close, c.Quit},
	}, contentLength, nil
}

func OpenFileSFTP(ctx context.Context, db *sqlx.DB, filePath string) (io.ReadCloser, int64, error) {
	urL, err := url.Parse(filePath)
	if err != nil {
		return nil, 0, err
	}

	var dest StorageDestination
	err = db.GetContext(ctx, &dest, `
		SELECT * FROM dahua_storage_destinations
		WHERE server_address = ? AND storage = ?
	`, urL.Host, StorageFTP)
	if err != nil {
		return nil, 0, err
	}

	conn, err := ssh.Dial("tcp", core.Address(dest.Server_Address, int(dest.Port)), &ssh.ClientConfig{
		User: dest.Username,
		Auth: []ssh.AuthMethod{ssh.Password(dest.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// TODO: check public key
			return nil
		},
	})
	if err != nil {
		return nil, 0, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, 0, err
	}

	username := "/" + dest.Username
	path, _ := strings.CutPrefix(urL.Path, username)

	var contentLength int64
	if stat, err := client.Stat(path); err == nil {
		contentLength = stat.Size()
	}

	rd, err := client.Open(path)
	if err != nil {
		client.Close()
		return nil, 0, err
	}

	return core.MultiReadCloser{
		Reader:  rd,
		Closers: []func() error{rd.Close, client.Close},
	}, contentLength, nil
}

func OpenFileLocal(ctx context.Context, client Client, filePath string) (io.ReadCloser, int64, error) {
	v, err := client.File.Do(ctx, dahuarpc.LoadFileURL(client.URL, filePath), dahuarpc.Cookie(client.RPC.Session(ctx)))
	if err != nil {
		return nil, 0, err
	}
	return v, v.ContentLength, nil
}

type File struct {
	ID           string
	Device_ID    int64
	Channel      int
	Start_Time   types.Time
	End_Time     types.Time
	Length       int64
	Type         string
	File_Path    string
	Duration     int64
	Disk         int64
	Video_Stream string
	Flags        types.Slice[string]
	Events       types.Slice[string]
	Cluster      int64
	Partition    int64
	Pic_Index    int64
	Repeat       int64
	Work_Dir     string
	Work_Dir_Sn  bool
	Storage      Storage
	Updated_At   types.Time
}

type FileScanCursor struct {
	Device_ID     int64
	Quick_Cursor  types.Time
	Full_Cursor   types.Time
	Full_Epoch    types.Time
	Full_Complete bool
	Updated_At    types.Time
}

// FileScanEpoch is the oldest a file can be.
var FileScanEpoch time.Time = core.Must2(time.ParseInLocation(time.DateTime, "2009-12-31 00:00:00", time.UTC))

const (
	// fileScanVolatilePeriod is latest time period that the device could still be writing files to disk.
	fileScanVolatilePeriod = 8 * time.Hour
	// fileScanMaxPeriod is the maximum time period a device can handle when scanning files before they give weird results.
	fileScanMaxPeriod = 30 * 24 * time.Hour
)

func NewFileCursorQuick(end time.Time) time.Time {
	return core.Oldest(end, time.Now().Add(-fileScanVolatilePeriod))
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

func NewScanRange(start, end time.Time) (FileScanRange, FileSubScan, bool) {
	scanRange := FileScanRange{
		Start: start,
		End:   end,
	}
	scanCursor, ok := scanRange.NextSubScan(start)
	return scanRange, scanCursor, ok
}

type FileScanRange struct {
	Start time.Time
	End   time.Time
}

type FileSubScan struct {
	ScanStart time.Time
	ScanEnd   time.Time
	Cursor    time.Time
}

func (r FileScanRange) NextSubScan(cursor time.Time) (FileSubScan, bool) {
	if r.Start.Equal(r.End) || cursor.Equal(r.End) {
		return FileSubScan{}, false
	}

	if r.Start.Before(r.End) {
		nextCursor := core.Oldest(cursor.Add(fileScanMaxPeriod), r.End)
		return FileSubScan{
			ScanStart: cursor,
			ScanEnd:   nextCursor,
			Cursor:    nextCursor,
		}, true
	} else {
		nextCursor := core.Newest(cursor.Add(-fileScanMaxPeriod), r.End)
		return FileSubScan{
			ScanStart: cursor,
			ScanEnd:   nextCursor,
			Cursor:    nextCursor,
		}, true
	}
}

func (r FileScanRange) Percent(cursor time.Time) float64 {
	top := cursor.Sub(r.Start).Abs().Hours()
	bottom := r.End.Sub(r.Start).Abs().Hours()
	percent := (top / bottom) * 100
	return percent
}

func ScanFilesManual(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64, start, end time.Time) (ScanFilesResult, error) {
	var result ScanFilesResult
	scanRange, scanCursor, ok := NewScanRange(start, end)
	for ok {
		scanResult, err := scanFiles(ctx, db, conn, deviceID, scanCursor.ScanStart, scanCursor.ScanEnd)
		if err != nil {
			return ScanFilesResult{}, nil
		}
		result.add(scanResult)

		scanCursor, ok = scanRange.NextSubScan(scanCursor.Cursor)
	}
	return result, nil
}

func ScanFilesFull(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64) (ScanFilesResult, error) {
	var fileScanCursor FileScanCursor
	err := db.GetContext(ctx, &fileScanCursor, `
		SELECT * FROM dahua_file_cursors WHERE device_id = ?
	`, deviceID)
	if err != nil {
		return ScanFilesResult{}, nil
	}

	var result ScanFilesResult
	scanRange, scanCursor, ok := NewScanRange(fileScanCursor.Full_Cursor.Time, fileScanCursor.Full_Epoch.Time)
	for ok {
		scanResult, err := scanFiles(ctx, db, conn, deviceID, scanCursor.ScanStart, scanCursor.ScanEnd)
		if err != nil {
			return ScanFilesResult{}, nil
		}
		result.add(scanResult)

		fullCursor := types.NewTime(NewFileCursorFull(scanCursor.ScanStart, fileScanCursor.Full_Epoch.Time))

		_, err = db.ExecContext(ctx, `
			UPDATE dahua_file_cursors SET full_cursor = ? WHERE device_id = ?
		`, fullCursor, deviceID)
		if err != nil {
			return ScanFilesResult{}, nil
		}

		scanCursor, ok = scanRange.NextSubScan(scanCursor.Cursor)
	}

	_, err = db.ExecContext(ctx, `
		UPDATE dahua_file_cursors SET full_cursor = full_epoch WHERE device_id = ?
	`, deviceID)
	if err != nil {
		return ScanFilesResult{}, nil
	}

	return result, nil
}

func ScanFilesQuick(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64) (ScanFilesResult, error) {
	var fileScanCursor FileScanCursor
	err := db.GetContext(ctx, &fileScanCursor, `
		SELECT * FROM dahua_file_cursors WHERE device_id = ?
	`, deviceID)
	if err != nil {
		return ScanFilesResult{}, nil
	}

	var result ScanFilesResult
	scanRange, scanCursor, ok := NewScanRange(fileScanCursor.Quick_Cursor.Time, time.Now())
	for ok {
		scanResult, err := scanFiles(ctx, db, conn, deviceID, scanCursor.ScanStart, scanCursor.ScanEnd)
		if err != nil {
			return ScanFilesResult{}, nil
		}
		result.add(scanResult)

		quickCursor := types.NewTime(NewFileCursorQuick(scanCursor.ScanEnd))

		_, err = db.ExecContext(ctx, `
			UPDATE dahua_file_cursors SET quick_cursor = ? WHERE device_id = ?
		`, quickCursor, deviceID)
		if err != nil {
			return ScanFilesResult{}, nil
		}

		scanCursor, ok = scanRange.NextSubScan(scanCursor.Cursor)
	}

	return result, nil
}

type ScanFilesResult struct {
	CreatedCount int64 `json:"created_count"`
	UpdatedCount int64 `json:"updated_count"`
	DeletedCount int64 `json:"deleted_count"`
}

func (lhs *ScanFilesResult) add(rhs ScanFilesResult) {
	lhs.CreatedCount += rhs.CreatedCount
	lhs.UpdatedCount += rhs.UpdatedCount
	lhs.DeletedCount += rhs.DeletedCount
}

// scanFiles files on a device by a time range.
func scanFiles(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64, scanStart, scanEnd time.Time) (ScanFilesResult, error) {
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
		return ScanFilesResult{}, err
	}

	dahuaStart, dahuaEnd, dahuaOrder := scanStart, scanEnd, mediafilefind.ConditionOrderAscent
	if dahuaEnd.Before(dahuaStart) {
		dahuaStart, dahuaEnd, dahuaOrder = scanEnd, scanStart, mediafilefind.ConditionOrderDescent
	}

	condition := mediafilefind.NewCondtion(
		dahuarpc.NewTimestamp(dahuaStart, data.Location.Location),
		dahuarpc.NewTimestamp(dahuaEnd, data.Location.Location),
		dahuaOrder,
	)

	updatedAt := types.NewTime(time.Now())

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
				return ScanFilesResult{}, err
			}
		case "video":
			stream, err = mediafilefind.OpenStream(ctx, conn, condition.Video())
			if err != nil {
				return ScanFilesResult{}, err
			}
		default:
			panic(fmt.Sprintf("invalid kind %s", kind))
		}

		// Until stream is empty
		for {
			files, next, err := stream.Next(ctx)
			if err != nil {
				return ScanFilesResult{}, err
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
					return ScanFilesResult{}, err
				}
				count, err := result.RowsAffected()
				if err != nil {
					return ScanFilesResult{}, err
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
					return ScanFilesResult{}, err
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
	`, deviceID, types.NewTime(scanStart), types.NewTime(scanEnd), updatedAt)
	if err != nil {
		return ScanFilesResult{}, err
	}

	deletedCount, err := result.RowsAffected()
	if err != nil {
		return ScanFilesResult{}, err
	}

	return ScanFilesResult{
		CreatedCount: createdCount,
		UpdatedCount: updatedCount,
		DeletedCount: deletedCount,
	}, err
}
