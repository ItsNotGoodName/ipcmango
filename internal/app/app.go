package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"slices"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/coaxialcontrolio"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
)

func useDevice(ctx context.Context, db *sqlx.DB, uuid string) (dahua.DahuaDevice, error) {
	var device dahua.DahuaDevice
	err := db.GetContext(ctx, &device, `
		SELECT * FROM dahua_devices WHERE uuid = ?
	`, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dahua.DahuaDevice{}, huma.Error404NotFound("Device not found")
		}
		return dahua.DahuaDevice{}, err
	}

	return device, nil
}

func useClient(ctx context.Context, dahuaStore *dahua.Store, device dahua.DahuaDevice) (dahua.Client, error) {
	client, err := dahuaStore.GetClient(ctx, dahua.NewConn(device))
	if err != nil {
		return dahua.Client{}, huma.Error404NotFound("Device not found")
	}
	return client, nil
}

func NewHumaConfig() huma.Config {
	return huma.DefaultConfig("IPCManView API", "1.0.0")
}

func Register(api huma.API, db *sqlx.DB, afs afero.Fs, dahuaStore *dahua.Store) {
	// Devices
	huma.Register(api, huma.Operation{
		Summary: "List devices",
		Method:  http.MethodGet,
		Path:    "/api/devices",
	}, func(ctx context.Context, input *struct{}) (*ListDevicesOutput, error) {
		rows, err := db.QueryxContext(ctx, `
			SELECT * FROM dahua_devices
		`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		body := []Device{}
		for rows.Next() {
			var v dahua.DahuaDevice
			if err := rows.StructScan(&v); err != nil {
				return nil, err
			}
			body = append(body, NewDevice(v))
		}

		return &ListDevicesOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*GetDeviceOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		return &GetDeviceOutput{
			Body: NewDevice(device),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Create device",
		Method:  http.MethodPost,
		Path:    "/api/devices",
	}, func(ctx context.Context, input *CreateDeviceInput) (*CreateDevicesOutput, error) {
		device, err := dahua.CreateDevice(ctx, db, dahua.CreateDeviceArgs{
			Name:            input.Body.Name,
			IP:              input.Body.IP,
			Username:        input.Body.Username,
			Password:        input.Body.Password,
			Location:        input.Body.Location,
			Features:        types.NewSlice(input.Body.Features),
			Email:           core.NullToSQLNull(input.Body.Email),
			Latitude:        core.NullToSQLNull(input.Body.Latitude),
			Longitude:       core.NullToSQLNull(input.Body.Longitude),
			SunriseOffset:   input.Body.SunriseOffset,
			SunsetOffset:    input.Body.SunsetOffset,
			SyncVideoInMode: core.NullToSQLNull(input.Body.SyncVideoInMode),
		})
		if err != nil {
			return nil, err
		}

		return &CreateDevicesOutput{
			Body: NewDevice(device),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Update device",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
		Body UpdateDevice
	}) (*UpdateDevicesOutput, error) {
		device, err := dahua.UpdateDevice(ctx, db, dahua.UpdateDeviceArgs{
			UUID:          input.UUID,
			Name:          input.Body.Name,
			IP:            input.Body.IP,
			Username:      input.Body.Username,
			Password:      core.NullToSQLNull(input.Body.Password),
			Location:      input.Body.Location,
			Features:      types.NewSlice(input.Body.Features),
			Email:         input.Body.Email,
			Latitude:      core.NullToSQLNull(input.Body.Latitude),
			Longitude:     core.NullToSQLNull(input.Body.Longitude),
			SunriseOffset: input.Body.SunriseOffset,
			SunsetOffset:  input.Body.SunsetOffset,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, huma.Error404NotFound("Device not found")
			}
			return nil, err
		}

		return &UpdateDevicesOutput{
			Body: NewDevice(device),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Delete device",
		Method:  http.MethodDelete,
		Path:    "/api/devices/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*struct{}, error) {
		return &struct{}{}, dahua.DeleteDevice(ctx, db, input.UUID)
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device coaxial caps",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/coaxial/caps",
	}, func(ctx context.Context, input *struct {
		UUID    string `path:"uuid" format:"uuid"`
		Channel int    `query:"channel"`
	}) (*GetDeviceCoaxialCapsOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetDeviceCoaxialCaps(ctx, client.RPC, input.Channel)
		if err != nil {
			return nil, err
		}

		return &GetDeviceCoaxialCapsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device coaxial status",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/coaxial/status",
	}, func(ctx context.Context, input *struct {
		UUID    string `path:"uuid" format:"uuid"`
		Channel int    `query:"channel"`
	}) (*GetDeviceCoaxialStatusOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetCoaxialStatus(ctx, client.RPC, input.Channel)
		if err != nil {
			return nil, err
		}

		return &GetDeviceCoaxialStatusOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device detail",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/detail",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*GetDeviceDetailOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetDetail(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &GetDeviceDetailOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device licenses",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/licenses",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*GetDeviceLicensesOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListLicenses(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &GetDeviceLicensesOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device ptz presets",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/ptz/presets",
	}, func(ctx context.Context, input *struct {
		UUID    string `path:"uuid" format:"uuid"`
		Channel int    `query:"channel"`
	}) (*ListDevicePTZPresetsOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListPTZPresets(ctx, client.PTZ, input.Channel)
		if err != nil {
			return nil, err
		}

		return &ListDevicePTZPresetsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device software versions",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/software",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*GetDeviceSoftwareOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetSoftwareVersion(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &GetDeviceSoftwareOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device storage",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/storage",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*ListDeviceStorageOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListStorage(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &ListDeviceStorageOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device users",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/users",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*ListDeviceUsersOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		position, err := dahua.GetDevicePosition(ctx, db, device.ID)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListUsers(ctx, client.RPC, position.Location.Location)
		if err != nil {
			return nil, err
		}

		return &ListDeviceUsersOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device active users",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/users/active",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*ListDeviceActiveUsersOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		position, err := dahua.GetDevicePosition(ctx, db, device.ID)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListActiveUsers(ctx, client.RPC, position.Location.Location)
		if err != nil {
			return nil, err
		}

		return &ListDeviceActiveUsersOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device groups",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/groups",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*ListDeviceGroupsOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListGroups(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &ListDeviceGroupsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device uptime",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/uptime",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*GetDeviceUptimeOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetUptime(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &GetDeviceUptimeOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device status",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/status",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*GetDeviceStatusOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body := dahua.GetStatus(ctx, client.RPC)

		return &GetDeviceStatusOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device snapshot",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/snapshot",
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Current snapshot of camera",
				Content:     map[string]*huma.MediaType{"image/jpeg": {}},
			},
		},
	}, func(ctx context.Context, input *struct {
		UUID    string `path:"uuid" format:"uuid"`
		Channel int    `query:"channel"`
		Type    int    `query:"type"`
	}) (*huma.StreamResponse, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		return &huma.StreamResponse{
			Body: func(ctx huma.Context) {
				snapshot, err := dahuacgi.SnapshotGet(ctx.Context(), client.CGI, input.Channel, input.Type)
				if err != nil {
					return
				}
				defer snapshot.Close()

				ctx.SetHeader("Content-Type", snapshot.ContentType)
				ctx.SetHeader("Content-Length", snapshot.ContentLength)

				io.Copy(ctx.BodyWriter(), snapshot)
			},
		}, nil
	})
	{
		eventHub := bus.NewHub[bus.Event]().Register()

		sse.Register(api, huma.Operation{
			Summary: "Listen for events from devices",
			Method:  http.MethodGet,
			Path:    "/api/events",
		}, map[string]any{
			"message": DeviceEventsOutput{},
		}, func(ctx context.Context, input *struct {
			DeviceUUIDs []string `query:"device-uuids"`
			Codes       []string `query:"codes"`
		}, send sse.Sender) {
			eventC, unsubscribeEventC := eventHub.Subscribe(ctx)
			defer unsubscribeEventC()

			for event := range eventC {
				if len(input.DeviceUUIDs) != 0 && !slices.Contains(input.DeviceUUIDs, event.DeviceKey.UUID) {
					continue
				}
				if len(input.Codes) != 0 && !slices.Contains(input.Codes, event.Event.Code) {
					continue
				}
				send.Data(DeviceEventsOutput{
					ID:         event.EventID,
					DeviceUUID: event.DeviceKey.UUID,
					Code:       event.Event.Code,
					Action:     event.Event.Action,
					Index:      int64(event.Event.Index),
					Data:       event.Event.Data,
					CreatedAt:  event.CreatedAt,
				})
			}
		})
	}
	huma.Register(api, huma.Operation{
		Summary: "Download device file",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/file",
		Responses: map[string]*huma.Response{
			"200": {
				Description: "File from camera",
				Content: map[string]*huma.MediaType{
					"image/jpeg":               {},
					"application/octet-stream": {},
				},
			},
		},
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
		// TODO: this should be path param wildcard but OpenAPI is stupid
		Name string `query:"name" required:"true"`
	}) (*huma.StreamResponse, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		return &huma.StreamResponse{
			Body: func(ctx huma.Context) {
				rd, err := client.File.Do(ctx.Context(), dahuarpc.LoadFileURL(client.Conn.URL, input.Name), dahuarpc.Cookie(client.RPC.Session(ctx.Context())))
				if err != nil {
					return
				}
				defer rd.Close()

				ctx.SetHeader("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filepath.Base(input.Name)))

				io.Copy(ctx.BodyWriter(), rd)
			},
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device VideoInMode",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/video-in-mode",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*DeviceVideoInModeOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetVideoInMode(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &DeviceVideoInModeOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Sync device VideoInMode",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}/video-in-mode/sync",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
		Body DeviceVideoInModeSchedule
	}) (*DeviceVideoInModeOutput, error) {
		device, err := useDevice(ctx, db, input.UUID)
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, dahuaStore, device)
		if err != nil {
			return nil, err
		}

		position, err := dahua.GetDevicePosition(ctx, db, device.ID)
		if err != nil {
			return nil, err
		}

		body, err := dahua.SyncVideoInMode(ctx, client.RPC, dahua.SyncVideoInModeArgs{
			Location:      core.Optional(input.Body.Location, position.Location).Location,
			Latitude:      core.Optional(input.Body.Latitude, position.Latitude),
			Longitude:     core.Optional(input.Body.Longitude, position.Longitude),
			SunriseOffset: core.Optional(input.Body.SunriseOffset, position.Sunrise_Offset).Duration,
			SunsetOffset:  core.Optional(input.Body.SunsetOffset, position.Sunset_Offset).Duration,
		})
		if err != nil {
			return nil, err
		}

		return &DeviceVideoInModeOutput{
			Body: body,
		}, nil
	})
	{
		setDeviceCoaxialState := func(ctx context.Context, uuid string, channel int, typE coaxialcontrolio.Type, action string) error {
			device, err := useDevice(ctx, db, uuid)
			if err != nil {
				return err
			}

			client, err := useClient(ctx, dahuaStore, device)
			if err != nil {
				return err
			}

			control := coaxialcontrolio.ControlRequest{
				Type:        typE,
				IO:          coaxialcontrolio.Off,
				TriggerMode: coaxialcontrolio.TriggerModeManual,
			}
			switch action {
			case "on":
				control.IO = coaxialcontrolio.On
			case "off":
				control.IO = coaxialcontrolio.Off
			default:
				status, err := dahua.GetCoaxialStatus(ctx, client.RPC, channel)
				if err != nil {
					return err
				}
				if status.WhiteLight {
					control.IO = coaxialcontrolio.Off
				} else {
					control.IO = coaxialcontrolio.On
				}
			}

			return coaxialcontrolio.Control(ctx, client.RPC, channel, control)
		}
		huma.Register(api, huma.Operation{
			Summary: "Set device white light state",
			Method:  http.MethodPost,
			Path:    "/api/devices/{uuid}/coaxial/white-light",
		}, func(ctx context.Context, input *struct {
			UUID    string `path:"uuid" format:"uuid"`
			Channel int    `query:"channel"`
			Action  string `query:"action" enum:"on,off,toggle"`
		}) (*struct{}, error) {
			return &struct{}{}, setDeviceCoaxialState(ctx, input.UUID, input.Channel, coaxialcontrolio.TypeWhiteLight, input.Action)
		})
		huma.Register(api, huma.Operation{
			Summary: "Set device speaker state",
			Method:  http.MethodPost,
			Path:    "/api/devices/{uuid}/coaxial/speaker",
		}, func(ctx context.Context, input *struct {
			UUID    string `path:"uuid" format:"uuid"`
			Channel int    `query:"channel"`
			Action  string `query:"action" enum:"on,off,toggle"`
		}) (*struct{}, error) {
			return &struct{}{}, setDeviceCoaxialState(ctx, input.UUID, input.Channel, coaxialcontrolio.TypeSpeaker, input.Action)
		})
	}
	huma.Register(api, huma.Operation{
		Summary: "Get setting",
		Method:  http.MethodGet,
		Path:    "/api/settings",
	}, func(ctx context.Context, i *struct{}) (*SettingOutput, error) {
		var settings system.Settings
		err := db.GetContext(ctx, &settings, `
			SELECT * FROM settings
		`)
		if err != nil {
			return nil, err
		}

		return &SettingOutput{
			Body: NewSettings(settings),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Update setting",
		Method:  http.MethodPut,
		Path:    "/api/settings",
	}, func(ctx context.Context, input *struct {
		Body UpdateSettingsInput
	}) (*SettingOutput, error) {
		settings, err := system.UpdateSettings(ctx, db, system.UpdateSettingsArgs{
			Location:        input.Body.Location,
			Latitude:        input.Body.Latitude,
			Longitude:       input.Body.Longitude,
			SunriseOffset:   input.Body.SunriseOffset,
			SunsetOffset:    input.Body.SunsetOffset,
			SyncVideoInMode: input.Body.SyncVideoInMode,
		})
		if err != nil {
			return nil, err
		}

		return &SettingOutput{
			Body: NewSettings(settings),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Default setting",
		Method:  http.MethodDelete,
		Path:    "/api/settings",
	}, func(ctx context.Context, i *struct{}) (*SettingOutput, error) {
		settings, err := system.DefaultSettings(ctx, db)
		if err != nil {
			return nil, err
		}

		return &SettingOutput{
			Body: NewSettings(settings),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List endpoints",
		Method:  http.MethodGet,
		Path:    "/api/endpoints",
	}, func(ctx context.Context, i *struct{}) (*ListEndpointsOutput, error) {
		rows, err := db.QueryxContext(ctx, `
				SELECT * FROM endpoints
		`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		body := []Endpoint{}
		for rows.Next() {
			var v system.Endpoint
			if err := rows.StructScan(&v); err != nil {
				return nil, err
			}
			body = append(body, NewEndpoint(v))
		}

		return &ListEndpointsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Create endpoint",
		Method:  http.MethodPost,
		Path:    "/api/endpoints",
	}, func(ctx context.Context, input *struct {
		Body CreateEndpointInput
	}) (*CreateEndpointsOutput, error) {
		endpoint, err := system.CreateEndpoint(ctx, db, system.CreateEndpointArgs{
			GoriseURL: input.Body.GoriseURL,
		})
		if err != nil {
			return nil, err
		}

		return &CreateEndpointsOutput{
			Body: NewEndpoint(endpoint),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Delete endpoint",
		Method:  http.MethodDelete,
		Path:    "/api/endpoints/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUID string `path:"uuid" format:"uuid"`
	}) (*struct{}, error) {
		return &struct{}{}, system.DeleteEndpoints(ctx, db, input.UUID)
	})
}

func NewEndpoint(v system.Endpoint) Endpoint {
	return Endpoint{
		UUID:      v.UUID,
		GoriseURL: v.Gorise_URL,
		CreatedAt: v.Created_At.Time,
		UpdatedAt: v.Updated_At.Time,
	}
}

type Endpoint struct {
	UUID      string
	GoriseURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSettings(v system.Settings) Settings {
	return Settings{
		Location:        v.Location,
		Latitude:        v.Latitude,
		Longitude:       v.Longitude,
		SunriseOffset:   v.Sunrise_Offset,
		SunsetOffset:    v.Sunset_Offset,
		UpdatedAt:       v.Updated_At.Time,
		SyncVideoInMode: v.Sync_Video_In_Mode,
	}
}

type Settings struct {
	Location        types.Location `json:"location"`
	Latitude        float64        `json:"latitude"`
	Longitude       float64        `json:"longitude"`
	SunriseOffset   types.Duration `json:"sunrise_offset"`
	SunsetOffset    types.Duration `json:"sunset_offset"`
	UpdatedAt       time.Time      `json:"updated_at"`
	SyncVideoInMode bool           `json:"sync_video_in_mode"`
}

type UpdateSettingsInput struct {
	Location        types.Location `json:"location"`
	Latitude        float64        `json:"latitude"`
	Longitude       float64        `json:"longitude"`
	SunriseOffset   types.Duration `json:"sunrise_offset"`
	SunsetOffset    types.Duration `json:"sunset_offset"`
	SyncVideoInMode bool           `json:"sync_video_in_mode"`
}

type SettingOutput struct {
	Body Settings
}

type DeviceCoaxialControlInput struct {
	Channel  int                           `json:"int"`
	Controls []DeviceCoaxialControlRequest `json:"controls"`
}

type DeviceCoaxialControlRequest struct {
	Type        coaxialcontrolio.Type        `json:"type"`
	IO          coaxialcontrolio.IO          `json:"io"`
	TriggerMode coaxialcontrolio.TriggerMode `json:"trigger_mode"`
}

type DeviceVideoInModeOutput struct {
	Body dahua.DeviceVideoInMode
}

type ListDevicesOutput struct {
	Body []Device
}

type DeviceVideoInModeSchedule struct {
	Location      *types.Location `json:"location,omitempty"`
	Latitude      *float64        `json:"latitude,omitempty"`
	Longitude     *float64        `json:"longitude,omitempty"`
	SunriseOffset *types.Duration `json:"sunrise_offset,omitempty"`
	SunsetOffset  *types.Duration `json:"sunset_offset,omitempty"`
}

type Device struct {
	UUID            string          `json:"uuid"`
	Name            string          `json:"name"`
	IP              net.IP          `json:"ip"`
	Username        string          `json:"username"`
	Location        *types.Location `json:"location"`
	Features        []dahua.Feature `json:"features"`
	Email           string          `json:"email"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Latitude        *float64        `json:"latitude"`
	Longitude       *float64        `json:"longitude"`
	SunriseOffset   *types.Duration `json:"sunrise_offset"`
	SunsetOffset    *types.Duration `json:"sunset_offset"`
	SyncVideoInMode *bool           `json:"sync_video_in_mode"`
}

func NewDevice(v dahua.DahuaDevice) Device {
	ip := net.ParseIP(v.IP)
	return Device{
		UUID:            v.UUID,
		Name:            v.Name,
		IP:              ip,
		Username:        v.Username,
		Location:        core.SQLNullToNull(v.Location),
		Features:        v.Features.V,
		Email:           v.Email.String,
		CreatedAt:       v.Created_At.Time,
		UpdatedAt:       v.Updated_At.Time,
		Latitude:        core.SQLNullToNull(v.Latitude),
		Longitude:       core.SQLNullToNull(v.Longitude),
		SunriseOffset:   core.SQLNullToNull(v.Sunrise_Offset),
		SunsetOffset:    core.SQLNullToNull(v.Sunset_Offset),
		SyncVideoInMode: core.SQLNullToNull(v.Sync_Video_In_Mode),
	}
}

type CreateDevice struct {
	Name            string          `json:"name,omitempty"`
	IP              string          `json:"ip,omitempty" format:"ipv4"`
	Username        string          `json:"username,omitempty"`
	Password        string          `json:"password,omitempty"`
	Location        *types.Location `json:"location,omitempty"`
	Features        []dahua.Feature `json:"features,omitempty"`
	Email           *string         `json:"email,omitempty"`
	Latitude        *float64        `json:"latitude,omitempty"`
	Longitude       *float64        `json:"longitude,omitempty"`
	SunriseOffset   *types.Duration `json:"sunrise_offset,omitempty"`
	SunsetOffset    *types.Duration `json:"sunset_offset,omitempty"`
	SyncVideoInMode *bool           `json:"sync_video_in_mode,omitempty"`
}

type UpdateDevice struct {
	Name            string          `json:"name"`
	IP              string          `json:"ip" format:"ipv4"`
	Username        string          `json:"username"`
	Password        *string         `json:"password,omitempty"`
	Location        *types.Location `json:"location,omitempty"`
	Features        []dahua.Feature `json:"features,omitempty"`
	Email           *string         `json:"email,omitempty"`
	Latitude        *float64        `json:"latitude"`
	Longitude       *float64        `json:"longitude"`
	SunriseOffset   *types.Duration `json:"sunrise_offset,omitempty"`
	SunsetOffset    *types.Duration `json:"sunset_offset,omitempty"`
	SyncVideoInMode *bool           `json:"sync_video_in_mode"`
}

type CreateDeviceInput struct {
	Body CreateDevice
}

func (i *CreateDeviceInput) Resolve(ctx huma.Context) []error {
	if i.Body.Name == "" {
		i.Body.Name = i.Body.IP
	}
	return nil
}

type CreateDevicesOutput struct {
	Body Device
}

type GetDeviceOutput struct {
	Body Device
}

type UpdateDevicesOutput struct {
	Body Device
}

type PatchDevicesOutput struct {
	Body Device
}

type DeleteDeviceOutput struct {
	UUID string `format:"uuid"`
}

type GetDeviceCoaxialCapsOutput struct {
	Body dahua.DeviceCoaxialCaps
}

type GetDeviceCoaxialStatusOutput struct {
	Body dahua.DeviceCoaxialStatus
}

type GetDeviceDetailOutput struct {
	Body dahua.DeviceDetail
}

type GetDeviceLicensesOutput struct {
	Body []dahua.DeviceLicense
}

type ListDevicePTZPresetsOutput struct {
	Body []dahua.DevicePTZPreset
}

type GetDeviceSoftwareOutput struct {
	Body dahua.DeviceSoftwareVersion
}

type ListDeviceStorageOutput struct {
	Body []dahua.DeviceStorage
}

type ListDeviceUsersOutput struct {
	Body []dahua.DeviceUser
}

type ListDeviceActiveUsersOutput struct {
	Body []dahua.DeviceActiveUser
}

type ListDeviceGroupsOutput struct {
	Body []dahua.DeviceGroup
}

type GetDeviceUptimeOutput struct {
	Body dahua.DeviceUptime
}

type GetDeviceStatusOutput struct {
	Body dahua.DeviceStatus
}

type GetDeviceSnapshotOutput struct {
	Body []byte
}

type DeviceEventsOutput struct {
	ID         string          `json:"id"`
	DeviceUUID string          `json:"device_uuid"`
	Code       string          `json:"code"`
	Action     string          `json:"action"`
	Index      int64           `json:"index"`
	Data       json.RawMessage `json:"data"`
	CreatedAt  time.Time       `json:"created_at"`
}

type CreateEndpointInput struct {
	GoriseURL string `json:"gorise_url"`
}

type ListEndpointsOutput struct {
	Body []Endpoint
}

type CreateEndpointsOutput struct {
	Body Endpoint
}
