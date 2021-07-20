package main

const schemaSQL = `
-- time is stored as Unix Time (int)

DROP TABLE IF EXISTS sent_msgs;
CREATE TABLE sent_msgs (
    time    int NOT NULL,
    dst     text NOT NULL,
    code    int NOT NULL,

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

    FOREIGN KEY (dst) 
        REFERENCES units (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE
);

DROP TABLE IF EXISTS received_msgs;
CREATE TABLE received_msgs (
    time    int NOT NULL,
    src     text NOT NULL,
    code    int NOT NULL,

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

    FOREIGN KEY (src) 
        REFERENCES units (ip) 
            ON DELETE CASCADE 
            ON UPDATE CASCADE
);
`
