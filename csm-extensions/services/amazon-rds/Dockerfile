FROM sidecar:latest

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    mysql-client \ 
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*
    
COPY ./SIDECAR_HOME/setup /catalog-service-manager/setup
COPY ./SIDECAR_HOME/connection /catalog-service-manager/connection
COPY ./SIDECAR_HOME/workspace /catalog-service-manager/workspace

COPY ./SIDECAR_RDS_BIN/amazon-rds-mysql /catalog-service-manager/bin/amazon-rds-mysql

#required for mysql-client to work correctly
ENV TERM=xterm-256color

ENTRYPOINT ["/catalog-service-manager/bin/catalog-service-manager"]
