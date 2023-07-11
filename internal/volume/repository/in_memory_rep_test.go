package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/t1d333/docker-volume-ftp-driver/internal/models"
)

var GetTestsSimple = map[string]struct {
	name     string
	in       volume.Volume
	expected volume.Volume
}{
	"Simple get negative test": {
		"test",
		volume.Volume{Name: "test", Mountpoint: "/test", CreatedAt: time.Date(2023, time.July, 10, 12, 12, 10, 10, time.UTC).Format(time.RFC3339Nano)},
		volume.Volume{Name: "test", Mountpoint: "/test", CreatedAt: time.Date(2023, time.July, 10, 12, 12, 10, 10, time.UTC).Format(time.RFC3339Nano)},
	},
}

var GetTestsNegative = map[string]struct {
	name string
	in   *volume.Volume
}{
	"Simple get negative test": {
		"test1",
		&volume.Volume{Name: "test", Mountpoint: "/test", CreatedAt: time.Date(2023, time.July, 10, 12, 12, 10, 10, time.UTC).Format(time.RFC3339Nano)},
	},
}

var CreateTestsNegative = map[string]struct {
	vol *volume.Volume
	opt *models.VolumeOptions
}{
	"Create volume with nil vol": {
		nil,
		&models.VolumeOptions{},
	},

	"Create volume with nil opt": {
		&volume.Volume{Name: "test"},
		nil,
	},
}

var ListTests = map[string]struct {
	in []*volume.Volume
}{
	"Empty list": {
		in: []*volume.Volume{},
	},

	"One item": {
		in: []*volume.Volume{
			{Name: "test", Mountpoint: "/test", CreatedAt: time.Date(2023, time.July, 10, 12, 12, 10, 10, time.UTC).Format(time.RFC3339Nano)},
		},
	},

	"Some items": {
		in: []*volume.Volume{
			{Name: "test", Mountpoint: "/test", CreatedAt: time.Date(2023, time.July, 10, 12, 12, 10, 10, time.UTC).Format(time.RFC3339Nano)},

			{Name: "test2", Mountpoint: "/test", CreatedAt: time.Date(2023, time.July, 10, 12, 12, 10, 10, time.UTC).Format(time.RFC3339Nano)},
		},
	},
}

func TestGetSimple(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	rep := CreateInMemoryRepository(logger)
	for name, test := range GetTestsSimple {
		err := rep.Create(&test.expected, &models.VolumeOptions{})
		assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: %s \n Input: %v \n Expected: %v \n Error: %v", name, test.in, test.expected, err))

		got, err := rep.Get(test.name)

		assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: %s \n Input: %v \n Expected: %v \n Error: %v", name, test.in, test.expected, err))

		assert.Equal(t, test.expected.Name, got.Name)
		assert.Equal(t, test.expected.Mountpoint, got.Mountpoint)
		assert.Equal(t, test.expected.CreatedAt, got.CreatedAt)
		assert.Equal(t, test.expected.Status, got.Status)
	}
}

func TestGetNegative(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	rep := CreateInMemoryRepository(logger)
	for name, test := range GetTestsNegative {
		err := rep.Create(test.in, &models.VolumeOptions{})
		assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: %s \n Input: %v \n Error: %v", name, test.in, err))

		_, err = rep.Get(test.name)
		assert.Error(t, err, "Get volume with not existing name")
	}
}

func TestCreateNegative(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	rep := CreateInMemoryRepository(logger)
	for name, test := range CreateTestsNegative {
		err := rep.Create(test.vol, test.opt)
		assert.Error(t, err, name)
	}
}

func TestCreateExistsVol(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	rep := CreateInMemoryRepository(logger)
	vol := &volume.Volume{Name: "test"}
	err := rep.Create(vol, &models.VolumeOptions{})
	assert.Nil(t, err, "Unexpected error on test: create exists volume")
	assert.Error(t, rep.Create(vol, &models.VolumeOptions{}))
}

func TestList(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	for name, test := range ListTests {
		rep := CreateInMemoryRepository(logger)
		for _, volume := range test.in {
			err := rep.Create(volume, &models.VolumeOptions{})
			assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: %s \n Input: %v \n Error: %v", name, test.in, err))
		}

		list, err := rep.List()

		assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: %s \n Input: %v \n Error: %v", name, test.in, err))

		assert.Equal(t, len(test.in), len(list))

		for _, volume := range test.in {
			flag := false
			for _, listItem := range list {
				if listItem == volume {
					flag = true
					break
				}
			}
			assert.True(t, flag)
		}
	}
}

