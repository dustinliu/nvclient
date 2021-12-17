# nvclient
open a file for editing in an already running nvim

## How it works
nvclient connect to the existing nvim instance by looking for the
environment variables `NVIM_LISTEN_ADDRESS`, if the `NVIM_LISTEN_ADDRESS`
is not found when nvclient start, it will start a new nvim instance.

If `NVIM_LISTEN_ADDRESS` is set when nvim start, nvim will use its vaule
to create the socket, otherwise `NVIM_LISTEN_ADDRESS` is only can be found
in the nvim terminal. In this case, you can check the value by
`:echo $NVIM_LISTEN_ADDRESS`

you can set the `NVIM_LISTEN_ADDRESS` in your shell's .dotfile,
so everytime you open a file with nvcllient, it will be sent to the same nvim
instance.

This is my .dotfile setting
```bash
export NVIM_LISTEN_ADDRESS=/tmp/nvim.socket
aliase vim=nvclient
```

## Installation
Download the binary at <https://github.com/dustinliu/nvclient/releases>
and copy the nvclient to the $PATH

for go user
```
go install github.com/dustinliu/nvclient@latest
```
