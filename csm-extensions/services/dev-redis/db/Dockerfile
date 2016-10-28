FROM docker:1.12.2-dind

ADD redis-dind /usr/local/bin/
COPY image/redis.tgz /root/

ENTRYPOINT ["redis-dind"]
CMD []
