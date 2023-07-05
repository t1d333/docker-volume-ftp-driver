package volume

import "github.com/docker/go-plugins-helpers/volume"

type VolumeRepository interface {
	Create(id int, name string) error
	List() ([]*volume.Volume, error)
	Get(name string) (*volume.Volume, error)
	Remove(name string) error
	Path(name string) (string, error)
	Mount(id int, name string) (string, error)
	Unmount(id int, name string) error
	Capabilities() volume.Capability
}
