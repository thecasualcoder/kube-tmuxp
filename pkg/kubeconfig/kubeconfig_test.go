package kubeconfig_test

import (
	"fmt"
	"testing"

	"github.com/arunvelsriram/kube-tmuxp/pkg/filesystem"
	"github.com/arunvelsriram/kube-tmuxp/pkg/internal/mock"
	"github.com/arunvelsriram/kube-tmuxp/pkg/kubeconfig"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var fs filesystem.FileSystem
	cfg := kubeconfig.New(fs)

	assert.NotNil(t, cfg)
}

func TestDelete(t *testing.T) {
	t.Run("should delete given kubeconfig file from filesystem", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().Remove("/Users/arunvelsriram/.kube/configs/context-name").Return(nil)

		cfg := kubeconfig.New(mockFS)
		err := cfg.Delete("context-name")

		assert.Nil(t, err)
	})

	t.Run("should return error if kubeconfig file cannot be deleted", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().Remove("/Users/arunvelsriram/.kube/configs/context-name").Return(fmt.Errorf("some error"))

		cfg := kubeconfig.New(mockFS)
		err := cfg.Delete("context-name")

		assert.EqualError(t, err, "some error")
	})
}
