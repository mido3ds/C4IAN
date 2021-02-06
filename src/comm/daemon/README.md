# Command Daemon
## Build Executable
```
$ go build
```

## Run
```
# store data at "comms.db", read public keys from "unit-keys.db", main address is 127.0.0.1:9180 and UI address is 127.0.0.1:7061
$ ./daemon --store comms.db --keys unit-keys.db --adr 127.0.0.1:9180 --ui-adr 127.0.0.1:7061

# for more options
$ ./daemon --help
```