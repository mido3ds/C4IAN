## Edit Toplogies
Use `edit` to create a new toplogy, edit an existing one connect to mininet-wifi and start nodes mobility.
```
$ ./edit
```
## Start Network Virtualisation
```
# First, install `watchgod`:
$ sudo pip3 install simple_websocket_server
```
```
## (1st terminal)
# start mininet-wifi with a topology of 2 nodes
$ sudo ./start mn topos/2nodes.topo

## (2nd terminal) 
# start routers in all nodes
$ sudo ./start routers

## (3rd terminal) 
# u1 ping u2
$ sudo ./start u1 ping 10.0.0.2 -c5
# u2 ping u1
$ sudo ./start u2 ping 10.0.0.1 -c5
```

## Start UNITs/CMDs in nodes
```sh
$ sudo ./start units
$ sudo ./start cmds
$ sudo ./start hals -- --video /path/to/video --audios-dir /path/to/audio/dir
```

## Start single UNIT/CMD/HAL locally
```sh
$ sudo ./start local unit
$ sudo ./start local cmd
$ sudo ./start local hal -- --video /path/to/video --audios-dir /path/to/audio/dir
```

## Watch mode
Use watch mode to automatically rerun routers/units/cmds whenever a change happens to their executable file.

First, install `watchgod`:
```
$ sudo pip3 install watchgod
```

Then enter watch mode by passing `-w` or `--watch` like this:
```
$ sudo ./start routers -w
$ sudo ./start units -w
$ sudo ./start cmds -w
$ sudo ./start local unit -w
$ sudo ./start local cmd -w
```

## Pass more options
With `start` script you can always pass more options to routers/units/cmds, just put them after `--`:
```sh
$ sudo ./start local cmd -- <options to cmd executable...>
```

## Note
You may need to run the following command to enable writing to `/tmp/`:
```
sudo sysctl fs.protected_regular=0
```
