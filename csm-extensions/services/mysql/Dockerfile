FROM csm:latest

RUN apk --no-cache add mysql-client

COPY ./CSM_HOME/setup /catalog-service-manager/setup
COPY ./CSM_HOME/connection /catalog-service-manager/connection
COPY ./CSM_HOME/workspace /catalog-service-manager/workspace
COPY ./CSM_HOME/status    /catalog-service-manager/status

#required for mysql-client to work correctly
ENV TERM=xterm-256color

ENV CSM_PARAMETERS "MYSQL_ROOT_PASSWORD"

ENTRYPOINT ["/catalog-service-manager/bin/catalog-service-manager"]
