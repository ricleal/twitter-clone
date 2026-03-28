#!/bin/sh

set -e

echo "Running E2E Tests"

# launch external requests
/e2e/requests.sh > /tmp/results.txt

# compare results
echo "Comparing Results"
if ! diff /tmp/results.txt /e2e/expected.txt > /tmp/e2e_diff.txt; then
  echo "***************"
  echo "FAIL: E2E Tests"
  cat /tmp/e2e_diff.txt
  exit 1
else
  echo "***************"
  echo "PASS: E2E Tests"
fi
