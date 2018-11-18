package tmuxp_test

import (
	"fmt"
	"testing"

	"github.com/arunvelsriram/kube-tmuxp/pkg/internal/mock"
	"github.com/arunvelsriram/kube-tmuxp/pkg/tmuxp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should create a tmuxp config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		tmuxpCfg, err := tmuxp.New("session", tmuxp.Windows{}, tmuxp.Environment{}, mockFS)

		assert.Nil(t, err)
		assert.NotNil(t, tmuxpCfg)
	})

	t.Run("should return error in home dir cannot be determined", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("", fmt.Errorf("some error"))
		_, err := tmuxp.New("session", tmuxp.Windows{}, tmuxp.Environment{}, mockFS)

		assert.EqualError(t, err, "some error")
	})
}

func TestTmuxpConfigsDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFS := mock.NewFileSystem(ctrl)
	mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
	tmuxpCfg, _ := tmuxp.New("session", tmuxp.Windows{}, tmuxp.Environment{}, mockFS)

	assert.Equal(t, "/Users/test/.tmuxp", tmuxpCfg.TmuxpConfigsDir())
}
