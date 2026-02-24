package mqtt

import (
	"context"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "localhost", config.Broker)
	assert.Equal(t, 1883, config.Port)
	assert.Equal(t, "iot-gateway", config.ClientID)
	assert.False(t, config.CleanSession)
	assert.True(t, config.AutoReconnect)
	assert.Equal(t, 30*time.Second, config.ConnectTimeout)
	assert.Equal(t, 60*time.Second, config.KeepAlive)
	assert.Equal(t, byte(1), config.QoS)
}

func TestNewClient(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		client := NewClient(nil, nil)
		require.NotNil(t, client)
		assert.NotNil(t, client.config)
		assert.Equal(t, "localhost", client.config.Broker)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &Config{
			Broker:   "mqtt.example.com",
			Port:     8883,
			ClientID: "test-client",
		}
		logger := logrus.NewEntry(logrus.StandardLogger())
		client := NewClient(config, logger)

		require.NotNil(t, client)
		assert.Equal(t, "mqtt.example.com", client.config.Broker)
		assert.Equal(t, 8883, client.config.Port)
		assert.Equal(t, "test-client", client.config.ClientID)
	})
}

func TestClient_SetHandlers(t *testing.T) {
	client := NewClient(nil, nil)

	t.Run("set message handler", func(t *testing.T) {
		client.SetMessageHandler(func(c *Client, msg mqtt.Message) {
			// Handler set
		})
		assert.NotNil(t, client.messageHandler)
	})

	t.Run("set connect handler", func(t *testing.T) {
		client.SetConnectHandler(func(c *Client) {})
		assert.NotNil(t, client.connectHandler)
	})

	t.Run("set connection lost handler", func(t *testing.T) {
		client.SetConnectionLostHandler(func(c *Client, err error) {})
		assert.NotNil(t, client.lostHandler)
	})
}

func TestClient_IsConnected(t *testing.T) {
	client := NewClient(nil, nil)

	// Initially not connected
	assert.False(t, client.IsConnected())
}

func TestClient_ConnectWithoutBroker(t *testing.T) {
	config := &Config{
		Broker:         "nonexistent.broker.local",
		Port:           1883,
		ClientID:       "test-client",
		ConnectTimeout: 1 * time.Second,
		AutoReconnect:  false,
	}
	client := NewClient(config, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	assert.Error(t, err)
}

func TestClient_PublishWithoutConnection(t *testing.T) {
	client := NewClient(nil, nil)

	err := client.Publish("test/topic", 1, false, "test message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not connected")
}

func TestClient_SubscribeWithoutConnection(t *testing.T) {
	client := NewClient(nil, nil)

	err := client.Subscribe("test/topic", 1, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not connected")
}

func TestClient_UnsubscribeWithoutConnection(t *testing.T) {
	client := NewClient(nil, nil)

	err := client.Unsubscribe("test/topic")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client not connected")
}

func TestClient_GetConfig(t *testing.T) {
	config := &Config{
		Broker:   "test.broker.com",
		Port:     1883,
		ClientID: "test-client",
	}
	client := NewClient(config, nil)

	retrievedConfig := client.GetConfig()
	assert.Equal(t, config, retrievedConfig)
}
