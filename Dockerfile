FROM unanet-docker.jfrog.io/alpine-base

ENV EVE_PORT 8080
ENV EVE_METRICS_PORT 3000
ENV EVE_SERVICE_NAME eve-api

ADD ./bin/eve-api /app/eve-api
ADD ./migrations /app/migrations
WORKDIR /app
CMD ["/app/eve-api"]

HEALTHCHECK --interval=1m --timeout=2s --start-period=60s \
    CMD curl -f http://localhost:${EVE_METRICS_PORT}/ || exit 1
