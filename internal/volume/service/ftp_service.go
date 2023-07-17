package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
	"github.com/t1d333/docker-volume-ftp-driver/internal/statemngr"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
)

type service struct {
	rep        pkgVolume.VolumeRepository
	mng        statemngr.StateManager
	logger     *logrus.Logger
	mountpoint string
}

func getURL(opt models.FTPConnectionOpt) string {
	return fmt.Sprintf("%s:%d", opt.Host, opt.Port)
}

func CreateFTPService(mountpoint string, mng statemngr.StateManager, rep pkgVolume.VolumeRepository, logger *logrus.Logger) (pkgVolume.VolumeService, error) {
	serv := &service{logger: logger, rep: rep, mountpoint: mountpoint, mng: mng}

	if err := mng.SyncState(); err != nil {
		return serv, err
	}
	return serv, nil
}

func (s *service) Create(name string, opt map[string]string) error {
	path, ok := opt["remotepath"]

	if !ok {
		path = "/"
	}

	ftpOpt := models.FTPConnectionOpt{}

	if host, ok := opt["host"]; !ok {
		return errors.New("Not specified from required one of the options")
	} else {
		ftpOpt.Host = host
	}

	if user, ok := opt["user"]; !ok {
		return errors.New("Not specified from required one of the options")
	} else {
		ftpOpt.User = user
	}

	if port, ok := opt["port"]; !ok {
		return errors.New("Not specified from required one of the options")
	} else {
		port, err := strconv.Atoi(port)
		if err != nil {
			return errors.New("Not a valid port")
		}
		ftpOpt.Port = port
	}

	if password, ok := opt["password"]; !ok {
		return errors.New("Not specified from required one of the options")
	} else {
		ftpOpt.Password = password
	}

	conn, err := ftp.Dial(getURL(ftpOpt), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		s.logger.WithField("Error", err).Error("Unable to connect to ftp server")
		return errors.New("Unable to connect to ftp server")
	}

	if err := conn.Login(ftpOpt.User, ftpOpt.Password); err != nil {
		s.logger.WithField("Error", err).Error("Unable to connect to ftp server")
		return errors.New("Unable to connect to ftp server. Failed authentication")
	}

	vol := &volume.Volume{
		Name:       name,
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
		Mountpoint: filepath.Join(s.mountpoint, name),
	}

	volumeOpt := &models.VolumeOptions{
		RemotePath:       path,
		FTPConnectionOpt: ftpOpt,
	}

	// TODO: добавить обработку параметра remotepath, если каталог по remotepath не существует

	if err := s.rep.Create(vol, volumeOpt); err != nil {
		return err
	}

	if err := s.mng.SaveState(); err != nil {
		s.logger.WithFields(logrus.Fields{"Error": err}).Error("Failed to update state data file")
	}

	return nil
}

func (s *service) List() ([]*volume.Volume, error) {
	return s.rep.List()
}

func (s *service) Get(name string) (*volume.Volume, error) {
	return s.rep.Get(name)
}

func (s *service) Remove(name string) error {
	volume, err := s.rep.Get(name)
	if err != nil {
		s.logger.WithFields(logrus.Fields{"Name": name, "Error": err}).Error("Failed to get volume for remove")
		return err
	}

	if isMount := s.rep.IsMount(name); isMount {
		s.logger.WithFields(logrus.Fields{"Name": name}).Error("Volume with is currently used")
		return fmt.Errorf("Volume with name '%s' is currently used", name)
	}

	if err := s.rep.Remove(name); err != nil {
		return err
	}

	if err := os.RemoveAll(volume.Mountpoint); err != nil {
		s.logger.WithFields(logrus.Fields{"Error": err}).Error("Failed to remove volume directory")
		return err
	}

	if err := s.mng.SaveState(); err != nil {
		s.logger.Error("Failed to update state data file")
	}

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

	ftpPath := fmt.Sprintf("%s:%d%s", opt.Host, opt.Port, opt.RemotePath)

	cmd := exec.Command("curlftpfs", ftpPath, volume.Mountpoint, "-o", fmt.Sprintf("user=%s:%s", opt.User, opt.Password), "-o", "nonempty")
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
