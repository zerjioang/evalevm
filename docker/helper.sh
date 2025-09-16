#!/usr/bin/env bash

set -euo pipefail

# bytecode argument
MODE="${1:-}"
BYTECODE="${2:-}"

# Only run if input is "creator"
if [[ "$MODE" == "ethersolve_creator" ]]; then
    # 1️⃣ Run the measure.sh script inside Docker
    ./measure.sh java -jar /opt/ethersolve/artifact/EtherSolve.jar --creation --tx-origin --re-entrancy --dot "$BYTECODE"
    # 3️⃣ Loop through Analysis_* files
    for f in Analysis_*; do
        echo ">>> $f"
        cat "$f"
        echo ""
        echo "<<<"
    done
elif [[ "$MODE" == "ethersolve_runtime" ]]; then
    # 1️⃣ Run the measure.sh script inside Docker
    ./measure.sh java -jar /opt/ethersolve/artifact/EtherSolve.jar --runtime --tx-origin --re-entrancy --dot "$BYTECODE"
    # 3️⃣ Loop through Analysis_* files
    for f in Analysis_*; do
        echo ">>> $f"
        cat "$f"
        echo ""
        echo "<<<"
    done
else
    # Fallback: pass all arguments to measure.sh
    ./measure.sh "$@"
fi
