package system

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/gorise"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Endpoint struct {
	core.Key
	Gorise_URL string
	Created_At types.Time
	Updated_At types.Time
}

type CreateEndpointArgs struct {
	GoriseURL string
}

func CreateEndpoint(ctx context.Context, db *sqlx.DB, args CreateEndpointArgs) (Endpoint, error) {
	_, err := gorise.Build(args.GoriseURL)
	if err != nil {
		return Endpoint{}, err
	}

	endpointUUID := uuid.NewString()
	createdAt := types.NewTime(time.Now())
	updatedAt := types.NewTime(time.Now())

	var endpoint Endpoint
	err = db.GetContext(ctx, &endpoint, `
		INSERT INTO endpoints (
			uuid,
			gorise_url,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?)
		RETURNING *
	`, endpointUUID, args.GoriseURL, createdAt, updatedAt)
	if err != nil {
		return endpoint, err
	}

	return endpoint, nil
}

func DeleteEndpoints(ctx context.Context, db *sqlx.DB, uuid string) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM endpoints WHERE uuid = ?
	`, uuid)
	return err
}
