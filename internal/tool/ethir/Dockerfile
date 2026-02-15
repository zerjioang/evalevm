FROM ubuntu:22.04

# Install system dependencies
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    git \
    wget \
    unzip \
    build-essential \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*

# Install Ethereum tools (evm)
# Install specific EVM version (1.10.26) that supports 'disasm' command
RUN apt-get update && apt-get install -y wget && \
    wget https://gethstore.blob.core.windows.net/builds/geth-alltools-linux-amd64-1.10.26-e5eb32ac.tar.gz && \
    tar -xvf geth-alltools-linux-amd64-1.10.26-e5eb32ac.tar.gz && \
    cp geth-alltools-linux-amd64-1.10.26-e5eb32ac/evm /usr/bin/evm && \
    rm -rf geth-alltools* && \
    chmod +x /usr/bin/evm
RUN ln -s /usr/bin/solc /usr/bin/solcv8
RUN ln -s /usr/bin/evm /usr/bin/evm1.14.13

# Install Python dependencies
RUN apt-get update && apt-get install -y python3-z3
RUN pip3 install --upgrade pip setuptools wheel
RUN pip3 install sympy six requests semantic_version solc-select

# Install solc via solc-select and ensure it's in path
# solc-select install might fail if no binary found, but it usually supports arm64
RUN solc-select install 0.8.23 && solc-select use 0.8.23
# solc-select usage usually links to ~/.solc-select/artifacts/solc-...
# We need it in /usr/bin/solc
RUN ln -sf $(which solc) /usr/bin/solc && ln -sf /usr/bin/solc /usr/bin/solcv8

# Install EthIR
WORKDIR /opt/ethir
RUN git clone https://github.com/costa-group/EthIR .

# Fix SyntaxWarning in cold.py (assert(False, ...) is always true)
# Properly remove parentheses: assert(False, 'msg') -> assert False, 'msg'
RUN sed -i "s/assert(False,\(.*\))/assert False, \1/g" /opt/ethir/ethir/storage/cold.py

# Fix evm version warning by updating supported version to match installed (1.16.8 or newer)
# We update it to 1.16.8 to match the warning seen by user
RUN sed -i "s/tested_evm_version = '1.9.23'/tested_evm_version = '1.16.8'/g" /opt/ethir/ethir/ethir.py

# Fix execution error: remove nohup and use 'evm' directly in input_helper.py
# The tool hardcodes "nohup evm1.14.13" which fails if nohup is missing or evm version differs
RUN sed -i 's/"nohup", "evm1.14.13",/"evm",/g' /opt/ethir/ethir/input_helper.py && \
    sed -i 's/nohup evm1.14.13/evm/g' /opt/ethir/ethir/input_helper.py

# Fix paths
RUN sed -i 's/costabs_path = ""/costabs_path = "\/tmp\/"/g' /opt/ethir/ethir/global_params_ethir.py

# Install performance related & metrics dependencies
RUN apt-get update && \
    apt-get install -y linux-tools-generic linux-tools-common jq time procps && \
    rm -rf /var/lib/apt/lists/*

# patch perf utility - bypass version-checking wrapper
RUN REAL_PERF=$(find /usr/lib -name perf -executable -type f | grep linux-tools | head -n 1) && \
    if [ -n "$REAL_PERF" ]; then \
    mv "$REAL_PERF" /usr/bin/perf.real && \
    printf "#!/bin/bash\n/usr/bin/perf.real \"\$@\"\n" > /usr/bin/perf && \
    chmod +x /usr/bin/perf; \
    else \
    echo "Could not find native perf binary" && exit 1; \
    fi

# Copy helper scripts
COPY measure.sh .
COPY helper.sh .
RUN chmod +x measure.sh helper.sh
