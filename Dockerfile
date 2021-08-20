FROM golang:1.16-alpine AS build
ADD . /market
ENV CGO_ENABLED=0
WORKDIR /market
RUN go build -o market ./cmd/market4

FROM alpine:latest
COPY --from=build /market/market /market/market
ENTRYPOINT ["/market/market"]
EXPOSE 9999