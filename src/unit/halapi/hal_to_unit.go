package halapi

// VideoPart is sent to unit periodically
// as configured in StartVideoStream sent by unit
type VideoPart struct {
	Video []byte
}

// HeartBeat is sent to unit periodically
// after connection establishment between HAL and unit
// sent each 1 second
type HeartBeat struct {
	// from 0 (dead or no sensor) to 100
	BeatsPerMinut int
}

// Location is a gps coordinates
// which are sent to unit periodically
// after connection establishment between HAL and unit
// sent each 1 second
type Location struct {
	Lon, Lat float64
}
