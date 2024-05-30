package dahua

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/gorise"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
)

// -------------------- EmailEndpoint CRUD
type EmailEndpoint struct {
	core.Key
	Global         bool
	Expression     string
	URLs           types.Slice[string]
	Title_Template string
	Body_Template  string
	Attachments    bool
	Created_At     types.Time
	Updated_At     types.Time
	Disabled_At    sql.Null[types.Time]
}

func GetEmailEndpointDeviceUUIDs(ctx context.Context, db *sqlx.DB, endpointKey core.Key) ([]string, error) {
	var deviceUUIDs []string
	err := db.SelectContext(ctx, &deviceUUIDs, `
		SELECT d.uuid FROM dahua_devices_to_email_endpoints AS t
		LEFT JOIN dahua_devices AS d ON t.device_id = d.id
		WHERE t.email_endpoint_id = ?;
	`, endpointKey.ID)
	return deviceUUIDs, err
}

type CreateEmailEndpointArgs struct {
	UUID          string
	Global        bool
	Expression    string
	TitleTemplate string
	BodyTemplate  string
	Attachments   bool
	URLs          types.Slice[string]
	DeviceUUIDs   []string
	Disabled      bool
}

func CreateEmailEndpoint(ctx context.Context, db *sqlx.DB, args CreateEmailEndpointArgs) (core.Key, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return core.Key{}, err
	}
	defer tx.Rollback()

	key, err := createEmailEndpoint(ctx, tx, args)
	if err != nil {
		return core.Key{}, err
	}

	return key, tx.Commit()
}

func PutEmailEndpoints(ctx context.Context, db *sqlx.DB, args []CreateEmailEndpointArgs) ([]core.Key, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		DELETE FROM dahua_email_endpoints
	`)
	if err != nil {
		return nil, err
	}

	var keys []core.Key
	for _, arg := range args {
		key, err := createEmailEndpoint(ctx, tx, arg)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, tx.Commit()
}

func createEmailEndpoint(ctx context.Context, tx *sqlx.Tx, args CreateEmailEndpointArgs) (core.Key, error) {
	for _, urL := range args.URLs.V {
		_, err := gorise.Build(urL)
		if err != nil {
			return core.Key{}, err
		}
	}

	createdAt := types.NewTime(time.Now())
	updatedAt := types.NewTime(time.Now())
	var disabledAt *types.Time
	if args.Disabled {
		t := types.NewTime(time.Now())
		disabledAt = &t
	}

	var key core.Key
	err := tx.GetContext(ctx, &key, `
		INSERT INTO dahua_email_endpoints (
			uuid,
			global,
			expression,
			title_template,
			body_template,
			attachments,
			urls,
			created_at,
			updated_at,
			disabled_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, uuid
	`,
		args.UUID,
		args.Global,
		args.Expression,
		args.TitleTemplate,
		args.BodyTemplate,
		args.Attachments,
		args.URLs,
		createdAt,
		updatedAt,
		disabledAt,
	)
	if err != nil {
		return core.Key{}, err
	}

	if len(args.DeviceUUIDs) != 0 {
		query, queryArgs, err := sqlx.In(`
			INSERT INTO dahua_devices_to_email_endpoints (device_id, email_endpoint_id)
			SELECT id, ? FROM dahua_devices WHERE uuid IN (?)
		`, key.ID, args.DeviceUUIDs)
		_, err = tx.ExecContext(ctx, tx.Rebind(query), queryArgs...)
		if err != nil {
			return core.Key{}, err
		}
	}

	return key, nil
}

func DeleteEndpoint(ctx context.Context, db *sqlx.DB, uuid string) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM dahua_email_endpoints WHERE uuid = ?
	`, uuid)
	return err
}

// -------------------- Email CRUD
type EmailMessage struct {
	core.Key
	Device_ID           int64
	Date                types.Time
	From                string
	To                  types.Slice[string]
	Subject             string
	Text                string
	Alarm_Event         string
	Alarm_Input_Channel string
	Alarm_Name          string
	Created_At          types.Time
}

type EmailAttachment struct {
	core.Key
	Message_ID sql.Null[int64]
	File_Name  string
	Size       int64
	Mime_Type  string
}

type CreateEmailArgs struct {
	DeviceKey         core.Key
	Date              types.Time
	From              string
	To                types.Slice[string]
	Subject           string
	Text              string
	AlarmEvent        string
	AlarmInputChannel int
	AlarmName         string
	Attachments       []CreateEmailParamsAttachment
}

type CreateEmailParamsAttachment struct {
	FileName string
	Content  []byte
}

