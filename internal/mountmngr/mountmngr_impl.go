package mountmngr

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"
)

type mountmngr struct {
	logger pkgLogger.Logger
}

// errors

var (
	CreateMountpoinError = errors.New("Failed to create mount point")
	MountError           = errors.New("Failed to mount directory")
	UnmountError         = errors.New("Failed to unmount directory")
	RemoveError          = errors.New("Failed to remove volume directory")
)

func NewMountManager(logger pkgLogger.Logger) MountManager {
	return &mountmngr{logger: logger}
}

func (mngr *mountmngr) Mount(vol *volume.Volume, opt *models.VolumeOptions) (string, error) {
	if err := os.MkdirAll(vol.Mountpoint, 0755); err != nil {
		return "", CreateMountpoinError
	}

	ftpPath := fmt.Sprintf("%s:%d%s", opt.Host, opt.Port, opt.RemotePath)

	cmd := exec.Command("curlftpfs", ftpPath, vol.Mountpoint, "-o", fmt.Sprintf("user=%s:%s", opt.User, opt.Password), "-o", "nonempty")

	if err := cmd.Run(); err != nil {
		return "", MountError
	}

	return vol.Mountpoint, nil
}

func (mngr *mountmngr) Unmount(volume *volume.Volume) error {
	cmd := exec.Command("umount", volume.Mountpoint)
	if err := cmd.Run(); err != nil {
		return UnmountError
	}

	if err := os.RemoveAll(volume.Mountpoint); err != nil {
		return UnmountError
	}

	return nil
}

func (mngr *mountmngr) Remove(volume *volume.Volume) error {
	if err := os.RemoveAll(volume.Mountpoint); err != nil {
		return RemoveError
	}

	return nil
}
