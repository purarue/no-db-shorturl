#!/bin/bash
# Modified from: https://gist.github.com/DimaKoz/06b7475317b12e7ffa724ef0e115a4ec

THIS_DIR="$(dirname "${BASH_SOURCE[0]}")"
cd "$THIS_DIR"

package=$(basename "$(pwd)")
if [[ -z "$package" ]]; then
	echo "usage: $0 <package-name>"
	exit 1
fi
package_name=$package

#the full list of the platforms: https://golang.org/doc/install/source#environment
platforms=(
  "darwin/amd64"
#  "darwin/arm64"
	"freebsd/amd64"
	"freebsd/arm"
	"linux/amd64"
	"linux/arm"
	"linux/arm64"
	"netbsd/amd64"
	"netbsd/arm"
	"openbsd/amd64"
	"openbsd/arm"
	"solaris/amd64"
	"windows/amd64"
	"windows/386")

BUILD_DIR="${THIS_DIR}/build"
[ -d "$BUILD_DIR" ] && rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"
cd "$BUILD_DIR"

for platform in "${platforms[@]}"; do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name="${package_name}-${GOOS}-${GOARCH}"
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

  echo "building ${platform}..."

	env GOOS=$GOOS GOARCH=$GOARCH go build -o "${output_name}" ..
	if [ $? -ne 0 ]; then
		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done
