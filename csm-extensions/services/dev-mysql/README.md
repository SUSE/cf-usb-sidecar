# Hello to the mysql sidecar

## Building it

First build the foundation images, from the toplevel directory of the
checkout, via `make build-image`.  Publish them via `make
publish-image`. Do not forget to set `DOCKER_REPOSITORY` and
`DOCKER_ORGANIZATION` to suitable values before doing that.

```
export DOCKER_REPOSITORY=...
export DOCKER_ORGANIZATION=...

make build-image
make publish-image
```

Then build and publish the sidecar image itself, in this directory,
using the same script as above. At last generate the helm chart to
deploy sidecars, via `make helm`.

```
export DOCKER_REPOSITORY=...
export DOCKER_ORGANIZATION=...

make build-image
make publish-image
make helm
```

## Parameters of the chart

|Name			|Example	|Description					|
|---			|---		|---						|
|CF_DOMAIN		|?		|The SCF base domain				|
|CF_CA_CERT		|?		|The SCF CA cert				|
|SERVICE_LOCATION	|http://...	|Url to kube service `cf-usb-sidecar-mysql`	|
|SERVICE_TYPE		|mysql		|The name used to register the sidecar with SCF	|
|SERVICE_MYSQL_PORT	|3306		|The port the mysql server listens on		|
|SERVICE_MYSQL_HOST	|?		|The host of the mysql database to use		|
|SERVICE_MYSQL_USER	|?		|User to access the mysql database		|
|SERVICE_MYSQL_PASS	|?		|Credentials for the user above			|
