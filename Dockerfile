FROM    golang:1.23 AS builder
WORKDIR  /src
COPY    ./     /src
RUN     go build cmd/*.go

FROM    debian:12-slim
COPY    --from=builder /src/zbproxy //src/zbproxy.json /app/
WORKDIR /app
CMD    ["/app/zbproxy"]


