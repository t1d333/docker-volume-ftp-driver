package repository

import (
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"

	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
)

type repository struct {
	data   *sync.Map
	logger *logrus.Logger
}

func NewInMemoryRepository(logger *logrus.Logger) pkgVolume.VolumeRepository {
	return &repository{data: new(sync.Map), logger: logger}
}

func (r *repository) Create(id int, name string) error {
	return nil
}

func (r *repository) List() ([]*volume.Volume, error) {
	return []*volume.Volume{}, nil
}

func (r *repository) Get(name string) (*volume.Volume, error) {
	return nil, nil
}

func (r *repository) Remove(name string) error {
	return nil
}

func (r *repository) Path(name string) (string, error) {
	return "", nil
}

func (r *repository) Mount(id int, name string) (string, error) {
	return "", nil
}

func (r *repository) Unmount(id int, name string) error {
	return nil
}

func (r *repository) Capabilities() volume.Capability {
	return volume.Capability{}
}
