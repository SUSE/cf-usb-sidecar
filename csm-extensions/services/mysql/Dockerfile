FROM sidecar:latest

RUN apk --no-cache add mysql-client

COPY ./SIDECAR_HOME/setup /catalog-service-manager/setup
COPY ./SIDECAR_HOME/connection /catalog-service-manager/connection
COPY ./SIDECAR_HOME/workspace /catalog-service-manager/workspace
COPY ./SIDECAR_HOME/status    /catalog-service-manager/status

#required for mysql-client to work correctly
ENV TERM=xterm-256color

ENV SIDECAR_PARAMETERS "MYSQL_ROOT_PASSWORD"

ENTRYPOINT ["/catalog-service-manager/bin/catalog-service-manager"]
