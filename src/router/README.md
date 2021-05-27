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

## Deployment
To enter adhoc mode and run the router locally.
```
# use defaults
$ sudo ./deploy

# override defaults
$ sudo IP=10.0.0.1 IFACE=wlan0 PASS=somePass SSID=someSSIDName ./deploy
```

## Run
Assuming you have adhoc mode enabled and you only want to run the router.
```
# start router and attach to interface "sta1-wlan0"
$ sudo ./router -i sta1-wlan0 -p passphrase

# list all options
$ ./router --help
```