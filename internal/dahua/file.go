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

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/ItsNotGoodName/ipcmanview/pkg/jobs"
	"github.com/jlaffaye/ftp"
	"github.com/jmoiron/sqlx"
	"github.com/maragudk/goqite"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

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

func fileScan(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64, start, end time.Time) error {
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
					result, err := db.ExecContext(ctx, `
						UPDATE dahua_files SET updated_at = ? WHERE device_id = ? AND file_path = ?
					`, updatedAt, deviceID, v.FilePath)
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
	`, deviceID, types.NewTime(start), types.NewTime(end), updatedAt)
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
	DeviceID  int64
	StartTime time.Time
	EndTime   time.Time
}

func NewFileScanJobClient(db *sqlx.DB) core.JobClient {
	name := "dahua_file_scan_jobs"
	timeout := 30 * time.Second
	queue := goqite.New(goqite.NewOpts{
		DB:      db.DB,
		Name:    name,
		Timeout: timeout,
	})
	runner := jobs.NewRunner(jobs.NewRunnerOpts{
		Extend:       timeout,
		Limit:        5,
		Log:          slog.Default().With("name", name),
		PollInterval: time.Second,
		Queue:        queue,
	})
	return core.NewJobClient(queue, runner)
}

func RegisterFileScanJob(client core.JobClient, db *sqlx.DB, dahuaStore *Store) core.Job[FileScanJob] {
	return core.NewJob(client, func(ctx context.Context, data FileScanJob) error {
		conn, err := dahuaStore.GetClient(ctx, types.Key{ID: data.DeviceID})
		if err != nil {
			return err
		}

		return fileScan(ctx, db, conn.RPC, data.DeviceID, data.StartTime, data.EndTime)
	})
}

func CreateFileScanJob(ctx context.Context, db *sqlx.DB, fileScanJob core.Job[FileScanJob], args FileScanJob) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	goqiteID, err := fileScanJob.CreateAndGetIDTx(ctx, tx.Tx, args)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO dahua_file_scan_jobs (device_id, goqite_id) VALUES (?, ?)
	`, args.DeviceID, goqiteID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

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
