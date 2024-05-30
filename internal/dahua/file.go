package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
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
			return false
		}
	} else {
		cursor = r.cursor.Add(-r.period)
		if cursor.Before(r.start) {
			return false
		}
	}

	r.cursor = cursor

	return true
}

func NewFileScanner(ctx context.Context, conn dahuarpc.Conn, location *time.Location, scanRange *FileScanRange) *FileScanner {
	return &FileScanner{
		conn:      conn,
		location:  location,
		scanRange: scanRange,
		streams:   [2]*mediafilefind.Stream{},
		closed:    false,
	}
}

type FileScanner struct {
	conn     dahuarpc.Conn
	location *time.Location

	scanRange *FileScanRange
	streams   [2]*mediafilefind.Stream
	closed    bool
}

func (s *FileScanner) Next(ctx context.Context) ([]mediafilefind.FindNextFileInfo, bool, error) {
	if s.closed {
		return nil, false, nil
	}

	for i := range s.streams {
		// Check if stream exists
		if s.streams[i] == nil {
			continue
		}

		files, next, err := s.streams[i].Next(ctx)
		if err != nil {
			s.Close()
			return nil, false, err
		}

		// Check if stream has more files
		if !next {
			s.streams[i] = nil
			continue
		}

		return files, true, nil
	}

	// Check if done scanning
	if !s.scanRange.Next() {
		s.Close()
		return nil, false, nil
	}

	// Get next scan range
	start, end := s.scanRange.Range()
	startTS, endTS := dahuarpc.NewTimestamp(start, s.location), dahuarpc.NewTimestamp(end, s.location)
	condition := mediafilefind.NewCondtion(startTS, endTS)
	if s.scanRange.ascending {
		condition.Order = mediafilefind.ConditionOrderAscent
	} else {
		condition.Order = mediafilefind.ConditionOrderDescent
	}

	// Open picture files stream
	pictureStream, err := mediafilefind.NewStream(ctx, s.conn, condition.Picture())
	if err != nil {
		s.Close()
		return nil, false, err
	}
	s.streams[0] = pictureStream

	// Open video files stream
	videoStream, err := mediafilefind.NewStream(ctx, s.conn, condition.Video())
	if err != nil {
		s.Close()
		return nil, false, err
	}
	s.streams[1] = videoStream

	return s.Next(ctx)
}

func (s *FileScanner) Close() {
	if s.closed {
		return
	}
	s.closed = true

	for _, stream := range s.streams {
		if stream == nil {
			continue
		}
		stream.Close()
	}
}
