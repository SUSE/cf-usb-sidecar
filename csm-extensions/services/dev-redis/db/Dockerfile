FROM docker:1.11.2-dind

ADD redis-dind /usr/local/bin/
COPY image/redis.tgz /root/

ENTRYPOINT ["redis-dind"]
CMD []
