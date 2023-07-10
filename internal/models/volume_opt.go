package models

import "github.com/jlaffaye/ftp"

type FTPConnectionOpt struct {
	User     string
	Host     string
	Port     int
	Password string
	Conn     *ftp.ServerConn
}

type VolumeOptions struct {
	RemotePath string
	FTPConnectionOpt
}
