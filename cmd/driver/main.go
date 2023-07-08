package main

import (
	"flag"
	"os/user"
	"strconv"

	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"

	"github.com/docker/go-plugins-helpers/volume"
	pkgVolume "github.com/t1d333/docker-volume-ftp-driver/internal/volume"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/repository"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/service"
)

func main() {
	ftpUser := flag.String("u", "", "FTP server user")
	ftpPassword := flag.String("P", "", "Password for FTP server user")
	ftpPort := flag.Int("p", 20, "Port of FTP server")
	ftpHost := flag.String("h", "", "Host of FTP server")

	flag.Parse()
	if ftpUser == nil || ftpHost == nil || ftpPassword == nil {
		return
	}

	logger := pkgLogger.InitializeNewLogger()
	rep := repository.CreateInMemoryRepository(logger)
	serv, err := service.CreateFTPService(service.FTPServiceOpt{
		Password: *ftpPassword,
		User:     *ftpUser,
		Host:     *ftpHost,
		Port:     *ftpPort,
	}, rep, logger)
	if err != nil {
		return
	}

	driver := pkgVolume.InitializeNewFTPDriver(serv, logger)
	handler := volume.NewHandler(driver)
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)

	if err := handler.ServeUnix("ftp-driver", gid); err != nil {
		logger.WithField("Error", err).Fatal("Can't start plugin")
	}
	logger.Info("Serve unix")
}
