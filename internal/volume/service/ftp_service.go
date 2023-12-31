package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/ftpmngr"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
	"github.com/t1d333/docker-volume-ftp-driver/internal/mountmngr"
	"github.com/t1d333/docker-volume-ftp-driver/internal/statemngr"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"
)

type service struct {
	rep          pkgVolume.VolumeRepository
	stateManager statemngr.StateManager
	mountManager mountmngr.MountManager
	ftpManager   ftpmngr.FTPManager
	logger       pkgLogger.Logger
	mountpoint   string
}

func CreateFTPService(mountpoint string, ftpManager ftpmngr.FTPManager, mountManager mountmngr.MountManager, stateManager statemngr.StateManager, rep pkgVolume.VolumeRepository, logger pkgLogger.Logger) (pkgVolume.VolumeService, error) {
	serv := &service{
		logger:       logger,
		rep:          rep,
		mountpoint:   mountpoint,
		stateManager: stateManager,
		mountManager: mountManager,
		ftpManager:   ftpManager,
	}

	if err := stateManager.SyncState(); err != nil {
		switch {
		case errors.Is(statemngr.OptionsInfoFileNotFoundError, err):
			return serv, nil
		case errors.Is(statemngr.VolumeInfoFileNotFoundError, err):
			return serv, nil
		default:
			return serv, err
		}
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

	if err := s.ftpManager.CheckConnection(&ftpOpt); err != nil {
		return fmt.Errorf("failed to ftpManager.CheckConnection in service.Create: %w", err)
	}

	if err := s.ftpManager.CheckRemoteDir(path, &ftpOpt); err != nil {
		return fmt.Errorf("failed to ftpManager.CheckRemoteDir in service.Create: %w", err)
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

	if err := s.rep.Create(vol, volumeOpt); err != nil {
		return fmt.Errorf("failed to repository.Create in service.Create: %w", err)
	}

	if err := s.stateManager.SaveState(); err != nil {
		s.logger.Errorf("Failed to update state data file: %s", err.Error())
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
		s.logger.Errorf("failed to get volume for remove with name: %s, err : %s", name, err.Error())

		return fmt.Errorf("failed repository.Get in service.Remove: %w", err)
	}

	if isMount := s.rep.IsMount(name); isMount {
		s.logger.Errorf("volume with name '%s' is currently used", name)
		return fmt.Errorf("volume with name '%s' is currently used", name)
	}

	if err := s.rep.Remove(name); err != nil {
		return fmt.Errorf("failed repository.Remove in service.Remove: %w", err)
	}

	if err := s.mountManager.Remove(volume); err != nil {
		return fmt.Errorf("failed mountmngr.Remove in service.Remove: %w", err)
	}

	if err := s.stateManager.SaveState(); err != nil {
		return fmt.Errorf("failed statemngr.SaveState in service.Remove: %w", err)
	}

	return nil
}

func (s *service) Path(name string) (string, error) {
	return s.rep.Path(name)
}

func (s *service) Mount(id, name string) (string, error) {
	volume, err := s.Get(name)
	if err != nil {
		return "", fmt.Errorf("failed service.Get in service.Mount: %w", err)
	}

	if s.rep.IsMount(volume.Name) {
		return volume.Mountpoint, nil
	}

	if err := s.rep.Mount(id, volume); err != nil {
		return "", fmt.Errorf("failed repository.Mount in service.Mount: %w", err)
	}

	opt := s.rep.GetVolumeOptions(volume.Name)

	path, err := s.mountManager.Mount(volume, opt)
	if err != nil {
		return path, fmt.Errorf("failed mountmngr.Mount in service.Mount: %w", err)
	}

	return path, nil
}

func (s *service) Unmount(id, name string) error {
	volume, err := s.Get(name)
	if err != nil {
		return fmt.Errorf("failed service.Get in service.Unmount: %w", err)
	}

	if !s.rep.IsMount(volume.Name) {
		return fmt.Errorf("volume with name: '%s' is not mounted", name)
	}

	if err := s.rep.Unmount(id, name); err != nil {
		return fmt.Errorf("failed repository.Unmount in service.Unmount: %w", err)
	}

	if list := s.rep.GetMountedIdsList(name); len(list) != 0 {
		return nil
	}

	if err := s.mountManager.Unmount(volume); err != nil {
		return fmt.Errorf("failed mountmngr.Mount in service.Mount: %w", err)
	}

	return nil
}

func (s *service) Capabilities() volume.Capability {
	return volume.Capability{Scope: "local"}
}
