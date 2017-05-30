DDGO
========

This cli create Datadog monitoring settings semi-automatically.

```bash
# Build golang image for CI
$ docker build -t ci-golang:1.8 -f Dockerfile .

# make xbuild with mounting working dir and specifying current UID and GID
# work dir of DinD Jenkins may be different from ${PWD}
$ docker run -v ${PWD}:/go/src/github.com/nntsugu/ddgo -w /go/src/github.com/nntsugu/ddgo ci-golang:1.8 make xbuild
```
