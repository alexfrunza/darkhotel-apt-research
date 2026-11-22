FROM debian:trixie-slim

ARG SUDO_VERSION=1.9.17

RUN apt-get update && apt-get install -y --no-install-recommends \
        supervisor python3 \
        zip unzip vim build-essential ca-certificates curl \
        iproute2 iputils-ping net-tools dnsutils traceroute netcat-openbsd tcpdump \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fsSLO https://www.sudo.ws/dist/sudo-${SUDO_VERSION}.tar.gz \
    && tar xzf sudo-${SUDO_VERSION}.tar.gz \
    && cd sudo-${SUDO_VERSION} \
    && ./configure --prefix=/usr --libexecdir=/usr/libexec \
                   --sysconfdir=/etc --with-secure-path --without-pam \
    && make -j"$(nproc)" \
    && make install \
    && cd / && rm -rf sudo-${SUDO_VERSION}*

# Non-root user "dev" with password "foobar", allowed to sudo (password required)
RUN useradd -m -s /bin/bash dev \
    && echo 'dev:foobar' | chpasswd \
    && echo 'dev ALL=(ALL) ALL' > /etc/sudoers.d/dev \
    && chmod 0440 /etc/sudoers.d/dev

COPY --chown=dev:dev go_script/unzipper /home/dev/unzipper
COPY --chown=dev:dev top_secret.zip /home/dev/top_secret.zip
COPY --chown=dev:dev passwords.txt /home/dev/passwords.txt

WORKDIR /home/dev

# supervisord as PID 1 (root), foreground. Loads /etc/supervisor/conf.d/*.conf
CMD ["bash", "-c", "supervisord -c /etc/supervisor/supervisord.conf && exec runuser -l dev"]
