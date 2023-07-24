package statemngr

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"
)

type statemanager struct {
	rep             pkgVolume.VolumeRepository
	logger          pkgLogger.Logger
	mountpoint      string
	volumesInfoPath string
	optionsInfoPath string
}

// errors

var (
	VolumeInfoFileNotFoundError  = errors.New("Volumes info file not found")
	OptionsInfoFileNotFoundError = errors.New("Options info file not found")
)

func NewStateManager(mountpoint string, logger pkgLogger.Logger, rep pkgVolume.VolumeRepository) (StateManager, error) {
	volumesPath := filepath.Join(mountpoint, "state", "volumes.json")
	optionsPath := filepath.Join(mountpoint, "state", "options.json")
	if err := os.MkdirAll(filepath.Join(mountpoint, "state"), 0755); err != nil {
		return nil, err
	}

	return &statemanager{
		rep:             rep,
		mountpoint:      mountpoint,
		volumesInfoPath: volumesPath,
		optionsInfoPath: optionsPath,
		logger:          logger,
	}, nil
}

func (mng *statemanager) SyncState() error {
	volumes, options, err := mng.readState()
	if err != nil {
		return err
	}

	for key, volume := range volumes {
		opt, ok := options[key]
		if ok {
			if err := mng.rep.Create(&volume, &opt); err != nil {
				return err
			}
		}
	}

	return nil
}

func (mnr *statemanager) SaveState() error {
	volumesList, _ := mnr.rep.List()

	volumesMap := make(map[string]volume.Volume)
	optionsMap := make(map[string]models.VolumeOptions)

	for _, vol := range volumesList {
		volumesMap[vol.Name] = *vol
		options := mnr.rep.GetVolumeOptions(vol.Name)
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

	if err := os.WriteFile(mnr.volumesInfoPath, volumesJson, 0644); err != nil {
		return err
	}

	if err := os.WriteFile(mnr.optionsInfoPath, optionsJson, 0644); err != nil {
		return err
	}

	return nil
}

func (mng *statemanager) readState() (map[string]volume.Volume, map[string]models.VolumeOptions, error) {
	data, err := os.ReadFile(mng.volumesInfoPath)

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

	data, err = os.ReadFile(mng.optionsInfoPath)
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
