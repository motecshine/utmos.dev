package aircraft

// RCOSD represents Remote Controller OSD (Operational Status Data)
// Push frequency: 0.5Hz
// Topic: thing/product/{device_sn}/osd
type RCOSD struct {
	// Basic information
	CapacityPercent *int     `json:"capacity_percent"` // Battery percentage (0-100%)
	Height          *float64 `json:"height"`           // Ellipsoid height in meters
	Latitude        *float64 `json:"latitude"`         // Latitude (-90 to 90)
	Longitude       *float64 `json:"longitude"`        // Longitude (-180 to 180)
	Country         *string  `json:"country"`          // Country code

	// Wireless link information
	WirelessLink *RCWirelessLink `json:"wireless_link"`
}

// RCWirelessLink represents the wireless link information for RC
type RCWirelessLink struct {
	DongleNumber  *int     `json:"dongle_number"`   // Number of dongles on aircraft
	Link4GState   *int     `json:"4g_link_state"`   // 4G link state: 0=disconnected, 1=connected
	SDRLinkState  *int     `json:"sdr_link_state"`  // SDR link state: 0=disconnected, 1=connected
	LinkWorkmode  *int     `json:"link_workmode"`   // Link work mode: 0=SDR, 1=4G fusion
	SDRQuality    *int     `json:"sdr_quality"`     // SDR signal quality (0-5)
	Link4GQuality *int     `json:"4g_quality"`      // Overall 4G signal quality (0-5)
	UAV4GQuality  *int     `json:"4g_uav_quality"`  // UAV-side 4G signal quality (0-5)
	GND4GQuality  *int     `json:"4g_gnd_quality"`  // Ground-side 4G signal quality (0-5)
	SDRFreqBand   *float64 `json:"sdr_freq_band"`   // SDR frequency band
	Link4GFreqBand *float64 `json:"4g_freq_band"`   // 4G frequency band
}
