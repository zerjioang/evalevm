FROM franck44/tacas25@sha256:23d3d8c76c3d62ce3d7d1b250a73aa1e7f006658503a30f81ce71b1397fd04a6

WORKDIR /tacas25/evm-dis

COPY measure.sh measure.sh
COPY helper.sh helper.sh

RUN chmod +x /tacas25/evm-dis/measure.sh && \
    chmod +x /tacas25/evm-dis/helper.sh

# install performance related & metrics dependencies
RUN apt-get update && \
    apt-get install -y gawk jq procps time linux-tools-generic linux-tools-common

# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/lib/linux-tools/6.8.0-83-generic/perf /usr/bin/perf
