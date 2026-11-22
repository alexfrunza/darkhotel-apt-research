FROM debian:trixie-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        supervisor python3 ca-certificates \
        iproute2 iputils-ping net-tools curl dnsutils traceroute netcat-openbsd tcpdump \
    && rm -rf /var/lib/apt/lists/*

COPY server.py /opt/heartbeat/server.py
COPY heartbeat-server.conf /etc/supervisor/conf.d/heartbeat-server.conf

# Defaults; override at run time, e.g. -e HEARTBEAT_TOKEN=secret
ENV PYTHONUNBUFFERED=1 \
    HEARTBEAT_HOST=0.0.0.0 \
    HEARTBEAT_PORT=8080 \
    HEARTBEAT_TOKEN="secret" \
    HEARTBEAT_LOG_FILE=/var/log/heartbeat/requests.log

RUN mkdir -p /var/log/heartbeat
VOLUME /var/log/heartbeat
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD python3 -c "import os,urllib.request as u; u.urlopen('http://127.0.0.1:%s/health' % os.environ['HEARTBEAT_PORT'], timeout=3)" || exit 1

# Runs as root. supervisord is PID 1, in the foreground, and starts the server.
CMD ["supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
