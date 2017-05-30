DDGO
========

This cli create Datadog monitoring settings semi-automatically.

# How to build on Docker

```bash
# Build golang image for CI
$ docker build -t ci-golang:1.8 -f Dockerfile .

# make xbuild with mounting working dir
$ docker run -v ${PWD}:/go/src/github.com/nntsugu/ddgo -w /go/src/github.com/nntsugu/ddgo ci-golang:1.8 make xbuild
```
