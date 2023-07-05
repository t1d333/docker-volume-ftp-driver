package volume

import "github.com/docker/go-plugins-helpers/volume"

type VolumeService interface {
	Create(name string, opt map[string]string) error
	List() ([]*volume.Volume, error)
	Get(name string) (*volume.Volume, error)
	Remove(name string) error
	Path(name string) (string, error)
	Mount(id int, name string) (string, error)
	Unmount(id int, name string) error
	Capabilities() volume.Capability
}
