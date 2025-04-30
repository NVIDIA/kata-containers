FROM ubuntu:22.04

RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y systemd
COPY kata-ops-logs.sh /usr/local/bin
RUN chmod +x /usr/local/bin/kata-ops-logs.sh
ENV JOURNAL_DIR="/host/var/log/journal/"
