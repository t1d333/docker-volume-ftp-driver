package main

import (
	"os/user"
	"strconv"

	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"

	"github.com/docker/go-plugins-helpers/volume"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
)

func main() {
	logger := pkgLogger.InitializeNewLogger()
	driver := pkgVolume.InitializeNewFTPDriver(logger)
	handler := volume.NewHandler(driver)
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)
    logger.Info("Serve unix")
	handler.ServeUnix("ftp-driver", gid)
}