func CreateEmail(ctx context.Context, db *sqlx.DB, afs afero.Fs, args CreateEmailArgs) (string, error) {
	messageUUID := uuid.NewString()
	createdAt := types.NewTime(time.Now())

	result, err := db.ExecContext(ctx, `
		INSERT INTO dahua_email_messages (
			uuid,
			device_id,
			date,
			'from',
			'to',
			subject,
			'text',
			alarm_event,
			alarm_input_channel,
			alarm_name,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		messageUUID,
		args.DeviceKey.ID,
		args.Date,
		args.From,
		args.To,
		args.Subject,
		args.Text,
		args.AlarmEvent,
		args.AlarmInputChannel,
		args.AlarmName,
		createdAt,
	)
	if err != nil {
		return "", err
	}
	messageID, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	for _, v := range args.Attachments {
		if err := createEmailAttachment(ctx, afs, db, messageID, v.FileName, v.Content); err != nil {
			return "", err
		}
	}

	bus.Publish(bus.EmailCreated{
		DeviceKey: args.DeviceKey,
		MessageKey: core.Key{
			UUID: messageUUID,
			ID:   messageID,
		},
	})

	return messageUUID, nil
}

func createEmailAttachment(ctx context.Context, afs afero.Fs, tx sqlx.ExecerContext, messageID int64, fileName string, content []byte) error {
	attachmentUUID := uuid.NewString()

	mimeType := http.DetectContentType(content)

	_, err := tx.ExecContext(ctx, `
		INSERT INTO dahua_email_attachments (
			uuid,
			message_id,
			file_name,
			size,
			mime_type
		) VALUES (?, ?, ?, ?, ?)
	`, attachmentUUID, messageID, fileName, len(content), mimeType)
	if err != nil {
		return err
	}

	file, err := afs.Create(attachmentUUID)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return err
	}

	return nil
}

func DeleteOrphanEmailAttachments(ctx context.Context, afs afero.Fs, db *sqlx.DB) error {
	var attachments []EmailAttachment
	err := db.Select(&attachments, `
		SELECT * FROM dahua_email_attachments WHERE message_id IS NULL LIMIT 10
	`)
	if err != nil {
		return err
	}
	if len(attachments) == 0 {
		return nil
	}

	for _, v := range attachments {
		if err := afs.Remove(v.UUID); err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, `
			DELETE FROM dahua_email_attachments WHERE id = ?
		`, v.ID)
		if err != nil {
			return nil
		}
	}

	return DeleteOrphanEmailAttachments(ctx, afs, db)
}

// ---------- END

type EmailContent struct {
	AlarmEvent        string
	AlarmInputChannel int
	AlarmDeviceName   string
	AlarmName         string
	IPAddress         string
}

func ParseEmailContent(text string) EmailContent {
	var content EmailContent
	for _, line := range strings.Split(text, "\n") {
		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "Alarm Event":
			content.AlarmEvent = value
		case "Alarm Input Channel":
			channel, _ := strconv.Atoi(value)
			content.AlarmInputChannel = channel
		case "Alarm Device Name":
			content.AlarmDeviceName = value
		case "Alarm Name":
			content.AlarmName = value
		case "IP Address":
			content.IPAddress = value
		default:
		}
	}

	return content
}

func NewDeleteOrphanEmailAttachmentsJob(db *sqlx.DB, afs afero.Fs) DeleteOrphanEmailAttachmentsJob {
	return DeleteOrphanEmailAttachmentsJob{
		db:  db,
		afs: afs,
	}
}

type DeleteOrphanEmailAttachmentsJob struct {
	db  *sqlx.DB
	afs afero.Fs
}

func (w DeleteOrphanEmailAttachmentsJob) Description() string {
	return "dahua.DeleteOrphanEmailAttachmentsJob"
}

func (w DeleteOrphanEmailAttachmentsJob) Execute(ctx context.Context) error {
	return DeleteOrphanEmailAttachments(ctx, w.afs, w.db)
}

func HandleEmail(ctx context.Context, db *sqlx.DB, afs afero.Fs, messageID int64) error {
	// Load message
	var message EmailMessage
	err := sqlx.GetContext(ctx, db, &message, `
		SELECT * FROM dahua_email_messages WHERE id = ?
	`, messageID)
	if err != nil {
		return err
	}

	// Load endpoints that match the message
	var endpoints []EmailEndpoint
	err = sqlx.Select(db, &endpoints, `
		SELECT t.* FROM dahua_email_endpoints AS t
		WHERE (t.global IS TRUE OR t.id IN (SELECT id FROM dahua_devices_to_email_endpoints AS r WHERE r.device_id = ?))
		AND t.disabled_at IS NULL
	`, message.Device_ID)
	if err != nil {
		return err
	}

	// Load attachments
	var attachments []EmailAttachment
	err = sqlx.SelectContext(ctx, db, &attachments, `
		SELECT * FROM dahua_email_attachments WHERE message_id = ?
	`, messageID)
	if err != nil {
		return err
	}

	email := NewEmailPayload(message, attachments)

	wg := sync.WaitGroup{}

	var errs error

	for _, endpoint := range endpoints {
		// Check if email is allowed to endpoint
		rule, err := NewEmailRule(endpoint.Expression)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if ok, err := rule.Match(email); !ok {
			if err != nil {
				errs = errors.Join(err)
			}
			continue
		}

		sender, err := NewSender(email, endpoint.Title_Template, endpoint.Body_Template, endpoint.Attachments)
		if err != nil {
			errs = errors.Join(err)
			continue
		}

		for _, endpointURL := range endpoint.URLs.V {
			wg.Add(1)
			go func(endpoint EmailEndpoint, sender EmailSender, endpointURL string) {
				defer wg.Done()
				if err := sender.Send(ctx, afs, endpointURL); err != nil {
					slog.Error("Failed to send to endpoint", "error", err)
				}
			}(endpoint, sender, endpointURL)
		}
	}

	wg.Wait()

	return errs
}

func NewEmailWorker(db *sqlx.DB, afs afero.Fs) EmailWorker {
	return EmailWorker{
		db:         db,
		afs:        afs,
		messageIDC: make(chan int64, 100),
	}
}

type EmailWorker struct {
	db         *sqlx.DB
	afs        afero.Fs
	messageIDC chan int64
}

func (w EmailWorker) String() string {
	return "dahua.EmailWorker"
}

func (w EmailWorker) Serve(ctx context.Context) error {
	return sutureext.SanitizeError(ctx, w.serve(ctx))
}

func (w EmailWorker) serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case messageID := <-w.messageIDC:
			err := HandleEmail(ctx, w.db, w.afs, messageID)
			if err != nil {
				return err
			}
		}
	}
}

func (w EmailWorker) Register() EmailWorker {
	bus.Subscribe(w.String(), func(ctx context.Context, event bus.EmailCreated) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case w.messageIDC <- event.MessageKey.ID:
			return nil
		}
	})
	return w
}

func NewEmailRule(expression string) (EmailRule, error) {
	if expression == "" {
		expression = "true"
	}

	t, err := template.New("").Parse("{{" + expression + "}}")
	if err != nil {
		return EmailRule{}, err
	}

	return EmailRule{
		Template: t,
	}, nil
}

type EmailRule struct {
	Template *template.Template
}

func (r EmailRule) Match(email EmailPayload) (bool, error) {
	var buffer bytes.Buffer
	err := r.Template.Execute(&buffer, email)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(buffer.String())
}

func NewEmailPayload(message EmailMessage, attachments []EmailAttachment) EmailPayload {
	v := []EmailPayloadAttachment{}
	for _, a := range attachments {
		v = append(v, EmailPayloadAttachment{
			UUID:     a.UUID,
			FileName: a.File_Name,
			Size:     a.Size,
			Mime:     a.Mime_Type,
		})

	}
	return EmailPayload{
		Message: EmailPayloadMessage{
			UUID:              message.UUID,
			Date:              message.Date.Time,
			From:              message.From,
			To:                message.To.V,
			Subject:           message.Subject,
			Text:              message.Text,
			AlarmEvent:        message.Alarm_Event,
			AlarmInputChannel: message.Alarm_Input_Channel,
			AlarmName:         message.Alarm_Name,
			CreatedAt:         message.Created_At.Time,
		},
		Attachments: v,
	}
}

type EmailPayload struct {
	Message     EmailPayloadMessage
	Attachments []EmailPayloadAttachment
}

type EmailPayloadMessage struct {
	UUID              string
	Date              time.Time
	From              string
	To                []string
	Subject           string
	Text              string
	AlarmEvent        string
	AlarmInputChannel string
	AlarmName         string
	CreatedAt         time.Time
}

type EmailPayloadAttachment struct {
	UUID     string
	FileName string
	Size     int64
	Mime     string
}

func EmailRenderTemplate(tmpl string, data EmailPayload) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func NewSender(email EmailPayload, titleTemplate string, bodyTemplate string, attachments bool) (EmailSender, error) {
	title, err := EmailRenderTemplate(titleTemplate, email)
	if err != nil {
		return EmailSender{}, err
	}

	body, err := EmailRenderTemplate(bodyTemplate, email)
	if err != nil {
		return EmailSender{}, err
	}

	return EmailSender{
		Email:       email,
		Title:       title,
		Body:        body,
		Attachments: attachments,
	}, nil
}

type EmailSender struct {
	Email       EmailPayload
	Title       string
	Body        string
	Attachments bool
}

func (s EmailSender) Send(ctx context.Context, afs afero.Fs, goriseURL string) error {
	sender, err := gorise.Build(goriseURL)
	if err != nil {
		return err
	}

	var closers []io.Closer
	defer func() {
		for _, closer := range closers {
			closer.Close()
		}
	}()

	attachments := []gorise.Attachment{}
	if s.Attachments {
		for _, v := range s.Email.Attachments {
			file, err := afs.Open(v.UUID)
			if err != nil {
				slog.Warn("Failed to open attachment", "error", err, "attachment-uuid", v.UUID)
				continue
			}
			closers = append(closers, file)

			attachments = append(attachments, gorise.Attachment{
				Name:   v.FileName,
				Mime:   v.Mime,
				Reader: file,
			})
		}
	}

	return sender.Send(ctx, gorise.Message{
		Title:       s.Title,
		Body:        s.Body,
		Attachments: attachments,
	})
}
