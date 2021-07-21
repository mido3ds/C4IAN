import sys
import sqlite3
from sqlite3 import Error

schema = \
"""
-- Time is stored as Unix Time (int)

DROP TABLE IF EXISTS forwarding;
CREATE TABLE forwarding (
    ip          text,
    dst         text,
    packet_hash blob,
    lat         real,
    lon         real
);
"""

db_file = sys.argv[1]

conn = sqlite3.connect(db_file)
conn.executescript(schema)
conn.close()
