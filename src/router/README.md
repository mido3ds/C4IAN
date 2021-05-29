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
# assuming interface=wlan0 and ip=10.0.0.1
$ sudo ./deploy wlan0 10.0.0.1

# override default env variables
$ sudo PASS=somePass SSID=someSSIDName ./deploy wlan0 10.0.0.1
```

## Run
Assuming you have adhoc mode enabled and you only want to run the router.
```
# start router and attach to interface "sta1-wlan0"
$ sudo ./router -i sta1-wlan0 -p passphrase

# list all options
$ ./router --help
```