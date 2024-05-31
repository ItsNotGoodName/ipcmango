package smtp

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"net/mail"
	"slices"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
)

type Server struct {
	server *smtp.Server
	db     *sqlx.DB
	afs    afero.Fs
}

func NewServer(db *sqlx.DB, afs afero.Fs, address string) Server {
	server := smtp.NewServer(nil)

	server.Addr = address
	server.Domain = "localhost"
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second
	server.MaxMessageBytes = 25 * 1024 * 1024
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true

	return Server{
		server: server,
		db:     db,
		afs:    afs,
	}
}

func (Server) String() string {
	return "smtp.Server"
}

func (s Server) Serve(ctx context.Context) error {
	slog.Info("Starting SMTP server", "address", s.server.Addr)

	s.server.Backend = smtp.BackendFunc(func(c *smtp.Conn) (smtp.Session, error) {
		address := c.Conn().RemoteAddr().String()
		log := slog.With("address", address)

		return &Session{
			ctx:     ctx,
			log:     log,
			db:      s.db,
			afs:     s.afs,
			address: address,
			from:    "",
			to:      "",
		}, nil
	})

	errC := make(chan error, 1)

	go func() { errC <- s.server.ListenAndServe() }()

	select {
	case err := <-errC:
		return err
	case <-ctx.Done():
	}

	s.server.Close()

	return nil
}

// Session is returned after EHLO.
type Session struct {
	ctx     context.Context
	log     *slog.Logger
	db      *sqlx.DB
	afs     afero.Fs
	address string
	from    string
	to      string
}

func (s *Session) AuthPlain(username, password string) error {
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from

	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = to

	return nil
}

func (s *Session) Data(r io.Reader) error {
	ctx := s.ctx
	slog := s.log
	db := s.db
	afs := s.afs

	// Get device by IP or email
	ip, _ := core.SplitAddress(s.address)
	var device struct {
		types.Key
		Name string
	}
	err := db.GetContext(ctx, &device, `
		SELECT id, uuid, name FROM dahua_devices WHERE ip = ? OR email = ? LIMIT 1
	`, ip, s.from)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		slog.Error("Failed to get device", "error", err)
		return err
	}
	slog = slog.With("device", device.Name)

	// Read
	e, err := enmime.ReadEnvelope(r)
	if err != nil {
		slog.Error("Failed to read envelope", "error", err)
		return err
	}

	// Parse to
	to := []string{s.to}
	if addresses, err := e.AddressList("To"); err == nil {
		for _, t := range addresses {
			to = append(to, t.Address)
		}
	} else {
		slog.Warn("Failed to get 'To' from address list", "error", err)
	}
	to = slices.Compact(to)

	// Parse date
	date, err := e.Date()
	if err != nil && !errors.Is(err, mail.ErrHeaderNotPresent) {
		slog.Warn("Failed to parse date", "error", err, "date", e.GetHeader("Date"))
	}

	// Parse subject
	subject := e.GetHeader("Subject")

	// Parse content
	content := dahua.ParseEmailContent(e.Text)

	// Create email
	attachments := make([]dahua.CreateEmailParamsAttachment, 0, len(e.Attachments))
	for _, a := range e.Attachments {
		attachments = append(attachments, dahua.CreateEmailParamsAttachment{
			FileName: a.FileName,
			Content:  a.Content,
		})
	}
	arg := dahua.CreateEmailArgs{
		DeviceKey:         device.Key,
		Date:              types.NewTime(date),
		From:              s.from,
		To:                types.NewSlice(to),
		Subject:           subject,
		Text:              e.Text,
		AlarmEvent:        content.AlarmEvent,
		AlarmName:         content.AlarmName,
		Attachments:       attachments,
		AlarmInputChannel: content.AlarmInputChannel,
	}
	messageID, err := dahua.CreateEmail(ctx, db, afs, arg)
	if err != nil {
		slog.Error("Failed to create email", "error", err)
		return err
	}
	slog.Info("Created email", "message-id", messageID)

	return nil
}

func (s *Session) Reset() {
}

func (s *Session) Logout() error {
	return nil
}

func (s *Session) AuthMechanisms() []string {
	return []string{
		sasl.Login,
	}
}

func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewLoginServer(func(username, password string) error {
		return s.AuthPlain(username, password)
	}), nil
}
