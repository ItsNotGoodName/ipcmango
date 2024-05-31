package peripheralchip

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Type string

var TypeBLOB Type = "BLOB"

type Cap struct {
	ChipType string `json:"string"`
	Count    int    `json:"number"`
	Model    string `json:"model"`
	Type     Type   `json:"type"`
}

func GetCaps(ctx context.Context, c dahuarpc.Conn) ([]Cap, error) {
	res, err := dahuarpc.Send[struct {
		Caps []Cap `json:"caps"`
	}](ctx, c, dahuarpc.New("PeripheralChip.getCaps"))
	if err != nil {
		return nil, err
	}

	return res.Params.Caps, nil
}

type Version struct {
	HardwareVersion string `json:"hardwareVersion"`
	IsMatching      int    `json:"isMatching"`
	SoftwareVersion string `json:"softwareVersion"`
}

func GetVersion(ctx context.Context, c dahuarpc.Conn, typ Type) ([]Version, error) {
	res, err := dahuarpc.Send[struct {
		Versions []Version `json:"versions"`
	}](ctx, c, dahuarpc.New("PeripheralChip.getVersion").
		Params(struct {
			Type Type `json:"type"`
		}{
			Type: typ,
		}))
	if err != nil {
		return nil, err
	}

	return res.Params.Versions, nil
}
