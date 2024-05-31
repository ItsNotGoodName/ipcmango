package dahuarpc

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
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

func NewTimestamp(date time.Time, deviceLocation *time.Location) Timestamp {
	if date.IsZero() {
		return ""
	}
	return Timestamp(date.In(deviceLocation).Format("2006-01-02 15:04:05"))
}

func (t Timestamp) Parse(deviceLocation *time.Location) (time.Time, error) {
	if t == "" {
		return time.Time{}.UTC(), nil
	}

	format := "2006-01-02 15:04:05"
	if strings.HasSuffix(string(t), "PM") || strings.HasSuffix(string(t), "AM") {
		format = "2006-01-02 03:04:05 PM"
	}

	date, err := time.ParseInLocation(format, string(t), deviceLocation)
	if err != nil {
		return date, err
	}

	return date.UTC(), nil
}

// ExtractFilePathTags extracts tags from the file path.
// Tags are strings surrounded by brackets
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

// Integer is for types that are supposed to be integers but for some reason the device returns a float.
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

// URL is the HTTP RPC API URL.
func URL(u *url.URL) string {
	return fmt.Sprintf("%s://%s/RPC2", u.Scheme, u.Hostname())
}

// LoginURL is the HTTP RPC API URL for login.
func LoginURL(u *url.URL) string {
	return fmt.Sprintf("%s://%s/RPC2_Login", u.Scheme, u.Hostname())
}

// LoadFileURL is the HTTP URL for accessing files.
// The file path must be absolute.
func LoadFileURL(u *url.URL, filePath string) string {
	return fmt.Sprintf("%s://%s/RPC_Loadfile%s", u.Scheme, u.Hostname(), filePath)
}

// Cookie creates a session cookie.
func Cookie(session string) string {
	return fmt.Sprintf("WebClientSessionID=%s; DWebClientSessionID=%s; DhWebClientSessionID=%s", session, session, session)
}

func DefaultTimeSection() TimeSection {
	return TimeSection{
		Number: 0,
		Start:  0,
		End:    24 * time.Hour,
	}
}

// NewTimeSectionFromString (e.g. "1 08:01:45-16:16:22").
func NewTimeSectionFromString(s string) (TimeSection, error) {
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

	number, err := strconv.Atoi(splitBySpace[0])
	if err != nil {
		return TimeSection{}, err
	}

	return TimeSection{
		Number: number,
		Start:  start,
		End:    end,
	}, nil
}

func NewTimeSectionFromRange(number int, start, end time.Time) TimeSection {
	return TimeSection{
		Number: number,
		Start:  time.Duration(start.Hour())*time.Hour + time.Duration(start.Minute())*time.Minute + time.Duration(start.Second())*time.Second,
		End:    time.Duration(end.Hour())*time.Hour + time.Duration(end.Minute())*time.Minute + time.Duration(end.Second())*time.Second,
	}
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
	Number int
	Start  time.Duration
	End    time.Duration
}

func (s *TimeSection) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	res, err := NewTimeSectionFromString(str)
	if err != nil {
		return err
	}

	s.Number = res.Number
	s.Start = res.Start
	s.End = res.End

	return nil
}

func (s TimeSection) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s TimeSection) String() string {
	return fmt.Sprintf(
		"%d %02d:%02d:%02d-%02d:%02d:%02d",
		s.Number,
		int(s.Start.Hours()),
		int(s.Start.Minutes())%60,
		int(s.Start.Seconds())%60,
		int(s.End.Hours()),
		int(s.End.Minutes())%60,
		int(s.End.Seconds())%60,
	)
}
