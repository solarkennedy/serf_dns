# Serf_DNS

A program to serve DNS requests based on the Serf toplology.

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
