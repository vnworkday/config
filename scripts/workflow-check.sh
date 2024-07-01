#!/usr/bin/env bash

set -euo pipefail

for workflow in $(git diff --cached --name-only --diff-filter=ACM | grep -E "^\.github/workflows/.*\.yml$"); do
    if npm exec action-validator "$workflow"; then
        echo "✅ $workflow is good to go!"
    else
        echo "❌ $workflow is invalid. Please fix the errors and try again."
        exit 1
    fi
done

echo "👍 All workflows are good to go!"