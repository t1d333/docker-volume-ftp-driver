package service

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	ftpMock "github.com/t1d333/docker-volume-ftp-driver/internal/ftpmngr/mocks"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
	mountMock "github.com/t1d333/docker-volume-ftp-driver/internal/mountmngr/mocks"
	stateMock "github.com/t1d333/docker-volume-ftp-driver/internal/statemngr/mocks"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/repository"
)

func TestGet(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	ftpmngr := ftpMock.NewFTPManager(t)
	mountmngr := mountMock.NewMountManager(t)
	statemngr := stateMock.NewStateManager(t)
	rep := repository.CreateInMemoryRepository(logger)
	mountpoint := "/test"

	inVolume := &volume.Volume{
		Name:       "test",
		Mountpoint: "/test/abc",
		Status:     make(map[string]interface{}),
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
	}

	err := rep.Create(inVolume, &models.VolumeOptions{})

	require.Nil(t, err)

	statemngr.On("SyncState").Return(nil)

	t.Run("succsess get volume", func(t *testing.T) {
		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		got, err := serv.Get("test")

		require.Nil(t, err)

		assert.Equal(t, inVolume.Name, got.Name)
		assert.Equal(t, inVolume.CreatedAt, got.CreatedAt)
		assert.Equal(t, inVolume.Mountpoint, got.Mountpoint)
		assert.Equal(t, inVolume.Status, got.Status)
	})

	t.Run("failed get volume", func(t *testing.T) {
		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		_, err = serv.Get("test1")

		assert.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	ftpmngr := ftpMock.NewFTPManager(t)
	mountmngr := mountMock.NewMountManager(t)
	statemngr := stateMock.NewStateManager(t)
	rep := repository.CreateInMemoryRepository(logger)
	mountpoint := "/test"

	statemngr.On("SyncState").Return(nil)
	statemngr.On("SaveState").Return(nil)
	ftpmngr.On("CheckConnection", mock.Anything).Return(nil).Once()

	t.Run("succsess creation", func(t *testing.T) {
		name := "volume"
		expectedPath := "/test/volume"
		opt := map[string]string{
			"user":     "admin",
			"host":     "localhost",
			"password": "password",
			"port":     "21",
		}

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		err = serv.Create(name, opt)
		require.Nil(t, err)

		got, _ := rep.Get(name)

		require.NotNil(t, got)

		assert.Equal(t, name, got.Name)
		assert.Equal(t, expectedPath, got.Mountpoint)
	})

	t.Run("creation without password option", func(t *testing.T) {
		name := "withoutPassword"
		opt := map[string]string{
			"user": "admin",
			"host": "localhost",
			"port": "21",
		}

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		err = serv.Create(name, opt)
		require.Error(t, err)
	})

	t.Run("creation without user option", func(t *testing.T) {
		name := "withoutUser"
		opt := map[string]string{
			"host":     "localhost",
			"password": "pswd",
			"port":     "21",
		}

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		err = serv.Create(name, opt)
		require.Error(t, err)
	})

	t.Run("creation without host option", func(t *testing.T) {
		name := "withoutHost"
		opt := map[string]string{
			"user":     "user",
			"password": "pswd",
			"port":     "21",
		}

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		err = serv.Create(name, opt)
		require.Error(t, err)
	})

	t.Run("creation without port option", func(t *testing.T) {
		name := "withoutHost"
		opt := map[string]string{
			"user":     "user",
			"password": "pswd",
			"host":     "host",
		}

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		err = serv.Create(name, opt)
		require.Error(t, err)
	})

	t.Run("creation with invalid port option", func(t *testing.T) {
		name := "withoutHost"
		opt := map[string]string{
			"user":     "user",
			"password": "pswd",
			"host":     "host",
			"port":     "abc",
		}

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		err = serv.Create(name, opt)
		require.Error(t, err)
	})

	ftpmngr.On("CheckConnection", mock.Anything).Return(nil).Once()
	stateError := errors.New("Failed to save state")
	statemngr.On("SaveState").Return(stateError)

	t.Run("get error from state manager", func(t *testing.T) {
		name := "stateError"
		opt := map[string]string{
			"user":     "user",
			"password": "pswd",
			"host":     "host",
			"port":     "21",
		}

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		err = serv.Create(name, opt)
		require.Nil(t, err)
	})

	ftpError := errors.New("Failed to connect to ftp server")
	ftpmngr.On("CheckConnection", mock.Anything).Return(ftpError)
	statemngr.On("SaveState").Return(nil)

	t.Run("get error from ftp manager", func(t *testing.T) {
		name := "ftpError"
		opt := map[string]string{
			"user":     "user",
			"password": "pswd",
			"host":     "host",
			"port":     "21",
		}

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		err = serv.Create(name, opt)
		require.Error(t, err)
	})
}

func TestList(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	ftpmngr := ftpMock.NewFTPManager(t)
	mountmngr := mountMock.NewMountManager(t)
	statemngr := stateMock.NewStateManager(t)
	rep := repository.CreateInMemoryRepository(logger)
	mountpoint := "/test"

	inVolume := &volume.Volume{
		Name:       "test",
		Mountpoint: "/test/abc",
		Status:     make(map[string]interface{}),
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
	}

	err := rep.Create(inVolume, &models.VolumeOptions{})

	require.Nil(t, err)

	statemngr.On("SyncState").Return(nil)

	t.Run("succsess get list", func(t *testing.T) {
		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		got, err := serv.List()

		require.Nil(t, err)

		require.Equal(t, 1, len(got))
		assert.Equal(t, inVolume.Name, got[0].Name)
		assert.Equal(t, inVolume.CreatedAt, got[0].CreatedAt)
		assert.Equal(t, inVolume.Mountpoint, got[0].Mountpoint)
		assert.Equal(t, inVolume.Status, got[0].Status)
	})
}

func TestRemove(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	ftpmngr := ftpMock.NewFTPManager(t)
	mountmngr := mountMock.NewMountManager(t)
	statemngr := stateMock.NewStateManager(t)
	mountpoint := "/test"

	inVolume := &volume.Volume{
		Name:       "test",
		Mountpoint: "/test/abc",
		Status:     make(map[string]interface{}),
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
	}

	statemngr.On("SyncState").Return(nil)
	statemngr.On("SaveState").Return(nil).Once()
	mountmngr.On("Remove", inVolume).Return(nil).Once()

	t.Run("succsess remove", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)
		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		err = serv.Remove("test")
		assert.Nil(t, err)
	})

	t.Run("remove not exists volume", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)
		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		err = serv.Remove("notExists")
		assert.Error(t, err)
	})

	t.Run("remove mounted volume", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)
		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		err = rep.Mount(uuid.NewString(), inVolume)

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		err = serv.Remove("test")
		assert.Error(t, err)
	})

	mountmngr.On("Remove", mock.Anything).Return(errors.New("Unexpected")).Once()

	t.Run("get error from mount manager", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)
		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		err = serv.Remove("test")

		assert.Error(t, err)
	})

	statemngr.On("SaveState").Return(errors.New("Unexpected")).Once()
	mountmngr.On("Remove", mock.Anything).Return(nil).Once()

	t.Run("get error from mount manager", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)
		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		err = serv.Remove("test")

		assert.Nil(t, err)
	})
}

