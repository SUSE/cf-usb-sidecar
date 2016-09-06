FROM gcr.io/google_containers/ubuntu-slim:0.3

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    nginx \
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

# forward nginx access and error logs to stdout and stderr of the daemon
# controller process
RUN ln -sf /proc/1/fd/1 /var/log/nginx/access.log \
	&& ln -sf /proc/1/fd/2 /var/log/nginx/error.log

COPY nginx.conf /etc/nginx/nginx.conf
COPY configure-nginx.bash /configure-nginx.bash

ENTRYPOINT ["/configure-nginx.bash"]
