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
