#!/bin/env python3
import socket
import json
import math

EARTH_RAD = 6371 * 1000


def _to_degrees(x):
    # x in meters, output in degrees (lon or lat)
    sign = +1 if x > 0 else -1
    x *= sign
    d = (180/math.pi) * (x/EARTH_RAD) % 360
    if d > 180:
        sign = -1*sign
        d = 360 - d
    return d*sign


def _to_cartesian(d):
    # d in degress, output in meters
    if d < 0:
        d += 360
    x = EARTH_RAD * d * math.pi/180
    return x


def to_gps_coords(x, y):
    return _to_degrees(x), _to_degrees(y)


def to_mn_coords(lon, lat):
    return _to_cartesian(lon), _to_cartesian(lat)


def send_location(sock, lon, lat):
    try:
        client = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)
        client.connect(sock)
        client.send(json.dumps({'lon': lon, 'lat': lat}).encode('ASCII'))
    except:
        pass


def kmhour_to_msec(speed: float) -> float:
    return speed * 1000 / (60*60)
