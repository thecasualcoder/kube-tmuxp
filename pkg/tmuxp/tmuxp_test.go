package tmuxp_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/arunvelsriram/kube-tmuxp/pkg/internal/mock"
	"github.com/arunvelsriram/kube-tmuxp/pkg/tmuxp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("should create a tmuxp config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		tmuxpCfg, err := tmuxp.NewConfig("session", tmuxp.Windows{}, tmuxp.Environment{}, mockFS)

		assert.Nil(t, err)
		assert.NotNil(t, tmuxpCfg)
	})

	t.Run("should return error in home dir cannot be determined", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("", fmt.Errorf("some error"))
		_, err := tmuxp.NewConfig("session", tmuxp.Windows{}, tmuxp.Environment{}, mockFS)

		assert.EqualError(t, err, "some error")
	})
}

func TestTmuxpConfigsDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFS := mock.NewFileSystem(ctrl)
	mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
	tmuxpCfg, _ := tmuxp.NewConfig("session", tmuxp.Windows{}, tmuxp.Environment{}, mockFS)

	assert.Equal(t, "/Users/test/.tmuxp", tmuxpCfg.TmuxpConfigsDir())
}

func TestSave(t *testing.T) {
	t.Run("should save the tmuxp config as file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		var writer bytes.Buffer
		mockFS.EXPECT().Create("tmuxp-config.yaml").Return(&writer, nil)
		tmuxpCfg, _ := tmuxp.NewConfig("session", tmuxp.Windows{{Name: "window"}}, tmuxp.Environment{"TEST_ENV": "value", "ANOTHER_TEST_ENV": "another-value"}, mockFS)

		err := tmuxpCfg.Save("tmuxp-config.yaml")

		actualContent := writer.String()
		expectedContent := `session_name: session
windows:
- window_name: window
  panes: []
environment:
  ANOTHER_TEST_ENV: another-value
  TEST_ENV: value
`
		assert.Nil(t, err)
		assert.Equal(t, expectedContent, actualContent)
	})

	t.Run("should return error if tmuxp config cannot be saved", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
		mockFS.EXPECT().Create("tmuxp-config.yaml").Return(nil, fmt.Errorf("some error"))
		tmuxpCfg, _ := tmuxp.NewConfig("session", tmuxp.Windows{{Name: "window"}}, tmuxp.Environment{"TEST_ENV": "value", "ANOTHER_TEST_ENV": "another-value"}, mockFS)

		err := tmuxpCfg.Save("tmuxp-config.yaml")

		assert.EqualError(t, err, "some error")
	})
}
