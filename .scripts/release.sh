#!/bin/bash
set -euo pipefail

# Install gh CLI
type -p curl >/dev/null || sudo apt install curl -y
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list >/dev/null
sudo apt update
sudo apt install gh -y

# Detect version label from the most recently merged PR
PR_NUMBER=$(gh pr list --state merged --base main --limit 1 --json number --jq ".[0].number")
echo "Last merged PR number: $PR_NUMBER"

if [ -z "$PR_NUMBER" ]; then
  echo "No merged PR found. Skipping release."
  exit 0
fi

PR_LABELS=$(gh pr view "$PR_NUMBER" --json labels --jq ".labels[].name" | tr "\n" " ")
echo "PR labels: $PR_LABELS"

VERSION_LABEL=""
for label in major minor patch; do
  if echo "$PR_LABELS" | grep -qw "$label"; then
    VERSION_LABEL="$label"
    break
  fi
done

echo "Version label selected: $VERSION_LABEL"

if [ -z "$VERSION_LABEL" ]; then
  echo "No version label found. Skipping release."
  exit 0
fi

# Release current version
VERSION=$(cat VERSION | tr -d '[:space:]')
echo "Releasing version: $VERSION"

git config user.name "NeetoBot"
git config user.email "bot@neeto.com"
git config core.hooksPath /dev/null

if git rev-parse "v${VERSION}" >/dev/null 2>&1; then
  echo "Tag v${VERSION} already exists. Skipping tag creation."
else
  git tag -a "v${VERSION}" -m "Release v${VERSION}"
  echo "Created tag v${VERSION}"
  git push origin "v${VERSION}"
  echo "Pushed tag v${VERSION}"
fi

export GORELEASER_CURRENT_TAG="v${VERSION}"

echo "Running tests..."
go test ./...
echo "Tests passed."

echo "GoReleaser version:"
goreleaser --version || true

echo "Running goreleaser release..."
goreleaser release --clean
echo "GoReleaser release complete."

echo "Uploading to S3 versioned directory..."
aws s3 cp dist/ "s3://neeto-downloads/cli/NeetoCal/v${VERSION}/" --recursive --exclude "*" --include "*.tar.gz" --include "*.zip" --include "checksums.txt"

echo "Uploading to S3 latest directory..."
aws s3 rm s3://neeto-downloads/cli/NeetoCal/latest/ --recursive
aws s3 cp "dist/neeto-cal-cli_${VERSION}_linux_amd64.tar.gz" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_linux_amd64.tar.gz
aws s3 cp "dist/neeto-cal-cli_${VERSION}_linux_arm64.tar.gz" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_linux_arm64.tar.gz
aws s3 cp "dist/neeto-cal-cli_${VERSION}_darwin_amd64.tar.gz" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_macos_amd64.tar.gz
aws s3 cp "dist/neeto-cal-cli_${VERSION}_darwin_arm64.tar.gz" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_macos_arm64.tar.gz
aws s3 cp "dist/neeto-cal-cli_${VERSION}_windows_amd64.zip" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_windows_amd64.zip
aws s3 cp "dist/neeto-cal-cli_${VERSION}_windows_arm64.zip" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_windows_arm64.zip
aws s3 cp dist/checksums.txt s3://neeto-downloads/cli/NeetoCal/latest/checksums.txt
echo "S3 upload complete."

# Bump version for next release
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"

case "$VERSION_LABEL" in
  major)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
  minor)
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  patch)
    PATCH=$((PATCH + 1))
    ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
echo "$NEW_VERSION" > VERSION
echo "Bumped version: $VERSION -> $NEW_VERSION"

git push origin --delete bump-version || true
git checkout -b bump-version
git add VERSION
git commit -m "Bump version to $NEW_VERSION"
git push --set-upstream origin bump-version
gh pr create -B main -H bump-version -t "Bump version to $NEW_VERSION" -b "" -l instant-mergepr
echo "Version bump PR created."
