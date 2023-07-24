package ftpmngr

import (
	"errors"
	"fmt"

	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
)

type ftpmngr struct {
	logger *logrus.Logger
}

// errors
var (
	ConnectionError = errors.New("Unable to connect to ftp server")
	AuthError       = errors.New("Unable to connect to ftp server. Failed authentication")
)

func getURL(opt *models.FTPConnectionOpt) string {
	return fmt.Sprintf("%s:%d", opt.Host, opt.Port)
}

func NewFTPManager(logger *logrus.Logger) FTPManager {
	return &ftpmngr{logger: logger}
}

func (mngr *ftpmngr) CheckConnection(opt *models.FTPConnectionOpt) error {
	conn, err := mngr.getConnection(opt)
	conn.Quit()
	return err
}

func (mngr *ftpmngr) getConnection(opt *models.FTPConnectionOpt) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(getURL(opt))
	if err != nil {
		mngr.logger.WithField("Error", err).Error("Unable to connect to ftp server")
		return nil, ConnectionError
	}

	if err := conn.Login(opt.User, opt.Password); err != nil {
		mngr.logger.WithField("Error", err).Error("Unable to connect to ftp server")
		return nil, AuthError
	}

	return conn, nil
}

func (mngr *ftpmngr) CheckRemoteDir(remotepath string, opt *models.FTPConnectionOpt) error {
	conn, err := mngr.getConnection(opt)
	defer conn.Quit()
	if err != nil {
		mngr.logger.Debug(err)
		return err
	}

	err = conn.ChangeDir(remotepath)
	if err != nil {
		mngr.logger.Debug(err)
		return errors.New("Remote dir not found")
	}

	return nil
}
