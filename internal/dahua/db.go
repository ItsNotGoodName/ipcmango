package dahua

import (
	"context"
	"database/sql"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	sq "github.com/Masterminds/squirrel"
)

// authFilter applies an authorization filter to a select query.
func authFilter(ctx context.Context, sb sq.SelectBuilder, deviceIDField string, level models.DahuaPermissionLevel) sq.SelectBuilder {
	actor := core.UseActor(ctx)

	if actor.Admin {
		return sb
	}

	return sb.
		Where(sq.Expr(deviceIDField+` IN (
			SELECT
				device_id
			FROM
				dahua_permissions
			WHERE
				dahua_permissions.level > ?
				AND (
					dahua_permissions.user_id = ?
					OR dahua_permissions.group_id IN (
						SELECT
							group_id
						FROM
							group_users
						WHERE
							group_users.user_id = ?
					)
				)
			)
		`, level, actor.UserID, actor.UserID))
}

func GetConn(ctx context.Context, id int64) (Conn, error) {
	actor := core.UseActor(ctx)
	v, err := app.DB.C().DahuaGetConn(ctx, repo.DahuaGetConnParams{
		ID:     id,
		Admin:  actor.Admin,
		UserID: core.NewNullInt64(actor.UserID),
	})
	if err != nil {
		return Conn{}, err
	}

	return Conn{
		ID:       v.ID,
		URL:      v.Url.URL,
		Username: v.Username,
		Password: v.Password,
		Location: v.Location.Location,
		Feature:  v.Feature,
		Seed:     int(v.Seed),
	}, nil
}

func ListConn(ctx context.Context) ([]Conn, error) {
	actor := core.UseActor(ctx)
	vv, err := app.DB.C().DahuaListConn(ctx, repo.DahuaListConnParams{
		Admin:  actor.Admin,
		UserID: core.NewNullInt64(actor.UserID),
	})
	if err != nil {
		return nil, err
	}

	conns := make([]Conn, 0, len(vv))
	for _, v := range vv {
		conns = append(conns, Conn{
			ID:       v.ID,
			URL:      v.Url.URL,
			Username: v.Username,
			Password: v.Password,
			Location: v.Location.Location,
			Feature:  v.Feature,
			Seed:     int(v.Seed),
		})
	}

	return conns, nil
}

