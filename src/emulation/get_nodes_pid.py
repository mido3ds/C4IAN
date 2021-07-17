import subprocess
import sys


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

assert(kind in ['r', 'u', 'c'])

n = 0

if kind == 'r':
    for name, pid in get_nodes():
        n += 1
        # TODO: load groups file
        print(f'-t {pid} -n ./router/router -i '
              f'{name}-wlan0 -p pass -l /tmp/{name}.router.locsock')
elif kind == 'u':
    for name, pid in get_nodes():
        if 'u' in name:
            n += 1
            print(f'-t {pid} -n ./unit/daemon/daemon '
                  f'-s /var/lib/caian/{name}.store.sqllite')

elif kind == 'c':
    for name, pid in get_nodes():
        if 'c' in name:
            n += 1
            # TODO: create units/groups files
            print(f'-t {pid} -n ./cmd/daemon/daemon '
                  f'-s /var/lib/caian/{name}.store.sqllite --videos-path /var/lib/caian/{name}.videos --units-path /tmp/units.json --groups-path /tmp/groups.json')

if n == 0:
    exit(1)
