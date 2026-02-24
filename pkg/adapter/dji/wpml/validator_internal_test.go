package wpml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests aim to trigger internal validator functions that are currently at 0% coverage

func TestValidator_InternalFunctions(t *testing.T) {
	validator, _ := NewValidator()

	// Test isRequiredForDrone by creating a validation scenario
	t.Run("Test drone validation patterns", func(t *testing.T) {
		// Create a test struct that would use drone validation
		result := validator.isRequiredForDrone("M3|M30", DroneM3Series)
		assert.True(t, result)

		result = validator.isRequiredForDrone("M300", DroneM3Series)
		assert.False(t, result)
	})

	t.Run("Test payload validation patterns", func(t *testing.T) {
		// Test isRequiredForPayload
		result := validator.isRequiredForPayload("H20|H30", PayloadH20)
		assert.True(t, result)

		result = validator.isRequiredForPayload("M30", PayloadH20)
		assert.False(t, result)
	})

	t.Run("Test drone pattern matching", func(t *testing.T) {
		// Test matchesDronePattern
		result := validator.matchesDronePattern("M3", DroneM3Series)
		assert.True(t, result)

		result = validator.matchesDronePattern("M300", DroneM3Series)
		assert.False(t, result)

		result = validator.matchesDronePattern("M300", DroneM300RTK)
		assert.True(t, result)

		result = validator.matchesDronePattern("UNKNOWN", DroneM3Series)
		assert.False(t, result)
	})

	t.Run("Test payload pattern matching", func(t *testing.T) {
		// Test matchesPayloadPattern
		result := validator.matchesPayloadPattern("H20", PayloadH20)
		assert.True(t, result)

		result2 := validator.matchesPayloadPattern("H20", PayloadH20T)
		assert.True(t, result2)

		result = validator.matchesPayloadPattern("H30", PayloadH20)
		assert.False(t, result)

		result = validator.matchesPayloadPattern("UNKNOWN", PayloadH20)
		assert.False(t, result)
	})
}

func TestValidator_RequiredForFields(t *testing.T) {
	validator, _ := NewValidator()

	// Call validateRequiredForDrone and validateRequiredForPayload indirectly
	// by creating a validation context - these are typically called during struct validation

	// These functions return true for now but we want to ensure they're called
	t.Run("validateRequiredForDrone coverage", func(t *testing.T) {
		// The functions are simple and return true, but we need to ensure they're covered
		// They're usually called as part of field validation logic

		// Create a scenario that might trigger these validations
		waylines := createValidWaylines("Validation Test")
		err := validator.ValidateStruct(waylines)
		// The validation should pass with our valid waylines
		assert.NoError(t, err)
	})
}

func TestCreateActionFromType_MoreCoverage(t *testing.T) {
	// Test CreateActionFromType function with more action types to improve coverage

	tests := []struct {
		name       string
		actionType string
		expected   ActionInterface
	}{
		{
			name:       "GimbalAngleLock",
			actionType: ActionTypeGimbalAngleLock,
			expected:   &GimbalAngleLockAction{},
		},
		{
			name:       "StartTimeLapse",
			actionType: ActionTypeStartTimeLapse,
			expected:   &StartTimeLapseAction{},
		},
		{
			name:       "StopTimeLapse",
			actionType: ActionTypeStopTimeLapse,
			expected:   &StopTimeLapseAction{},
		},
		{
			name:       "Unknown type",
			actionType: "unknownAction",
			expected:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateActionFromType(tt.actionType)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.IsType(t, tt.expected, result)
			}
		})
	}
}
