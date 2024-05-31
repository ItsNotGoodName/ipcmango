package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/jmoiron/sqlx"
	"github.com/k0kubun/pp/v3"
)

type StorageDestination struct {
	types.Key
	types.Timestamp
	Name             string
	Storage          string
	Server_Address   string
	Port             int
	Username         string
	Password         string
	Remote_Directory string
}

type CreateStorageDestinationArgs struct {
	UUID            string
	Name            string
	Storage         string
	ServerAddress   string
	Port            int
	Username        string
	Password        string
	RemoteDirectory string
}

func CreateStorageDestination(ctx context.Context, db *sqlx.DB, args CreateStorageDestinationArgs) (types.Key, error) {
	return createStorageDestination(ctx, db, args)
}

func PutStorageDestinations(ctx context.Context, db *sqlx.DB, args []CreateStorageDestinationArgs) ([]types.Key, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = db.ExecContext(ctx, `
		DELETE FROM dahua_storage_destinations
	`)
	if err != nil {
		return nil, err
	}

	var keys []types.Key
	for _, arg := range args {
		key, err := createStorageDestination(ctx, tx, arg)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return keys, nil
}

func createStorageDestination(ctx context.Context, db sqlx.QueryerContext, args CreateStorageDestinationArgs) (types.Key, error) {
	createdAt := types.NewTime(time.Now())
	updatedAt := types.NewTime(time.Now())

	var key types.Key
	err := sqlx.GetContext(ctx, db, &key, `
		INSERT INTO dahua_storage_destinations (
			uuid,
			name,
			storage,
			server_address,
			port,
			username,
			password,
			remote_directory,
			created_at,
			updated_at
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
		RETURNING id, uuid;
	`,
		args.UUID,
		args.Name,
		args.Storage,
		args.ServerAddress,
		args.Port,
		args.Username,
		args.Password,
		args.RemoteDirectory,
		createdAt,
		updatedAt,
	)
	if err != nil {
		pp.Println(err)
		return key, err
	}

	return key, nil
}

type UpdateStorageDestinationArgs struct {
	UUID            string
	Name            string
	Storage         string
	ServerAddress   string
	Port            int
	Username        string
	Password        string
	RemoteDirectory string
}

func UpdateStorageDestination(ctx context.Context, db *sqlx.DB, args UpdateStorageDestinationArgs) (types.Key, error) {
	updatedAt := types.NewTime(time.Now())

	var key types.Key
	err := db.GetContext(ctx, &key, `
		UPDATE dahua_storage_destinations SET
			name = ?,
			storage = ?,
			server_address = ?,
			port = ?,
			username = ?,
			password = ?,
			remote_directory = ?,
			updated_at = ?
		WHERE uuid = ?
		RETURNING id, uuid;
	`,
		args.Name,
		args.Storage,
		args.ServerAddress,
		args.Port,
		args.Username,
		args.Password,
		args.RemoteDirectory,
		updatedAt,
		args.UUID,
	)
	if err != nil {
		return key, err
	}

	return key, nil
}

func DeleteStorageDestination(ctx context.Context, db *sqlx.DB, uuid string) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM StorageDestination WHERE uuid = ?
	`, uuid)
	return err
}
