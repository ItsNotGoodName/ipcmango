package mediafilefind

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Stream struct {
	conn   dahuarpc.Conn
	object int64
	count  int
	closed bool
}

func NewStream(ctx context.Context, c dahuarpc.Conn, condtion Condition) (*Stream, error) {
	object, err := Create(ctx, c)
	if err != nil {
		return nil, err
	}

	var closed bool
	ok, err := FindFile(ctx, c, object, condtion)
	if err != nil {
		var resErr *dahuarpc.Error
		if !errors.As(err, &resErr) {
			return nil, err
		}

		if resErr.Type != dahuarpc.ErrorTypeNoData {
			return nil, err
		}

		closed = true
	} else {
		closed = !ok
	}

	return &Stream{
		conn:   c,
		object: object,
		count:  64,
		closed: closed,
	}, nil
}

func (s *Stream) Next(ctx context.Context) ([]FindNextFileInfo, error) {
	if s.closed {
		return nil, nil
	}

	files, err := FindNextFile(ctx, s.conn, s.object, s.count)
	if err != nil {
		s.Close()
		return nil, err
	}

	if files.Infos == nil {
		s.Close()
		return nil, nil
	}

	if files.Found < s.count {
		s.Close()
	}

	return files.Infos, nil
}

func (s *Stream) Close() {
	if s.closed {
		return
	}

	s.closed = true

	ctx := context.Background()

	if _, err := Close(ctx, s.conn, s.object); err != nil {
		return
	}

	Destroy(ctx, s.conn, s.object)
}
