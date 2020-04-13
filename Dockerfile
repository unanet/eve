FROM unanet-docker.jfrog.io/alpine-base
RUN apk --no-cache add ca-certificates
ADD ./bin/eve-api /app/eve-api
ADD ./bin/eve-scheduler /app/eve-scheduler
WORKDIR /app
CMD ["/app/eve-api"]
