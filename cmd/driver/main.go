package main

import (
	"os/user"
	"strconv"

	"github.com/sirupsen/logrus"
	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/statemngr"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/repository"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/service"
)

const (
	mountpoint = "/var/run/docker/ftp-driver/"
)

func main() {
	logger := pkgLogger.InitializeNewLogger()
	rep := repository.CreateInMemoryRepository(logger)
	mng, err := statemngr.NewStateManager(mountpoint, rep)
	serv, err := service.CreateFTPService(mountpoint, mng, rep, logger)
	if err != nil {
		logger.WithFields(logrus.Fields{"Error": err}).Fatal("Failed to create service")
		return
	}

	driver := pkgVolume.InitializeNewFTPDriver(serv, logger)
	handler := volume.NewHandler(driver)
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)

	logger.Info("Serve unix")
	if err := handler.ServeUnix("ftp-driver", gid); err != nil {
		logger.WithField("Error", err).Fatal("Can't start plugin")
	}
}
