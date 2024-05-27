package dahua

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/gorise"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
)

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
}

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

	_, err := tx.ExecContext(ctx, `
			INSERT INTO dahua_email_attachments (
				uuid,
				message_id,
				file_name,
				size
			) VALUES (?, ?, ?, ?)
		`, attachmentUUID, messageID, fileName, len(content))
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

func SendEmailToEndpoints(ctx context.Context, db *sqlx.DB, afs afero.Fs, messageKey core.Key) error {
	var endpoints []system.Endpoint
	err := sqlx.Select(db, &endpoints, `
		SELECT * FROM endpoints
	`)
	if err != nil {
		return err
	}

	var message EmailMessage
	err = sqlx.GetContext(ctx, db, &message, `
		SELECT * FROM dahua_email_messages WHERE id = ?
	`, messageKey.ID)
	if err != nil {
		return err
	}

	var attachments []EmailAttachment
	err = sqlx.SelectContext(ctx, db, &attachments, `
		SELECT * FROM dahua_email_attachments WHERE message_id = ?
	`, messageKey.ID)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(endpoint system.Endpoint) {
			defer wg.Done()

			slog := slog.With("endpoint-uuid", endpoint.UUID)

			sender, err := gorise.Build(endpoint.Gorise_URL)
			if err != nil {
				slog.Error("Failed to build gorise url", "error", err)
				return
			}

			var closers []io.Closer
			defer func() {
				for _, closer := range closers {
					closer.Close()
				}
			}()

			goriseAttachments := []gorise.Attachment{}
			for _, v := range attachments {
				file, err := afs.Open(v.UUID)
				if err != nil {
					slog.Error("Failed to open attachment", "error", err, "attachment-uuid", v.UUID)
					continue
				}
				closers = append(closers, file)

				goriseAttachments = append(goriseAttachments, gorise.Attachment{
					Name:   v.File_Name,
					Mime:   "image/jpeg",
					Reader: file,
				})
			}

			err = sender.Send(ctx, gorise.Message{
				Title:       message.Subject,
				Body:        message.Text,
				Attachments: goriseAttachments,
			})
			if err != nil {
				slog.Error("Failed to send to endpoint", "error", err)
				return
			}

		}(endpoint)
	}

	wg.Wait()

	return nil
}

func RegisterEmailToEndpoints(db *sqlx.DB, afs afero.Fs) {
	bus.Subscribe("RegisterEmailToEndpoints", func(ctx context.Context, event bus.EmailCreated) error {
		go func() {
			err := SendEmailToEndpoints(ctx, db, afs, event.MessageKey)
			if err != nil {
				slog.Error("Failed to send email to endpoints", "error", err)
			}
		}()
		return nil
	})
}
