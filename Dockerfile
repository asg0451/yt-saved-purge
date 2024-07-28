FROM golang:1.22-bullseye AS builder
WORKDIR /build
RUN apt-get update && apt-get install -y make gcc musl-dev

COPY Makefile .

COPY go.mod go.sum  ./
RUN go mod download

COPY ./src ./src
RUN make build

FROM debian:bullseye-slim
RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates
RUN apt-get install -y curl htop dnsutils file wget && rm -rf /var/lib/apt/lists/* # for debugging

COPY --from=builder /build/bin/yt-saved-purge /usr/local/bin/yt-saved-purge

ENTRYPOINT ["/usr/local/bin/yt-saved-purge"]
