# Use a base image with JDK 11
FROM eclipse-temurin:11-jdk-jammy

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
WORKDIR /opt/evmlisa

# Clone the EVMLiSA repository (latest stable version)
# RUN git clone https://github.com/lisa-analyzer/evm-lisa.git .

# Copy the entire EVMLiSA project into the container (latest version)
COPY . .

RUN chmod +x /opt/evmlisa/measure.sh && \
    chmod +x /opt/evmlisa/helper.sh

RUN apt-get update && apt-get install -y gawk jq procps time linux-tools-generic linux-tools-common
# patch perf utility
RUN rm /usr/bin/perf && ln -s /usr/lib/linux-tools/*/perf /usr/bin/perf

# Build the project with Gradle
RUN gradle assemble

# Create results directory
RUN mkdir -p /opt/evmlisa/execution/results

# Set the entry point of the container
ENTRYPOINT ["java", "-jar", "/opt/evmlisa/build/libs/evm-lisa-all.jar"]
