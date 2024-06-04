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
	"github.com/jlaffaye/ftp"
	"github.com/jmoiron/sqlx"
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

type FileScanJob struct {
	Command   FileScanCommand
	StartTime time.Time
	EndTime   time.Time
	Data      []FileScanData
}

// ENUM(quick,full,manual)
type FileScanCommand string

type FileScanData struct {
	DeviceID  int64
	StartTime time.Time
	EndTime   time.Time
}

func NewFileScanService(db *sqlx.DB, store *Store, quickScanInterval time.Duration) FileScanService {
	return FileScanService{
		db:                db,
		store:             store,
		jobC:              make(chan FileScanJob),
		quickScanInterval: quickScanInterval,
	}
}

type FileScanService struct {
	db                *sqlx.DB
	store             *Store
	jobC              chan FileScanJob
	quickScanInterval time.Duration
}

func (w FileScanService) String() string {
	return "dahua.FileScanService"
}

func (w FileScanService) Serve(ctx context.Context) error {
	slog := slog.With("service", w.String())
	slog.Info("Started service")

	t := time.NewTicker(w.quickScanInterval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			err := w.handle(ctx, slog, FileScanJob{
				Command: FileScanCommandQuick,
			})
			if err != nil {
				return err
			}
		case job := <-w.jobC:
			if err := w.handle(ctx, slog, job); err != nil {
				return err
			}
		}
	}
}

