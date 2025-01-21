# Preface

UDP forwarding with encryption and FEC.

before:

```
client --> 8.8.8.8:53
```

after:

```
client --> (uboost client) --> (uboost server) --> 8.8.8.8:53
```

## Quick start

Assuming that you want to forward UDP traffic from local port 53 to 8.8.8.8:53, you are required to start a uboost server to forward packet:

```bash
./uboost -side server -src 0.0.0.0:5300 -dst 8.8.8.8:53 -mask secret -fec 1
```

and also a uboost client to serve on local port 53:

```bash
./uboost -side client -src 0.0.0.0:53 -dst 120.50.20.7:5300 -mask secret -fec 1
```

(PS: 120.50.20.7 is the remote server running uboost server)

After that, traffic between uboost client and server will be accelerated and encrypted.

Quick test:

```bash
nslookup www.google.com 127.0.0.1
```

The dns request to 127.0.0.1:53 will be accepted by uboost client and forward to 8.8.8.8:53 via uboost server. 

## Download

Download the binaries in the [releases](https://github.com/skywind3000/uboost/releases) page.

## Usage

And execute it in the command line:

```text
Usage of ./uboost:
  -side string
        forward side: client/server
  -src string
        local address, eg: 0.0.0.0:8080
  -dst string
        destination address, eg: 8.8.8.8:443
  -mask string
        encryption/decryption key
  -fec int
        fec redundancy
  -mark uint
        fwmark value
```

To see the usage documentation.

## Credit

TODO
