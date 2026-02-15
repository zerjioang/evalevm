#!/usr/bin/env bash

set -euo pipefail

# bytecode argument
MODE="${1:-}"
BYTECODE="${2:-}"

if [[ "$MODE" == "ethersolve" ]]; then
    RUN_MODE="${2:-}"
    AUDIT="${3:-}"
    BYTECODE="${4:-}"
    
    FLAGS="--dot"
    if [[ "$RUN_MODE" == "creator" ]]; then
        FLAGS="--creation $FLAGS "
    else
        # default to runtime for 'runtime' or 'any'
        FLAGS="--runtime $FLAGS "
    fi

    if [[ "$AUDIT" == "true" ]]; then
        FLAGS=" --tx-origin --re-entrancy $FLAGS"
    fi

    # run the measure.sh script inside Docker
    ./measure.sh java -jar /opt/ethersolve/artifact/EtherSolve.jar $FLAGS "$BYTECODE"
    # Loop through Analysis_* files
    for f in Analysis_*; do
        echo ">>> $f"
        cat "$f"
        echo ""
        echo "<<<"
    done
elif [[ "$MODE" == "evmlisa" ]]; then
    # Shift the first two arguments (MODE and BYTECODE) to get the rest
    shift 2
    EXTRA_ARGS="$@"
    # run the measure.sh script inside Docker
    ./measure.sh java -jar /opt/evmlisa/build/libs/evm-lisa-all.jar --show-all-instructions-in-cfg --paper-stats --bytecode "$BYTECODE" $EXTRA_ARGS
    if [ -d "execution/results" ]; then
      find execution/results -type f | sort | while read -r f; do
        echo ">>> $f"
        cat "$f"
        echo ""
        echo "<<<"
      done
    fi
elif [[ "$MODE" == "rattle" ]]; then
    # run the measure.sh script inside Docker
    echo "$BYTECODE" > code.evm
    ./measure.sh python3 /opt/rattle/rattle-cli.py --no-split-functions --optimize --input code.evm
    rm code.evm
    for f in *.dot; do
      if [ -f "$f" ]; then
        echo ">>> $f"
        cat "$f"
        echo ""
        echo "<<<"
      fi
    done
elif [[ "$MODE" == "vandal" ]]; then
    # run the measure.sh script inside Docker
    echo "$BYTECODE" > code.evm
    ./measure.sh python3 /vandal/bin/decompile -n -v -g graph.dot code.evm
    rm code.evm
    for f in *.dot; do
      if [ -f "$f" ]; then
        echo ">>> $f"
        cat "$f"
        echo ""
        echo "<<<"
      fi
    done
else
    # Fallback: pass all arguments to measure.sh
    ./measure.sh "$@"
fi
