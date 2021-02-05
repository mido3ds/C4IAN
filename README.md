# Development
## Common Requirements
### golang 1.15.8
```
$ sudo apt update && sudo apt install -y wget && (
    set -e
    cd /tmp
    wget -c https://golang.org/dl/go1.15.8.linux-amd64.tar.gz
    tar xvf go1.15.8.linux-amd64.tar.gz
    sudo chown -R root:root ./go
    sudo mv go /usr/local
    mkdir -p $HOME/.config/golang/bin
    echo >>"$HOME"/.bash_profile
    echo 'export GOPATH="$HOME/.config/golang"' >>"$HOME"/.bash_profile
    echo 'export PATH="$PATH:/usr/local/go/bin:$GOPATH/bin"' >>"$HOME"/.bash_profile
    echo >>"$HOME"/.bash_profile
    . "$HOME"/.bash_profile
)
```
### yarn
### nodejs

## Issues
- `electron: error while loading shared libraries: libgconf-2.so.4: cannot open shared object file: No such file or directory`
```
$ sudo apt install -y libgconf-2-4
```

Follow the rest of `README`s.