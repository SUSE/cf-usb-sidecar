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
using a similar set of commands and the same docker configuration. At last
generate the helm chart to deploy sidecars, via `make helm`.

```
make build-image
make build-service-image ;# Needed only for publish to work
make publish-image
make helm
```

## Parameters of the chart

|Name			|Example	|Description					|
|---			|---		|---						|
|CF_ADMIN_PASSWORD	|?		|SCF cluster admin password			|
|CF_ADMIN_USER		|admin		|SCF cluster admin user				|
|CF_CA_CERT		|?		|The SCF CA cert				|
|CF_DOMAIN		|cf-dev.io	|The SCF base domain				|
|SERVICE_LOCATION	|http://...	|Url to kube service `cf-usb-sidecar-mysql`	|
|SERVICE_MYSQL_HOST	|mysql		|The host of the mysql database to use		|
|SERVICE_MYSQL_PASS	|?		|Credentials for the user above			|
|SERVICE_MYSQL_PORT	|3306		|The port the mysql server listens on		|
|SERVICE_MYSQL_USER	|root		|User to access the mysql database		|
|SERVICE_TYPE		|mysql		|The name used to register the sidecar with SCF	|
|UAA_CA_CERT		|?		|The UAA CA cert   				|


## Mysql version compatibility note

Currently, the mysql sidecar will only work with deployments which use
`mysql_native_password` as their authentication plugin `--default-auth`. This is the
default for myql versions 8.0.3 and earlier, but later versions will need to be started
with `--default-auth=mysql_native_password` *before* any user creation, in order to work.

See [go-sql-driver/mysql Issue #785](https://github.com/go-sql-driver/mysql/issues/785) for more information.
