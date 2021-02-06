# Unit Daemon
## Build Executable
```
$ go build
```

## Run
```
# run in virtual mode, main port is 9080, UI port is 8090, store files in "unit.db", read public keys at "comms-keys.db" and private key in /path/to/mykey
$ ./daemon virt --port 9080 --ui-port 8090 --store unit.db --keys comms-keys.db --priv-key /path/to/mykey

# for more options
$ ./daemon --help
```
