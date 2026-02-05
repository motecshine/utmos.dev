package wpml

// createValidWaylines creates a fully valid Waylines object for testing
func createValidWaylines(name string) *Waylines {
	return &Waylines{
		Name:                    name,
		Description:             "Test Description",
		DroneModel:              DroneM3Series,
		PayloadModel:            PayloadMatrice3TD,
		TemplateType:            TemplateTypeWaypoint,
		GlobalHeight:            50.0,
		GlobalSpeed:             15.0,
		ClimbMode:               "vertical",
		TakeOffSecurityHeight:   120.0,
		GlobalRTHHeight:         100.0,
		AircraftYawMode:         "followWayline",
		GimbalPitchMode:         "usePointSetting",
		GlobalTransitionalSpeed: 10.0,
		Waypoints: []WaylinesWaypoint{
			{
				Latitude:    39.9093,
				Longitude:   116.3974,
				Height:      50.0,
				Speed:       15.0,
				TriggerType: TriggerTypeReachPoint,
			},
		},
	}
}
