package volume

import (
	"github.com/docker/go-plugins-helpers/volume"
	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"
)

type FTPDriver struct {
	logger pkgLogger.Logger
	serv   VolumeService
}

func InitializeNewFTPDriver(serv VolumeService, logger pkgLogger.Logger) *FTPDriver {
	return &FTPDriver{logger: logger, serv: serv}
}

func (d *FTPDriver) Create(req *volume.CreateRequest) error {
	d.logger.Infow("create request", "name", req.Name, "opt", req.Options)
	return d.serv.Create(req.Name, req.Options)
}

func (d *FTPDriver) List() (*volume.ListResponse, error) {
	d.logger.Info("List request")
	list, err := d.serv.List()
	return &volume.ListResponse{Volumes: list}, err
}

func (d *FTPDriver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
	d.logger.Infow("get request", "name", req.Name)
	vol, err := d.serv.Get(req.Name)
	return &volume.GetResponse{Volume: vol}, err
}

func (d *FTPDriver) Remove(req *volume.RemoveRequest) error {
	d.logger.Infow("remove request", "name", req.Name)
	return d.serv.Remove(req.Name)
}

func (d *FTPDriver) Path(req *volume.PathRequest) (*volume.PathResponse, error) {
	d.logger.Infow("path request", "name", req.Name)
	path, err := d.serv.Path(req.Name)
	return &volume.PathResponse{Mountpoint: path}, err
}

func (d *FTPDriver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
	d.logger.Infow("mount request", "name", req.Name, "id", req.ID)
	path, err := d.serv.Mount(req.ID, req.Name)
	return &volume.MountResponse{Mountpoint: path}, err
}

func (d *FTPDriver) Unmount(req *volume.UnmountRequest) error {
	d.logger.Infow("unmount request", "name", req.Name, "id", req.ID)
	return d.serv.Unmount(req.ID, req.Name)
}

func (d *FTPDriver) Capabilities() *volume.CapabilitiesResponse {
	d.logger.Info("capabilities request")
	tmp := d.serv.Capabilities()
	return &volume.CapabilitiesResponse{Capabilities: tmp}
}
