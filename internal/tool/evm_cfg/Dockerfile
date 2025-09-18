FROM rust:1.87 AS build
RUN cargo install --git https://github.com/plotchy/evm-cfg --root /installed

FROM python:3.12-slim
WORKDIR /app
RUN pip3 install pydot==3.0.4

WORKDIR /opt/evm-cfg
COPY --from=build /installed/bin/evm-cfg /opt/evm-cfg/evm-cfg

# install perf script
COPY measure.sh measure.sh
COPY helper.sh helper.sh

RUN chmod +x /opt/evm-cfg/measure.sh && \
    chmod +x /opt/evm-cfg/helper.sh

RUN apt-get update && apt-get install -y --no-install-recommends \
    jq \
    time \
    linux-perf \
    procps \
    && apt-get clean && rm -rf /var/lib/apt/lists/*
