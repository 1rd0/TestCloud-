package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestNewConfig_Defaults(t *testing.T) {
	t.Run("Should return default values", func(t *testing.T) {
		cfg, err := New("")
		require.NoError(t, err)

		assert.Equal(t, ":8080", cfg.Listen)
		assert.Equal(t, 100, cfg.Rate.Capacity)
		assert.Equal(t, 10, cfg.Rate.RPS)
		assert.Empty(t, cfg.LB.Backends)

	})
}

func TestNewConfig_FromYAML(t *testing.T) {
	t.Run("Should load config from YAML file", func(t *testing.T) {
		// Получаем абсолютный путь к тестовому файлу

		cfg, err := New("configTest.yaml")
		require.NoError(t, err)

		assert.Equal(t, ":8080", cfg.Listen)
		assert.Equal(t, []string{"http://localhost:9001", "http://localhost:9002"}, cfg.LB.Backends)
		assert.Equal(t, 100, cfg.Rate.Capacity)
		assert.Equal(t, 10, cfg.Rate.RPS)
		assert.Equal(t, 5*time.Second, cfg.Health.Interval)
	})

	t.Run("Should return error for invalid YAML", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "config-*.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString("invalid: yaml: [")
		require.NoError(t, err)
		tmpFile.Close()

		_, err = New(tmpFile.Name())
		assert.Error(t, err)
	})
}
