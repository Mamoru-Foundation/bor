FROM golang:1.19 as builder

RUN apt-get update -y && apt-get upgrade -y \
    && apt install build-essential git -y

WORKDIR /app
COPY . .
RUN make bor

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && apt-get clean && rm -r /var/lib/apt/lists/*

COPY --from=builder /app/build/bin/bor /usr/bin/

ENTRYPOINT ["bor"]
