# Command Daemon
## Build Executable
```
$ go build
```

## Run
```
# store data at "comms.db", read public keys from "unit-keys.db" and bind to port 7061
$ ./daemon --store comms.db --keys unit-keys.db --port 7061

# for more options
$ ./daemon --help
```