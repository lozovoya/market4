FROM golang:1.16-alpine AS build
ADD . /market
ENV CGO_ENABLED=0
WORKDIR /market
RUN go build -o market ./cmd/market4

FROM alpine:latest
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

COPY --from=build /market/market /market/market
ENTRYPOINT ["/market/market"]
EXPOSE 9999

