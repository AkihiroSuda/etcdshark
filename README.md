# etcdshark: packet capturing tool for etcd raft

[![Build Status](https://travis-ci.org/AkihiroSuda/etcdshark.svg?branch=master)](https://travis-ci.org/AkihiroSuda/etcdshark)
[![Go Report Card](https://goreportcard.com/badge/github.com/AkihiroSuda/etcdshark)](https://goreportcard.com/report/github.com/AkihiroSuda/etcdshark)

Still work-in progress.

## Usage

Install:

    $ go get github.com/AkihiroSuda/etcdshark

Live mode:

    $ sudo ./etcdshark -i eth0

Off-line mode:

    $ ./etcdshark -r a.pcap


## TODO
### Raft packets
  - [X] /raft/probing
  - [X] /raft/stream/message
  - [ ] /raft/stream/msgapp
  - [ ] /raft/snapshot

### C/S packets

  - [X] v2
  - [ ] v3

### Output

  - [X] colorful output for each of the stream
  - [ ] JSON output (friendly to `jq`)

### Others

  - [ ] unit test
