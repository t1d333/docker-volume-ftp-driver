package volume

import (
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"
)

type FTPDriver struct {
	logger *logrus.Logger
}

func InitializeNewFTPDriver(logger *logrus.Logger) *FTPDriver {
	return &FTPDriver{logger: logger}
}

func (d *FTPDriver) Create(req *volume.CreateRequest) error {
	d.logger.WithFields(logrus.Fields{"Name": req.Name, "Opt":  req.Options}).Info("Create request")
	return nil
}

func (d *FTPDriver) List() (*volume.ListResponse, error) {
	d.logger.Info("List request")
	return nil, nil
}

func (d *FTPDriver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
	d.logger.WithFields(logrus.Fields{"Name": req.Name}).Info("Get request")
	return nil, nil
}

func (d *FTPDriver) Remove(req *volume.RemoveRequest) error {
	d.logger.WithFields(logrus.Fields{"Name": req.Name}).Info("Remove request")
	return nil
}

func (d *FTPDriver) Path(req *volume.PathRequest) (*volume.PathResponse, error) {
	d.logger.WithFields(logrus.Fields{"Name": req.Name}).Info("Path request")
	return nil, nil
}

func (d *FTPDriver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
	d.logger.WithFields(logrus.Fields{"ID": req.ID, "Name": req.Name}).Info("Mount request")
	return nil, nil
}

func (d *FTPDriver) Unmount(req *volume.UnmountRequest) error {
	d.logger.WithFields(logrus.Fields{"ID": req.ID, "Name": req.Name}).Info("Unmount request")
	return nil
}

func (d *FTPDriver) Capabilities() *volume.CapabilitiesResponse {
	d.logger.Info("Capabilities request")
	return nil
}
