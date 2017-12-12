#!/bin/sh

mkdir -p /etc/nginx/conf.d/tcp/

cat <<EOF > /etc/nginx/conf.d/tcp/service.conf
upstream service {
    server ${TARGET_SERVICE_HOST}.${KUBE_SERVICE_DOMAIN_SUFFIX}:${TARGET_SERVICE_PORT};
}

server {
    listen ${TARGET_SERVICE_PORT};
    proxy_pass  service;
}

EOF

echo "Updated /etc/nginx/conf.d/tcp/service.conf"
nginx
