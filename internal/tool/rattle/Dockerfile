FROM python:3.10-slim

# Avoid interactive prompts during package install
ENV DEBIAN_FRONTEND=noninteractive

# Install required system and build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    build-essential \
    python3-dev \
    graphviz \
    gawk \
    jq \
    procps \
    time \
    linux-perf \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /opt/rattle

# Copy project files
COPY . .

# Ensure helper scripts are executable
RUN chmod +x /opt/rattle/measure.sh /opt/rattle/helper.sh

# 2025: patch cbor2 version
RUN sed -i.bak 's/^cbor2==.*/cbor2>=5.4.6/' requirements.txt

# Install Python dependencies (ensure latest build tools)
RUN pip install --upgrade pip setuptools wheel \
    && pip install --no-cache-dir -r requirements.txt

# Default entrypoint to run the Python CLI
ENTRYPOINT ["python3", "rattle-cli.py"]
