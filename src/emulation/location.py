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

SCALE=95596.88361256948

def to_gps_coords(x, y):
    return _to_degrees(x), _to_degrees(y)
    # return x/SCALE, y/SCALE


def to_mn_coords(lon, lat):
    return _to_cartesian(lon), _to_cartesian(lat)
    # return lon*SCALE, lat*SCALE


def send_location(sock, lon, lat):
    try:
        client = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)
        client.connect(sock)
        client.send(json.dumps({'lon': lon, 'lat': lat}).encode('ASCII'))
    except:
        pass


def __lonlat_to_xyz(lon, lat):
    # lon, lat in degrees
    assert lon >= -180 and lon <= 180
    lon += 180
    lon %= 360
    assert lon >= 0 and lon < 360
    # [0, 360) -> [0, 0x10000)
    x = int(lon / 360.0 * 0x10000)
    assert x >= 0 and x < 0x10000

    assert lat >= -90 and lat <= 90
    lat += 90
    lat %= 180
    assert lat >= 0 and lat < 180
    # [0, 180) -> [0, 0x10000)
    y = int(lat / 180.0 * 0x10000)
    assert x >= 0 and x < 0x10000

    return x, y


def __dist(x1, y1, x2, y2):
    return math.sqrt((x1-x2)**2 + (y1-y2)**2)
    