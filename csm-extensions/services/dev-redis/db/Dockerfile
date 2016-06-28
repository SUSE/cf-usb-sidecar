FROM docker:dind

ADD redis-dind /usr/local/bin/
COPY image/redis.tgz /root/

ENTRYPOINT ["redis-dind"]
CMD []
