FROM openjdk:8-jdk

WORKDIR /opt/ethersolve
COPY . .

RUN chmod +rwx measure.sh
RUN apt-get update && apt-get install -y gawk jq procps time linux-perf
RUN ./gradlew -v && \
    ./gradlew jar

# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/bin/perf_5.10 /usr/bin/perf