func TestMount(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	ftpmngr := ftpMock.NewFTPManager(t)
	mountmngr := mountMock.NewMountManager(t)
	statemngr := stateMock.NewStateManager(t)
	mountpoint := "/test"

	inVolume := &volume.Volume{
		Name:       "test",
		Mountpoint: "/test/abc",
		Status:     make(map[string]interface{}),
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
	}

	statemngr.On("SyncState").Return(nil)

	mountmngr.On("Mount", mock.Anything, mock.Anything).Return(inVolume.Mountpoint, nil).Once()

	t.Run("succsess mount", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)
		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		path, err := serv.Mount(uuid.NewString(), "test")

		assert.Nil(t, err)
		assert.Equal(t, inVolume.Mountpoint, path)
	})

	t.Run("mount already mounted volume", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)
		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		err = rep.Mount(uuid.NewString(), inVolume)

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		path, err := serv.Mount(uuid.NewString(), "test")

		assert.Nil(t, err)

		assert.Equal(t, inVolume.Mountpoint, path)
	})

	mountmngr.On("Mount", mock.Anything, mock.Anything).Return(inVolume.Mountpoint, errors.New("Unexpected")).Once()

	t.Run("get error from mount manager", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)
		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		_, err = serv.Mount(uuid.NewString(), "test")

		assert.Error(t, err)
	})

	t.Run("mount not exists volume", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		_, err = serv.Mount(uuid.NewString(), "test")

		assert.Error(t, err)
	})
}

