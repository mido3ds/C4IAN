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

	// ODMRP
	// Time to live of ODMRP packet
	ODMRPDefaultTTL = 100
	// The age of entries in the members table
	MembersTableTimeout = 2 * time.Second // 960 * time.Millisecond
	// The age of entries in the forward table
	ForwardTableTimeout = 2 * time.Second // 960 * time.Millisecond
	// The age of entries of ODMRP cache
	ODMRPCacheTimeout = 2 * time.Second // 960 * time.Millisecond
	// The age of entries in the multi forward table
	MultiForwardTableTimeout = 2 * time.Second // 960 * time.Millisecond
	// The delay between consequent joinquery to maintain the multicast graph
	JQRefreshTime = 400 * time.Millisecond
	// The Timeout to fill the forward table to start sending the message
	FillForwardTableTimeout = 5 * time.Second
)
