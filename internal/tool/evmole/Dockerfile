FROM python:3.12-slim@sha256:f8081b61393cc16b2339d377bb523a5646cdd28fc6105612b3893c405d27f730

RUN apt update && \
    apt-get install -y build-essential gcc

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
