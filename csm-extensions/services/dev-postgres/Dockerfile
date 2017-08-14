FROM cf-usb-sidecar:latest

COPY ./SIDECAR_HOME /catalog-service-manager
RUN chmod +x /catalog-service-manager/bin/run.sh

ENTRYPOINT ["sh", "/catalog-service-manager/bin/run.sh"]
