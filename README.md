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

## Why?
I was previously managing my kubeconfigs using the `$KUBECONFIG` environment variable. However, in order to automate this process, you have to do add something like this to your rc files:
```bash
KUBECONFIG=$(find $HOME/.kube -type f -name '*.conf' 2> /dev/null | sed ':a;N;$!ba;s/\n/:/g')
```
... that gets you a `$KUBECONFIG` variable with all your kubeconfigs separated by colons. The problem is that if you're frequently working with new/modified kubeconfigs, you'd have to trigger this command again.

With the `kconf` command, there's no need for `$KUBECONFIG` since `kubectl` already looks at `$HOME/.kube/config` by default. Additionally, as soon as you have a new kubeconfig, you can `add` it pretty easily and quickly.
