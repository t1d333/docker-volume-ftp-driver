package ftpmngr

import (
	"errors"
	"fmt"

	"github.com/jlaffaye/ftp"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
	pkgLogger "github.com/t1d333/docker-volume-ftp-driver/pkg/logger"
)

type ftpmngr struct {
	logger pkgLogger.Logger
}

func getURL(opt *models.FTPConnectionOpt) string {
	return fmt.Sprintf("%s:%d", opt.Host, opt.Port)
}

func NewFTPManager(logger pkgLogger.Logger) FTPManager {
	return &ftpmngr{logger: logger}
}

func (mngr *ftpmngr) CheckConnection(opt *models.FTPConnectionOpt) error {
	conn, err := mngr.getConnection(opt)
	if err == nil {
		if err := conn.Quit(); err != nil {
			mngr.logger.Errorf("failed to close ftp connection: %s", err.Error())
		}
	} else {
		return fmt.Errorf("failed to connect to ftp server: %w", err)
	}

	return nil
}

func (mngr *ftpmngr) getConnection(opt *models.FTPConnectionOpt) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(getURL(opt))
	if err != nil {
		mngr.logger.Errorf("unable to connect to ftp server: %s", err.Error())
		return nil, fmt.Errorf("unable to connect to ftp server in ftpmngr.getConnection: %w", err)
	}

	if err := conn.Login(opt.User, opt.Password); err != nil {
		mngr.logger.Errorf("unable to login to ftp server: %s", err.Error())
		return nil, fmt.Errorf("unable to login to ftp server in ftpmngr.getConnection: %w", err)
	}

	return conn, nil
}

func (mngr *ftpmngr) CheckRemoteDir(remotepath string, opt *models.FTPConnectionOpt) error {
	conn, err := mngr.getConnection(opt)
	if err != nil {
		mngr.logger.Errorf("unable to check remote dir: %s", err.Error())
		return fmt.Errorf("unable to connect to ftp server in ftpmngr.CheckRemoteDir: %w", err)
	}

	err = conn.ChangeDir(remotepath)
	if err != nil {
		mngr.logger.Errorf("unable to find remote dir: %s", err.Error())
		return errors.New("remote dir not found")
	}

	if err := conn.Quit(); err != nil {
		mngr.logger.Errorf("Failed to close ftp connection: %s", err.Error())
	}

	return nil
}
