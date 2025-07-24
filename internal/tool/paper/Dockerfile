FROM paper-algorithm:latest

COPY measure.sh measure.sh
RUN chmod +rwx measure.sh
RUN apt-get update && apt-get install -y gawk jq procps time linux-perf
