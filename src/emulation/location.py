#!/bin/env python3
import socket
import json
import math

EARTH_RAD = 6371 * 1000
CAIRO_UNI_LONLAT = (31.2108959, 30.0272552)


def _to_degrees(x):
    # x in meters, output in degrees (lon or lat)
    sign = +1 if x > 0 else -1
    x *= sign
    d = (180/math.pi) * (x/EARTH_RAD) % 360
    if d > 180:
        sign = -1*sign
        d = 360 - d
    return d*sign


def to_gps_coords(x, y):
    lon = _to_degrees(x) + CAIRO_UNI_LONLAT[0]
    lat = _to_degrees(y) + CAIRO_UNI_LONLAT[1]
    return lon, lat


def send_location(sock, lon, lat):
    try:
        client = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)
        client.connect(sock)
        client.send(json.dumps({'lon': lon, 'lat': lat}).encode('ASCII'))
    except:
        pass
