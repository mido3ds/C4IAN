# Unit Daemon
## Build Executable
```
$ go build
```

## Run
```
# run in virtual mode, listen at port 8090, store files in "unit.db" and read public keys at "comms-keys.db"
$ ./daemon virt --port 8090 --store unit.db --keys comms-keys.db

# for more options
$ ./daemon --help
```
