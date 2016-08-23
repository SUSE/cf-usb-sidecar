FROM sidecar:latest

RUN apk --no-cache add --update mysql-client
COPY ./SIDECAR_HOME/setup /catalog-service-manager/setup
COPY ./SIDECAR_HOME/connection /catalog-service-manager/connection
COPY ./SIDECAR_HOME/workspace /catalog-service-manager/workspace

COPY ./SIDECAR_RDS_BIN/amazon-rds-mysql /catalog-service-manager/bin/amazon-rds-mysql

#required for mysql-client to work correctly
ENV TERM=xterm-256color

ENTRYPOINT ["/catalog-service-manager/bin/catalog-service-manager"]
