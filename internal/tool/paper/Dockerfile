FROM paper-algorithm:latest

WORKDIR /opt

RUN apt-get update && \
    apt-get install -y --no-install-recommends gawk jq procps time linux-perf && \
    rm -rf /var/lib/apt/lists/*

# Copy scripts after dependencies (these change more often)
COPY measure.sh measure.sh
COPY helper.sh helper.sh

# Set executable permissions in a single layer
RUN chmod +x /opt/measure.sh /opt/helper.sh
