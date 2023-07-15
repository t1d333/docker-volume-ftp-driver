package service

import (
	"encoding/json"
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
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
)

const (
	mountpoint = "/var/run/docker/ftp-driver/"
)

// errors

var (
	VolumeInfoFileNotFoundError  = errors.New("Volumes info file not found")
	OptionsInfoFileNotFoundError = errors.New("Options info file not found")
)

type service struct {
	rep             pkgVolume.VolumeRepository
	logger          *logrus.Logger
	mountpoint      string
	volumesInfoPath string
	optionsInfoPath string
}

func getURL(opt models.FTPConnectionOpt) string {
	return fmt.Sprintf("%s:%d", opt.Host, opt.Port)
}

func CreateFTPService(rep pkgVolume.VolumeRepository, logger *logrus.Logger) (pkgVolume.VolumeService, error) {
	volumesPath := filepath.Join(mountpoint, "state", "volumes.json")
	optionsPath := filepath.Join(mountpoint, "state", "options.json")

	if err := os.MkdirAll(filepath.Join(mountpoint, "state"), 0755); err != nil {
		logger.WithFields(logrus.Fields{"Error": err}).Fatalf("Failed to create state directory")
	}

	serv := &service{logger: logger, rep: rep, mountpoint: mountpoint, volumesInfoPath: volumesPath, optionsInfoPath: optionsPath}
	if err := serv.syncData(); err != nil {
		if !errors.Is(err, VolumeInfoFileNotFoundError) && !errors.Is(err, OptionsInfoFileNotFoundError) {
			return serv, err
		}
	}

	return serv, nil
}

func (s *service) syncData() error {
	volumes, options, err := s.readState()
	if err != nil {
		return err
	}

	for key, volume := range volumes {
		opt, ok := options[key]
		if ok {
			if err := s.rep.Create(&volume, &opt); err != nil {
				s.logger.Error("Failed to sync data")
				return err
			}
		} else {
			s.logger.Warnf("Failed to find options for volume %s", key)
		}
	}

	return nil
}

func (s *service) readState() (map[string]volume.Volume, map[string]models.VolumeOptions, error) {
	data, err := os.ReadFile(s.volumesInfoPath)

	volumes := make(map[string]volume.Volume, 0)
	options := make(map[string]models.VolumeOptions, 0)

	if err != nil {
		if os.IsNotExist(err) {
			return volumes, options, VolumeInfoFileNotFoundError
		} else {
			return volumes, options, err
		}
	}

	if err := json.Unmarshal(data, &volumes); err != nil {
		return volumes, options, err
	}

	data, err = os.ReadFile(s.optionsInfoPath)
	if err != nil {
		if os.IsNotExist(err) {
			return volumes, options, OptionsInfoFileNotFoundError
		} else {
			return volumes, options, err
		}
	}

	if err := json.Unmarshal(data, &options); err != nil {
		return volumes, options, err
	}

	return volumes, options, nil
}

func (s *service) saveState() error {
	volumesList, _ := s.List()
	volumesMap := make(map[string]volume.Volume)
	optionsMap := make(map[string]models.VolumeOptions)

	for _, vol := range volumesList {
		volumesMap[vol.Name] = *vol
		options := s.rep.GetVolumeOptions(vol.Name)
		if options != nil {
			optionsMap[vol.Name] = *options
		}
	}

	volumesJson, err := json.Marshal(volumesMap)
	if err != nil {
		return err
	}

	optionsJson, err := json.Marshal(optionsMap)
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.volumesInfoPath, volumesJson, 0644); err != nil {
		return err
	}

	if err := os.WriteFile(s.optionsInfoPath, optionsJson, 0644); err != nil {
		return err
	}

	return nil
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

	if err := s.saveState(); err != nil {
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

	if err := s.saveState(); err != nil {
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
