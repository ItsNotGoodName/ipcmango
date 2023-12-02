package dahua

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func NewMemCameraStore() *MemCameraStore {
	return &MemCameraStore{
		mu:   sync.Mutex{},
		data: make(map[int64]models.DahuaCamera),
	}
}

type MemCameraStore struct {
	mu   sync.Mutex
	data map[int64]models.DahuaCamera
}

func (s *MemCameraStore) Save(ctx context.Context, camera ...models.DahuaCamera) error {
	s.mu.Lock()
	for _, camera := range camera {
		s.data[camera.ID] = camera
	}
	s.mu.Unlock()
	return nil
}

func (s *MemCameraStore) Get(ctx context.Context, id int64) (models.DahuaCamera, bool, error) {
	s.mu.Lock()
	camera, found := s.data[id]
	s.mu.Unlock()
	return camera, found, nil
}

func (s *MemCameraStore) List(ctx context.Context) ([]models.DahuaCamera, error) {
	s.mu.Lock()
	cameras := make([]models.DahuaCamera, 0, len(s.data))
	for _, c := range s.data {
		cameras = append(cameras, c)
	}
	s.mu.Unlock()
	slices.SortFunc(cameras, func(a, b models.DahuaCamera) int { return cmp.Compare(a.ID, b.ID) })
	return cameras, nil
}
