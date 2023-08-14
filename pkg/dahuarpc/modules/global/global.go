package global

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

func FirstLogin(ctx context.Context, c dahuarpc.ClientLogin, username string) (dahuarpc.Response[dahuarpc.AuthParam], error) {
	return dahuarpc.SendRaw[dahuarpc.AuthParam](ctx, c.
		RPCLogin().
		Method("global.login").
		Params(dahuarpc.JSON{
			"userName":   username,
			"password":   "",
			"loginType":  "Direct",
			"clientType": "Web3.0",
		}))
}

func SecondLogin(ctx context.Context, c dahuarpc.ClientLogin, username, password, loginType, authorityType string) error {
	_, err := dahuarpc.Send[any](ctx, c.
		RPCLogin().
		Method("global.login").
		Params(dahuarpc.JSON{
			"userName":      username,
			"password":      password,
			"clientType":    "Web3.0",
			"loginType":     loginType,
			"authorityType": authorityType,
		}))

	return err
}

func GetCurrentTime(ctx context.Context, c dahuarpc.Client) (string, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahuarpc.Send[struct {
		Time string `json:"time"`
	}](ctx, rpc.Method("global.getCurrentTime"))

	return res.Params.Time, err
}

func KeepAlive(ctx context.Context, c dahuarpc.Client) (int, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return 0, err
	}

	res, err := dahuarpc.Send[struct {
		Timeout int `json:"timeout"`
	}](ctx, rpc.Method("global.keepAlive"))

	return res.Params.Timeout, err
}

func Logout(ctx context.Context, c dahuarpc.Client) (bool, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return false, err
	}

	res, err := dahuarpc.Send[bool](ctx, rpc.Method("global.logout"))

	return res.Params, err
}
