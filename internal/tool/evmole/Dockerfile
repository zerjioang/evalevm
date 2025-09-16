FROM python:3.12-slim

RUN apt update && \
    apt-get install -y build-essential gcc

WORKDIR /opt/evmole

RUN pip install evmole
COPY run.py run.py

# install perf script
COPY measure.sh measure.sh
COPY helper.sh helper.sh

RUN chmod +x /opt/measure.sh && \
    chmod +x /opt/helper.sh

RUN apt-get update && apt-get install -y --no-install-recommends \
    jq \
    time \
    linux-perf \
    procps \
    && apt-get clean && rm -rf /var/lib/apt/lists/*
