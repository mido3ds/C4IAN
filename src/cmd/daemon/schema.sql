-- time is stored as Unix Time (int)
DROP TABLE IF EXISTS units;
CREATE TABLE units (
    ip text PRIMARY KEY
);

DROP TABLE IF EXISTS groups;
CREATE TABLE groups (
    ip text PRIMARY KEY
);

DROP TABLE IF EXISTS members;
CREATE TABLE members (
    group_ip    text,
    unit_ip     text,

    PRIMARY KEY(group_ip, unit_ip),
    FOREIGN KEY (group_ip) 
        REFERENCES groups (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE,
    FOREIGN KEY (unit_ip) 
        REFERENCES units (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE
);

DROP TABLE IF EXISTS sent_msgs;
CREATE TABLE sent_msgs (
    time    int NOT NULL,
    dst     text NOT NULL,
    body    text NOT NULL,

    PRIMARY KEY (time, dst),
    FOREIGN KEY (dst) 
        REFERENCES units (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE
);

DROP TABLE IF EXISTS sent_audios;
CREATE TABLE sent_audios (
    time    int NOT NULL,
    dst     text NOT NULL,
    body    blob NOT NULL,

    PRIMARY KEY (time, dst),
    FOREIGN KEY (dst) 
        REFERENCES units (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE
);

DROP TABLE IF EXISTS received_msgs;
CREATE TABLE received_msgs (
    time    int NOT NULL,
    src     text NOT NULL,
    body    text NOT NULL,

    PRIMARY KEY (time, src),
    FOREIGN KEY (src) 
        REFERENCES units (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE
);

DROP TABLE IF EXISTS received_audios;
CREATE TABLE received_audios (
    time    int NOT NULL,
    src     text NOT NULL,
    body    blob NOT NULL,

    PRIMARY KEY (time, src),
    FOREIGN KEY (src) 
        REFERENCES units (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE
);

DROP TABLE IF EXISTS received_videos;
CREATE TABLE received_videos (
    time    int NOT NULL,
    src     text NOT NULL,
    path    text NOT NULL,

    PRIMARY KEY (time, src),
    FOREIGN KEY (src) 
        REFERENCES units (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE
);

DROP TABLE IF EXISTS received_sensor_data;
CREATE TABLE received_sensor_data (
    time        int NOT NULL,
    src         text NOT NULL,
    heartbeat   int NOT NULL,
    loc_x       int NOT NULL,
    loc_y       int NOT NULL,

    PRIMARY KEY (time, src),
    FOREIGN KEY (src) 
        REFERENCES units (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE
);