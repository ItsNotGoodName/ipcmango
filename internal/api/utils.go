package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
)

// ---------- Stream

func useStream(c echo.Context) *json.Encoder {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response())
}

type StreamPayload struct {
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	OK      bool    `json:"ok"`
}

func sendStreamError(c echo.Context, enc *json.Encoder, err error) error {
	str := err.Error()
	if encodeErr := enc.Encode(StreamPayload{
		OK:      false,
		Message: &str,
	}); encodeErr != nil {
		return errors.Join(encodeErr, err)
	}

	c.Response().Flush()

	return err
}

func sendStream(c echo.Context, enc *json.Encoder, data any) error {
	err := enc.Encode(StreamPayload{
		OK:   true,
		Data: data,
	})
	if err != nil {
		return sendStreamError(c, enc, err)
	}

	c.Response().Flush()

	return nil
}

// ---------- Queries

func queryInts(c echo.Context, key string) ([]int64, error) {
	ids := make([]int64, 0)
	idsStr := c.QueryParams()[key]
	for _, v := range idsStr {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, echo.ErrBadRequest.WithInternal(err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func queryIntOptional(c echo.Context, key string) (int, error) {
	str := c.QueryParam(key)
	if str == "" {
		return 0, nil
	}

	number, err := strconv.Atoi(str)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}

func queryBoolOptional(c echo.Context, key string) (bool, error) {
	str := c.QueryParam(key)
	if str == "" {
		return false, nil
	}

	bool, err := strconv.ParseBool(str)
	if err != nil {
		return false, echo.ErrBadRequest.WithInternal(err)
	}

	return bool, nil
}

func queryDahuaScanRange(startStr, endStr string) (models.TimeRange, error) {
	end := time.Now()
	start := end.Add(-dahuacore.MaxScanPeriod)
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
		}
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
		}
	}

	res, err := core.NewTimeRange(start, end)
	if err != nil {
		return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
	}

	return res, nil
}

var encoder = schema.NewEncoder()

var decoder = schema.NewDecoder()

func ParseForm(c echo.Context, form any) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	if err := decoder.Decode(form, c.Request().PostForm); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	return nil
}

func ParseLocation(location string) (*time.Location, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return nil, echo.ErrBadRequest.WithInternal(err)
	}
	return loc, nil
}

func DecodeQuery(c echo.Context, dst any) error {
	if err := decoder.Decode(dst, c.Request().URL.Query()); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	return nil
}

func EncodeQuery(src any) url.Values {
	query := make(url.Values)
	err := encoder.Encode(src, query)
	if err != nil {
		panic(err)
	}
	return query
}

func ValidateStruct(src any) error {
	err := validate.Validate.Struct(src)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	return err
}

func FormatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}

func PathID(c echo.Context) (int64, error) {
	number, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}

func QueryInt(c echo.Context, key string) (int64, error) {
	str := c.QueryParam(key)
	if str == "" {
		return 0, echo.ErrBadRequest
	}

	number, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}

func UseTimeRange(start, end string) (models.TimeRange, error) {
	var startTime, endTime time.Time
	if start != "" {
		var err error
		startTime, err = time.ParseInLocation("2006-01-02T15:04", start, time.Local)
		if err != nil {
			return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
		}
	}

	if end != "" {
		var err error
		endTime, err = time.ParseInLocation("2006-01-02T15:04", end, time.Local)
		if err != nil {
			return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
		}
	} else if start != "" {
		endTime = time.Now()
	}

	r, err := core.NewTimeRange(startTime, endTime)
	if err != nil {
		return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
	}

	return r, nil
}
