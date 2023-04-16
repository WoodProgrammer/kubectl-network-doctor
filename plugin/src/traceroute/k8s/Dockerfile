FROM debian:bullseye-slim
WORKDIR /opt/traceroute
COPY main.sh .

RUN apt-get update && apt-get install -y \
    traceroute wget \
    && rm -rf /var/lib/apt/lists/*

RUN wget -q https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O yq
RUN mv yq /usr/local/bin/ && chmod +x /usr/local/bin/yq

RUN apt-get purge -y wget

CMD ["./main.sh"]