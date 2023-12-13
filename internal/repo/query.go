package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	sq "github.com/Masterminds/squirrel"
)

type CreateDahuaCameraParams = createDahuaCameraParams

type CreateDahuaFileCursorParams = createDahuaFileCursorParams

func (db DB) CreateDahuaCamera(ctx context.Context, arg CreateDahuaCameraParams, args2 CreateDahuaFileCursorParams) (int64, error) {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback()

	id, err := tx.createDahuaCamera(ctx, arg)
	if err != nil {
		return 0, err
	}

	// TODO: sql.NullInt64 should just be int64...
	err = tx.allocateDahuaSeed(ctx, sql.NullInt64{
		Valid: true,
		Int64: id,
	})
	if err != nil {
		return 0, err
	}

	args2.CameraID = id
	err = tx.createDahuaFileCursor(ctx, args2)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db DB) UpsertDahuaFiles(ctx context.Context, args CreateDahuaFileParams) (int64, error) {
	id, err := db.UpdateDahuaFile(ctx, UpdateDahuaFileParams{
		CameraID:    args.CameraID,
		Channel:     args.Channel,
		StartTime:   args.StartTime,
		EndTime:     args.EndTime,
		Length:      args.Length,
		Type:        args.Type,
		FilePath:    args.FilePath,
		Duration:    args.Duration,
		Disk:        args.Disk,
		VideoStream: args.VideoStream,
		Flags:       args.Flags,
		Events:      args.Events,
		Cluster:     args.Cluster,
		Partition:   args.Partition,
		PicIndex:    args.PicIndex,
		Repeat:      args.Repeat,
		WorkDir:     args.WorkDir,
		WorkDirSn:   args.WorkDirSn,
		UpdatedAt:   args.UpdatedAt,
	})
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	return db.CreateDahuaFile(ctx, args)
}

type ListDahuaEventParams struct {
	pagination.Page
	Code      []string
	Action    []string
	CameraID  []int64
	Start     types.Time
	End       types.Time
	Ascending bool
}

type ListDahuaEventResult struct {
	pagination.PageResult
	Data []DahuaEvent
}

func (db DB) ListDahuaEvent(ctx context.Context, arg ListDahuaEventParams) (ListDahuaEventResult, error) {
	where := sq.And{}

	eq := sq.Eq{}
	if len(arg.Code) != 0 {
		eq["code"] = arg.Code
	}
	if len(arg.Action) != 0 {
		eq["action"] = arg.Action
	}
	if len(arg.CameraID) != 0 {
		eq["camera_id"] = arg.CameraID
	}
	where = append(where, eq)

	and := sq.And{}
	if !arg.Start.IsZero() {
		and = append(and, sq.GtOrEq{"created_at": arg.Start})
	}
	if !arg.End.IsZero() {
		and = append(and, sq.Lt{"created_at": arg.End})
	}
	where = append(where, and)

	order := "created_at DESC"
	if arg.Ascending {
		order = "created_at ASC"
	}

	var res []DahuaEvent
	err := ssq.Query(ctx, db, &res, sq.
		Select("*").
		From("dahua_events").
		Where(where).
		OrderBy(order).
		Limit(uint64(arg.Page.Limit())).
		Offset(uint64(arg.Page.Offset())))
	if err != nil {
		return ListDahuaEventResult{}, err
	}

	var count int
	err = ssq.QueryOne(ctx, db, &count, sq.
		Select("COUNT(*)").
		From("dahua_events").
		Where(where))
	if err != nil {
		return ListDahuaEventResult{}, err
	}

	return ListDahuaEventResult{
		PageResult: arg.Page.Result(count),
		Data:       res,
	}, nil
}

type ListDahuaFileParams struct {
	pagination.Page
	Type      []string
	CameraID  []int64
	Start     types.Time
	End       types.Time
	Ascending bool
}

type ListDahuaFileResult struct {
	pagination.PageResult
	Data []DahuaFile
}

func (db DB) ListDahuaFile(ctx context.Context, arg ListDahuaFileParams) (ListDahuaFileResult, error) {
	where := sq.And{}

	eq := sq.Eq{}
	if len(arg.Type) != 0 {
		eq["type"] = arg.Type
	}
	if len(arg.CameraID) != 0 {
		eq["camera_id"] = arg.CameraID
	}
	where = append(where, eq)

	and := sq.And{}
	if !arg.Start.IsZero() {
		and = append(and, sq.GtOrEq{"start_time": arg.Start})
	}
	if !arg.End.IsZero() {
		and = append(and, sq.Lt{"start_time": arg.End})
	}
	where = append(where, and)

	order := "start_time DESC"
	if arg.Ascending {
		order = "start_time ASC"
	}

	var res []DahuaFile
	err := ssq.Query(ctx, db, &res, sq.
		Select("*").
		From("dahua_files").
		Where(where).
		OrderBy(order).
		Limit(uint64(arg.Page.Limit())).
		Offset(uint64(arg.Page.Offset())))
	if err != nil {
		return ListDahuaFileResult{}, err
	}

	var count int
	err = ssq.QueryOne(ctx, db, &count, sq.
		Select("COUNT(*)").
		From("dahua_files").
		Where(where))
	if err != nil {
		return ListDahuaFileResult{}, err
	}

	return ListDahuaFileResult{
		PageResult: arg.Page.Result(count),
		Data:       res,
	}, nil
}

func (db DB) GetDahuaEventRuleByEvent(ctx context.Context, event models.DahuaEvent) (models.DahuaEventRule, error) {
	res, err := db.getDahuaEventRuleByEvent(ctx, getDahuaEventRuleByEventParams{
		CameraID: event.CameraID,
		Code:     event.Code,
	})
	if err != nil {
		return models.DahuaEventRule{}, err
	}
	if len(res) == 0 {
		return models.DahuaEventRule{}, nil
	}

	return models.DahuaEventRule{
		IgnoreDB:   res[0].IgnoreDb,
		IgnoreLive: res[0].IgnoreLive,
		IgnoreMQTT: res[0].IgnoreMqtt,
	}, nil
}
