#!/bin/sh

echo "Running E2E Tests"

# launch external requests
/e2e/requests.sh > /e2e/results.txt

# compare results
echo "Comparing Results"
# compare results, ignoring id/user_id fields (values differ per run)
# Uses grep -v instead of diff -I, which is not supported in BusyBox diff
grep -v '".*id":' /e2e/results.txt > /tmp/results_no_id.txt
grep -v '".*id":' /e2e/expected.txt > /tmp/expected_no_id.txt
errors=$(diff /tmp/results_no_id.txt /tmp/expected_no_id.txt)

# exit with error code if there are differences
if [ -n "$errors" ]; then
  echo "***************"
  echo "FAIL: E2E Tests"
  echo "$errors"
else
  echo "***************"
  echo "PASS: E2E Tests"
fi
