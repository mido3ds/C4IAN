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
$ sudo iptables -A OUTPUT -j NFQUEUE --queue-num 0

# start router and attach to interface "wlan0"
$ sudo ./router --iface wlan0 --queue-num 0

# you must undo the redirect, otherwise all packets will stuck
$ sudo iptables -D OUTPUT -j NFQUEUE --queue-num 0

# list all options
$ ./router --help
```