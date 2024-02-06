package migrations

import (
	"context"
	_ "embed"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

func Normalize(ctx context.Context, db sqlite.DB) error {
	_, err := db.ExecContext(ctx, `
WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
INSERT OR IGNORE INTO dahua_seeds (seed) SELECT value from generate_series;
INSERT OR IGNORE INTO dahua_event_rules (code) VALUES ('');
	`)
	if err != nil {
		return err
	}

	{
		c := dahua.NewFileCursor()
		err := db.C().DahuaNormalizeFileCursors(context.Background(), repo.DahuaNormalizeFileCursorsParams{
			QuickCursor: c.QuickCursor,
			FullCursor:  c.FullCursor,
			FullEpoch:   c.FullEpoch,
			Scan:        c.Scan,
			ScanPercent: c.ScanPercent,
			ScanType:    c.ScanType,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
