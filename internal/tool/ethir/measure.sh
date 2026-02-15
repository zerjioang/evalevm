#!/bin/bash
export LC_NUMERIC=C

if [ $# -lt 1 ]; then
  echo '{"error":"Usage: ./measure.sh <command> [args...]"}'
  exit 1
fi

TMP_OUT=$(mktemp)
TMP_PERF=$(mktemp)
trap 'rm -f "$TMP_OUT" "$TMP_PERF"' EXIT

# Check if perf is available and working
PERF_CMD=""
if command -v perf >/dev/null 2>&1 && perf stat true >/dev/null 2>&1; then
  PERF_CMD="perf stat -x, -o $TMP_PERF --"
else
  echo "WARNING: perf not found or not working, CPU metrics will be missing" >&2
fi

# Run command with time (and optional perf)
/usr/bin/time -f 'max_ram_kb=%M\nuser_seconds=%U\nsys_seconds=%S\nreal_seconds=%e\nexit_code=%x' -o "$TMP_OUT" \
  $PERF_CMD "$@"

# Parse time output
MAX_RAM=$(grep '^max_ram_kb=' "$TMP_OUT" | cut -d= -f2)
USER_SEC=$(grep '^user_seconds=' "$TMP_OUT" | cut -d= -f2)
SYS_SEC=$(grep '^sys_seconds=' "$TMP_OUT" | cut -d= -f2)
REAL_SEC=$(grep '^real_seconds=' "$TMP_OUT" | cut -d= -f2)
EXIT_CODE=$(grep '^exit_code=' "$TMP_OUT" | cut -d= -f2)

# Calculate execution time in microseconds and milliseconds
EXEC_TIME_MS=$(awk -v r="$REAL_SEC" 'BEGIN { printf("%.0f", r * 1000) }')
EXEC_TIME_US=$(awk -v ms="$EXEC_TIME_MS" 'BEGIN { printf("%.0f", ms * 1000) }')

# Calculate average CPU usage percentage
CPU_TIME=$(awk -v u="$USER_SEC" -v s="$SYS_SEC" 'BEGIN { print u + s }')
AVG_CPU_PERCENT=$(awk -v cpu="$CPU_TIME" -v real="$REAL_SEC" 'BEGIN { if (real>0) printf("%.2f", (cpu/real)*100); else print 0 }')

# Helper to parse perf stat values
parse_perf_val() {
  grep ",$1" "$TMP_PERF" | awk -F, '{gsub(/[[:space:]]/, "", $1); print $1}'
}

INSTRUCTIONS=$(parse_perf_val instructions)
CYCLES=$(parse_perf_val cycles)
CONTEXT_SWITCHES=$(parse_perf_val 'context-switches')
PAGE_FAULTS=$(parse_perf_val 'minor-faults')
BRANCH_MISSES=$(parse_perf_val 'branch-misses')

# Ensure values are numeric
sanitize_number() {
  [[ "$1" =~ ^[0-9]+(\.[0-9]+)?$ ]] && echo "$1" || echo 0
}

# Print metrics report to stdout with a title header
echo "evalevm_perf_metrics_start"
jq -n \
  --argjson max_ram "$(sanitize_number "$MAX_RAM")" \
  --argjson exec_time_us "$(sanitize_number "$EXEC_TIME_US")" \
  --argjson exec_time_ms "$(sanitize_number "$EXEC_TIME_MS")" \
  --argjson exec_time_s "$(sanitize_number "$REAL_SEC")" \
  --argjson exit_status "$(sanitize_number "$EXIT_CODE")" \
  --argjson avg_cpu_percent "$(sanitize_number "$AVG_CPU_PERCENT")" \
  --argjson instructions "$(sanitize_number "$INSTRUCTIONS")" \
  --argjson cpu_cycles "$(sanitize_number "$CYCLES")" \
  --argjson context_switches "$(sanitize_number "$CONTEXT_SWITCHES")" \
  --argjson page_faults "$(sanitize_number "$PAGE_FAULTS")" \
  --argjson branch_misses "$(sanitize_number "$BRANCH_MISSES")" \
  '{
    max_ram_kb: $max_ram,
    exec_time_us: $exec_time_us,
    exec_time_ms: $exec_time_ms,
    exec_time_s: $exec_time_s,
    exit_status: $exit_status,
    avg_cpu_percent: $avg_cpu_percent,
    instructions: $instructions,
    cpu_cycles: $cpu_cycles,
    context_switches: $context_switches,
    page_faults: $page_faults,
    branch_misses: $branch_misses
  }'
echo "evalevm_perf_metrics_end"
