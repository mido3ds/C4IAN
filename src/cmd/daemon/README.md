# Command Daemon
## Build Executable
```
$ go build
```

## Run
```
# main port is 9180, UI port is 7061, store data at "comms.db", read public keys from "unit-keys.db" and private key in /path/to/mykey
$ ./daemon --port 9180 --ui-port 7061 --store comms.db --keys unit-keys.db --priv-key /path/to/mykey

# for more options
$ ./daemon --help
```