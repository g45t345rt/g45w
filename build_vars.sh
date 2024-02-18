VERSION="v0.10.1"
NAME=G45W

FLAGS="-X github.com/g45t345rt/g45w/settings.Version=$VERSION"
FLAGS="$FLAGS -X github.com/g45t345rt/g45w/settings.BuildTime=$(date +%s)"
FLAGS="$FLAGS -X github.com/g45t345rt/g45w/settings.GitVersion=$(git describe --tags --dirty --always)"

APPID="com.github.g45t345rt.g45w"

export FLAGS APPID VERSION