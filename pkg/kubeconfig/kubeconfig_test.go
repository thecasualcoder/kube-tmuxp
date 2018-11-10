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

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)

		kubeCfg, err := kubeconfig.New(mockFS)

		assert.Nil(t, err)
		assert.NotNil(t, kubeCfg)
	})

	t.Run("should return error if home dir cannot be determined", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("", fmt.Errorf("some error"))
		_, err := kubeconfig.New(mockFS)

		assert.EqualError(t, err, "some error")
	})
}

func TestDelete(t *testing.T) {
	t.Run("should delete given kubeconfig file from filesystem", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		mockFS.EXPECT().Remove("/Users/test/.kube/configs/context-name").Return(nil)

		kubeCfg, _ := kubeconfig.New(mockFS)
		err := kubeCfg.Delete("context-name")

		assert.Nil(t, err)
	})

	t.Run("should skip error arises when deleting kubeconfig file that does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		mockFS.EXPECT().Remove("/Users/test/.kube/configs/context-name").Return(&os.PathError{Op: "remove", Path: "/Users/test/.kube/configs/context-name", Err: os.ErrNotExist})

		kubeCfg, _ := kubeconfig.New(mockFS)
		err := kubeCfg.Delete("context-name")

		assert.Nil(t, err)
	})

	t.Run("should return error if error occurs when deleting the kubeconfig file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		mockFS.EXPECT().Remove("/Users/test/.kube/configs/context-name").Return(fmt.Errorf("some error"))

		kubeCfg, _ := kubeconfig.New(mockFS)
		err := kubeCfg.Delete("context-name")

		assert.EqualError(t, err, "some error")
	})
}
