package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"

	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"
)

type repository struct {
	volumes        *sync.Map
	options        *sync.Map
	mountedVolumes *sync.Map
	logger         pkgLogger.Logger
}

func CreateInMemoryRepository(logger pkgLogger.Logger) pkgVolume.VolumeRepository {
	return &repository{volumes: new(sync.Map), options: new(sync.Map), mountedVolumes: new(sync.Map), logger: logger}
}

func (r *repository) Create(v *volume.Volume, opt *models.VolumeOptions) error {
	if v == nil || opt == nil {
		return errors.New("Volume or options is nil")
	}
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
	if _, err := r.Get(name); err != nil {
		return err
	}

	if r.IsMount(name) {
		return fmt.Errorf("Volume with name '%s' is currently used", name)
	}

	r.volumes.Delete(name)

	return nil
}

func (r *repository) Path(name string) (string, error) {
	tmp, ok := r.volumes.Load(name)
	if !ok {
		return "", errors.New("Volume with this name not found")
	}

	vol := tmp.(*volume.Volume)
	return vol.Mountpoint, nil
}

func (r *repository) Mount(id string, volume *volume.Volume) error {
	ids, ok := r.mountedVolumes.Load(volume.Name)
	if ok {
		ids := ids.(*sync.Map)
		_, ok := ids.LoadOrStore(id, true)

		if ok {
			return errors.New("Volume with this name already mounted")
		}

	} else {
		ids := new(sync.Map)
		ids.Store(id, true)
		r.mountedVolumes.Store(volume.Name, ids)
	}

	return nil
}

func (r *repository) IsMount(name string) bool {
	_, ok := r.mountedVolumes.Load(name)
	if !ok {
		return ok
	}

	list := r.GetMountedIdsList(name)
	return len(list) > 0
}

func (r *repository) Unmount(id, name string) error {
	ids, ok := r.mountedVolumes.Load(name)
	if !ok {
		return errors.New("Volume is not mounted")
	} else {
		ids := ids.(*sync.Map)
		_, ok := ids.Load(id)
		if ok {
			ids.Delete(id)
		} else {
			return errors.New("Volume is not mounted")
		}
	}

	return nil
}

func (r *repository) GetMountedIdsList(name string) []string {
	ids, ok := r.mountedVolumes.Load(name)
	if ok {
		ids := ids.(*sync.Map)
		list := make([]string, 0)

		ids.Range(func(key, value any) bool {
			list = append(list, key.(string))
			return true
		})
		return list
	}

	return []string{}
}
