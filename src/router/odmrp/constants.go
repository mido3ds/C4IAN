package odmrp

const (
	MULTICAST_PATTERN = "^(?:2[23][4-9])\\." +
		"(?:(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.){2}" +
		"(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])$"
	BROADCAST_PATTERN      = "255.255.255.255"
	DEFAULT_TIME_TO_LIVE   = 50
	MAX_NUM_PREV_HOP_PAIRS = 100

	JOIN_QUERY_SIZE             = 16
	JOIN_REPLY_SIZE             = 16
	PACKET_BUFFER_SIZE          = 20
	JOIN_SRC_HASH_SIZE          = 10
	JOIN_SOURCE_HASH_SIZE       = 10
	JOIN_HASH_SIZE              = 20
	JOIN_QUERY_REFRESH_INTERVAL = 3
	FLAG_TIMEOUT                = 3 * JOIN_QUERY_REFRESH_INTERVAL
	CACHE_SIZE                  = 500
	HASH_SIZE                   = 20
	ARP_TIMEOUT                 = 30e-3
	MAX_NUM_RET                 = 3 // max number of rrep retransmissions
	INF_NUM_RET                 = 100
	PACK_TIMEOUT                = 0.035 // passive ack timeout
	JQ_REFRESH_INTERVAL         = 3     // interval between join query floods
)

/*	ARP_TIMEOUT
	ARP broadcasts a request packet to all the machines on the LAN and asks
	if any of the machines know they are using that particular IP address.
*/

type Time float64

type IsFound int

const (
	NOT_FOUND IsFound = iota
	FOUND
	FOUND_LONGER
)
