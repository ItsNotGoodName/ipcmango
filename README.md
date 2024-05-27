# IPCManView

**🚧 WORK IN PROGRESS, EACH COMMIT MOST LIKELY BREAKS THE LAST 🚧**

Application for managing and viewing Dahua devices.

TODO: explain why this program exists

https://github.com/ItsNotGoodName/ipcmanview/assets/35015993/5c9a5482-a031-49ae-a56d-5242bd471505

# Features

- Single binary
- View device information (e.g. software version, license, storage, …)
- Subscribe to device events
- View snapshot of cameras
- Receive and view emails from devices
- Send emails to other messaging endpoints (e.g Telegram, ntfy, ...)
- Sync VideoInMode with sunrise and sunset

# Usage

```
ipcmanview
```

## Configuration

| Environment Variable    | Default           | Description                                            |
| ----------------------- | ----------------- | ------------------------------------------------------ |
| `SERVICE_DIR`           | "ipcmanview_data" | Directory path for storing data.                       |
| `SERVICE_HOST`          |                   | Host to listen on (e.g. "127.0.0.1").                  |
| `SERVICE_PORT`          | 8080              | Port to listen on.                                     |
| `SERVICE_SMTP_HOST`     |                   | SMTP host to listen on (e.g. "127.0.0.1").             |
| `SERVICE_SMTP_PORT`     | 1025              | SMTP port to listen on.                                |
| `SERVICE_MQTT_ADDRESS`  |                   | MQTT server address (e.g. "mqtt://192.168.1.20:1883"). |
| `SERVICE_MQTT_TOPIC`    | "ipcmanview"      | MQTT server topic to publish messages.                 |
| `SERVICE_MQTT_USERNAME` |                   | MQTT server username for authentication.               |
| `SERVICE_MQTT_PASSWORD` |                   | MQTT server password for authentication.               |
| `SERVICE_MQTT_HA`       | false             | Enable Home Assistant MQTT discovery.                  |
| `SERVICE_MQTT_HA_TOPIC` | "homeassistant"   | Home Assistant MQTT discover topic.                    |

# Roadmap

Roadmap is in order of importance.

- View files on devices
- Support editing the config `General` (System > General)
- Support editing the config `Email` (Network > SMTP(Email))
- Support editing the config `VideoAnalyseRule`
- View DAV files in local storage via RTSP (see 4.1.3 in the Dahua HTTP API PDF)
- Create and cache thumbnails for files
- Act as a HomeKit bridge for viewing cameras
- Support two-way talk on cameras that support it (see [./pkg/dahuacgi/audio.go](./pkg/dahuacgi/audio.go))
