package ftpmngr

import (
	"errors"
	"fmt"
	"time"

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
	conn, err := ftp.Dial(getURL(opt), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		mngr.logger.WithField("Error", err).Error("Unable to connect to ftp server")
		return ConnectionError
	}

	if err := conn.Login(opt.User, opt.Password); err != nil {
		mngr.logger.WithField("Error", err).Error("Unable to connect to ftp server")
		return AuthError
	}

	return nil
}
