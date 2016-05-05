FROM catalog-service-manager:base

RUN apk --no-cache add --update mysql-client
COPY ./CSM_HOME/setup /catalog-service-manager/setup
COPY ./CSM_HOME/connection /catalog-service-manager/connection
COPY ./CSM_HOME/workspace /catalog-service-manager/workspace

COPY ./CSM_RDS_BIN/amazon-rds-mysql /catalog-service-manager/bin/amazon-rds-mysql

#required for mysql-client to work correctly
ENV TERM=xterm-256color

ENTRYPOINT /catalog-service-manager/bin/catalog-service-manager