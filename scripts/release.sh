#!/usr/bin/env bash

set -euo pipefail

usage() {
    echo "📋 Usage: $0 -t/--type <major|minor|patch>"
    echo "✅ Example: $0 -t patch"
    exit 0
}

check_branch() {
    echo "🔍 Checking branch..."
    echo "🌿 Current branch: $current_branch"
    if [[ ! ("$current_branch" == "main") ]]; then
        echo "⚠️ You must be on the release branch."
        exit 1
    fi
}

check_version() {
    echo "🔍 Checking version..."
    echo "📦 Package version: $(node -p "require('./package.json').version")"
    latest_git_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
    pkg_version="v$(node -p "require('./package.json').version")"
    if [[ $latest_git_tag != "$pkg_version" ]]; then
        echo "⚠️ Latest tag is not the same as the package version. Please check the package version."
        exit 1
    fi
}

check_dirty() {
    echo "🔍 Checking git status..."
    if [[ -n $(git status -s) ]]; then
        echo "⚠️ Git is dirty. Please commit your changes before releasing."
        exit 1
    fi
}

release() {
    echo "🚀 Releasing with type: $release_type"

    new_tag="$(npm version "$release_type" --no-commit-hooks --no-git-tag-version)"
    echo "🔖 Latest tag: $latest_git_tag"
    echo "🏷️ New tag: $new_tag"

    echo "📦 Staging files and committing..."
    git add .
    git commit -m "chore: prepare to release Go proto module: $new_tag"
    git push origin "$current_branch"

    echo -e "\n🚛 Pushing changes to the remote repository...\n"
    git tag -a "$new_tag" -m "Release $new_tag"
    git push origin "$new_tag"

    echo -e "\n🎉 Done! Released $new_tag"
}

release_type="patch"

while [ "$#" -ge 2 ]; do
    case "$1" in
        -t|--type) release_type="$2"; shift 2;;
        *) echo "❗️ Invalid argument: $1"; usage;;
    esac
done

current_branch=$(git rev-parse --abbrev-ref HEAD)

check_branch
check_version
check_dirty
release