package volume

import (
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
)

type VolumeRepository interface {
	Create(v *volume.Volume, opt *models.VolumeOptions) error
	List() ([]*volume.Volume, error)
	Get(name string) (*volume.Volume, error)
	Remove(name string) error
	Path(name string) (string, error)
	Mount(id string, volume *volume.Volume) error
	Unmount(id, name string) error
	IsMount(name string) bool
	GetMountedIdsList(name string) []string
	GetVolumeOptions(name string) *models.VolumeOptions
}
