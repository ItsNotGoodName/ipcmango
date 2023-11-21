package webserver

import (
	"bytes"
	"cmp"
	"context"
	"net/http"
	"slices"
	"strconv"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	webcore "github.com/ItsNotGoodName/ipcmanview/internal/web/core"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func RegisterMiddleware(e *echo.Echo) {
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: web.AssetFS(),
	}))
}

func RegisterRoutes(e *echo.Echo, w Server) {
	e.DELETE("/dahua/cameras/:id", w.DahuaCamerasIDDelete)
	e.GET("/", w.Index)
	e.GET("/dahua", w.Dahua)
	e.GET("/dahua/cameras", w.DahuaCameras)
	e.GET("/dahua/cameras/create", w.DahuaCamerasCreate)
	e.GET("/dahua/data", w.DahuaData)
	e.GET("/dahua/snapshots", w.DahuaSnapshots)
	e.GET("/dahua/events/stream", w.DahuaEventStream)
	e.GET("/dahua/events", w.DahuaEvent)
	e.POST("/dahua/cameras/create", w.DahuaCamerasCreatePOST)
}

type Server struct {
	db         *sqlc.Queries
	dahuaStore *dahua.Store
	pubSub     api.PubSub
}

func New(db *sqlc.Queries, dahuaStore *dahua.Store, pubSub api.PubSub) Server {
	return Server{
		db:         db,
		dahuaStore: dahuaStore,
		pubSub:     pubSub,
	}
}

func (s Server) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (s Server) Dahua(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua", nil)
}

func (s Server) DahuaEvent(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-events", nil)
}

func (s Server) DahuaEventStream(c echo.Context) error {
	w := c.Response()

	w.Header().Set(echo.HeaderContentType, "text/event-stream")
	w.Header().Set(echo.HeaderCacheControl, "no-cache")
	w.Header().Set(echo.HeaderConnection, "keep-alive")

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	eventsC, err := s.pubSub.SubscribeDahuaEvents(ctx, []string{})
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-eventsC:
			if err := c.Echo().Renderer.Render(buf, "dahua-events", TemplateBlock{
				"event-row",
				Data{
					"Event": event.Event,
				},
			}, c); err != nil {
				return err
			}
			w.Write(formatSSE("message", buf.String()))
			buf.Reset()
			w.Flush()
		}
	}
}

func (s Server) DahuaCamerasIDDelete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := s.db.DeleteDahuaCamera(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s Server) DahuaCameras(c echo.Context) error {
	cameras, err := s.db.ListDahuaCamera(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-cameras", Data{
		"Cameras": cameras,
	})
}

func (s Server) DahuaCamerasCreate(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-cameras-create", nil)
}

func (s Server) DahuaCamerasCreatePOST(c echo.Context) error {
	var form struct {
		Name     string
		Address  string
		Username string
		Password string
		Location string
	}
	if err := parseForm(c, &form); err != nil {
		return err
	}

	location, err := core.NewLocation(form.Location)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)

	}
	dto, err := dahua.NewDahuaCamera("", models.DTODahuaCamera{
		Address:  form.Address,
		Username: form.Username,
		Password: form.Password,
		Location: location,
	})

	ctx := c.Request().Context()

	_, err = s.db.CreateDahuaCamera(ctx, sqlc.CreateDahuaCameraParams{
		Name:      form.Name,
		Username:  dto.Username,
		Password:  dto.Password,
		Address:   dto.Address,
		Location:  dto.Location,
		CreatedAt: dto.CreatedAt,
	})
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/dahua/cameras")
}

func (s Server) DahuaSnapshots(c echo.Context) error {
	cameras, err := s.db.ListDahuaCamera(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-snapshots", Data{
		"Cameras": cameras,
	})
}

