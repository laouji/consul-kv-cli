# Consul KV CLI

A cli tool that makes it easy to pipe formatted text from stdout into consul's kv, preserving
linebreaks and other formatting

## Features

* put - used to set a new key or overwrite an existing one setting a unique name for each node

## Example

Store the output of crontab -l in the kvs

```sh
user@host01 $ consul-kv-stash put mykeysuffix crontab -l
```

Another host can then fetch the contents by prepending the hostname to the key suffix

```sh
user@host02 $ curl http://localhost:8500/v1/kv/host01/mykeysuffix
```
