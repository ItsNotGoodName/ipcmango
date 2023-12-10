package webserver

import (
	"bytes"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	webdahua "github.com/ItsNotGoodName/ipcmanview/internal/web/dahua"
	"github.com/ItsNotGoodName/ipcmanview/pkg/htmx"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.GET("/dahua/cameras/:id/update", w.DahuaCamerasUpdate)
	e.GET("/dahua/cameras/create", w.DahuaCamerasCreate)
	e.GET("/dahua/events", w.DahuaEvents)
	e.GET("/dahua/events/live", w.DahuaEventsLive)
	e.GET("/dahua/events/stream", w.DahuaEventStream)
	e.GET("/dahua/events/:id/data", w.DahuaEventsIDData)
	e.GET("/dahua/files", w.DahuaFiles)
	e.GET("/dahua/snapshots", w.DahuaSnapshots)

	e.POST("/dahua/cameras/create", w.DahuaCamerasCreatePOST)
	e.POST("/dahua/cameras/:id/update", w.DahuaCamerasUpdatePOST)
}

type Server struct {
	db         sqlc.DB
	pub        pubsub.Pub
	dahuaStore *dahua.Store
	dahuaBus   *dahua.Bus
}

func New(db sqlc.DB, pub pubsub.Pub, dahuaStore *dahua.Store, dahuaBus *dahua.Bus) Server {
	return Server{
		db:         db,
		pub:        pub,
		dahuaStore: dahuaStore,
		dahuaBus:   dahuaBus,
	}
}

func (s Server) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (s Server) Dahua(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua", nil)
}

func (s Server) DahuaEvents(c echo.Context) error {
	ctx := c.Request().Context()

	params := struct {
		CameraID  []int64
		Code      []string
		Action    []string
		Page      int64 `validate:"gt=0"`
		PerPage   int64
		Data      bool
		Start     string
		End       string
		Ascending bool
	}{
		Page:    1,
		PerPage: 10,
	}
	if err := api.ParseQuery(c, &params); err != nil {
		return err
	}
	if err := api.ValidateStruct(params); err != nil {
		return err
	}

	var start, end time.Time
	if params.Start != "" {
		var err error
		start, err = time.ParseInLocation("2006-01-02T15:04", params.Start, time.Local)
		if err != nil {
			return echo.ErrBadRequest.WithInternal(err)
		}
	}

	if params.End != "" {
		var err error
		end, err = time.ParseInLocation("2006-01-02T15:04", params.End, time.Local)
		if err != nil {
			return echo.ErrBadRequest.WithInternal(err)
		}
	}

	events, err := s.db.ListDahuaEvent(ctx, sqlc.ListDahuaEventParams{
		CameraID: params.CameraID,
		Code:     params.Code,
		Action:   params.Action,
		Page: pagination.Page{
			Page:    int(params.Page),
			PerPage: int(params.PerPage),
		},
		Start:     types.NewTime(start),
		End:       types.NewTime(end),
		Ascending: params.Ascending,
	})
	if err != nil {
		return err
	}

	eventCodes, err := s.db.ListDahuaEventCodes(ctx)
	if err != nil {
		return err
	}

	eventActions, err := s.db.ListDahuaEventActions(ctx)
	if err != nil {
		return err
	}

	cameras, err := s.db.ListDahuaCamera(ctx)
	if err != nil {
		return err
	}

	data := Data{
		"Params":      params,
		"Cameras":     cameras,
		"Events":      events,
		"EventCodes":  eventCodes,
		"EventAction": eventActions,
	}

	if htmx.GetRequest(c.Request()) && !htmx.GetBoosted(c.Request()) {
		htmx.SetReplaceURL(c.Response(), "/dahua/events?"+api.NewQuery(params).Encode())
		return c.Render(http.StatusOK, "dahua-events", TemplateBlock{"htmx", data})
	}

	return c.Render(http.StatusOK, "dahua-events", data)
}

func (s Server) DahuaEventsIDData(c echo.Context) error {
	id, err := api.PathID(c)
	if err != nil {
		return err
	}

	data, err := s.db.GetDahuaEventData(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-events", TemplateBlock{
		"hx-event-data",
		data,
	})
}

func (s Server) DahuaFiles(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-files", nil)
}

func (s Server) DahuaEventsLive(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-events-live", nil)
}

