FROM golang:1.14.2-alpine AS builder
WORKDIR /src/
COPY cmd /src/cmd
COPY internal /src/internal
COPY pkg /src/pkg
COPY go.mod /src/
COPY go.sum /src/ 
RUN go build -o eve-api ./cmd/eve-api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/eve-api /app/eve-api
CMD ["/app/eve-api"]