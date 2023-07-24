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

// errors
var (
	ConnectionError = errors.New("Unable to connect to ftp server")
	AuthError       = errors.New("Unable to connect to ftp server. Failed authentication")
)

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
			mngr.logger.Errorf("Failed to close ftp connection: %s", err.Error())
		}
	}
	return err
}

func (mngr *ftpmngr) getConnection(opt *models.FTPConnectionOpt) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(getURL(opt))
	if err != nil {
		mngr.logger.Errorf("Unable to connect to ftp server: %s", err.Error())
		return nil, ConnectionError
	}

	if err := conn.Login(opt.User, opt.Password); err != nil {
		mngr.logger.Errorf("Unable to login to ftp server: %s", err.Error())
		return nil, AuthError
	}

	return conn, nil
}

func (mngr *ftpmngr) CheckRemoteDir(remotepath string, opt *models.FTPConnectionOpt) error {
	conn, err := mngr.getConnection(opt)
	if err != nil {
		mngr.logger.Errorf("Unable to check remote dir: %s", err.Error())
		return err
	}

	err = conn.ChangeDir(remotepath)
	if err != nil {
		mngr.logger.Errorf("Unable to find remote dir: %s", err.Error())
		return errors.New("Remote dir not found")
	}

	if err := conn.Quit(); err != nil {
		mngr.logger.Errorf("Failed to close ftp connection: %s", err.Error())
	}
	
	return nil
}
