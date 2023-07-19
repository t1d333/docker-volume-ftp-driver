# Docker volume plugin for ftp

![Coverage](https://img.shields.io/badge/Coverage-68.9%25-yellow)
[![Go Report Card](https://goreportcard.com/badge/github.com/vieux/docker-volume-sshfs)](https://goreportcard.com/report/github.com/t1d333/docker-volume-ftp-driver)

This plugin allows you to mount remotes directories from ftp server

## Usage

### Install plugin

```
$ docker plugin install t1d333/ftp-driver:latest
```

or using GitHub

```
# clone repository
$ git clone https://github.com/t1d333/docker-volume-ftp-driver.github

$ cd docker-volume-ftp-driver

# build plugin
$ make

# enable plugin
$ make enable
```

### Create a volume

Options

- `host` - host name or IP address of the ftp server
- `user` - username for authentication
- `password` - password for authentication
- `port` - port for connection
- `remotepath` - path on ftp server for mount. The path must be **_absolute_** and **_start with a /_**

**_The source path on the ftp server must exist_**, otherwise an error will occur when creating the volume.

All options except `remotepath` are **_required_**

```
$ docker volume create -d t1d333/ftp-driver \
	-o host=<hostname or ip address> \
	-o user=<username> \
	-o password=<password> \
	-o port=<port> [-o remotepath=<remotepath>] \
	--name ftpvolume

$ docker volume ls
DRIVER                     VOLUME NAME
local                      e83cd4dda625e94150533511447050b3d7b29f194e37777cf0c3d8b547696b45
local                      ea4cfbe0c50202f1616b555c6e3d2087563b102ae02092102a00470bed8dc921
local                      ea8ebc180669eeca9b238b84d3a72a13b36bce00c68cf3d2caf6ae15415cfe30
local                      eda644e6c01842da1a59e746b774804016417129a05d210f99f54a0008f3d4eb
t1d333/ftp-driver:latest   ftpvolume
```

## Use the volume

```
$ docker run -it -v ftpvolume:<path> debian ls <path>
```
