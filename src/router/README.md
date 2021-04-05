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
# redirect all packets to queue 0
$ sudo iptables -t filter -A OUTPUT -j NFQUEUE --queue-num 0

# start router and attach to interface "sta1-wlan0"
$ sudo ./router -i sta1-wlan0

# you must undo the redirect, otherwise all packets will stuck
$ sudo iptables -t filter -D OUTPUT -j NFQUEUE --queue-num 0

# list all options
$ ./router --help
```