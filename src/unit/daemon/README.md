# Unit Daemon
## Build Executable
```
$ go build
```

## Run
```
# run in virtual mode, main address is 127.0.0.1:9080, UI address is 127.0.0.1:8090, store files in "unit.db" and read public keys at "comms-keys.db"
$ ./daemon virt --adr 127.0.0.1:9080 --ui-adr 127.0.0.1:8090 --store unit.db --keys comms-keys.db

# for more options
$ ./daemon --help
```
