VERSION="0.4.1"
VERSION_INCREMENT=1
NAME=G45W

FLAGS="-X github.com/g45t345rt/g45w/settings.Version=v$VERSION"
FLAGS="$FLAGS -X github.com/g45t345rt/g45w/settings.BuildTime=$(date +%s)"
FLAGS="$FLAGS -X github.com/g45t345rt/g45w/settings.GitVersion=$(git describe --tags --dirty --always)"

APPID="com.github.g45t345rt.g45w"

export FLAGS APPID VERSION VERSION_INCREMENT