VERSION="0.2.1"
VERSION_INCREMENT=1
NAME=secret-wallet

FLAGS="-X github.com/secretnamebasis/secret-wallet/settings.Version=v$VERSION"
FLAGS="$FLAGS -X github.com/secretnamebasis/secret-wallet/settings.BuildTime=$(date +%s)"
FLAGS="$FLAGS -X github.com/secretnamebasis/secret-wallet/settings.GitVersion=$(git describe --tags --dirty --always)"

APPID="com.github.secretnamebasis.secret_wallet"

export FLAGS APPID VERSION VERSION_INCREMENT
