# nvclient
open a file for editing in an already running nvim

## How it works
nvclient connect to the existing nvim instance by looking for the
environment variables `NVIM_LISTEN_ADDRESS`, if the `NVIM_LISTEN_ADDRESS`
is not found when nvclient start, it will start a new nvim instance.

If `NVIM_LISTEN_ADDRESS` is already set when nvim start, nvim will use its vaule
to create the socket, otherwise `NVIM_LISTEN_ADDRESS` is only can be found
in the nvim terminal. In this case, you can check the value by
`:echo $NVIM_LISTEN_ADDRESS`

you can do the trick by setting the `NVIM_LISTEN_ADDRESS` in your shell's .dotfile,
so everytime you open a file with nvcllient, it will be sent to the same nvim
instance, but beware any new nvim instance will overwrite the socket file.

For tmux user, nvclient read the tmux environment for `NVIM_LISTEN_ADDRESS`
if it is running in tmux, but you have to set the tmux `update-environment`
correctly in .tmux.conf to make it work.

Example:
```
 set -g update-environment "SSH_ASKPASS WINDOWID XAUTHORITY NVIM_LISTEN_ADDRESS"
```

## Installation
Download the binary at <https://github.com/dustinliu/nvclient/releases>
and copy the nvclient to the $PATH

for go user
```
go install github.com/dustinliu/nvclient@latest
```
