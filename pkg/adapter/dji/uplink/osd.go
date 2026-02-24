package uplink

// OSDData represents DJI OSD telemetry data
type OSDData struct {
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Altitude        float64 `json:"altitude"`
	Height          float64 `json:"height"`
	Speed           float64 `json:"speed"`
	Heading         float64 `json:"heading"`
	Pitch           float64 `json:"pitch"`
	Roll            float64 `json:"roll"`
	Yaw             float64 `json:"yaw"`
	BatteryPercent  int     `json:"battery_percent"`
	FlightMode      string  `json:"flight_mode"`
	GPSSatellites   int     `json:"gps_satellites"`
	SignalStrength  int     `json:"signal_strength"`
	HomeDistance    float64 `json:"home_distance"`
	HorizontalSpeed float64 `json:"horizontal_speed"`
	VerticalSpeed   float64 `json:"vertical_speed"`
	WindSpeed       float64 `json:"wind_speed"`
	WindDirection   float64 `json:"wind_direction"`
}

// ParseOSDData parses DJI OSD data from properties
func ParseOSDData(properties map[string]any) (*OSDData, error) {
	data := &OSDData{}

	if v, ok := properties["latitude"].(float64); ok {
		data.Latitude = v
	}
	if v, ok := properties["longitude"].(float64); ok {
		data.Longitude = v
	}
	if v, ok := properties["altitude"].(float64); ok {
		data.Altitude = v
	}
	if v, ok := properties["height"].(float64); ok {
		data.Height = v
	}
	if v, ok := properties["speed"].(float64); ok {
		data.Speed = v
	}
	if v, ok := properties["heading"].(float64); ok {
		data.Heading = v
	}
	if v, ok := properties["pitch"].(float64); ok {
		data.Pitch = v
	}
	if v, ok := properties["roll"].(float64); ok {
		data.Roll = v
	}
	if v, ok := properties["yaw"].(float64); ok {
		data.Yaw = v
	}
	if v, ok := properties["battery_percent"].(float64); ok {
		data.BatteryPercent = int(v)
	}
	if v, ok := properties["flight_mode"].(string); ok {
		data.FlightMode = v
	}
	if v, ok := properties["gps_satellites"].(float64); ok {
		data.GPSSatellites = int(v)
	}
	if v, ok := properties["signal_strength"].(float64); ok {
		data.SignalStrength = int(v)
	}
	if v, ok := properties["home_distance"].(float64); ok {
		data.HomeDistance = v
	}
	if v, ok := properties["horizontal_speed"].(float64); ok {
		data.HorizontalSpeed = v
	}
	if v, ok := properties["vertical_speed"].(float64); ok {
		data.VerticalSpeed = v
	}
	if v, ok := properties["wind_speed"].(float64); ok {
		data.WindSpeed = v
	}
	if v, ok := properties["wind_direction"].(float64); ok {
		data.WindDirection = v
	}

	return data, nil
}

