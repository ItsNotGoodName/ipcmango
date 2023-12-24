package dahuarpc

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AuthParam struct {
	Encryption string `json:"encryption"`
	Random     string `json:"random"`
	Realm      string `json:"realm"`
}

// HashPassword runs the hashing algorithm for the password.
func (a AuthParam) HashPassword(username, password string) string {
	switch a.Encryption {
	case "Basic":
		return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	case "Default":
		return strings.ToUpper(fmt.Sprintf("%x",
			md5.Sum([]byte(fmt.Sprintf(
				"%s:%s:%s",
				username,
				a.Random,
				strings.ToUpper(fmt.Sprintf(
					"%x",
					md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", username, a.Realm, password))))))))))
	default:
		return password
	}
}

type Timestamp string

// NewTimestamp converts the given UTC time to the given location and returns the timestamp.
func NewTimestamp(date time.Time, cameraLocation *time.Location) Timestamp {
	return Timestamp(date.In(cameraLocation).Format("2006-01-02 15:04:05"))
}

// Parse returns the UTC time for the given timestamp and camera location.
func (t Timestamp) Parse(cameraLocation *time.Location) (time.Time, error) {
	if strings.HasSuffix(string(t), "PM") || strings.HasSuffix(string(t), "AM") {
		date, err := time.ParseInLocation("2006-01-02 03:04:05 PM", string(t), cameraLocation)
		if err != nil {
			return date, err
		}

		return date.UTC(), nil
	} else {
		date, err := time.ParseInLocation("2006-01-02 15:04:05", string(t), cameraLocation)
		if err != nil {
			return date, err
		}

		return date.UTC(), nil
	}
}

// ExtractFilePathTags extracts tags that are surrounded by brackets from the given file path.
func ExtractFilePathTags(filePath string) []string {
	search := filePath
	idx := strings.LastIndex(filePath, "/")
	if idx != -1 {
		search = filePath[idx:]
	}

	var tags []string
	tokens := strings.Split(search, "[")
	for i := 1; i < len(tokens); i++ {
		if end := strings.Index(tokens[i], "]"); end != -1 {
			tags = append(tags, tokens[i][:end])
		}
	}

	return tags
}

// Integer is for types that are supposed to integer but for some reason the camera returns a float.
type Integer int64

func (s *Integer) UnmarshalJSON(data []byte) error {
	var number float64
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	*s = Integer(number)

	return nil
}

func (s Integer) Integer() int64 {
	return int64(s)
}

func URL(httpAddress string) string {
	return fmt.Sprintf("%s/RPC2", httpAddress)
}

func LoginURL(httpAddress string) string {
	return fmt.Sprintf("%s/RPC2_Login", httpAddress)
}

func LoadFileURL(httpAddress, path string) string {
	return fmt.Sprintf("%s/RPC_Loadfile%s", httpAddress, path)
}

func Cookie(session string) string {
	return fmt.Sprintf("WebClientSessionID=%s; DWebClientSessionID=%s; DhWebClientSessionID=%s", session, session, session)
}

// NewTimeSection (e.g. "1 08:01:45-16:16:22")
func NewTimeSection(s string) (TimeSection, error) {
	splitBySpace := strings.Split(s, " ")
	if len(splitBySpace) != 2 {
		return TimeSection{}, fmt.Errorf("invalid number of spaces: %d", len(splitBySpace))
	}

	splitByDash := strings.Split(splitBySpace[1], "-")
	if len(splitByDash) != 2 {
		return TimeSection{}, fmt.Errorf("invalid number of dashes: %d", len(splitByDash))
	}

	start, err := durationFromTimeString(splitByDash[0])
	if err != nil {
		return TimeSection{}, err
	}

	end, err := durationFromTimeString(splitByDash[1])
	if err != nil {
		return TimeSection{}, err
	}

	return TimeSection{
		Enable: splitBySpace[0] == "1",
		Start:  start,
		End:    end,
	}, nil
}

// durationFromTimeString (e.g. "08:01:45")
func durationFromTimeString(s string) (time.Duration, error) {
	arr := strings.Split(s, ":")
	if len(arr) != 3 {
		return 0, fmt.Errorf("invalid number of colons: %d", len(arr))
	}

	var numbers [3]int
	for i := range arr {
		var err error
		numbers[i], err = strconv.Atoi(arr[i])
		if err != nil {
			return 0, err
		}
	}

	return time.Duration(numbers[0])*time.Hour + time.Duration(numbers[1])*time.Minute + time.Duration(numbers[2])*time.Second, nil
}

type TimeSection struct {
	Enable bool
	Start  time.Duration
	End    time.Duration
}

func (s *TimeSection) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	res, err := NewTimeSection(str)
	if err != nil {
		return err
	}

	s.Enable = res.Enable
	s.Start = res.Start
	s.End = res.End

	return nil
}

func (s TimeSection) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s TimeSection) String() string {
	var enable int
	if s.Enable {
		enable = 1
	}

	return fmt.Sprintf(
		"%d %02d:%02d:%02d-%02d:%02d:%02d",
		enable,
		int(s.Start.Hours()),
		int(s.Start.Minutes())%60,
		int(s.Start.Seconds())%60,
		int(s.End.Hours()),
		int(s.End.Minutes())%60,
		int(s.End.Seconds())%60,
	)
}
