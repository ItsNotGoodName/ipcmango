package storage

import (
	"context"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

func GetDeviceAllInfo(ctx context.Context, c dahuarpc.Client) ([]Storage, error) {
	var object int64
	{
		rpc, err := c.RPC(ctx)
		if err != nil {
			return []Storage{}, err
		}

		res, err := dahuarpc.Send[any](ctx, rpc.Method("storage.factory.instance"))
		if err != nil {
			return []Storage{}, err
		}

		object = res.Result.Integer()
	}

	rpc, err := c.RPC(ctx)
	if err != nil {
		return []Storage{}, err
	}

	res, err := dahuarpc.Send[getDeviceAllInfoResult](ctx, rpc.Method("storage.getDeviceAllInfo").Object(object))
	if err != nil {
		return []Storage{}, err
	}

	return res.Params.Info, err
}

type Storage struct {
	Name   string          `json:"Name"`
	State  string          `json:"State"`
	Detail []StorageDetail `json:"Detail"`
}

type StorageDetail struct {
	Path       string           `json:"Path"`
	Type       string           `json:"Type"`
	TotalBytes dahuarpc.Integer `json:"TotalBytes"`
	UsedBytes  dahuarpc.Integer `json:"UsedBytes"`
	IsError    bool             `json:"IsError"`
}

type getDeviceAllInfoResult struct {
	Info []Storage
}

func (g *getDeviceAllInfoResult) UnmarshalJSON(data []byte) error {
	{
		res := struct {
			Info []Storage `json:"info"`
		}{}

		if err := json.Unmarshal(data, &res); err == nil {
			g.Info = res.Info
			return nil
		}
	}

	var storages []Storage
	if err := json.Unmarshal(data, &storages); err != nil {
		return err
	}

	g.Info = storages

	return nil
}
