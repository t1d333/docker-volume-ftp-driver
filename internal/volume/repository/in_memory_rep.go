package repository

import (
	"errors"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"

	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
)

type repository struct {
	volumes        *sync.Map
	options        *sync.Map
	mountedVolumes *sync.Map
	logger         *logrus.Logger
}

func CreateInMemoryRepository(logger *logrus.Logger) pkgVolume.VolumeRepository {
	return &repository{volumes: new(sync.Map), options: new(sync.Map), mountedVolumes: new(sync.Map), logger: logger}
}

func (r *repository) Create(v *volume.Volume, opt *models.VolumeOptions) error {
	_, ok := r.volumes.LoadOrStore(v.Name, v)
	if ok {
		return errors.New("Volume arleady exists")
	}

	r.options.Store(v.Name, opt)
	return nil
}

func (r *repository) GetVolumeOptions(name string) *models.VolumeOptions {
	opt, _ := r.options.Load(name)
	return opt.(*models.VolumeOptions)
}

func (r *repository) List() ([]*volume.Volume, error) {
	res := make([]*volume.Volume, 0)
	r.volumes.Range(func(key any, value any) bool {
		res = append(res, value.(*volume.Volume))
		return true
	})

	return res, nil
}

func (r *repository) Get(name string) (*volume.Volume, error) {
	vol, ok := r.volumes.Load(name)
	if !ok {
		return nil, errors.New("Volume not found")
	}

	return vol.(*volume.Volume), nil
}

func (r *repository) Remove(name string) error {
	return nil
}

func (r *repository) Path(name string) (string, error) {
	tmp, ok := r.volumes.Load("name")
	if !ok {
		return "", errors.New("Volume with this name not found")
	}

	vol := tmp.(*volume.Volume)
	return vol.Mountpoint, nil
}

func (r *repository) Mount(volume *volume.Volume) error {
	_, ok := r.mountedVolumes.LoadOrStore(volume.Name, volume)
	if ok {
		return errors.New("Volume with this name already mounted")
	}

	return nil
}

func (r *repository) IsMount(name string) bool {
	_, ok := r.mountedVolumes.Load(name)
	return ok
}

func (r *repository) Unmount(id, name string) error {
	return nil
}
