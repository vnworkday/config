#!/usr/bin/env bash

set -eufo pipefail

check_verify() {
  go mod tidy
  go mod verify
  exit_code=$?

  if [[ $exit_code -ne 0 ]]; then
    echo "🚫 go mod tidy or go mod verify failed."
    exit 1
  fi

  echo "✅  Code complies with go mod requirements."
}

check_import() {
  # We only require goimports to have been run on files that were changed
  # relative to the main branch, so that we can gradually create more consistency
  # rather than bulk-changing everything at once.
  # If we seem to be running inside a GitHub Actions pull request check
  # then we'll use the PR's target branch from this variable instead.
  base_branch=${GITHUB_BASE_REF:-"origin/main"}

  # Get the list of files that were changed relative to the main branch
  declare -a target_files
  while IFS= read -r line; do
    target_files+=("$line")
  done < <(git diff --name-only origin/main --diff-filter=MA | grep "\.go" | grep -v ".pb.go" | grep -v ".go-version" | grep -v ".golangci.yaml")

  if [[ ${#target_files[@]} -eq 0 ]]; then
    echo "🚫 No Go files changed relative to $base_branch, skipping import check."
    return
  fi

#  echo "🔍 Checking the following files:"
#  for file in "${target_files[@]}"; do
#    echo "  - $file"
#  done

  declare -a bad_files
  bad_files=()
  for file in "${target_files[@]}"; do
    output=$(go run golang.org/x/tools/cmd/goimports -l "${file}")
    exit_code=$?

    [[ $exit_code -ne 0 ]] && echo "🚫 Failed to run goimports on $file" && exit 1
    [[ -n "$output" ]] && bad_files+=("$file")
  done

  if [[ ${#bad_files[@]} -gt 1 ]]; then
    echo "⚠️ The following files import statements that disagree with \"goimports\":"
    for file in "${bad_files[@]}"; do
      echo "  - $file"
    done
    # Read the auto-fix flag from the input. if it's set to "true" then we'll run goimports with the -w flag
    # to automatically fix the import statements in the files that need it.
    if [[ $# -ge 1 && "$1" == "true" ]]; then
      echo "🔧 Running goimports with the -w flag to automatically fix import statements."
      go run golang.org/x/tools/cmd/goimports -w -l "${bad_files[@]}"
    else
      echo "🚫 Please run \"scripts/import-check.sh true\" on these files to fix their imports."
      exit 1
    fi
  fi

  echo "✅  Code complies with go imports requirements."
}

check_golangci() {
  golangci-lint run --fix
  exit_code=$?

  if [[ $exit_code -ne 0 ]]; then
    echo "🚫 Linting failed."
    return
  fi

  echo "✅  Code complies with linting requirements."
}

check_verify
check_import true
check_golangci

echo "==============================="
echo "🎉 All checks passed."