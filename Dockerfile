# FROM unanet-docker.jfrog.io/alpine-base
FROM alpine

ENV EVE_PORT 3000
ENV EVE_METRICS_PORT 3001
ENV EVE_SERVICE_NAME eve-api
ENV EVE_SERVER_FLAG true
ENV EVE_MIGRATE_FLAG true
ENV VAULT_ADDR https://vault.unanet.io
ENV VAULT_ROLE k8s-devops

ADD ./bin/eve-api /app/eve-api
ADD ./migrations /app/migrations
WORKDIR /app
CMD ["/app/eve-api"]

HEALTHCHECK --interval=1m --timeout=2s --start-period=60s \
    CMD curl -f http://localhost:${EVE_METRICS_PORT}/ || exit 1
