FROM paper-algorithm:latest

WORKDIR /opt

COPY measure.sh measure.sh
COPY helper.sh helper.sh

RUN chmod +x /opt/measure.sh && \
    chmod +x /opt/helper.sh

RUN apt-get update && apt-get install -y gawk jq procps time linux-perf
