set -e

COMMIT_MSG=$(git log -1 --pretty=format:"%s")
echo "Latest commit message: $COMMIT_MSG"

if [ "$COMMIT_MSG" != "Bump version" ]; then
  echo "Not a version bump commit. Skipping release."
  exit 0
fi

VERSION=$(cat VERSION | tr -d '[:space:]')
echo "Releasing version: $VERSION"

git config user.name "NeetoBot"
git config user.email "bot@neeto.com"
git tag -a "v${VERSION}" -m "Release v${VERSION}"
git push origin "v${VERSION}"

go test ./...
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
