package volume

import (
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"
)

type FTPDriver struct {
	logger *logrus.Logger
	serv   VolumeService
}

func InitializeNewFTPDriver(serv VolumeService, logger *logrus.Logger) *FTPDriver {
	return &FTPDriver{logger: logger, serv: serv}
}

func (d *FTPDriver) Create(req *volume.CreateRequest) error {
	d.logger.WithFields(logrus.Fields{"Name": req.Name, "Opt": req.Options}).Info("Create request")
	return d.serv.Create(req.Name, req.Options)
}

func (d *FTPDriver) List() (*volume.ListResponse, error) {
	d.logger.Info("List request")
	list, err := d.serv.List()
	return &volume.ListResponse{Volumes: list}, err
}

func (d *FTPDriver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
	d.logger.WithFields(logrus.Fields{"Name": req.Name}).Info("Get request")
	vol, err := d.serv.Get(req.Name)
	return &volume.GetResponse{Volume: vol}, err
}

func (d *FTPDriver) Remove(req *volume.RemoveRequest) error {
	d.logger.WithFields(logrus.Fields{"Name": req.Name}).Info("Remove request")
	return d.serv.Remove(req.Name)
}

func (d *FTPDriver) Path(req *volume.PathRequest) (*volume.PathResponse, error) {
	d.logger.WithFields(logrus.Fields{"Name": req.Name}).Info("Path request")
	path, err := d.serv.Path(req.Name)
	return &volume.PathResponse{Mountpoint: path}, err
}

func (d *FTPDriver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
	d.logger.WithFields(logrus.Fields{"ID": req.ID, "Name": req.Name}).Info("Mount request")
	path, err := d.serv.Mount(req.ID, req.Name)
	return &volume.MountResponse{Mountpoint: path}, err
}

func (d *FTPDriver) Unmount(req *volume.UnmountRequest) error {
	d.logger.WithFields(logrus.Fields{"ID": req.ID, "Name": req.Name}).Info("Unmount request")
	return d.serv.Unmount(req.ID, req.Name)
}

func (d *FTPDriver) Capabilities() *volume.CapabilitiesResponse {
	d.logger.Info("Capabilities request")
	return nil
}
