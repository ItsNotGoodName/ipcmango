package models

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager/config"
)

const DahuaFileTypeJPG = "jpg"
const DahuaFileTypeDAV = "dav"

type DahuaEventWorkerState string

const (
	DahuaEventWorkerStateConnecting   DahuaEventWorkerState = "connecting"
	DahuaEventWorkerStateConnected    DahuaEventWorkerState = "connected"
	DahuaEventWorkerStateDisconnected DahuaEventWorkerState = "disconnected"
)

type DahuaError struct {
	Error string `json:"error"`
}

type DahuaStatus struct {
	DeviceID     int64     `json:"device_id"`
	Address      string    `json:"address"`
	Username     string    `json:"username"`
	Location     string    `json:"location"`
	Seed         int       `json:"seed"`
	RPCError     string    `json:"rpc_error"`
	RPCState     string    `json:"rpc_state"`
	RPCLastLogin time.Time `json:"rpc_last_login"`
}

type DahuaDevice struct {
	ID        int64
	Name      string
	Address   *url.URL
	Username  string
	Password  string
	Location  *time.Location
	Feature   DahuaFeature
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DahuaConn struct {
	ID       int64
	Address  *url.URL
	Username string
	Password string
	Location *time.Location
	Feature  DahuaFeature
	Seed     int
}

type DahuaDeviceConn struct {
	DahuaDevice
	DahuaConn
}

type DahuaFeature int

func (f DahuaFeature) EQ(feature DahuaFeature) bool {
	return feature != 0 && f&feature == feature
}

const (
	// DahuaFeatureCamera means the device has at least 1 camera.
	DahuaFeatureCamera DahuaFeature = 1 << iota
)

type DahuaDetail struct {
	DeviceID         int64  `json:"device_id"`
	SN               string `json:"sn"`
	DeviceClass      string `json:"device_class"`
	DeviceType       string `json:"device_type"`
	HardwareVersion  string `json:"hardware_version"`
	MarketArea       string `json:"market_area"`
	ProcessInfo      string `json:"process_info"`
	Vendor           string `json:"vendor"`
	OnvifVersion     string `json:"onvif_version"`
	AlgorithmVersion string `json:"algorithm_version"`
}

type DahuaSoftwareVersion struct {
	DeviceID                int64  `json:"device_id"`
	Build                   string `json:"build"`
	BuildDate               string `json:"build_date"`
	SecurityBaseLineVersion string `json:"security_base_line_version"`
	Version                 string `json:"version"`
	WebVersion              string `json:"web_version"`
}

type DahuaLicense struct {
	DeviceID      int64     `json:"device_id"`
	AbroadInfo    string    `json:"abroad_info"`
	AllType       bool      `json:"all_type"`
	DigitChannel  int       `json:"digit_channel"`
	EffectiveDays int       `json:"effective_days"`
	EffectiveTime time.Time `json:"effective_time"`
	LicenseID     int       `json:"license_id"`
	ProductType   string    `json:"product_type"`
	Status        int       `json:"status"`
	Username      string    `json:"username"`
}

type DahuaCoaxialStatus struct {
	DeviceID   int64 `json:"device_id"`
	WhiteLight bool  `json:"white_light"`
	Speaker    bool  `json:"speaker"`
}

type DahuaCoaxialCaps struct {
	DeviceID                     int64 `json:"device_id"`
	SupportControlFullcolorLight bool  `json:"support_control_fullcolor_light"`
	SupportControlLight          bool  `json:"support_control_light"`
	SupportControlSpeaker        bool  `json:"support_control_speaker"`
}

type DahuaFile struct {
	ID          int64     `json:"id"`
	DeviceID    int64     `json:"device_id"`
	Channel     int       `json:"channel"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Length      int       `json:"length"`
	Type        string    `json:"type"`
	FilePath    string    `json:"file_path"`
	Duration    int       `json:"duration"`
	Disk        int       `json:"disk"`
	VideoStream string    `json:"video_stream"`
	Flags       []string  `json:"flags"`
	Events      []string  `json:"events"`
	Cluster     int       `json:"cluster"`
	Partition   int       `json:"partition"`
	PicIndex    int       `json:"pic_index"`
	Repeat      int       `json:"repeat"`
	WorkDir     string    `json:"work_dir"`
	WorkDirSN   bool      `json:"work_dir_sn"`
	Local       bool      `json:"local"`
}

type DahuaEvent struct {
	ID        int64           `json:"id"`
	DeviceID  int64           `json:"device_id"`
	Code      string          `json:"code"`
	Action    string          `json:"action"`
	Index     int             `json:"index"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
}

type DahuaEventRule struct {
	IgnoreDB   bool
	IgnoreLive bool
	IgnoreMQTT bool
}

type DahuaStorage struct {
	DeviceID   int64  `json:"device_id"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Path       string `json:"path"`
	Type       string `json:"type"`
	TotalBytes int64  `json:"total_bytes"`
	UsedBytes  int64  `json:"used_bytes"`
	IsError    bool   `json:"is_error"`
}

type DahuaUser struct {
	DeviceID      int64     `json:"device_id"`
	ClientAddress string    `json:"client_address"`
	ClientType    string    `json:"client_type"`
	Group         string    `json:"group"`
	ID            int       `json:"id"`
	LoginTime     time.Time `json:"login_time"`
	Name          string    `json:"name"`
}

type DahuaSunriseSunset struct {
	SwitchMode  config.SwitchMode    `json:"switch_mode"`
	TimeSection dahuarpc.TimeSection `json:"time_section"`
}

type DahuaStream struct {
	ID           int64
	DeviceID     int64
	Name         string
	Channel      int
	Subtype      int
	MediamtxPath string
	EmbedURL     string
}

type DahuaScanType string

var (
	DahuaScanTypeUnkown  DahuaScanType = ""
	DahuaScanTypeFull    DahuaScanType = "full"
	DahuaScanTypeQuick   DahuaScanType = "quick"
	DahuaScanTypeReverse DahuaScanType = "reverse"
)
