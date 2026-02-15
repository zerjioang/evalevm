FROM python:3.8-slim

WORKDIR /opt/octopus

RUN apt-get update && apt-get install -y \
    graphviz \
    xdg-utils \
    git \
    procps \
    time \
    gawk jq \
    linux-perf \
    && rm -rf /var/lib/apt/lists/*

# Install octopus
# Using git clone because pip install octopus might get the wrong package or outdated
RUN git clone https://github.com/FuzzingLabs/octopus . && \
    pip install . && \
    pip install graphviz

# Copy helper scripts
COPY measure.sh .
COPY helper.sh .
RUN chmod +x measure.sh helper.sh

# The tool entrypoint will be via octopus_eth_evm.py which is in the repo root
# We need to make sure it's executable or call with python3
