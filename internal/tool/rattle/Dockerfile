FROM python:3.10-slim

# Avoid interactive prompts during package install
ENV DEBIAN_FRONTEND=noninteractive

# Install required system packages
RUN apt-get update && apt-get install graphviz gcc -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /opt/rattle
COPY . .

COPY measure.sh measure.sh
COPY helper.sh helper.sh

RUN chmod +x /opt/measure.sh && \
    chmod +x /opt/helper.sh

RUN apt-get update && apt-get install -y gawk jq procps time linux-perf

# Install Python dependencies
RUN pip3 install --upgrade pip && pip install --no-cache-dir -r requirements.txt

# Default entrypoint to run the Python CLI
ENTRYPOINT ["python3", "rattle-cli.py"]
