package ftpmngr

import "github.com/t1d333/docker-volume-ftp-driver/internal/models"

type FTPManager interface {
	CheckConnection(opt *models.FTPConnectionOpt) error
	CheckRemoteDir(remotepath string, opt *models.FTPConnectionOpt) error
}
