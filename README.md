# Netmap

Netmap stands for Network Mapper, a visualizer for your inventory of network devices. Netmap starts collecting LLDP information with a single device credential and recursively collects neighbor of neighbors information. Built with love by Roopesh and friends in Go.

## Usage

```
roopesh:~/ $ netmap create --help                                                                                                                                                                         
Create topology diagram

Usage:
  netmap create

Flags:
  -h, --help              help for create
  -n, --hostname string   hostname to connect
  -p, --password string   password to connect to the host
  -u, --username string   username to connect to the host
roopesh:~/ $ netmap create -n ok270 -u admin -p password
```

<figure>
  <img src="./graph.png" alt="Graph created by netmap">
  <figcaption>Graph created by Netmap</figcaption>
</figure>