func (s Server) DahuaEventStream(c echo.Context) error {
	ctx := c.Request().Context()
	w := c.Response()

	w.Header().Set(echo.HeaderContentType, "text/event-stream")
	w.Header().Set(echo.HeaderCacheControl, "no-cache")
	w.Header().Set(echo.HeaderConnection, "keep-alive")

	buf := new(bytes.Buffer)
	sub, err := s.pub.Subscribe(ctx, func(event pubsub.Event) error {
		evt, ok := event.(models.EventDahuaCameraEvent)
		if !ok {
			return nil
		}

		if err := c.Echo().Renderer.Render(buf, "dahua-events-live", TemplateBlock{
			"event-row",
			Data{
				"Event": evt.Event,
			},
		}, c); err != nil {
			return err
		}
		_, err := w.Write(api.FormatSSE("message", buf.String()))
		if err != nil {
			return err
		}
		buf.Reset()
		w.Flush()
		return nil
	}, models.EventDahuaCameraEvent{})
	if err != nil {
		return err
	}
	defer sub.Close()

	return sub.Wait(ctx)
}

func (s Server) DahuaCamerasIDDelete(c echo.Context) error {
	id, err := api.PathID(c)
	if err != nil {
		return err
	}

	if err := s.db.DeleteDahuaCamera(c.Request().Context(), id); err != nil {
		return err
	}
	s.dahuaBus.CameraDeleted(id)

	return c.NoContent(http.StatusOK)
}

func (s Server) DahuaCameras(c echo.Context) error {
	if htmx.GetRequest(c.Request()) && !htmx.GetBoosted(c.Request()) {
		tables, err := useDahuaTables(c.Request().Context(), s.db, s.dahuaStore)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "dahua-cameras", TemplateBlock{
			"htmx",
			tables,
		})
	}

	cameras, err := s.db.ListDahuaCamera(c.Request().Context())
	if err != nil {
		return err
	}

	fileCursors, err := s.db.ListDahuaFileCursor(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-cameras", Data{
		"Cameras":     cameras,
		"FileCursors": fileCursors,
	})
}

func (s Server) DahuaCamerasCreate(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-cameras-create", Data{
		"Locations": webdahua.Locations,
	})
}

func (s Server) DahuaCamerasCreatePOST(c echo.Context) error {
	var form struct {
		Name     string
		Address  string
		Username string
		Password string
		Location string
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}
	if form.Name == "" {
		form.Name = form.Address
	}
	if form.Username == "" {
		form.Username = "admin"
	}
	location, err := core.NewLocation(form.Location)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	camera, err := dahua.NewDahuaCamera(models.DahuaCamera{
		Address:  form.Address,
		Username: form.Username,
		Password: form.Password,
		Location: location,
	})

	ctx := c.Request().Context()

	id, err := s.db.CreateDahuaCamera(ctx, sqlc.CreateDahuaCameraParams{
		Name:      form.Name,
		Username:  camera.Username,
		Password:  camera.Password,
		Address:   camera.Address,
		Location:  camera.Location,
		CreatedAt: types.NewTime(camera.CreatedAt),
		UpdatedAt: types.NewTime(camera.CreatedAt),
	}, webdahua.NewFileCursor())
	if err != nil {
		return err
	}
	dbCamera, err := s.db.GetDahuaCamera(ctx, id)
	if err != nil {
		return err
	}

	s.dahuaBus.CameraCreated(webdahua.ConvertGetDahuaCameraRow(dbCamera))

	return c.Redirect(http.StatusSeeOther, "/dahua/cameras")
}

func (s Server) DahuaCamerasUpdate(c echo.Context) error {
	camera, err := useDahuaCamera(c, s.db)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-cameras-update", Data{
		"Locations": webdahua.Locations,
		"Camera":    camera,
	})
}

func (s Server) DahuaCamerasUpdatePOST(c echo.Context) error {
	camera, err := useDahuaCamera(c, s.db)
	if err != nil {
		return err
	}

	var form struct {
		Name     string
		Address  string
		Username string
		Password string
		Location string
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}
	location, err := core.NewLocation(form.Location)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if form.Password == "" {
		form.Password = camera.Password
	}

	dto, err := dahua.UpdateDahuaCamera(models.DahuaCamera{
		ID:       camera.ID,
		Address:  form.Address,
		Username: form.Username,
		Password: form.Password,
		Location: location,
	})
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	_, err = s.db.UpdateDahuaCamera(ctx, sqlc.UpdateDahuaCameraParams{
		ID:        dto.ID,
		Name:      form.Name,
		Username:  dto.Username,
		Password:  dto.Password,
		Address:   dto.Address,
		Location:  dto.Location,
		UpdatedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		return err
	}

	dbCamera, err := s.db.GetDahuaCamera(ctx, camera.ID)
	if err != nil {
		return err
	}
	s.dahuaBus.CameraUpdated(webdahua.ConvertGetDahuaCameraRow(dbCamera))

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
