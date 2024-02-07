FROM golang:1.21 as builder

RUN apt-get update -y && apt-get upgrade -y \
    && apt install build-essential git -y

WORKDIR /app
COPY . .
RUN make bor

FROM debian:12

COPY docker/cron/cron.conf /etc/cron.d/cron.conf
COPY docker/cron/prune.sh /prune.sh
COPY docker/supervisord/gethlighthousebn.conf /etc/supervisor/conf.d/supervisord.conf

# Install Supervisor and create the Unix socket
RUN touch /var/run/supervisor.sock

RUN apt-get update && apt-get install cron supervisor ca-certificates tini -y \
    && apt-get clean && rm -r /var/lib/apt/lists/*

COPY --from=builder /app/build/bin/bor /usr/bin/

ENTRYPOINT ["/usr/bin/tini", "--", "supervisord", "-n", "-c",  "/etc/supervisor/conf.d/supervisord.conf"]
