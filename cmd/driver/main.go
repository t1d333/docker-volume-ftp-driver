package main

import (
	"os/user"
	"strconv"

	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"

	"github.com/docker/go-plugins-helpers/volume"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/repository"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/service"
)

func main() {
	logger := pkgLogger.InitializeNewLogger()
	rep := repository.CreateInMemoryRepository(logger)
	serv, err := service.CreateFTPService(rep, logger)
	if err != nil {
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
