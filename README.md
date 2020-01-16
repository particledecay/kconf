# kconf
An opinionated command line tool for managing multiple kubeconfigs.

## Description
kconf works by storing all kubeconfig information in a single file (`$HOME/.kube/config`). This file is looked at by default when using `kubectl`.

## Usage
To merge in a new kubeconfig file:
```bash
kconf add /path/to/kubeconfig.conf
```
To remove an existing kubeconfig:
```bash
kconf rm myContext
```

To view all saved contexts in the kubeconfig:
```bash
kconf list
```

## Why?
I was previously managing my kubeconfigs using the `$KUBECONFIG` environment variable. However, in order to automate this process, you have to do something like this in your rc files:
```bash
KUBECONFIG=$(find $HOME/.kube -type f -name '*.conf' 2> /dev/null | sed ':a;N;$!ba;s/\n/:/g')
```
... that gets you a `$KUBECONFIG` variable with all your kubeconfigs separated by colons. The problem is that if you're frequently working with new/modified kubeconfigs, you'd have to trigger this command each time something changed.

With the `kconf` command, there's no need for `$KUBECONFIG` since `kubectl` already looks at `$HOME/.kube/config` by default. Additionally, as soon as you have a new kubeconfig, you can `add` it pretty easily and quickly.

## Known Issues
Check out the [Issues](https://github.com/particledecay/kconf/issues) section or specifically [issues created by me](https://github.com/particledecay/kconf/issues?q=is:issue+is:open+sort:updated-desc+author:particledecay)