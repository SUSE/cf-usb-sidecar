FROM ubuntu:14.04

RUN apt-get update && apt-get install -y curl && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*
ADD test-image-up /usr/local/bin/

ENTRYPOINT ["test-image-up","curl --silent -k -H \"x-sidecar-token: $SIDECAR_API_KEY\" ${SIDECAR_ENDPOINT}/status | grep successful"]
CMD []