func TestGetOptionsSimpe(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	rep := CreateInMemoryRepository(logger)

	volume := &volume.Volume{Name: "test"}
	ftpOpt := models.FTPConnectionOpt{User: "admin", Host: "localhost", Port: 21, Password: "admin"}
	options := &models.VolumeOptions{RemotePath: "/test", FTPConnectionOpt: ftpOpt}

	err := rep.Create(volume, options)
	assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: get options simple \n Input: {volume: %v, \n options: %v } \n Error: %v", *volume, *options, err))

	got := rep.GetVolumeOptions(volume.Name)
	assert.NotNil(t, got, "Not found options for volume")

	assert.Equal(t, options.RemotePath, got.RemotePath)
	assert.Equal(t, options.User, got.User)
	assert.Equal(t, options.Host, got.Host)
	assert.Equal(t, options.Password, got.Password)
	assert.Equal(t, options.Port, got.Port)
}

func TestPathSimple(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	rep := CreateInMemoryRepository(logger)
	path := "/mnt/test"
	volume := &volume.Volume{Name: "test", Mountpoint: path}
	rep.Create(volume, &models.VolumeOptions{})
	got, err := rep.Path(volume.Name)
	assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: get path simple \n Input: {volume: %v} \n Error: %v", *volume, err))
	assert.Equal(t, path, got, "Paths not equal")
}

func TestPathNegative(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	rep := CreateInMemoryRepository(logger)
	_, err := rep.Path("abcde")
	assert.Error(t, err, "Get path for not existing volume")
}

func TestRemoveSimple(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	rep := CreateInMemoryRepository(logger)
	volume := &volume.Volume{Name: "test"}
	rep.Create(volume, &models.VolumeOptions{})
	err := rep.Remove(volume.Name)
	assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: remove simple \n Input: {volume: %v} \n Error: %v", *volume, err))
	_, err = rep.Get(volume.Name)
	assert.Error(t, err)
}

func TestRemoveNotExistingVolume(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	rep := CreateInMemoryRepository(logger)
	err := rep.Remove("not exists")
	assert.Error(t, err, "Removing not existing volume")
}

func TestRemoveMountedVolume(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	rep := CreateInMemoryRepository(logger)

	volume := &volume.Volume{Name: "test"}
	rep.Create(volume, &models.VolumeOptions{})

	rep.Mount(uuid.NewString(), volume)
	err := rep.Remove(volume.Name)

	assert.Error(t, err, "Removing mounted volume")
}

func TestMountSimple(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	rep := CreateInMemoryRepository(logger)

	vol := &volume.Volume{Name: "test"}

	rep.Create(vol, &models.VolumeOptions{})

	err := rep.Mount(uuid.NewString(), vol)
	assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: mount simple \n Input: {volume: %v} \n Error: %v", *vol, err))

	assert.True(t, rep.IsMount(vol.Name))

	volume2 := &volume.Volume{Name: "test2"}
	rep.Create(volume2, &models.VolumeOptions{})

	assert.False(t, rep.IsMount(volume2.Name))
}

func TestMountNegative(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	rep := CreateInMemoryRepository(logger)

	vol := &volume.Volume{Name: "test"}

	rep.Create(vol, &models.VolumeOptions{})

	id := uuid.NewString()
	err := rep.Mount(id, vol)
	assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: mount simple \n Input: {volume: %v} \n Error: %v", *vol, err))

	assert.True(t, rep.IsMount(vol.Name))

	assert.Error(t, rep.Mount(id, vol))
}

func TestUnmountSimple(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	rep := CreateInMemoryRepository(logger)

	vol := &volume.Volume{Name: "test"}

	rep.Create(vol, &models.VolumeOptions{})

	id := uuid.NewString()

	err := rep.Mount(id, vol)
	assert.Nil(t, err, fmt.Sprintf("Unexpected error on test: mount simple \n Input: {volume: %v} \n Error: %v", *vol, err))

	assert.True(t, rep.IsMount(vol.Name))
	assert.Nil(t, rep.Unmount(id, vol.Name))
	assert.False(t, rep.IsMount(vol.Name))
}

func TestUnmountNegative(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	rep := CreateInMemoryRepository(logger)

	vol := &volume.Volume{Name: "test"}

	rep.Create(vol, &models.VolumeOptions{})

	id := uuid.NewString()

	assert.Error(t, rep.Unmount(id, vol.Name))
	rep.Mount(id, vol)

	assert.Error(t, rep.Unmount(uuid.NewString(), vol.Name))
}

func TestGetMountedIdsList(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	rep := CreateInMemoryRepository(logger)
	vol := &volume.Volume{Name: "test"}

	rep.Create(vol, &models.VolumeOptions{})

	id := uuid.NewString()

	assert.Equal(t, 0, len(rep.GetMountedIdsList("test")))

	rep.Mount(id, vol)
	assert.Equal(t, 1, len(rep.GetMountedIdsList("test")))
	assert.Equal(t, id, rep.GetMountedIdsList("test")[0])
}
