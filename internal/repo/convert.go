package repo

import "github.com/ItsNotGoodName/ipcmanview/internal/models"

// this is stupid

func (c ListDahuaDeviceByIDsRow) Convert() models.DahuaDeviceConn {
	return models.DahuaDeviceConn{
		DahuaDevice: models.DahuaDevice{
			ID:        c.ID,
			Address:   c.Address,
			Username:  c.Username,
			Password:  c.Password,
			Location:  c.Location.Location,
			Name:      c.Name,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
		},
		DahuaConn: models.DahuaConn{
			ID:       c.ID,
			Address:  c.Address,
			Username: c.Username,
			Password: c.Password,
			Location: c.Location.Location,
			Seed:     int(c.Seed),
		},
	}
}

func (c GetDahuaDeviceRow) Convert() models.DahuaDeviceConn {
	return models.DahuaDeviceConn{
		DahuaDevice: models.DahuaDevice{
			ID:        c.ID,
			Address:   c.Address,
			Username:  c.Username,
			Password:  c.Password,
			Location:  c.Location.Location,
			Name:      c.Name,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
		},
		DahuaConn: models.DahuaConn{
			ID:       c.ID,
			Address:  c.Address,
			Username: c.Username,
			Password: c.Password,
			Location: c.Location.Location,
			Seed:     int(c.Seed),
		},
	}
}

func (c ListDahuaDeviceRow) Convert() models.DahuaDeviceConn {
	return models.DahuaDeviceConn{
		DahuaDevice: models.DahuaDevice{
			ID:        c.ID,
			Address:   c.Address,
			Username:  c.Username,
			Password:  c.Password,
			Location:  c.Location.Location,
			Name:      c.Name,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
		},
		DahuaConn: models.DahuaConn{
			ID:       c.ID,
			Address:  c.Address,
			Username: c.Username,
			Password: c.Password,
			Location: c.Location.Location,
			Seed:     int(c.Seed),
		},
	}
}

func (c DahuaFile) Convert() models.DahuaFile {
	return models.DahuaFile{
		ID:          c.ID,
		DeviceID:    c.DeviceID,
		Channel:     int(c.Channel),
		StartTime:   c.StartTime.Time,
		EndTime:     c.EndTime.Time,
		Length:      int(c.Length),
		Type:        c.Type,
		FilePath:    c.FilePath,
		Duration:    int(c.Duration),
		Disk:        int(c.Disk),
		VideoStream: c.VideoStream,
		Flags:       c.Flags.Slice,
		Events:      c.Events.Slice,
		Cluster:     int(c.Cluster),
		Partition:   int(c.Partition),
		PicIndex:    int(c.PicIndex),
		Repeat:      int(c.Repeat),
		WorkDir:     c.WorkDir,
		WorkDirSN:   int(c.WorkDirSn),
	}
}
