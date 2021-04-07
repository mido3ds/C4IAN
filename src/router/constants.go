package main

const (
	MULTICAST_PATTERN = "^(?:2[23][4-9])\\." +
		"(?:(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.){2}" +
		"(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])$"
	BROADCASTPATTERN     = "255.255.255.255"
	DEFAULT_TIME_TO_LIVE = 50
)
