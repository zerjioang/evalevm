FROM openjdk:8-jdk

WORKDIR /opt/ethersolve
COPY . .

RUN apt-get update && apt-get install -y gawk jq procps time linux-perf && \
    chmod +x /opt/ethersolve/measure.sh && \
    chmod +x /opt/ethersolve/helper.sh && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

RUN ./gradlew -v && \
    ./gradlew jar

# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/bin/perf_5.10 /usr/bin/perf
