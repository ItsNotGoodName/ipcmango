package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/jlaffaye/ftp"
	echo "github.com/labstack/echo/v4"
)

type DahuaFileCache interface {
	// Save(ctx context.Context, file models.DahuaFile, r io.ReadCloser) error
	Exists(ctx context.Context, file models.DahuaFile) (bool, error)
	Get(ctx context.Context, file models.DahuaFile) (io.ReadCloser, error)
}

func NewServer(
	pub pubsub.Pub,
	db repo.DB,
	dahuaStore *dahua.Store,
	dahuaFileCache DahuaFileCache,
) *Server {
	return &Server{
		pub:            pub,
		db:             db,
		dahuaStore:     dahuaStore,
		dahuaFileCache: dahuaFileCache,
	}
}

type Server struct {
	pub            pubsub.Pub
	db             repo.DB
	dahuaStore     *dahua.Store
	dahuaFileCache DahuaFileCache
}

func (s *Server) Register(e *echo.Echo) {
	e.GET("/v1/dahua", s.Dahua)
	e.GET("/v1/dahua-events", s.DahuaEvents)
	e.GET("/v1/dahua/:id/audio", s.DahuaIDAudio)
	e.GET("/v1/dahua/:id/coaxial/caps", s.DahuaIDCoaxialCaps)
	e.GET("/v1/dahua/:id/coaxial/status", s.DahuaIDCoaxialStatus)
	e.GET("/v1/dahua/:id/detail", s.DahuaIDDetail)
	e.GET("/v1/dahua/:id/error", s.DahuaIDError)
	e.GET("/v1/dahua/:id/events", s.DahuaIDEvents)
	e.GET("/v1/dahua/:id/files", s.DahuaIDFiles)
	e.GET("/v1/dahua/:id/files/*", s.DahuaIDFilesPath)
	e.GET("/v1/dahua/:id/licenses", s.DahuaIDLicenses)
	e.GET("/v1/dahua/:id/snapshot", s.DahuaIDSnapshot)
	e.GET("/v1/dahua/:id/software", s.DahuaIDSoftware)
	e.GET("/v1/dahua/:id/storage", s.DahuaIDStorage)
	e.GET("/v1/dahua/:id/users", s.DahuaIDUsers)

	e.POST("/v1/dahua/:id/ptz/preset", s.DahuaIDPTZPresetPOST)
	e.POST("/v1/dahua/:id/rpc", s.DahuaIDRPCPOST)
}

