package constants

import "time"

const (
	// The age of entries in global and zone flooding tables
	FloodingTableEntryAge = 60
	// Time allowed for sARP responses to arrive and neighborhood table to be updated
	SARPHoldTime = 500 * time.Millisecond
	// Time between consequent sARP requests (neighborhood discoveries)
	SARPDelay = 1 * time.Second
	// Time between sending interzone LSR packets
	InterzoneLSRDelay = 1 * time.Second
	// The age of entries in the destination zone discovery (DZD) cache
	DZCacheAge = 60 * time.Second
	// The delay between consequent trials to find destination zone
	DZDRetryTimeout = 3 * time.Second
	// The maximum number of attempts to find the destination zone
	DZDMaxRetry = 5

	// ODMRP
	// Time to live of ODMRP packet
	ODMRPDefaultTTL = 100
	// The age of entries in the members table
	MembersTableTimeout = 2500 * time.Millisecond
	// The age of entries in the forward table
	ForwardTableTimeout = 2500 * time.Millisecond
	// The age of entries of ODMRP cache
	ODMRPCacheTimeout = 2500 * time.Millisecond
	// The age of entries in the multi forward table
	MultiForwardTableTimeout = 2500 * time.Millisecond
	// The delay between consequent joinquery to maintain the multicast graph
	JQRefreshTime = 500 * time.Millisecond
	// The Timeout to fill the forward table to start sending the message
	FillForwardTableTimeout = 3 * time.Second
	// the multi forward table timeout value must be larger (e.g., 3 to 5 times) than the value of route refresh interval.
)
