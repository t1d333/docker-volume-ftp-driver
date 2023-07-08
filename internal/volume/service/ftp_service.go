package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
)

const (
	mountpoint = "/var/run/docker/ftp-driver/"
)

type service struct {
	conn       *ftp.ServerConn
	rep        pkgVolume.VolumeRepository
	logger     *logrus.Logger
	opt        FTPServiceOpt
	mountpoint string
}

type FTPServiceOpt struct {
	User     string
	Host     string
	Port     int
	Password string
}

func getURL(opt FTPServiceOpt) string {
	return fmt.Sprintf("%s:%d", opt.Host, opt.Port)
}

func CreateFTPService(opt FTPServiceOpt, rep pkgVolume.VolumeRepository, logger *logrus.Logger) (pkgVolume.VolumeService, error) {
	logger.Info("Connecting to ftp server...")
	conn, err := ftp.Dial(getURL(opt), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		logger.WithField("Error", err).Error("Unable to connect to ftp server")
		return nil, errors.New("Unable to connect to ftp server")
	}

	if err := conn.Login(opt.User, opt.Password); err != nil {
		logger.WithField("Error", err).Error("Unable to connect to ftp server")
		return nil, errors.New("Unable to connect to ftp server. Failed authentication")
	}

	logger.Info("Successful connection to ftp server")
	return &service{conn: conn, logger: logger, rep: rep, opt: opt, mountpoint: mountpoint}, nil
}

func (s *service) Create(name string, opt map[string]string) error {
	path, ok := opt["remotepath"]
	if !ok {
		path = "/"
	}

	if ok {
		if err := s.conn.MakeDir(strings.TrimPrefix(path, "/")); err != nil {
			s.logger.Error("Failed to create remote directory")
			return err
		}
	}

	vol := &volume.Volume{
		Name:       name,
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
		Mountpoint: filepath.Join(s.mountpoint, name),
	}

	volumeOpt := &models.VolumeOptions{
		RemotePath: path,
	}

	return s.rep.Create(vol, volumeOpt)
}

func (s *service) List() ([]*volume.Volume, error) {
	return s.rep.List()
}

func (s *service) Get(name string) (*volume.Volume, error) {
	return s.rep.Get(name)
}

func (s *service) Remove(name string) error {
	return nil
}

func (s *service) Path(name string) (string, error) {
	return s.rep.Path(name)
}

func (s *service) Mount(id, name string) (string, error) {
	volume, err := s.Get(name)
	if err != nil {
		return "", err
	}

	if s.rep.IsMount(volume.Name) {
		return volume.Mountpoint, nil
	}

	if err := s.rep.Mount(id, volume); err != nil {
		return "", err
	}

	if err := os.MkdirAll(volume.Mountpoint, 0755); err != nil {
		s.logger.WithField("Error", err).Error("Failed to create mount point")
		return "", errors.New("Failed to create mount point")
	}

	opt := s.rep.GetVolumeOptions(volume.Name)

	ftpPath := fmt.Sprintf("%s:%d%s", s.opt.Host, s.opt.Port, opt.RemotePath)

	cmd := exec.Command("curlftpfs", ftpPath, volume.Mountpoint, "-o", fmt.Sprintf("user=%s:%s", s.opt.User, s.opt.Password), "-o", "nonempty")
	if out, err := cmd.CombinedOutput(); err != nil {
		s.logger.WithFields(logrus.Fields{"Error": err, "Out": string(out)}).Error("Failed to mount directory")
		return "", errors.New("Failed to mount directory")
	}

	return volume.Mountpoint, nil
}

func (s *service) Unmount(id, name string) error {
	volume, err := s.Get(name)
	if err != nil {
		return err
	}

	if !s.rep.IsMount(volume.Name) {
		return fmt.Errorf("Volume with name: '%s' is not mounted", name)
	}

	if err := s.rep.Unmount(id, name); err != nil {
		return err
	}

	if list := s.rep.GetMountedIdsList(name); len(list) != 0 {
		s.logger.Info("Mounted list", list)
		return nil
	}

	cmd := exec.Command("umount", volume.Mountpoint)
	if out, err := cmd.CombinedOutput(); err != nil {
		s.logger.WithFields(logrus.Fields{"Error": err, "Out": string(out)}).Error("Failed to unmount directory")
		return errors.New("Failed to unmount directory")
	}

	if err := os.RemoveAll(volume.Mountpoint); err != nil {
		s.logger.WithFields(logrus.Fields{"Error": err}).Error("Failed to remove mounted directory")
		return errors.New("Failed to remove mounted directory")
	}

	return nil
}

func (s *service) Capabilities() volume.Capability {
	return volume.Capability{Scope: "local"}
}
