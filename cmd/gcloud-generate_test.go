package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_mergeEnvs(t *testing.T) {
	t.Run("merge simple envs", func(t *testing.T) {
		result := mergeEnvs(map[string]string{
			"key":  "value",
			"key1": "oldValue",
		}, map[string]string{
			"key1": "newValue",
			"key2": "value",
		})

		assert.Equal(t, map[string]string{
			"key":  "value",
			"key1": "newValue",
			"key2": "value",
		}, result)
	})

	t.Run("merge additional envs containing base envs", func(t *testing.T) {
		result := mergeEnvs(map[string]string{
			"key":  "someValue",
			"key1": "oldValue",
		}, map[string]string{
			"key1": "$key",
			"key2": "$key",
		})

		assert.Equal(t, map[string]string{
			"key":  "someValue",
			"key1": "someValue",
			"key2": "someValue",
		}, result)
	})

	t.Run("merge additional envs containing process/os envs", func(t *testing.T) {
		assert.NoError(t, os.Setenv("SOME_KEY", "someEnvValue"))

		result := mergeEnvs(map[string]string{
			"key": "someValue",
		}, map[string]string{
			"key1": "$SOME_KEY$SOME_KEY",
			"key2": "$key$SOME_KEY",
		})

		assert.Equal(t, map[string]string{
			"key":  "someValue",
			"key1": "someEnvValuesomeEnvValue",
			"key2": "someValuesomeEnvValue",
		}, result)
	})
}
