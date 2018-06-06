# PostgreSQL database sidecar

This chart has the following set of parameters:

|Name              |Example        |Description
|---               |---            |---
|CF_ADMIN_PASSWORD |hunter2        |SCF cluster admin password
|CF_ADMIN_USER     |admin          |SCF cluster admin user name
|CF_CA_CERT        |----- BEGIN... |The SCF CA cert
|CF_DOMAIN         |cf-dev.io      |The SCF base domain
|SERVICE_LOCATION  |http://...     |URL to Kubernetes service `cf-usb-sidecar-postgresql`
|PGHOST            |pg-staging     |The host of the postgres database to use
|PGPORT            |5432           |The port the postgres server listens on
|PGSSLMODE         |disable        |Connection to postgres server, one of `disable`, `require`, `verify-ca`, `verify-full`
|PGUSER            |root           |User to access the postgres database
|PGPASSWORD        |hunter2        |Credentials for the user above
|PGDATABASE        |postgres       |Name of database to connect to (optional)
|SERVICE_TYPE      |postgres       |The name used to register the sidecar with SCF
|SIDECAR_LOG_LEVEL |debug          |Logging level; more verbose than `info` is not recommended
|UAA_CA_CERT       |----- BEGIN... |The UAA CA certificate
