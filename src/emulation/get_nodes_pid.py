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


kind = sys.argv[1]

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

n = 0

if kind == 'r':
    for name, pid in get_nodes():
        n += 1
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
            print(f'-t {pid} -n '
                  f'./unit/daemon/daemon '
                  f'--store /var/lib/caian/{name}.store.sqllite '
                  f'--hal-socket-path /tmp/{name}.hal.sock '
                  f'--cmd-addr {unit_get_cmd_addr(name)} ')

elif kind == 'c':
    for name, pid in get_nodes():
        if 'c' in name:
            n += 1
            print(f'-t {pid} -n '
                  f'./cmd/daemon/daemon '
                  f'--store /var/lib/caian/{name}.store.sqllite '
                  f'--videos-path /var/lib/caian/{name}.videos '
                  f'--units-path {cmd_get_units_file(name)} '
                  f'--groups-path {cmd_get_groups_file(name)} ')

elif kind == 'h':
    for name, pid in get_nodes():
        if 'u' in name:
            n += 1
            print(f'-t {pid} -n '
                  f'./unit/halsimulation/halsimulation '
                  f'--hal-socket-path /tmp/{name}.hal.sock '
                  f'--location-socket /tmp/{name}.hal.locsock ')

if n == 0:
    exit(1)
