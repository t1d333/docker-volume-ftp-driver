package mountmngr

import (
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
)

type MountManager interface {
	Mount(vol *volume.Volume, opt *models.VolumeOptions) (string, error)
	Unmount(vol *volume.Volume) error
	Remove(vol *volume.Volume) error
}
