package dahua

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
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

var FileScanEpoch time.Time = core.Must2(time.ParseInLocation(time.DateTime, "2009-12-31 00:00:00", time.UTC))

const (
	fileScanVolatilePeriod = 8 * time.Hour
	fileScanMaxPeriod      = 30 * 24 * time.Hour
)

func NewFileCursorQuick(end, now time.Time) time.Time {
	return core.Oldest(end, now.Add(-fileScanVolatilePeriod))
}

func NewFileCursorFull(start, epoch time.Time) time.Time {
	return core.Newest(start, epoch)
}

func NewDefaultFileCursor() DefaultFileCursor {
	now := time.Now()
	return DefaultFileCursor{
		QuickCursor: types.NewTime(now.Add(-fileScanVolatilePeriod)),
		FullCursor:  types.NewTime(now),
		FullEpoch:   types.NewTime(FileScanEpoch),
	}
}

type DefaultFileCursor struct {
	QuickCursor types.Time
	FullCursor  types.Time
	FullEpoch   types.Time
}

func NewFileScanRange(start, end time.Time) *FileScanRange {
	if start.After(end) {
		panic("invalid time range")
	}

	cursor := end.Add(fileScanMaxPeriod)

	return &FileScanRange{
		start:  start,
		end:    end,
		cursor: cursor,
	}
}

type FileScanRange struct {
	start time.Time
	end   time.Time

	cursor time.Time
}

func (r *FileScanRange) Cursor() time.Time {
	return r.cursor
}

func (r *FileScanRange) Percent() float64 {
	if r.cursor.Equal(r.start) {
		return 100.0
	}
	return (r.end.Sub(r.cursor).Hours() / r.end.Sub(r.start).Hours()) * 100
}

func (r *FileScanRange) Range() (time.Time, time.Time) {
	start := r.cursor.Add(-fileScanMaxPeriod)
	if start.Before(r.start) {
		start = r.start
	}
	return start, r.cursor
}

func (r *FileScanRange) Next() bool {
	cursor := r.cursor.Add(-fileScanMaxPeriod)
	if cursor.Before(r.start) {
		r.cursor = r.start
		return false
	}

	r.cursor = cursor

	return true
}

func NewCondition(ctx context.Context, scanRange *FileScanRange, location *time.Location) mediafilefind.Condition {
	start, end := scanRange.Range()
	startTs, endTs := dahuarpc.NewTimestamp(start, location), dahuarpc.NewTimestamp(end, location)
	condition := mediafilefind.NewCondtion(startTs, endTs)
	condition.Order = mediafilefind.ConditionOrderDescent
	return condition
}

