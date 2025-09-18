FROM docker.io/smartbugs/vandal:latest@sha256:134f5e1660920ffd0e426829f484b4c7d86c5bdd2987c896e842845a2a9a9b41

# Set working directory
WORKDIR /opt/vandal

COPY . .

RUN chmod +x /opt/vandal/measure.sh && \
    chmod +x /opt/vandal/helper.sh

# Install required system and build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    gawk \
    jq \
    procps \
    time \
    linux-tools-generic linux-tools-common \
    && rm -rf /var/lib/apt/lists/*

# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/lib/linux-tools/5.4.0-216-generic/perf /usr/bin/perf
