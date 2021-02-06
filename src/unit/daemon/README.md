# Unit Daemon
## Build Executable
```
$ go build
```

## Run
```
# run in virtual mode, main port is 9080, UI port is 8090, store files in "unit.db" and read public keys at "comms-keys.db"
$ ./daemon virt --port 9080 --ui-port 8090 --store unit.db --keys comms-keys.db

# for more options
$ ./daemon --help
```
