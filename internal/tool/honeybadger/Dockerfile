FROM christoftorres/honeybadger

COPY datasets/honeypots honeypots
COPY honeybadger honeybadger
COPY measure.sh measure.sh
COPY helper.sh helper.sh

RUN chmod +x /opt/measure.sh && \
    chmod +x /opt/helper.sh

RUN apt-get update && apt-get install -y gawk jq procps time linux-tools-generic linux-tools-common
# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/lib/linux-tools/4.4.0-210-generic/perf /usr/bin/perf
