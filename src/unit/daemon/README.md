# Unit Daemon
## Build Executable
```
$ go build
```

## Run
```
# run in virtual mode, listen at port 8090, store files in "unit.db"
$ ./daemon --store unit.db virt --port 8090

# for more options
$ ./daemon --help
```
