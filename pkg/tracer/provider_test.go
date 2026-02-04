package tracer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/shared/config"
)

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.TracerConfig
		wantErr bool
	}{
		{
			name: "valid config with tracing enabled",
			cfg: &config.TracerConfig{
				Enabled:      true,
				Endpoint:     "http://localhost:4318/v1/traces",
				ServiceName:  "test-service",
				SamplingRate: 1.0,
			},
			wantErr: false,
		},
		{
			name: "tracing disabled",
			cfg: &config.TracerConfig{
				Enabled:     false,
				ServiceName: "test-service",
			},
			wantErr: false,
		},
		{
			name: "empty service name",
			cfg: &config.TracerConfig{
				Enabled:      true,
				Endpoint:     "http://localhost:4318/v1/traces",
				ServiceName:  "",
				SamplingRate: 1.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewProvider(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, provider)

			// Clean up
			if provider != nil {
				_ = provider.Shutdown(context.Background())
			}
		})
	}
}

func TestProvider_Tracer(t *testing.T) {
	cfg := &config.TracerConfig{
		Enabled:      true,
		Endpoint:     "http://localhost:4318/v1/traces",
		ServiceName:  "test-service",
		SamplingRate: 1.0,
	}

	provider, err := NewProvider(cfg)
	require.NoError(t, err)
	defer func() { _ = provider.Shutdown(context.Background()) }()

	tracer := provider.Tracer("test-tracer")
	assert.NotNil(t, tracer)
}

func TestProvider_Shutdown(t *testing.T) {
	cfg := &config.TracerConfig{
		Enabled:      true,
		Endpoint:     "http://localhost:4318/v1/traces",
		ServiceName:  "test-service",
		SamplingRate: 1.0,
	}

	provider, err := NewProvider(cfg)
	require.NoError(t, err)

	err = provider.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestNoopProvider(t *testing.T) {
	cfg := &config.TracerConfig{
		Enabled:     false,
		ServiceName: "test-service",
	}

	provider, err := NewProvider(cfg)
	require.NoError(t, err)
	assert.NotNil(t, provider)

	tracer := provider.Tracer("test-tracer")
	assert.NotNil(t, tracer)

	err = provider.Shutdown(context.Background())
	assert.NoError(t, err)
}
