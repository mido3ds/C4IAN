import subprocess
import sys
import json


def get_nodes():
    try:
        lines = [
            line.split()
            for line in subprocess.check_output('pgrep -u root -a bash | grep mininet:', shell=True)
            .decode('utf8')
            .splitlines()
        ]

        for line in lines:
            name = line[-1].split(':')[-1]
            pid = int(line[0])
            yield name, pid
    except:
        exit(1)


def node_to_tcp(name: str) -> int:
    kind = name[0]
    assert kind in ['u', 'c']

    num = int(name[1:])
    assert num >= 0 and num < 100

    if kind == 'c':
        return 3100 + num

    return 3200 + num

kind = sys.argv[1]
mode = sys.argv[2]

with open('/tmp/mn.metadata.json') as f:
    metadata = json.load(f)


def router_get_groups_file(name):
    return metadata['router']['groups_file'][name]


def unit_get_cmd_addr(name):
    return metadata['unit']['cmd_addr'][name]


def cmd_get_units_file(name):
    return metadata['cmd']['units_file'][name]


def cmd_get_groups_file(name):
    return metadata['cmd']['groups_file'][name]


assert(kind in ['r', 'u', 'c', 'h'])
assert(mode in ['nsenter', 'socat'])

n = 0

if kind == 'r':
    for name, pid in get_nodes():
        n += 1
        assert mode == 'nsenter'
        print(f'-t {pid} -n '
              f'./router/router '
              f'--iface {name}-wlan0 '
              f'--pass xXxHaCkEr_MaNxXx '  # very hard to guess password
              f'--location-socket /tmp/{name}.router.locsock '
              f'--mgroups-file {router_get_groups_file(name)} ')

elif kind == 'u':
    for name, pid in get_nodes():
        if 'u' in name:
            n += 1
            if mode == 'nsenter':
                print(f'-t {pid} -n '
                      f'./unit/daemon/daemon '
                      f'--iface {name}-wlan0 '
                      f'--ui-socket /tmp/{name}.unit.sock '
                      f'--store /var/lib/caian/{name}.store.sqllite '
                      f'--hal-socket-path /tmp/{name}.hal.sock '
                      f'--cmd-addr {unit_get_cmd_addr(name)} ')
            elif mode == 'socat':
                unix = f'/tmp/{name}.unit.sock'
                tcp = node_to_tcp(name)
                print(f'{unix} <--> {tcp}', file=sys.stderr)
                print(f'TCP4-LISTEN:{tcp},fork,reuseaddr '
                      f'UNIX-CONNECT:{unix}')

elif kind == 'c':
    for name, pid in get_nodes():
        if 'c' in name:
            n += 1
            if mode == 'nsenter':
                print(f'-t {pid} -n '
                      f'./cmd/daemon/daemon '
                      f'--iface {name}-wlan0 '
                      f'--ui-socket /tmp/{name}.cmd.sock '
                      f'--store /var/lib/caian/{name}.store.sqllite '
                      f'--videos-path /var/lib/caian/{name}.videos '
                      f'--units-path {cmd_get_units_file(name)} '
                      f'--groups-path {cmd_get_groups_file(name)} ')
            elif mode == 'socat':
                unix = f'/tmp/{name}.cmd.sock'
                tcp = node_to_tcp(name)
                print(f'{unix} <--> {tcp}', file=sys.stderr)
                print(f'TCP4-LISTEN:{tcp},fork,reuseaddr '
                      f'UNIX-CONNECT:{unix}')

elif kind == 'h':
    for name, pid in get_nodes():
        if 'u' in name:
            n += 1
            assert mode == 'nsenter'
            print(f'-t {pid} -n '
                  f'./unit/halsimulation/halsimulation '
                  f'--iface {name}-wlan0 '
                  f'--hal-socket-path /tmp/{name}.hal.sock '
                  f'--location-socket /tmp/{name}.hal.locsock ')

if n == 0:
    exit(1)
