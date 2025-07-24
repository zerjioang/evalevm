FROM franck44/tacas25:latest

WORKDIR /tacas25/evm-dis
COPY measure.sh measure.sh
RUN chmod +rwx measure.sh

# install performance related & metrics dependencies
RUN apt-get update && \
    apt-get install -y gawk jq procps time linux-tools-generic linux-tools-common

# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/lib/linux-tools/6.8.0-64-generic/perf /usr/bin/perf
