FROM python:3.10-slim

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update && \
    apt-get install -y build-essential gcc cargo pkg-config libssl-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /opt/evmole

RUN pip install evmole
COPY run.py run.py

# install perf script
COPY measure.sh measure.sh
COPY helper.sh helper.sh

RUN chmod +x /opt/evmole/measure.sh && \
    chmod +x /opt/evmole/helper.sh

RUN apt-get update && apt-get install -y --no-install-recommends \
    jq \
    time \
    linux-perf \
    procps \
    && apt-get clean && rm -rf /var/lib/apt/lists/*
