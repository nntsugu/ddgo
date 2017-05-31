DDGO
========

This cli create Datadog monitoring settings semi-automatically.

# How to build on Docker

```bash
# Build golang image
$ docker build -t build-golang:1.8 -f Dockerfile .

# make xbuild with mounting working dir
$ docker run -v ${PWD}:/go/src/github.com/nntsugu/ddgo -w /go/src/github.com/nntsugu/ddgo build-golang:1.8 make xbuild
```

After that you can see zip/tar files under the snapshot directory.
Please use the suitable one for your OS.

# How to use it

Run the following command then all monitoring settings in PATH_TO_MONITORING_CONF_DIR will be created on Datadog. 

```bash
$ ddgo -f PATH_TO_DATADOG_KEYS -m PATH_TO_MONITORING_CONF_DIR
```

## PATH_TO_DATADOG_KEYS

```bash
datadog:
  api_key: <<API KEY>>
  app_key: <<APPLICATION KEY>>
```

You can get the above keys by the following URL.
- https://app.datadoghq.com/account/settings#api


## PATH_TO_MONITORING_CONF_DIR
- see examples/monitoring_setting
- You can get json from existing setting via export function of Datadog.


