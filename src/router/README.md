# Router Daemon
## Requirements
```
$ sudo apt update &&
    sudo apt install -y libnetfilter-queue-dev
```

## Build Executable
```
$ go build
```

## Run
```
# start router and attach to interface "sta1-wlan0"
$ sudo ./router -i sta1-wlan0

# list all options
$ ./router --help
```