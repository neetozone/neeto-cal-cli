set -e

install_gh() {
  type -p curl >/dev/null || sudo apt install curl -y
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
  sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list >/dev/null
  sudo apt update
  sudo apt install gh -y
}

detect_version_label() {
  PR_NUMBER=$(gh pr list --state merged --base main --limit 1 --json number --jq ".[0].number")
  echo "Last merged PR number: $PR_NUMBER"

  [[ -z "$PR_NUMBER" ]] && echo "No merged PR found." && exit 0

  PR_LABELS=$(gh pr view $PR_NUMBER --json labels --jq ".labels[].name" | tr "\n" " " | cut -d " " -f 1-)
  echo "PR labels: $PR_LABELS"

  [[ -z "$PR_LABELS" ]] && exit 0

  ALL_VERSION_LABELS=(major minor patch)

  PR_VERSION_LABELS=()
  for version_label in "${ALL_VERSION_LABELS[@]}"; do
    if echo "$PR_LABELS" | grep -q -o "$version_label"; then
      PR_VERSION_LABELS+=("$version_label")
    fi
  done

  VERSION_LABEL=$(echo "${PR_VERSION_LABELS[0]}" | xargs)
  echo "Version label selected: $VERSION_LABEL"

  [[ -z "$VERSION_LABEL" ]] && exit 0
}

release() {
  VERSION=$(cat VERSION | tr -d '[:space:]')
  echo "Releasing version: $VERSION"

  git config user.name "NeetoBot"
  git config user.email "bot@neeto.com"
  git config core.hooksPath /dev/null

  if git rev-parse "v${VERSION}" >/dev/null 2>&1; then
    echo "Tag v${VERSION} already exists. Skipping tag creation."
  else
    git tag -a "v${VERSION}" -m "Release v${VERSION}"
    echo "Tag v${VERSION} created."
    git push origin "v${VERSION}"
    echo "Tag v${VERSION} pushed."
  fi

  export GORELEASER_CURRENT_TAG="v${VERSION}"
  echo "Running tests..."
  go test ./...
  echo "Tests passed."
  echo "GoReleaser version:"
  goreleaser --version || true
  echo "Running goreleaser release..."
  goreleaser release --clean

  aws s3 cp dist/ "s3://neeto-downloads/cli/NeetoCal/v${VERSION}/" --recursive --exclude "*" --include "*.tar.gz" --include "*.zip" --include "checksums.txt"

  aws s3 rm s3://neeto-downloads/cli/NeetoCal/latest/ --recursive
  aws s3 cp "dist/neeto-cal-cli_${VERSION}_linux_amd64.tar.gz" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_linux_amd64.tar.gz
  aws s3 cp "dist/neeto-cal-cli_${VERSION}_linux_arm64.tar.gz" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_linux_arm64.tar.gz
  aws s3 cp "dist/neeto-cal-cli_${VERSION}_darwin_amd64.tar.gz" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_macos_amd64.tar.gz
  aws s3 cp "dist/neeto-cal-cli_${VERSION}_darwin_arm64.tar.gz" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_macos_arm64.tar.gz
  aws s3 cp "dist/neeto-cal-cli_${VERSION}_windows_amd64.zip" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_windows_amd64.zip
  aws s3 cp "dist/neeto-cal-cli_${VERSION}_windows_arm64.zip" s3://neeto-downloads/cli/NeetoCal/latest/neetocal_windows_arm64.zip
  aws s3 cp dist/checksums.txt s3://neeto-downloads/cli/NeetoCal/latest/checksums.txt
}

bump_version() {
  CURRENT_VERSION=$(cat VERSION | tr -d '[:space:]')
  IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

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
  echo "Bumped version: $CURRENT_VERSION -> $NEW_VERSION"

  git add VERSION
  git commit -m "Bump version to $NEW_VERSION"
  git push origin main
}

install_gh
detect_version_label
release
bump_version