func TestUnmount(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	ftpmngr := ftpMock.NewFTPManager(t)
	mountmngr := mountMock.NewMountManager(t)
	statemngr := stateMock.NewStateManager(t)
	mountpoint := "/test"
	id := uuid.NewString()

	inVolume := &volume.Volume{
		Name:       "test",
		Mountpoint: "/test/abc",
		Status:     make(map[string]interface{}),
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
	}

	statemngr.On("SyncState").Return(nil)

	mountmngr.On("Unmount", mock.Anything).Return(nil).Once()

	t.Run("succsess unmount", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)
		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		err = rep.Mount(id, inVolume)

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		err = serv.Unmount(id, inVolume.Name)

		assert.Nil(t, err)
	})

	t.Run("unmount not exists volume", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		err = serv.Unmount(id, inVolume.Name)

		assert.Error(t, err)
	})

	t.Run("unmount not mounted volume", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)

		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		err = serv.Unmount(id, inVolume.Name)

		assert.Error(t, err)
	})

	t.Run("twice mount one volume", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)

		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		id2 := uuid.NewString()

		err = rep.Mount(id, inVolume)

		require.Nil(t, err)

		err = rep.Mount(id2, inVolume)

		require.Nil(t, err)

		err = serv.Unmount(id, inVolume.Name)

		require.Nil(t, err)

		assert.True(t, rep.IsMount(inVolume.Name))
	})

	mountmngr.On("Unmount", mock.Anything).Return(errors.New("Unexpected"))

	t.Run("get error from mount manager", func(t *testing.T) {
		rep := repository.CreateInMemoryRepository(logger)

		err := rep.Create(inVolume, &models.VolumeOptions{})

		require.Nil(t, err)

		err = rep.Mount(id, inVolume)

		require.Nil(t, err)

		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

		require.Nil(t, err)

		err = serv.Unmount(id, inVolume.Name)

		require.Error(t, err)
	})
}

func TestPath(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	ftpmngr := ftpMock.NewFTPManager(t)
	mountmngr := mountMock.NewMountManager(t)
	statemngr := stateMock.NewStateManager(t)
	rep := repository.CreateInMemoryRepository(logger)
	mountpoint := "/test"

	inVolume := &volume.Volume{
		Name:       "test",
		Mountpoint: "/test/abc",
		Status:     make(map[string]interface{}),
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
	}

	err := rep.Create(inVolume, &models.VolumeOptions{})

	require.Nil(t, err)

	statemngr.On("SyncState").Return(nil)

	t.Run("succsess get path", func(t *testing.T) {
		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		got, err := serv.Path("test")

		require.Nil(t, err)

		assert.Equal(t, inVolume.Mountpoint, got)
	})

	t.Run("failed get path", func(t *testing.T) {
		serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)
		require.Nil(t, err)

		_, err = serv.Path("test1")

		assert.Error(t, err)
	})
}

func TestCapabilities(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	ftpmngr := ftpMock.NewFTPManager(t)
	mountmngr := mountMock.NewMountManager(t)
	statemngr := stateMock.NewStateManager(t)
	rep := repository.CreateInMemoryRepository(logger)
	mountpoint := "/test"

	statemngr.On("SyncState").Return(nil).Once()

	serv, err := CreateFTPService(mountpoint, ftpmngr, mountmngr, statemngr, rep, logger)

	require.Nil(t, err)

	expected := "local"

	got := serv.Capabilities()
	assert.Equal(t, expected, got.Scope)
}
