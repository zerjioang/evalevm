# very time consuming to compile security in 2025
FROM smartbugs/securify:usolc

COPY measure.sh measure.sh
RUN chmod +rwx measure.sh
RUN apt-get update && apt-get install -y gawk jq procps time linux-tools-generic linux-tools-common

# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/lib/linux-tools/4.15.0-213-generic/perf /usr/bin/perf
