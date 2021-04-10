package odmrp

const (
	MULTICAST_PATTERN = "^(?:2[23][4-9])\\." +
		"(?:(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.){2}" +
		"(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])$"
	BROADCAST_PATTERN      = "255.255.255.255"
	DEFAULT_TIME_TO_LIVE   = 50
	MAX_NUM_PREV_HOP_PAIRS = 100
	JOIN_QUERY_SIZE        = 16
	JOIN_REPLY_SIZE        = 16
	PACKET_BUFFER_SIZE     = 20
)
