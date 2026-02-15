FROM eclipse-temurin:8-jdk

WORKDIR /opt/ethersolve
RUN apt-get update && apt-get install -y git gawk jq procps time linux-tools-generic && \
    rm -rf /var/lib/apt/lists/*

# Clone repo
RUN git clone https://github.com/SeUniVr/EtherSolve .

# Build
ENV GRADLE_OPTS="-Dorg.gradle.daemon=false -Xmx1g"
RUN chmod +x gradlew && ./gradlew --no-daemon jar

# Copy helper scripts
COPY measure.sh .
COPY helper.sh .
RUN chmod +x measure.sh helper.sh

# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/lib/linux-tools/*/perf /usr/bin/perf
