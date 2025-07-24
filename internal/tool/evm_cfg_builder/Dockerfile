# syntax=docker/dockerfile:1

FROM python:3.6-slim

# Avoid interactive prompts and install required system tools
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install gcc -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/evm_cfg_builder

COPY . .

RUN chmod +rwx measure.sh
RUN apt-get update && apt-get install -y gawk jq procps time linux-perf
RUN pip install --upgrade pip && \
    pip install .

# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/bin/perf_5.10 /usr/bin/perf

ENTRYPOINT ["python3", "evm_cfg_builder.py"]
