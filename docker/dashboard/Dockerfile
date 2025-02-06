ARG platform
FROM --platform=${platform} debian:stable

ENV TZ="Asia/Shanghai"

COPY docker/dashboard/entrypoint.sh /entrypoint.sh
COPY docker/dashboard/scripts /scripts

RUN export DEBIAN_FRONTEND="noninteractive" && \
	apt update && apt install -y ca-certificates tzdata && \
	update-ca-certificates && \
	ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
	dpkg-reconfigure tzdata && \
	chmod +x /entrypoint.sh

WORKDIR /opt/sleepy-dashboard

COPY examples/dashboard /opt/sleepy-dashboard/examples
COPY custom /opt/sleepy-dashboard/custom
COPY resources /opt/sleepy-dashboard/resources
COPY sleepy-dashboard /opt/sleepy-dashboard/sleepy-dashboard

EXPOSE 80 443
VOLUME [ "/data" ]

ENTRYPOINT ["/entrypoint.sh"]
CMD ["--config", "/data/config.json"]