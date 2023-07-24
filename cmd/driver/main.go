package main

import (
	"os/user"
	"strconv"

	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/ftpmngr"
	"github.com/t1d333/docker-volume-ftp-driver/internal/mountmngr"
	"github.com/t1d333/docker-volume-ftp-driver/internal/statemngr"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/repository"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/service"
)

const (
	mountpoint = "/var/run/docker/ftp-driver/"
)

func main() {
	logger := pkgLogger.NewLogger()
	rep := repository.CreateInMemoryRepository(logger)

	stateManager, err := statemngr.NewStateManager(mountpoint, logger, rep)
	if err != nil {
		logger.Fatalf("failed to create state manager: %s", err.Error())
		return
	}
	ftpManager := ftpmngr.NewFTPManager(logger)
	mountManager := mountmngr.NewMountManager(logger)
	serv, err := service.CreateFTPService(mountpoint, ftpManager, mountManager, stateManager, rep, logger)
	if err != nil {
		logger.Fatalf("failed to create service: %s", err.Error())
		return
	}

	driver := pkgVolume.InitializeNewFTPDriver(serv, logger)
	handler := volume.NewHandler(driver)
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)

	logger.Info("serve unix")
	if err := handler.ServeUnix("ftp-driver", gid); err != nil {
		logger.Fatalf("can't start plugin: %s", err.Error())
	}
}
