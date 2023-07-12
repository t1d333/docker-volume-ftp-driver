package models

type FTPConnectionOpt struct {
	User     string
	Host     string
	Port     int
	Password string
}

type VolumeOptions struct {
	RemotePath string
	FTPConnectionOpt
}
