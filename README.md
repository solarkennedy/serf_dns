# Serf_DNS
[![Build Status](https://travis-ci.org/solarkennedy/serf_dns.svg?branch=master)](https://travis-ci.org/solarkennedy/serf_dns)

A program to serve DNS requests based on the Serf topology.

## Description

This program uses the [Serf RPC](http://www.serfdom.io/docs/agent/rpc.html) 
to query Serf and inspect the current names and IPs.

From this data it builds an in-memory mapping of DNS records, and responds to
DNS requests.

## Building

You need at least Go 1.2

## Running

I wouldn't yet :( 

There is no error handling, it's basically crap.

## Examples

You could run this program, and then have BIND forward a particular zone 
to it:

```
zone "serf." in {
    type forward;
    ; forward to Serf_dns running on localhost!
    forwarders { 127.0.0.1 port 8053 ; };
};
```

## Demo

Once running, you can confirm the records are in place with dig:
```
$ # Using a "serf" fake TLD
$ # Lets query the local serf_dns running on 8053...
$ dig +short @localhost -p 8053  server1.xkyle.com.serf.
192.168.1.67
$ # Run serf members to verify...
$ serf members
server1.xkyle.com     192.168.1.67:7946    alive
server2.xkyle.com     192.168.1.69:7946    alive
$ echo 'It Works!'
It Works!
```