func (w FileScanService) handle(ctx context.Context, slog *slog.Logger, job FileScanJob) error {
	panic("not implemented")
	// type BatchJob struct {
	// 	DeviceID  int64
	// 	StartTime time.Time
	// 	EndTime   time.Time
	// 	FullEpoch time.Time
	// }
	//
	// workers := 3
	// var batchJobs []BatchJob
	// switch job.Command {
	// case FileScanCommandManual:
	// 	slog = slog.With("type", "manual")
	//
	// 	// Add batch jobs from job data
	// 	if len(job.Data) == 0 {
	// 		var deviceIDs []int64
	// 		err := w.db.SelectContext(ctx, &deviceIDs, `
	// 			SELECT id FROM dahua_devices
	// 		`)
	// 		if err != nil {
	// 			return err
	// 		}
	//
	// 		for _, deviceID := range deviceIDs {
	// 			batchJobs = append(batchJobs, BatchJob{
	// 				DeviceID:  deviceID,
	// 				StartTime: job.StartTime,
	// 				EndTime:   job.EndTime,
	// 			})
	// 		}
	// 	} else {
	// 		for _, data := range job.Data {
	// 			batchJobs = append(batchJobs, BatchJob{
	// 				DeviceID:  data.DeviceID,
	// 				StartTime: data.StartTime,
	// 				EndTime:   data.EndTime,
	// 			})
	// 		}
	// 	}
	// case FileScanCommandQuick, FileScanCommandFull:
	// 	// Extract devices ids
	// 	var deviceIDs []int64
	// 	for _, data := range job.Data {
	// 		deviceIDs = append(deviceIDs, data.DeviceID)
	// 	}
	//
	// 	// Get cursors by device ids
	// 	var (
	// 		query     string
	// 		queryArgs []any
	// 		err       error
	// 	)
	// 	if len(deviceIDs) == 0 {
	// 		query = `
	// 			SELECT * FROM dahua_file_cursors
	// 		`
	// 	} else {
	// 		query, queryArgs, err = sqlx.In(`
	// 			SELECT * FROM dahua_file_cursors
	// 			WHERE ? = 0 OR device_id IN (?)
	// 		`, len(deviceIDs), deviceIDs)
	// 	}
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	var cursors []FileScanCursor
	// 	err = sqlx.SelectContext(ctx, w.db, &cursors, query, queryArgs...)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	// Add batch jobs from cursors
	// 	switch job.Command {
	// 	case FileScanCommandFull:
	// 		slog = slog.With("type", "full")
	//
	// 		for _, cursor := range cursors {
	// 			batchJobs = append(batchJobs, BatchJob{
	// 				DeviceID:  cursor.Device_ID,
	// 				StartTime: cursor.Full_Epoch.Time,
	// 				EndTime:   cursor.Full_Cursor.Time,
	// 			})
	// 		}
	// 	case FileScanCommandQuick:
	// 		slog = slog.With("type", "quick")
	//
	// 		end := time.Now()
	// 		for _, cursor := range cursors {
	// 			batchJobs = append(batchJobs, BatchJob{
	// 				DeviceID:  cursor.Device_ID,
	// 				StartTime: cursor.Quick_Cursor.Time,
	// 				EndTime:   end,
	// 			})
	// 		}
	// 	}
	// default:
	// 	panic("invalid command")
	// }
	//
	// slog.Info("Started batch scan")
	// timer := time.Now()
	// wg := sync.WaitGroup{}
	//
	// sema := make(chan struct{}, workers)
	//
	// for _, batchJob := range batchJobs {
	// 	wg.Add(1)
	//
	// 	// Run batch job
	// 	go func(batchJob BatchJob) {
	// 		defer wg.Done()
	//
	// 		sema <- struct{}{}
	// 		defer func() { <-sema }()
	//
	// 		// Get client
	// 		client, err := w.store.GetClient(ctx, types.Key{ID: batchJob.DeviceID})
	// 		if err != nil {
	// 			slog.Error("Failed to get client", "error", err, "device-id", batchJob.DeviceID)
	// 			return
	// 		}
	// 		slog := slog.With("device", client.Conn.Name)
	// 		slog.Info("Starting scan")
	//
	// 		// Scan
	// 		scanRange := NewFileScanRange(batchJob.StartTime, batchJob.EndTime)
	// 		if err := fileScan(ctx, w.db, client.RPC, batchJob.DeviceID, scanRange, func() {
	// 			// Save file cursor when full scan fails
	// 			if job.Command == FileScanCommandFull {
	// 				cursor := types.NewTime(scanRange.Cursor())
	// 				_, err = w.db.ExecContext(ctx, `
	// 					UPDATE dahua_file_cursors SET full_cursor = ? WHERE device_id = ?
	// 				`, cursor, batchJob.DeviceID)
	// 			} else if job.Command == FileScanCommandQuick {
	// 				cursor := types.NewTime(scanRange.Cursor())
	// 				_, err = w.db.ExecContext(ctx, `
	// 					UPDATE dahua_file_cursors SET quick_cursor = ? WHERE device_id = ?
	// 				`, cursor, batchJob.DeviceID)
	// 			}
	// 		}); err != nil {
	// 			slog.Error("Failed to scan files", "error", err)
	// 			return
	// 		}
	//
	// 		switch job.Command {
	// 		case FileScanCommandManual:
	// 		case FileScanCommandQuick:
	// 			// Update quick cursor
	// 			quickCursor := types.NewTime(NewFileCursorQuick(batchJob.EndTime, time.Now()))
	// 			_, err = w.db.ExecContext(ctx, `
	// 				UPDATE dahua_file_cursors SET quick_cursor = ? WHERE device_id = ?
	// 			`, quickCursor, batchJob.DeviceID)
	// 			if err != nil {
	// 				slog.Error("Failed to update quick cursor", "error", err)
	// 				return
	// 			}
	// 		case FileScanCommandFull:
	// 			// Update full cursor
	// 			fullCursor := types.NewTime(NewFileCursorFull(batchJob.StartTime, batchJob.FullEpoch))
	// 			_, err = w.db.ExecContext(ctx, `
	// 				UPDATE dahua_file_cursors SET full_cursor = ? WHERE device_id = ?
	// 			`, fullCursor, batchJob.DeviceID)
	// 			if err != nil {
	// 				slog.Error("Failed to update full cursor", "error", err)
	// 				return
	// 			}
	// 		default:
	// 			panic("invalid command")
	// 		}
	// 	}(batchJob)
	// }
	//
	// wg.Wait()
	// slog.Info("Finished batch scan", "duration", time.Now().Sub(timer).String())
	// return nil
}

func (w FileScanService) Queue(ctx context.Context, job FileScanJob) error {
	select {
	case w.jobC <- job:
		return nil
	default:
		return fmt.Errorf("file scan job in progress")
	}
}