func (s Server) DahuaData(c echo.Context) error {
	ctx := c.Request().Context()

	cameras, err := s.db.ListDahuaCamera(ctx)
	if err != nil {
		return err
	}

	type cameraData struct {
		detail               models.DahuaDetail
		softwareVersion      models.DahuaSoftwareVersion
		licenses             []models.DahuaLicense
		storage              []models.DahuaStorage
		coaxialcontrolStatus []models.DahuaCoaxialStatus
	}

	cameraDataC := make(chan cameraData, len(cameras))
	wg := sync.WaitGroup{}
	for _, camera := range cameras {
		wg.Add(1)
		go func(camera sqlc.DahuaCamera) {
			defer wg.Done()

			conn := s.dahuaStore.ConnByCamera(ctx, webcore.ConvertDahuaCamera(camera))
			log := log.With().Str("id", conn.Camera.ID).Logger()

			var data cameraData

			{
				res, err := dahua.GetDahuaDetail(ctx, conn.Camera.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get detail")
				} else {
					data.detail = res
				}
			}

			{
				res, err := dahua.GetSoftwareVersion(ctx, conn.Camera.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get software version")
				} else {
					data.softwareVersion = res
				}
			}

			{
				res, err := dahua.GetLicenseList(ctx, conn.Camera.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get licenses")
				} else {
					data.licenses = res
				}
			}

			{
				res, err := dahua.GetStorage(ctx, conn.Camera.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get storage")
				} else {
					data.storage = res
				}
			}

			{
				caps, err := dahua.GetCoaxialCaps(ctx, conn.Camera.ID, conn.RPC, 0)
				if err != nil {
					log.Err(err).Msg("Failed to get coaxial caps")
				} else if caps.SupportControlLight || caps.SupportControlSpeaker || caps.SupportControlFullcolorLight {
					res, err := dahua.GetCoaxialStatus(ctx, conn.Camera.ID, conn.RPC, 0)
					if err != nil {
						log.Err(err).Msg("Failed to get coaxial status")
					} else {
						data.coaxialcontrolStatus = append(data.coaxialcontrolStatus, res)
					}
				}
			}

			cameraDataC <- data
		}(camera)
	}
	wg.Wait()
	close(cameraDataC)

	conns, err := s.dahuaStore.ConnListByCameras(ctx)
	if err != nil {
		return err
	}
	status := make([]models.DahuaStatus, 0, len(conns))
	for _, conn := range conns {
		status = append(status, dahua.GetDahuaStatus(conn.Camera, conn.RPC.Conn))
	}

	details := make([]models.DahuaDetail, 0, len(cameras))
	softwareVersions := make([]models.DahuaSoftwareVersion, 0, len(cameras))
	licenses := make([]models.DahuaLicense, 0, len(cameras))
	storage := make([]models.DahuaStorage, 0, len(cameras))
	coaxialStatus := make([]models.DahuaCoaxialStatus, 0, len(cameras))
	for data := range cameraDataC {
		details = append(details, data.detail)
		softwareVersions = append(softwareVersions, data.softwareVersion)
		licenses = append(licenses, data.licenses...)
		storage = append(storage, data.storage...)
		coaxialStatus = append(coaxialStatus, data.coaxialcontrolStatus...)
	}
	slices.SortFunc(details, func(a, b models.DahuaDetail) int { return cmp.Compare(a.CameraID, b.CameraID) })
	slices.SortFunc(softwareVersions, func(a, b models.DahuaSoftwareVersion) int { return cmp.Compare(a.CameraID, b.CameraID) })
	slices.SortFunc(licenses, func(a, b models.DahuaLicense) int { return cmp.Compare(a.CameraID, b.CameraID) })
	slices.SortFunc(storage, func(a, b models.DahuaStorage) int { return cmp.Compare(a.CameraID, b.CameraID) })
	slices.SortFunc(coaxialStatus, func(a, b models.DahuaCoaxialStatus) int { return cmp.Compare(a.CameraID, b.CameraID) })

	return c.Render(http.StatusOK, "dahua-data", Data{
		"Cameras":          cameras,
		"Status":           status,
		"Details":          details,
		"SoftwareVersions": softwareVersions,
		"Licenses":         licenses,
		"Storage":          storage,
		"CoaxialStatus":    coaxialStatus,
	})
}
