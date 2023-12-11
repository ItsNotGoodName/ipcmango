package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	echo "github.com/labstack/echo/v4"
)

func (s *Server) RegisterDahuaRoutes(e *echo.Echo) {
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

func useDahuaConn(c echo.Context, connStore DahuaConnStore, store *dahua.Store) (dahua.Conn, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return dahua.Conn{}, echo.ErrBadRequest.WithInternal(err)
	}

	ctx := c.Request().Context()

	camera, found, err := connStore.Get(ctx, id)
	if err != nil {
		return dahua.Conn{}, err
	}
	if !found {
		return dahua.Conn{}, echo.ErrNotFound.WithInternal(err)
	}

	client := store.Conn(ctx, camera)

	return client, nil
}

func (s *Server) Dahua(c echo.Context) error {
	ctx := c.Request().Context()

	cameras, err := s.dahuaConnStore.List(ctx)
	if err != nil {
		return err
	}

	conns := s.dahuaStore.ConnList(ctx, cameras)

	res := make([]models.DahuaStatus, 0, len(conns))
	for _, conn := range conns {
		res = append(res, dahua.GetDahuaStatus(conn.Camera, conn.RPC))
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDRPCPOST(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	var req struct {
		Method string          `json:"method"`
		Params json.RawMessage `json:"params"`
		Object int64           `json:"object"`
		Seq    int             `json:"seq"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	rpc, err := conn.RPC.RPC(ctx)
	if err != nil {
		return err
	}

	res, err := dahuarpc.SendRaw[json.RawMessage](ctx, rpc.
		Method(req.Method).
		Params(req.Params).
		Object(req.Object).
		Seq(req.Seq))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDDetail(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetDahuaDetail(ctx, conn.Camera.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDSoftware(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetSoftwareVersion(ctx, conn.Camera.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDLicenses(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetLicenseList(ctx, conn.Camera.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDError(c echo.Context) error {
	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	res := dahua.GetError(conn.RPC)

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaIDSnapshot(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
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

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	direct, err := queryBoolOptional(c, "direct")
	if err != nil {
		return err
	}

	if direct {
		// Get events directly from the camera

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

			data := dahua.NewDahuaEvent(conn.Camera.ID, event, time.Now())

			if err := sendStream(c, stream, data); err != nil {
				return err
			}
		}
	} else {
		// Get events from PubSub

		ctx := ctx
		stream := useStream(c)

		sub, err := s.pub.Subscribe(ctx, func(event pubsub.Event) error {
			evt, ok := event.(models.EventDahuaCameraEvent)
			if !ok {
				return nil
			}

			return sendStream(c, stream, evt.Event)
		}, models.EventDahuaCameraEvent{})
		if err != nil {
			return err
		}
		defer sub.Close()

		return sub.Wait(ctx)
	}
}

func (s *Server) DahuaEvents(c echo.Context) error {
	ctx := c.Request().Context()

	ids, err := queryInts(c, "id")
	if err != nil {
		return err
	}

	stream := useStream(c)

	sub, err := s.pub.Subscribe(ctx, func(event pubsub.Event) error {
		evt, ok := event.(models.EventDahuaCameraEvent)
		if !ok {
			return nil
		}

		if len(ids) != 0 && !slices.Contains(ids, evt.Event.CameraID) {
			return nil
		}

		return sendStream(c, stream, evt.Event)
	}, models.EventDahuaCameraEvent{})
	defer sub.Close()

	return sub.Wait(ctx)
}

func (s *Server) DahuaIDFiles(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	var form struct {
		Start string
		End   string
	}
	if err := ParseQuery(c, &form); err != nil {
		return err
	}

	scanRange, err := queryDahuaScanRange(form.Start, form.End)
	if err != nil {
		return err
	}

	iter := dahua.NewScanPeriodIterator(scanRange)

	filesC := make(chan []mediafilefind.FindNextFileInfo)
	stream := useStream(c)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for period, ok := iter.Next(); ok; period, ok = iter.Next() {
		errC := dahua.Scan(ctx, conn.RPC, period, conn.Camera.Location, filesC)

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
				res, err := dahua.NewDahuaFiles(conn.Camera.ID, files, dahua.GetSeed(conn.Camera), conn.Camera.Location)
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

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	filePath := c.Param("*")

	dahuaFile, err := s.db.GetDahuaFileByFilePath(ctx, filePath)
	if err != nil {
		return err
	}

	var rd io.ReadCloser
	if exists, err := s.dahuaFileCache.Exists(ctx, dahuaFile); err != nil {
		return err
	} else if !exists {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, dahuarpc.LoadFileURL(dahua.NewHTTPAddress(conn.Camera.Address), filePath), nil)
		if err != nil {
			return err
		}

		session, err := conn.RPC.RPCSession(ctx)
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
	defer rd.Close()

	_, err = io.Copy(c.Response().Writer, rd)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) DahuaIDAudio(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
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

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialStatus(ctx, conn.Camera.ID, conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *Server) DahuaIDCoaxialCaps(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialCaps(ctx, conn.Camera.ID, conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *Server) DahuaIDPTZPresetPOST(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
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

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	storage, err := dahua.GetStorage(ctx, conn.Camera.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, storage)
}

func (s *Server) DahuaIDUsers(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaConnStore, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetUsers(ctx, conn.Camera.ID, conn.RPC, conn.Camera.Location)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
