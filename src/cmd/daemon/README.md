# Command Daemon
## Build Executable
```
$ go build
```

## Run
```
# TCP/UDP Port to communicate with units is 9180, units are listening on port 6000, connect to the UI using port 7061, store data at "comms.db"
$ ./daemon --port 9180 --units-port 6000 --ui-port 7061 --store comms.db

# for more options
$ ./daemon --help
```