func (s *Server) Dahua(c echo.Context) error {
	ctx := c.Request().Context()

	dbDevices, err := s.db.ListDahuaDevice(ctx)
	if err != nil {
		return err
	}
	var conns []models.DahuaConn
	for _, v := range dbDevices {
		conns = append(conns, v.Convert().DahuaConn)
	}

	clients := s.dahuaStore.ClientList(ctx, conns)

	res := make([]models.DahuaStatus, 0, len(clients))
	for _, client := range clients {
		res = append(res, dahua.GetDahuaStatus(client.Conn, client.RPC))
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDRPCPOST(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	var req struct {
		Method string          `json:"method"`
		Params json.RawMessage `json:"params"`
		Object int64           `json:"object"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := dahuarpc.SendRaw[json.RawMessage](ctx, conn.RPC, dahuarpc.New(req.Method).
		Params(req.Params).
		Object(req.Object))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDDetail(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetDahuaDetail(ctx, conn.Conn.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDSoftware(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetSoftwareVersion(ctx, conn.Conn.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDLicenses(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetLicenseList(ctx, conn.Conn.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDError(c echo.Context) error {
	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res := dahua.GetError(conn.RPC)

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDSnapshot(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	snapshot, err := dahuacgi.SnapshotGet(ctx, conn.CGI, channel)
	if err != nil {
		return err
	}
	defer snapshot.Close()

	c.Response().Header().Set(echo.HeaderContentLength, snapshot.ContentLength)
	c.Response().Header().Set(echo.HeaderCacheControl, "no-store")

	_, err = io.Copy(c.Response().Writer, snapshot)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) DahuaIDEvents(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	direct, err := queryBoolOptional(c, "direct")
	if err != nil {
		return err
	}

	if direct {
		// Get events directly from the device

		manager, err := dahuacgi.EventManagerGet(ctx, conn.CGI, 0)
		if err != nil {
			return err
		}
		reader := manager.Reader()

		stream := useStream(c)

		for {
			err := reader.Poll()
			if err != nil {
				return sendStreamError(c, stream, err)
			}

			event, err := reader.ReadEvent()
			if err != nil {
				return sendStreamError(c, stream, err)
			}

			data := dahua.NewDahuaEvent(conn.Conn.ID, event)

			if err := sendStream(c, stream, data); err != nil {
				return err
			}
		}
	} else {
		// Get events from PubSub

		sub, eventsC, err := s.pub.SubscribeChan(ctx, 10, models.EventDahuaDeviceEvent{})
		if err != nil {
			return err
		}
		defer sub.Close()

		stream := useStream(c)

		for event := range eventsC {
			evt, ok := event.(models.EventDahuaDeviceEvent)
			if !ok {
				continue
			}

			err := sendStream(c, stream, evt.Event)
			if err != nil {
				return sendStreamError(c, stream, err)
			}
		}

		if err := sub.Error(); err != nil {
			return sendStreamError(c, stream, err)
		}

		return nil
	}
}

func (s *Server) DahuaEvents(c echo.Context) error {
	ctx := c.Request().Context()

	ids, err := queryInts(c, "id")
	if err != nil {
		return err
	}

	sub, eventsC, err := s.pub.SubscribeChan(ctx, 10, models.EventDahuaDeviceEvent{})
	if err != nil {
		return err
	}
	defer sub.Close()

	stream := useStream(c)

	for event := range eventsC {
		evt, ok := event.(models.EventDahuaDeviceEvent)
		if !ok {
			continue
		}

		if len(ids) != 0 && !slices.Contains(ids, evt.Event.DeviceID) {
			continue
		}

		if err := sendStream(c, stream, evt.Event); err != nil {
			return sendStreamError(c, stream, err)
		}
	}

	if err := sub.Error(); err != nil {
		return sendStreamError(c, stream, err)
	}

	return nil
}

func (s *Server) DahuaIDFiles(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	var form struct {
		Start string
		End   string
	}
	if err := DecodeQuery(c, &form); err != nil {
		return err
	}

	timeRange, err := UseTimeRange(form.Start, form.End)
	if err != nil {
		return err
	}

	iterator := dahua.NewScannerPeriodIterator(timeRange)

	filesC := make(chan []mediafilefind.FindNextFileInfo)
	stream := useStream(c)

	for period, ok := iterator.Next(); ok; period, ok = iterator.Next() {
		cancel, errC := dahua.Scanner(ctx, conn.RPC, period, conn.Conn.Location, filesC)
		defer cancel()

	inner:
		for {
			select {
			case <-ctx.Done():
				return sendStreamError(c, stream, ctx.Err())
			case err := <-errC:
				if err != nil {
					return sendStreamError(c, stream, err)
				}
				break inner
			case files := <-filesC:
				res, err := dahua.NewDahuaFiles(conn.Conn.ID, files, dahua.GetSeed(conn.Conn), conn.Conn.Location)
				if err != nil {
					return sendStreamError(c, stream, err)
				}

				if err := sendStream(c, stream, res); err != nil {
					return sendStreamError(c, stream, err)
				}
			}
		}
	}

	return nil
}

func (s *Server) DahuaIDFilesPath(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	// FIXME: vite dev proxy removes double slashes which messes up the file path
	filePath := c.Param("*")
	// if !strings.HasPrefix(filePath, "/") && !strings.Contains(strings.SplitN(filePath, "/", 2)[0], "://") {
	// 	filePath = "/" + filePath
	// }

	dbDahuaFile, err := s.db.GetDahuaFileByFilePath(ctx, repo.GetDahuaFileByFilePathParams{
		DeviceID: conn.Conn.ID,
		FilePath: filePath,
	})
	if err != nil {
		if repo.IsNotFound(err) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}
	dahuaFile := dbDahuaFile.Convert()

	storage := core.StorageFromFilePath(dahuaFile.FilePath)

	var rd io.ReadCloser
	switch storage {
	case models.StorageLocal:
		if exists, err := s.dahuaFileCache.Exists(ctx, dahuaFile); err != nil {
			return err
		} else if !exists {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, dahuarpc.LoadFileURL(dahua.NewHTTPAddress(conn.Conn.Address), filePath), nil)
			if err != nil {
				return err
			}

			session, err := conn.RPC.Session(ctx)
			if err != nil {
				return err
			}

			req.Header.Add("Cookie", dahuarpc.Cookie(session))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			rd = resp.Body
		} else {
			rd, err = s.dahuaFileCache.Get(ctx, dahuaFile)
			if err != nil {
				return err
			}
		}
	case models.StorageFTP:
		u, err := url.Parse(dahuaFile.FilePath)
		if err != nil {
			return err
		}

		cred, err := s.db.GetDahuaCredential(ctx, repo.GetDahuaCredentialParams{
			ServerAddress: u.Hostname(),
			Storage:       storage,
		})
		if err != nil {
			return echo.ErrNotFound.WithInternal(err)
		}

		c, err := ftp.Dial(cred.ServerAddress+":"+strconv.FormatInt(cred.Port, 10), ftp.DialWithContext(ctx))
		if err != nil {
			return err
		}

		err = c.Login(cred.Username, cred.Password)
		if err != nil {
			return err
		}
		defer c.Quit()

		username := "/" + cred.Username
		path, _ := strings.CutPrefix(u.Path, username)

		rd, err = c.Retr(path)
		if err != nil {
			return err
		}
	}
	defer rd.Close()

	_, err = io.Copy(c.Response().Writer, rd)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) DahuaIDAudio(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	audioStream, err := dahuacgi.AudioStreamGet(ctx, conn.CGI, channel, dahuacgi.HTTPTypeSinglePart)
	if err != nil {
		return err
	}

	c.Request().Header.Add("ContentType", audioStream.ContentType)

	_, err = io.Copy(c.Response().Writer, audioStream)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) DahuaIDCoaxialStatus(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialStatus(ctx, conn.Conn.ID, conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *Server) DahuaIDCoaxialCaps(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialCaps(ctx, conn.Conn.ID, conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *Server) DahuaIDPTZPresetPOST(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	index, err := queryIntOptional(c, "index")
	if err != nil {
		return err
	}

	err = dahua.SetPreset(ctx, conn.PTZ, channel, index)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (s *Server) DahuaIDStorage(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	storage, err := dahua.GetStorage(ctx, conn.Conn.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, storage)
}

func (s *Server) DahuaIDUsers(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetUsers(ctx, conn.Conn.ID, conn.RPC, conn.Conn.Location)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
