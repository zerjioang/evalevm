FROM franck44/tacas25@sha256:23d3d8c76c3d62ce3d7d1b250a73aa1e7f006658503a30f81ce71b1397fd04a6

WORKDIR /tacas25/evm-dis

COPY measure.sh measure.sh
COPY helper.sh helper.sh

RUN chmod +x /tacas25/evm-dis/measure.sh && \
    chmod +x /tacas25/evm-dis/helper.sh

# install performance related & metrics dependencies
RUN apt-get update && \
    apt-get install -y gawk jq procps time linux-tools-generic linux-tools-common

# patch perf utility - dynamically find the correct directory and overwrite wrapper
RUN REAL_PERF=$(find /usr/lib -name perf -executable -type f | grep linux-tools | head -n 1) && \
    if [ -n "$REAL_PERF" ]; then \
    mv "$REAL_PERF" /usr/bin/perf.real && \
    printf "#!/bin/bash\n/lib64/ld-linux-x86-64.so.2 /usr/bin/perf.real \"\$@\"\n" > /usr/bin/perf && \
    chmod +x /usr/bin/perf; \
    else \
    echo "Could not find native perf binary" && exit 1; \
    fi
