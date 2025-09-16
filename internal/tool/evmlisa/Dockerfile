# Use a base image with JDK 11
FROM openjdk:11-jdk-slim

# Specify the version of Gradle to install
ARG GRADLE_VERSION=8.0

# Install Gradle
RUN apt-get update && apt-get install -y wget unzip git && \
    wget https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip && \
    mkdir /opt/gradle && \
    unzip -d /opt/gradle gradle-${GRADLE_VERSION}-bin.zip && \
    rm gradle-${GRADLE_VERSION}-bin.zip && \
    ln -s /opt/gradle/gradle-${GRADLE_VERSION}/bin/gradle /usr/bin/gradle && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Clone the EVMLiSA repository (latest stable version)
# RUN git clone https://github.com/lisa-analyzer/evm-lisa.git .

# Copy the entire EVMLiSA project into the container (latest version)
COPY . .

COPY measure.sh measure.sh
COPY helper.sh helper.sh

RUN chmod +x /opt/measure.sh && \
    chmod +x /opt/helper.sh

RUN apt-get update && apt-get install -y gawk jq procps time linux-perf
# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/bin/perf_5.10 /usr/bin/perf

# Build the project with Gradle
RUN gradle assemble

# Create results directory
RUN mkdir -p /app/execution/results

# Set the entry point of the container
ENTRYPOINT ["java", "-jar", "/app/build/libs/evm-lisa-all.jar"]
