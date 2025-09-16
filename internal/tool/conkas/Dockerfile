FROM python:3.11-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      build-essential \
      curl \
      git \
      && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/conkas
COPY . .

RUN apt-get update && apt-get install -y gawk jq procps time linux-perf && \
    chmod +x /opt/conkas/measure.sh && \
    chmod +x /opt/conkas/helper.sh

RUN pip install --no-cache-dir -r requirements.txt
RUN pip install --no-cache-dir z3-solver --prefer-binary

ENTRYPOINT ["python3", "conkas.py"]
CMD ["-h"]