func ListDeviceIDs(ctx context.Context) ([]int64, error) {
	sb := sq.
		Select("id").
		From("dahua_devices")

	var res []int64
	err := ssq.Query(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_devices.id", levelDefault))
	return res, err
}

// ---------- Count

type dbCountRow struct {
	Count int64
}

func CountEvents(ctx context.Context) (int64, error) {
	sb := sq.
		Select("COUNT(*) AS count").
		From("dahua_events")

	var res dbCountRow
	err := ssq.QueryOne(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_events.device_id", levelDefault))
	return res.Count, err
}

func CountEmails(ctx context.Context) (int64, error) {
	sb := sq.
		Select("COUNT(*) AS count").
		From("dahua_email_messages")

	var res dbCountRow
	err := ssq.QueryOne(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_email_messages.device_id", levelEmail))
	return res.Count, err
}

type ListLatestEmailsResult struct {
	repo.DahuaEmailMessage
	AttachmentCount int64
}

func ListLatestEmails(ctx context.Context, count int) ([]ListLatestEmailsResult, error) {
	sb := sq.
		Select("dahua_email_messages.*", "COUNT(dahua_email_attachments.id) AS attachment_count").
		From("dahua_email_messages").
		LeftJoin("dahua_email_attachments ON dahua_email_attachments.message_id = dahua_email_messages.id").
		OrderBy("created_at DESC").
		GroupBy("dahua_email_messages.id").
		Limit(uint64(count))

	var res []ListLatestEmailsResult
	err := ssq.Query(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_email_messages.device_id", levelEmail))
	return res, err
}

func ListLatestFiles(ctx context.Context, count int) ([]repo.DahuaFile, error) {
	sb := sq.
		Select("*").
		From("dahua_files").
		OrderBy("start_time DESC").
		Limit(uint64(count))

	var res []repo.DahuaFile
	err := ssq.Query(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_files.device_id", levelEmail))
	return res, err
}

func GetDevice(ctx context.Context, id int64) (repo.DahuaDevice, error) {
	sb := sq.
		Select("*").
		From("dahua_devices").
		Where("id = ?", id)

	var res repo.DahuaDevice
	err := ssq.QueryOne(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_devices.id", levelDefault))
	return res, err
}

func ListDevices(ctx context.Context) ([]repo.DahuaDevice, error) {
	sb := sq.
		Select("*").
		From("dahua_devices")

	var res []repo.DahuaDevice
	err := ssq.Query(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_devices.id", levelDefault))
	return res, err
}

type EmailFilter struct {
	FilterDeviceIDs   []int64
	FilterAlarmEvents []string
}

func (arg EmailFilter) where() sq.Eq {
	where := sq.Eq{}
	if len(arg.FilterDeviceIDs) != 0 {
		where["dahua_email_messages.device_id"] = arg.FilterDeviceIDs
	}
	if len(arg.FilterAlarmEvents) != 0 {
		where["dahua_email_messages.alarm_event"] = arg.FilterAlarmEvents
	}
	return where
}

type ListEmailsParams struct {
	pagination.Page
	EmailFilter
	Ascending bool
}

type ListEmailsResult struct {
	pagination.PageResult
	Items []ListEmailsResultItems
}

type ListEmailsResultItems struct {
	repo.DahuaEmailMessage
	DeviceName      string
	AttachmentCount int
}

func ListEmails(ctx context.Context, arg ListEmailsParams) (ListEmailsResult, error) {
	where := arg.where()

	order := "dahua_email_messages.id"
	if arg.Ascending {
		order += " ASC"
	} else {
		order += " DESC"
	}
	sb := sq.
		Select(
			"dahua_email_messages.*",
			"COUNT(dahua_email_attachments.id) AS attachment_count",
			"dahua_devices.name AS device_name",
		).
		From("dahua_email_messages").
		LeftJoin("dahua_email_attachments ON dahua_email_attachments.message_id = dahua_email_messages.id").
		LeftJoin("dahua_devices ON dahua_devices.id = dahua_email_messages.device_id").
		Where(where).
		OrderBy(order).
		GroupBy("dahua_email_messages.id").
		Offset(uint64(arg.Offset())).
		Limit(uint64(arg.Limit()))

	var items []ListEmailsResultItems
	if err := ssq.Query(ctx, app.DB, &items, authFilter(ctx, sb, "dahua_email_messages.device_id", levelEmail)); err != nil {
		return ListEmailsResult{}, err
	}

	sb = sq.
		Select("COUNT(*) AS count").
		From("dahua_email_messages").
		Where(where)

	var count dbCountRow
	if err := ssq.QueryOne(ctx, app.DB, &count, authFilter(ctx, sb, "dahua_email_messages.device_id", levelEmail)); err != nil {
		return ListEmailsResult{}, err
	}

	return ListEmailsResult{
		PageResult: arg.Result(int(count.Count)),
		Items:      items,
	}, nil
}

type GetEmailParams struct {
	EmailFilter
	ID int64
}

type GetEmailResult struct {
	NextEmailID int64
	Message     repo.DahuaEmailMessage
	Attachments []repo.DahuaListEmailAttachmentsForMessageRow
	Filter      EmailFilter
}

type GetEmailResultAttachments struct {
	repo.DahuaEmailAttachment
	repo.DahuaAferoFile
}

func GetEmail(ctx context.Context, arg GetEmailParams) (GetEmailResult, error) {
	sb := sq.
		Select("*").
		From("dahua_email_messages").
		Where("id <= ?", arg.ID).
		Where(arg.where()).
		OrderBy("id DESC").
		Limit(2)
	var messages []repo.DahuaEmailMessage
	if err := ssq.Query(ctx, app.DB, &messages, authFilter(ctx, sb, "dahua_email_messages.device_id", levelEmail)); err != nil {
		return GetEmailResult{}, err
	}
	if len(messages) == 0 || messages[0].ID != arg.ID {
		return GetEmailResult{}, core.ErrNotFound
	}
	message := messages[0]
	nextEmailID := messages[0].ID
	if len(messages) == 2 {
		nextEmailID = messages[1].ID
	}

	attachments, err := app.DB.C().DahuaListEmailAttachmentsForMessage(ctx, arg.ID)
	if err != nil {
		return GetEmailResult{}, err
	}

	return GetEmailResult{
		NextEmailID: nextEmailID,
		Message:     message,
		Attachments: attachments,
	}, nil
}

type GetEmailAroundParams struct {
	ID int64
	EmailFilter
}

type GetEmailAroundResult struct {
	EmailSeen       int64
	PreviousEmailID int64
	Count           int64
}

func GetEmailAround(ctx context.Context, arg GetEmailAroundParams) (GetEmailAroundResult, error) {
	where := arg.where()

	sb := sq.
		Select(
			"MIN(id) AS previous_email_id",
			"COUNT(*) as email_seen",
		).
		From("dahua_email_messages").
		Where("id > ?", arg.ID).
		Where(where)
	var res struct {
		PreviousEmailID sql.NullInt64
		EmailSeen       int64
	}
	if err := ssq.QueryOne(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_email_messages.device_id", levelEmail)); err != nil {
		return GetEmailAroundResult{}, err
	}

	emailSeen := res.EmailSeen
	previousEmailID := arg.ID
	if res.PreviousEmailID.Valid {
		previousEmailID = res.PreviousEmailID.Int64
	}

	sb = sq.
		Select("COUNT(*) AS count").
		From("dahua_email_messages").
		Where(where)

	var count dbCountRow
	if err := ssq.QueryOne(ctx, app.DB, &count, authFilter(ctx, sb, "dahua_email_messages.device_id", levelEmail)); err != nil {
		return GetEmailAroundResult{}, err
	}

	return GetEmailAroundResult{
		EmailSeen:       emailSeen + 1,
		PreviousEmailID: previousEmailID,
		Count:           count.Count,
	}, nil
}

type EventFilter struct {
	FilterDeviceIDs []int64
	FilterCodes     []string
	FilterActions   []string
}

func (arg EventFilter) where() sq.Eq {
	where := sq.Eq{}
	if len(arg.FilterDeviceIDs) != 0 {
		where["dahua_events.device_id"] = arg.FilterDeviceIDs
	}
	if len(arg.FilterCodes) != 0 {
		where["dahua_events.code"] = arg.FilterCodes
	}
	if len(arg.FilterActions) != 0 {
		where["dahua_events.action"] = arg.FilterActions
	}
	return where
}

type ListEventsParams struct {
	pagination.Page
	Ascending bool
	EventFilter
}

type ListEventsResult struct {
	pagination.PageResult
	Items []ListEventsResultItems
}

type ListEventsResultItems struct {
	repo.DahuaEvent
	DeviceName string
}

func ListEvents(ctx context.Context, arg ListEventsParams) (ListEventsResult, error) {
	where := arg.where()

	order := "dahua_events.id"
	if arg.Ascending {
		order += " ASC"
	} else {
		order += " DESC"
	}
	sb := sq.
		Select(
			"dahua_events.*",
			"dahua_devices.name AS device_name",
		).
		From("dahua_events").
		LeftJoin("dahua_devices ON dahua_devices.id = dahua_events.device_id").
		Where(where).
		OrderBy(order).
		Offset(uint64(arg.Offset())).
		Limit(uint64(arg.Limit()))

	var items []ListEventsResultItems
	if err := ssq.Query(ctx, app.DB, &items, authFilter(ctx, sb, "dahua_events.device_id", levelDefault)); err != nil {
		return ListEventsResult{}, err
	}

	sb = sq.
		Select("COUNT(*) AS count").
		From("dahua_events").
		Where(where)

	var res dbCountRow
	if err := ssq.QueryOne(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_events.device_id", levelDefault)); err != nil {
		return ListEventsResult{}, err
	}

	return ListEventsResult{
		PageResult: arg.Result(int(res.Count)),
		Items:      items,
	}, nil
}

func ListEmailAlarmEvents(ctx context.Context) ([]string, error) {
	sb := sq.Select("DISTINCT alarm_event").From("dahua_email_messages")

	var res []string
	if err := ssq.Query(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_email_messages.device_id", levelEmail)); err != nil {
		return nil, err
	}

	return res, nil
}

func ListEventCodes(ctx context.Context) ([]string, error) {
	sb := sq.Select("DISTINCT code").From("dahua_events")

	var res []string
	if err := ssq.Query(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_events.device_id", levelEmail)); err != nil {
		return nil, err
	}

	return res, nil
}

func ListEventActions(ctx context.Context) ([]string, error) {
	sb := sq.Select("DISTINCT action").From("dahua_events")

	var res []string
	if err := ssq.Query(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_events.device_id", levelEmail)); err != nil {
		return nil, err
	}

	return res, nil
}

func ListStreams(ctx context.Context, deviceID int64) ([]repo.DahuaStream, error) {
	sb := sq.
		Select("*").
		From("dahua_streams").
		Where("device_id = ?", deviceID)

	var res []repo.DahuaStream
	sb = authFilter(ctx, sb, "dahua_streams.device_id", levelDefault)
	if err := ssq.Query(ctx, app.DB, &res, sb); err != nil {
		return nil, err
	}

	return res, nil
}

func ListEventRules(ctx context.Context) ([]repo.DahuaEventRule, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return nil, err
	}
	return app.DB.C().DahuaListEventRules(ctx)
}

type FileFilter struct {
	FilterDeviceIDs []int64
	FilterMonth     time.Time
}

func (arg FileFilter) where() sq.And {
	and := sq.And{}

	if !arg.FilterMonth.IsZero() {
		month := types.NewTime(arg.FilterMonth)
		and = append(and, sq.Expr(`(start_time >= datetime(?, 'start of month') AND start_time < datetime( ?, 'start of month', '+1 month'))`, month, month))
	}

	eq := sq.Eq{}

	if len(arg.FilterDeviceIDs) != 0 {
		eq["dahua_files.device_id"] = arg.FilterDeviceIDs
	}

	return append(and, eq)
}

func CountFiles(ctx context.Context) (int64, error) {
	sb := sq.
		Select("COUNT(*) AS count").
		From("dahua_files")

	var res dbCountRow
	err := ssq.QueryOne(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_files.device_id", levelDefault))
	return res.Count, err
}

type CountFilesByMonthResult struct {
	Month types.Time
	Count int64
}

func CountFilesByMonth(ctx context.Context, filter FileFilter) ([]CountFilesByMonthResult, error) {
	sb := sq.
		Select(
			"datetime(start_time, 'start of month') AS month",
			"count(id) as count",
		).
		From("dahua_files").
		Where(filter.where()).
		GroupBy("strftime('%Y-%m', substr(start_time, 1, 10))")

	var res []CountFilesByMonthResult
	err := ssq.Query(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_files.device_id", levelDefault))
	if err != nil {
		return nil, err
	}

	return res, nil
}

type ListFilesParams struct {
	pagination.Page
	Ascending bool
	Filter    FileFilter
}

type ListFilesResult struct {
	pagination.PageResult
	Items []ListFilesResultItems
}

type ListFilesResultItems struct {
	repo.DahuaFile
	DeviceName string
}

func ListFiles(ctx context.Context, arg ListFilesParams) (ListFilesResult, error) {
	where := arg.Filter.where()

	order := "dahua_files.start_time"
	if arg.Ascending {
		order += " ASC"
	} else {
		order += " DESC"
	}
	sb := sq.
		Select(
			"dahua_files.*",
			"dahua_devices.name AS device_name",
		).
		From("dahua_files").
		LeftJoin("dahua_devices ON dahua_devices.id = dahua_files.device_id").
		Where(where).
		OrderBy(order).
		Offset(uint64(arg.Offset())).
		Limit(uint64(arg.Limit()))

	var items []ListFilesResultItems
	if err := ssq.Query(ctx, app.DB, &items, authFilter(ctx, sb, "dahua_files.device_id", levelDefault)); err != nil {
		return ListFilesResult{}, err
	}

	sb = sq.
		Select("COUNT(*) AS count").
		From("dahua_files").
		Where(where)

	var res dbCountRow
	if err := ssq.QueryOne(ctx, app.DB, &res, authFilter(ctx, sb, "dahua_files.device_id", levelDefault)); err != nil {
		return ListFilesResult{}, err
	}

	return ListFilesResult{
		PageResult: arg.Result(int(res.Count)),
		Items:      items,
	}, nil
}
