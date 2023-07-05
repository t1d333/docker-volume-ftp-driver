package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
)

type service struct {
	conn   *ftp.ServerConn
	logger *logrus.Logger
}

type FTPServiceOpt struct {
	User     string
	Host     string
	Port     int
	Password string
}

func getURL(opt FTPServiceOpt) string {
	return fmt.Sprintf("ftp://%s:%s@%s:%d", opt.User, opt.Password, opt.Host, opt.Port)
}

func CreateFTPService(opt FTPServiceOpt, logger *logrus.Logger) (pkgVolume.VolumeService, error) {
	logger.Info("Connecting to ftp server...")
	conn, err := ftp.Dial(getURL(opt), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		logger.WithField("Error", err).Error("Unable to connect to ftp server")
		return nil, errors.New("Unable to connect to ftp server")
	}

	return &service{conn: conn, logger: logger}, nil
}

func (s *service) Create(name string, opt map[string]string) error {
	return nil
}

func (s *service) List() ([]*volume.Volume, error) {
	return []*volume.Volume{}, nil
}

func (s *service) Get(name string) (*volume.Volume, error) {
	return nil, nil
}

func (s *service) Remove(name string) error {
	return nil
}

func (s *service) Path(name string) (string, error) {
	return "", nil
}

func (s *service) Mount(id int, name string) (string, error) {
	return "", nil
}

func (s *service) Unmount(id int, name string) error {
	return nil
}

func (s *service) Capabilities() volume.Capability {
	return volume.Capability{}
}
