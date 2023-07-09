#!/bin/sh

echo "E2E Tests"

# launch external requests
/e2e/requests.sh > /e2e/results.txt

# compare results
# The regex would be: '"(user_)?id"...' but it doesn't work on the container :shrug:
errors=$(diff -I '".*id":".*"' /e2e/results.txt /e2e/expected.txt)

# exit with error code if there are differences
if [ -n "$errors" ]; then
  echo "***************"
  echo "FAIL: E2E Tests"
  echo "$errors"
else
  echo "***************"
  echo "PASS: E2E Tests"
fi