// fileScan assumes that no other goroutines are calling this function with the same deviceID.
func fileScan(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64, scanRange *FileScanRange) error {
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
		return err
	}

	updatedAt := types.NewTime(time.Now())

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
		bus.Publish(bus.FileScanProgressed{
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
					return err
				}
			case "video":
				stream, err = mediafilefind.OpenStream(ctx, conn, condition.Video())
				if err != nil {
					return err
				}
			default:
				panic(fmt.Sprintf("invalid kind %s", kind))
			}

			// Until stream is empty
			for {
				files, next, err := stream.Next(ctx)
				if err != nil {
					return err
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
						return err
					}
					count, err := result.RowsAffected()
					if err != nil {
						return err
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
						return err
					}

					createdCount++
				}
			}

			stream.Close()
		}
	}
	progress := scanRange.Percent()
	bus.Publish(bus.FileScanProgressed{
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
	`, deviceID, types.NewTime(scanRange.start), types.NewTime(scanRange.end), updatedAt)
	if err != nil {
		return err
	}

	deletedCount, err := result.RowsAffected()
	if err != nil {
		return err
	}

	bus.Publish(bus.FileScanFinished{
		DeviceKey:    data.Key,
		CreatedCount: createdCount,
		UpdatedCount: updatedCount,
		DeletedCount: deletedCount,
	})

	return err
}

type FileScanJob struct {
	Command FileScanCommand
	Data    []FileScanData
}

// ENUM(quick,full,manual)
type FileScanCommand string

type FileScanData struct {
	DeviceID  int64
	StartTime time.Time
	EndTime   time.Time
}

func NewFileScanService(db *sqlx.DB, store *Store) FileScanService {
	return FileScanService{
		db:    db,
		store: store,
		jobC:  make(chan FileScanJob),
	}
}

type FileScanService struct {
	db    *sqlx.DB
	store *Store
	jobC  chan FileScanJob
}

func (w FileScanService) String() string {
	return "dahua.FileScanService"
}

func (w FileScanService) Serve(ctx context.Context) error {
	slog := slog.With("service", w.String())
	slog.Info("Started service")

	quickScanInterval := 30 * time.Second
	t := time.NewTicker(quickScanInterval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			w.handle(ctx, slog, FileScanJob{
				Command: FileScanCommandQuick,
			})
		case job := <-w.jobC:
			w.handle(ctx, slog, job)
		}
	}
}

func (w FileScanService) handle(ctx context.Context, slog *slog.Logger, job FileScanJob) error {
	type BatchJob struct {
		DeviceID  int64
		StartTime time.Time
		EndTime   time.Time
		FullEpoch time.Time
	}

	workers := 3
	var batchJobs []BatchJob
	switch job.Command {
	case FileScanCommandManual:
		slog = slog.With("type", "manual")
		batchJobs = make([]BatchJob, 0, len(job.Data))

		// Add batch jobs from job data
		for _, data := range job.Data {
			batchJobs = append(batchJobs, BatchJob{
				DeviceID:  data.DeviceID,
				StartTime: data.StartTime,
				EndTime:   data.EndTime,
			})
		}
	case FileScanCommandQuick, FileScanCommandFull:
		// Extract devices ids
		var deviceIDs []int64
		for _, data := range job.Data {
			deviceIDs = append(deviceIDs, data.DeviceID)
		}

		// Get cursors by device ids
		var cursors []FileScanCursor
		query, queryArgs, _ := sqlx.In(`
			SELECT * FROM dahua_file_cursors
			WHERE ? = 0 OR device_id IN (?)
		`, len(deviceIDs), deviceIDs)
		err := sqlx.SelectContext(ctx, w.db, &cursors, query, queryArgs...)
		if err != nil {
			return err
		}

		// Add batch jobs from cursors
		batchJobs := make([]BatchJob, 0, len(cursors))
		switch job.Command {
		case FileScanCommandFull:
			slog = slog.With("type", "full")

			for _, cursor := range cursors {
				batchJobs = append(batchJobs, BatchJob{
					DeviceID:  cursor.Device_ID,
					StartTime: cursor.Full_Epoch.Time,
					EndTime:   cursor.Full_Cursor.Time,
				})
			}
		case FileScanCommandQuick:
			slog = slog.With("type", "quick")

			end := time.Now()
			for _, cursor := range cursors {
				batchJobs = append(batchJobs, BatchJob{
					DeviceID:  cursor.Device_ID,
					StartTime: cursor.Quick_Cursor.Time,
					EndTime:   end,
				})
			}
		}
	default:
		panic("invalid command")
	}

	slog.Info("Started file scan")
	timer := time.Now()
	wg := sync.WaitGroup{}

	sema := make(chan struct{}, workers)

	for _, batchJob := range batchJobs {
		wg.Add(1)

		// Run batch job
		go func(batchJob BatchJob) {
			defer wg.Done()

			sema <- struct{}{}
			defer func() { <-sema }()

			// Get client
			client, err := w.store.GetClient(ctx, types.Key{ID: batchJob.DeviceID})
			if err != nil {
				slog.Error("Failed to get client", "error", err, "device-id", batchJob.DeviceID)
				return
			}
			slog := slog.With("device", client.Conn.Name)

			// Scan
			scanRange := NewFileScanRange(batchJob.StartTime, batchJob.EndTime)
			if err := fileScan(ctx, w.db, client.RPC, batchJob.DeviceID, scanRange); err != nil {
				slog.Error("Failed to scan files", "error", err)

				// Save file cursor when full scan fails
				if job.Command == FileScanCommandFull {
					cursor := types.NewTime(scanRange.Cursor())
					_, err = w.db.ExecContext(ctx, `
						UPDATE dahua_file_cursors SET full_cursor = ? WHERE device_id = ?
					`, cursor, batchJob.DeviceID)
				}

				return
			}

			switch job.Command {
			case FileScanCommandManual:
			case FileScanCommandQuick:
				// Update quick cursor
				quickCursor := types.NewTime(NewFileCursorQuick(batchJob.EndTime, time.Now()))
				_, err = w.db.ExecContext(ctx, `
					UPDATE dahua_file_cursors SET quick_cursor = ? WHERE device_id = ?
				`, quickCursor, batchJob.DeviceID)
				if err != nil {
					slog.Error("Failed to update quick cursor", "error", err)
					return
				}
			case FileScanCommandFull:
				// Update full cursor
				fullCursor := types.NewTime(NewFileCursorFull(batchJob.StartTime, batchJob.FullEpoch))
				_, err = w.db.ExecContext(ctx, `
					UPDATE dahua_file_cursors SET full_cursor = ? WHERE device_id = ?
				`, fullCursor, batchJob.DeviceID)
				if err != nil {
					slog.Error("Failed to update full cursor", "error", err)
					return
				}
			default:
				panic("invalid command")
			}
		}(batchJob)
	}

	wg.Wait()
	slog.Info("Finished file scan", "duration", time.Now().Sub(timer).String())
	return nil
}

func (w FileScanService) Queue(ctx context.Context, job FileScanJob) error {
	select {
	case w.jobC <- job:
		return nil
	default:
		return fmt.Errorf("file scan job in progress")
	}
}

func ResetFileScanCursor(ctx context.Context, db sqlx.ExecerContext, deviceID int64) error {
	cursor := NewDefaultFileCursor()
	updatedAt := types.NewTime(time.Now())

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
		cursor.QuickCursor,
		cursor.FullCursor,
		cursor.FullEpoch,
		updatedAt,
	)
	return err
}
