package volume

import (
	"errors"
	"testing"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/t1d333/docker-volume-ftp-driver/internal/volume/mocks"
	"go.uber.org/zap"
)

func TestCreate(t *testing.T) {
	mockServ := mocks.NewVolumeService(t)
	name := "test"
	options := make(map[string]string, 0)
	conf := zap.NewDevelopmentConfig()
	conf.Level.SetLevel(zap.PanicLevel)
	log, _ := conf.Build()
	logger := log.Sugar()

	t.Run("Success creation", func(t *testing.T) {
		mockServ.On("Create", name, options).Return(nil).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		err := driver.Create(&volume.CreateRequest{Name: name, Options: options})

		assert.Nil(t, err)
	})

	expectedErr := errors.New("Failed to create")

	t.Run("Failed creation", func(t *testing.T) {
		mockServ.On("Create", name, options).Return(expectedErr)
		driver := InitializeNewFTPDriver(mockServ, logger)

		err := driver.Create(&volume.CreateRequest{Name: name, Options: options})

		assert.ErrorIs(t, expectedErr, err)
	})
}

func TestGet(t *testing.T) {
	mockServ := mocks.NewVolumeService(t)
	name := "test"
	date := time.Now().Format(time.RFC3339Nano)

	conf := zap.NewDevelopmentConfig()
	conf.Level.SetLevel(zap.PanicLevel)
	log, _ := conf.Build()
	logger := log.Sugar()

	expected := &volume.Volume{
		Name:       "test",
		Mountpoint: "/test",
		CreatedAt:  date,
		Status:     make(map[string]interface{}),
	}

	t.Run("Success get", func(t *testing.T) {
		mockServ.On("Get", name).Return(expected, nil).Once()
		driver := InitializeNewFTPDriver(mockServ, logger)

		got, err := driver.Get(&volume.GetRequest{Name: "test"})

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, expected.Name, got.Volume.Name)
		assert.Equal(t, expected.CreatedAt, got.Volume.CreatedAt)
		assert.Equal(t, expected.Mountpoint, got.Volume.Mountpoint)
	})

	expectedErr := errors.New("Failed to get volume")

	t.Run("Failed get", func(t *testing.T) {
		mockServ.On("Get", mock.Anything).Return(nil, expectedErr).Once()
		driver := InitializeNewFTPDriver(mockServ, logger)

		_, err := driver.Get(&volume.GetRequest{Name: "test"})

		assert.ErrorIs(t, expectedErr, err)
	})
}

func TestList(t *testing.T) {
	mockServ := mocks.NewVolumeService(t)
	conf := zap.NewDevelopmentConfig()
	conf.Level.SetLevel(zap.PanicLevel)
	log, _ := conf.Build()
	logger := log.Sugar()

	expectedList := make([]*volume.Volume, 0)

	t.Run("Success get list", func(t *testing.T) {
		mockServ.On("List").Return(expectedList, nil).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		got, err := driver.List()

		assert.Nil(t, err)

		assert.Equal(t, len(expectedList), len(got.Volumes))
	})

	expectedErr := errors.New("Failed to get list")

	t.Run("Failed to get list", func(t *testing.T) {
		mockServ.On("List").Return(make([]*volume.Volume, 0), expectedErr).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		_, err := driver.List()

		assert.ErrorIs(t, expectedErr, err)
	})
}

func TestPath(t *testing.T) {
	mockServ := mocks.NewVolumeService(t)

	conf := zap.NewDevelopmentConfig()
	conf.Level.SetLevel(zap.PanicLevel)
	log, _ := conf.Build()
	logger := log.Sugar()

	volumeName := "test"
	expectedPath := "/test/test1"

	t.Run("Success get path", func(t *testing.T) {
		mockServ.On("Path", volumeName).Return(expectedPath, nil).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		got, err := driver.Path(&volume.PathRequest{Name: volumeName})

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, expectedPath, got.Mountpoint)
	})

	expectedErr := errors.New("Failed to get path")

	t.Run("Failed to get path", func(t *testing.T) {
		mockServ.On("Path", volumeName).Return("", expectedErr).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		_, err := driver.Path(&volume.PathRequest{Name: volumeName})

		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestRemove(t *testing.T) {
	mockServ := mocks.NewVolumeService(t)
	conf := zap.NewDevelopmentConfig()
	conf.Level.SetLevel(zap.PanicLevel)
	log, _ := conf.Build()
	logger := log.Sugar()

	volumeName := "test"

	t.Run("Success remove", func(t *testing.T) {
		mockServ.On("Remove", volumeName).Return(nil).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		err := driver.Remove(&volume.RemoveRequest{Name: volumeName})

		assert.Nil(t, err)
	})

	expectedErr := errors.New("Failed to remove")

	t.Run("Failed remove", func(t *testing.T) {
		mockServ.On("Remove", volumeName).Return(expectedErr).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		err := driver.Remove(&volume.RemoveRequest{Name: volumeName})

		assert.ErrorIs(t, expectedErr, err)
	})
}

func TestMount(t *testing.T) {
	mockServ := mocks.NewVolumeService(t)
	conf := zap.NewDevelopmentConfig()
	conf.Level.SetLevel(zap.PanicLevel)
	log, _ := conf.Build()
	logger := log.Sugar()

	volumeName := "test"
	expectedPath := "/test/test1"
	id := uuid.NewString()

	t.Run("Success mount", func(t *testing.T) {
		mockServ.On("Mount", id, volumeName).Return(expectedPath, nil).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		got, err := driver.Mount(&volume.MountRequest{Name: volumeName, ID: id})

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, expectedPath, got.Mountpoint)
	})

	expectedErr := errors.New("Failed to mount volume")

	t.Run("Failed mount", func(t *testing.T) {
		mockServ.On("Mount", id, volumeName).Return(expectedPath, expectedErr).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		_, err := driver.Mount(&volume.MountRequest{Name: volumeName, ID: id})

		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestUnmount(t *testing.T) {
	mockServ := mocks.NewVolumeService(t)
	conf := zap.NewDevelopmentConfig()
	conf.Level.SetLevel(zap.PanicLevel)
	log, _ := conf.Build()
	logger := log.Sugar()

	volumeName := "test"
	id := uuid.NewString()

	t.Run("Success unmount", func(t *testing.T) {
		mockServ.On("Unmount", id, volumeName).Return(nil).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		err := driver.Unmount(&volume.UnmountRequest{Name: volumeName, ID: id})

		assert.Nil(t, err)
	})

	expectedErr := errors.New("Failed to unmount volume")

	t.Run("Failed unmount", func(t *testing.T) {
		mockServ.On("Unmount", id, volumeName).Return(expectedErr).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		err := driver.Unmount(&volume.UnmountRequest{Name: volumeName, ID: id})

		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestCapabilities(t *testing.T) {
	mockServ := mocks.NewVolumeService(t)
	conf := zap.NewDevelopmentConfig()
	conf.Level.SetLevel(zap.PanicLevel)
	log, _ := conf.Build()
	logger := log.Sugar()

	expected := volume.Capability{Scope: "local"}

	t.Run("Success get capabilities ", func(t *testing.T) {
		mockServ.On("Capabilities").Return(expected).Once()

		driver := InitializeNewFTPDriver(mockServ, logger)

		got := driver.Capabilities()

		assert.NotNil(t, got)

		assert.Equal(t, expected.Scope, got.Capabilities.Scope)
	})
}
