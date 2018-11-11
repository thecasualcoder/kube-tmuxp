package kubeconfig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/arunvelsriram/kube-tmuxp/pkg/internal/mock"
	"github.com/arunvelsriram/kube-tmuxp/pkg/kubeconfig"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should create kubeconfig", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCmdr := mock.NewCommander(ctrl)
		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)

		kubeCfg, err := kubeconfig.New(mockFS, mockCmdr)

		assert.Nil(t, err)
		assert.NotNil(t, kubeCfg)
	})

	t.Run("should return error if home dir cannot be determined", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCmdr := mock.NewCommander(ctrl)
		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("", fmt.Errorf("some error"))
		_, err := kubeconfig.New(mockFS, mockCmdr)

		assert.EqualError(t, err, "some error")
	})
}

func TestDelete(t *testing.T) {
	t.Run("should delete given kubeconfig file from filesystem", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCmdr := mock.NewCommander(ctrl)
		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		mockFS.EXPECT().Remove("/Users/test/.kube/configs/context-name").Return(nil)

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		err := kubeCfg.Delete("/Users/test/.kube/configs/context-name")

		assert.Nil(t, err)
	})

	t.Run("should skip error arises when deleting kubeconfig file that does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCmdr := mock.NewCommander(ctrl)
		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		mockFS.EXPECT().Remove("/Users/test/.kube/configs/context-name").Return(&os.PathError{Op: "remove", Path: "/Users/test/.kube/configs/context-name", Err: os.ErrNotExist})

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		err := kubeCfg.Delete("/Users/test/.kube/configs/context-name")

		assert.Nil(t, err)
	})

	t.Run("should return error if error occurs when deleting the kubeconfig file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCmdr := mock.NewCommander(ctrl)
		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		mockFS.EXPECT().Remove("/Users/test/.kube/configs/context-name").Return(fmt.Errorf("some error"))

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		err := kubeCfg.Delete("/Users/test/.kube/configs/context-name")

		assert.EqualError(t, err, "some error")
	})
}

func TestAddRegionalCluster(t *testing.T) {
	t.Run("should invoke command for adding regional cluster", func(*testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)

		mockCmdr := mock.NewCommander(ctrl)
		args := []string{
			"beta",
			"container",
			"clusters",
			"get-credentials",
			"test-cluster",
			"--region=test-region",
			"--project=test-project",
		}
		envs := []string{
			"CLOUDSDK_CONTAINER_USE_V1_API_CLIENT=false",
			"CLOUDSDK_CONTAINER_USE_V1_API=false",
			"KUBECONFIG=/Users/test/.kube/configs/test-context",
		}
		mockCmdr.EXPECT().Execute("gcloud", args, envs).Return("Context added successfully", nil)

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		err := kubeCfg.AddRegionalCluster("test-project", "test-cluster", "test-region", "/Users/test/.kube/configs/test-context")

		assert.Nil(t, err)
	})

	t.Run("should return error if command failed to execute", func(*testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)

		mockCmdr := mock.NewCommander(ctrl)
		args := []string{
			"beta",
			"container",
			"clusters",
			"get-credentials",
			"test-cluster",
			"--region=test-region",
			"--project=test-project",
		}
		envs := []string{
			"CLOUDSDK_CONTAINER_USE_V1_API_CLIENT=false",
			"CLOUDSDK_CONTAINER_USE_V1_API=false",
			"KUBECONFIG=/Users/test/.kube/configs/test-context",
		}
		mockCmdr.EXPECT().Execute("gcloud", args, envs).Return("", fmt.Errorf("some error"))

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		err := kubeCfg.AddRegionalCluster("test-project", "test-cluster", "test-region", "/Users/test/.kube/configs/test-context")

		assert.EqualError(t, err, "some error")
	})
}

func TestAddZonalCluster(t *testing.T) {
	t.Run("should invoke command for adding zonal cluster", func(*testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)

		mockCmdr := mock.NewCommander(ctrl)
		args := []string{
			"container",
			"clusters",
			"get-credentials",
			"test-cluster",
			"--zone=test-zone",
			"--project=test-project",
		}
		envs := []string{
			"KUBECONFIG=/Users/test/.kube/configs/test-context",
		}
		mockCmdr.EXPECT().Execute("gcloud", args, envs).Return("Context added successfully", nil)

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		err := kubeCfg.AddZonalCluster("test-project", "test-cluster", "test-zone", "/Users/test/.kube/configs/test-context")

		assert.Nil(t, err)
	})

	t.Run("should return error if command failed to execute", func(*testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)

		mockCmdr := mock.NewCommander(ctrl)
		args := []string{
			"container",
			"clusters",
			"get-credentials",
			"test-cluster",
			"--zone=test-zone",
			"--project=test-project",
		}
		envs := []string{
			"KUBECONFIG=/Users/test/.kube/configs/test-context",
		}
		mockCmdr.EXPECT().Execute("gcloud", args, envs).Return("", fmt.Errorf("some error"))

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		err := kubeCfg.AddZonalCluster("test-project", "test-cluster", "test-zone", "/Users/test/.kube/configs/test-context")

		assert.EqualError(t, err, "some error")
	})
}

func TestRenameContext(t *testing.T) {
	t.Run("should rename a context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)

		mockCmdr := mock.NewCommander(ctrl)
		args := []string{
			"config",
			"rename-context",
			"old-context-name",
			"new-context-name",
		}
		envs := []string{
			"KUBECONFIG=/Users/test/.kube/configs/new-context-name",
		}
		mockCmdr.EXPECT().Execute("kubectl", args, envs).Return("Context renamed", nil)

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		err := kubeCfg.RenameContext("old-context-name", "new-context-name", "/Users/test/.kube/configs/new-context-name")

		assert.Nil(t, err)
	})

	t.Run("should return error if renaming context fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)

		mockCmdr := mock.NewCommander(ctrl)
		args := []string{
			"config",
			"rename-context",
			"old-context-name",
			"new-context-name",
		}
		envs := []string{
			"KUBECONFIG=/Users/test/.kube/configs/new-context-name",
		}
		mockCmdr.EXPECT().Execute("kubectl", args, envs).Return("", fmt.Errorf("some error"))

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		err := kubeCfg.RenameContext("old-context-name", "new-context-name", "/Users/test/.kube/configs/new-context-name")

		assert.EqualError(t, err, "some error")
	})
}

func TestKubeCfgsDir(t *testing.T) {
	t.Run("should return the directory in which kube configs are stored", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)

		mockCmdr := mock.NewCommander(ctrl)

		kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
		dir := kubeCfg.KubeCfgsDir()

		assert.Equal(t, "/Users/test/.kube/configs", dir)
	})
}
