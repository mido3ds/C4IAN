package constants

import "time"

const (
	// The age of entries in global and zone flooding tables
	FloodingTableEntryAge = 60
	// Time allowed for sARP responses to arrive and neighborhood table to be updated
	SARPHoldTime = 500 * time.Millisecond
	// Time between consequent sARP requests (neighborhood discoveries)
	SARPDelay = 1 * time.Second
	// The age of entries in the destination zone discovery (DZD) cache
	DZCacheAge = 5 * time.Second
	// The delay between consequent trials to find destination zone
	DZDRetryTimeout = 3 * time.Second
	// The maximum number of attempts to find the destination zone
	DZDMaxRetry = 2

	// TODO: Add constants in odmrp & tables
)
