package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
)

type ScanType string

var (
	ScanTypeFull  ScanType = "full"
	ScanTypeQuick ScanType = "quick"
)

const scanVolatileDuration = 8 * time.Hour

func NewFileCursor() repo.CreateDahuaFileCursorParams {
	now := time.Now()
	return repo.CreateDahuaFileCursorParams{
		DeviceID:    0,
		QuickCursor: types.NewTime(now.Add(-scanVolatileDuration)),
		FullCursor:  types.NewTime(now),
		FullEpoch:   types.NewTime(ScannerEpoch),
		Percent:     0,
	}
}

func updateFileCursor(fileCursor repo.DahuaFileCursor, scanPeriod ScannerPeriod, scanType ScanType) repo.DahuaFileCursor {
	switch scanType {
	case ScanTypeFull:
		// Update FullCursor
		if scanPeriod.Start.Before(fileCursor.FullCursor.Time) {
			fileCursor.FullCursor = types.NewTime(scanPeriod.Start)
		}
	case ScanTypeQuick:
		// Update QuickCursor
		quickCursor := time.Now().Add(-scanVolatileDuration)
		if scanPeriod.End.Before(quickCursor) {
			fileCursor.QuickCursor = types.NewTime(scanPeriod.End)
		} else {
			fileCursor.QuickCursor = types.NewTime(quickCursor)
		}
	default:
		panic("unknown type")
	}

	return fileCursor
}

func getScanRange(fileCursor repo.DahuaFileCursor, scanType ScanType) models.TimeRange {
	switch scanType {
	case ScanTypeFull:
		return models.TimeRange{
			Start: fileCursor.FullEpoch.Time,
			End:   fileCursor.FullCursor.Time,
		}
	case ScanTypeQuick:
		return models.TimeRange{
			Start: fileCursor.QuickCursor.Time,
			End:   time.Now(),
		}
	default:
		panic("unknown type")
	}
}

// ScanReset cannot be called concurrently for the same device.
func ScanReset(ctx context.Context, db repo.DB, id int64) error {
	fileCursor := NewFileCursor()
	_, err := db.UpdateDahuaFileCursor(ctx, repo.UpdateDahuaFileCursorParams{
		QuickCursor: fileCursor.QuickCursor,
		FullCursor:  fileCursor.FullCursor,
		FullEpoch:   fileCursor.FullEpoch,
		DeviceID:    id,
		Percent:     0,
	})
	return err
}

// Scan cannot be called concurrently for the same device.
func Scan(ctx context.Context, db repo.DB, rpcClient dahuarpc.Conn, device models.DahuaConn, scanType ScanType) error {
	fileCursor, err := db.UpdateDahuaFileCursorPercent(ctx, repo.UpdateDahuaFileCursorPercentParams{
		DeviceID: device.ID,
		Percent:  0,
	})
	if err != nil {
		return err
	}

	updated_at := types.NewTime(time.Now())
	iterator := NewScannerPeriodIterator(getScanRange(fileCursor, scanType))
	mediaFilesC := make(chan []mediafilefind.FindNextFileInfo)

	for scannerPeriod, ok := iterator.Next(); ok; scannerPeriod, ok = iterator.Next() {
		cancel, errC := Scanner(ctx, rpcClient, scannerPeriod, device.Location, mediaFilesC)
		defer cancel()

	inner:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errC:
				if err != nil {
					return err
				}
				break inner
			case mediaFiles := <-mediaFilesC:
				files, err := NewDahuaFiles(device.ID, mediaFiles, int(device.Seed), device.Location)
				if err != nil {
					return err
				}

				for _, f := range files {
					_, err := db.UpsertDahuaFiles(ctx, repo.CreateDahuaFileParams{
						DeviceID:    device.ID,
						Channel:     int64(f.Channel),
						StartTime:   types.NewTime(f.StartTime),
						EndTime:     types.NewTime(f.EndTime),
						Length:      int64(f.Length),
						Type:        f.Type,
						FilePath:    f.FilePath,
						Duration:    int64(f.Duration),
						Disk:        int64(f.Disk),
						VideoStream: f.VideoStream,
						Flags:       types.StringSlice{Slice: f.Flags},
						Events:      types.StringSlice{Slice: f.Events},
						Cluster:     int64(f.Cluster),
						Partition:   int64(f.Partition),
						PicIndex:    int64(f.PicIndex),
						Repeat:      int64(f.Repeat),
						WorkDir:     f.WorkDir,
						WorkDirSn:   f.WorkDirSN,
						UpdatedAt:   updated_at,
						Storage:     core.StorageFromFilePath(f.FilePath),
					})
					if err != nil {
						return err
					}
				}
			}
		}

		err := db.DeleteDahuaFile(ctx, repo.DeleteDahuaFileParams{
			UpdatedAt: updated_at,
			DeviceID:  device.ID,
			Start:     types.NewTime(scannerPeriod.Start.UTC()),
			End:       types.NewTime(scannerPeriod.End.UTC()),
		})
		if err != nil {
			return err
		}

		fileCursor = updateFileCursor(fileCursor, scannerPeriod, scanType)

		fileCursor, err = db.UpdateDahuaFileCursor(ctx, repo.UpdateDahuaFileCursorParams{
			QuickCursor: fileCursor.QuickCursor,
			FullCursor:  fileCursor.FullCursor,
			FullEpoch:   fileCursor.FullEpoch,
			DeviceID:    device.ID,
			Percent:     iterator.Percent(),
		})
		if err != nil {
			return err
		}
	}
	fileCursor, err = db.UpdateDahuaFileCursor(ctx, repo.UpdateDahuaFileCursorParams{
		QuickCursor: fileCursor.QuickCursor,
		FullCursor:  fileCursor.FullCursor,
		FullEpoch:   fileCursor.FullEpoch,
		DeviceID:    device.ID,
		Percent:     iterator.Percent(),
	})
	if err != nil {
		return err
	}

	return nil
}
