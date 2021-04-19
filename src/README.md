## Edit Toplogies
To create a new toplogy or edit and existing one, use `edit`.

For example:
```
$ ./edit topos/simple.topo
```

## Start Network Virtualisation
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
```
