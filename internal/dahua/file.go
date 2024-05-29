package dahua

import (
	"time"
)

const MaxFileScanPeriod = 30 * 24 * time.Hour

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

// func NewScaner(ctx context.Context, conn dahuarpc.Conn, start, end time.Time) (*Scanner, error) {
// 	return &Scanner{
// 		maxScannerPeriod: 30 * 24 * time.Hour,
// 		conn:             conn,
// 		start:            start,
// 		end:              end,
// 		closed:           false,
// 		cursor:           end,
// 		streams:          []mediafilefind.Stream{},
// 	}, nil
// }
//
// type Scanner struct {
// 	maxScannerPeriod time.Duration
// 	conn             dahuarpc.Conn
//
// 	closed  bool
// 	cursor  time.Time
// 	streams []mediafilefind.Stream
// }
//
// func (s *Scanner) Cursor() time.Time {
// 	return s.cursor
// }
//
// func (s *Scanner) Next(ctx context.Context) ([]mediafilefind.FindNextFileInfo, error) {
// 	if s.closed {
// 		return nil, nil
// 	}
//
// 	// fetch
// 	for _, stream := range s.streams {
// 		files, err := stream.Next(ctx)
// 		if files == nil && err == nil {
// 			continue
// 		}
// 	}
//
// 	// next
// 	if s.start.Equal(s.cursor) {
// 		return ScannerPeriod{}, false
// 	}
//
// 	cursor := s.cursor.Add(-s.maxScannerPeriod)
// 	start := cursor.Add(-s.maxScannerPeriod)
// 	if start.Before(s.start) {
// 		start = s.start
// 	}
//
// 	s.cursor = cursor
//
// 	stream1, err := mediafilefind.NewStream(ctx, s.conn, mediafilefind.Condition{})
// 	stream2, err := mediafilefind.NewStream(ctx, s.conn, mediafilefind.Condition{})
//
// 	// Only mutation in this struct
// 	s.cursor = cursor
//
// 	s.closed = true
// 	return nil, nil
// }
//
// func (s *Scanner) Close() {
// 	if s.closed {
// 		return
// 	}
//
// 	for _, stream := range s.streams {
// 		stream.Close()
// 	}
// 	s.closed = true
// }
