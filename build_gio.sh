#!/bin/bash

if [[ $# -ne 2 ]]; then
  echo "Usage: $0 <OS> <ARCH>"
  echo "Example: $0 windows amd64"
  exit 1
fi

GOOS=$1
GOARCH=$2
OUTPUT="./build/g45w_${GOOS}_${GOARCH}"

source ./build_vars.sh

if [ $GOOS = "windows" ]; then
  OUTPUT+=".exe"
fi

if [ $GOOS = "android" ]; then
  OUTPUT+=".apk"
fi

if [ $GOOS = "macos" ]; then
  OUTPUT+=".app"
fi

## gogio commands
# target        = flag.String("target", "", "specify target (ios, tvos, android, js).\n")
#	archNames     = flag.String("arch", "", "specify architecture(s) to include (arm, arm64, amd64).")
#	minsdk        = flag.Int("minsdk", 0, "specify the minimum supported operating system level")
#	buildMode     = flag.String("buildmode", "exe", "specify buildmode (archive, exe)")
#	destPath      = flag.String("o", "", "output file or directory.\nFor -target ios or tvos, use the .app suffix to target simulators.")
#	appID         = flag.String("appid", "", "app identifier (for -buildmode=exe)")
#	name          = flag.String("name", "", "app name (for -buildmode=exe)")
#	version       = flag.Int("version", 1, "app version (for -buildmode=exe)")
#	printCommands = flag.Bool("x", false, "print the commands")
#	keepWorkdir   = flag.Bool("work", false, "print the name of the temporary work directory and do not delete it when exiting.")
#	linkMode      = flag.String("linkmode", "", "set the -linkmode flag of the go tool")
#	extraLdflags  = flag.String("ldflags", "", "extra flags to the Go linker")
#	extraTags     = flag.String("tags", "", "extra tags to the Go tool")
#	iconPath      = flag.String("icon", "", "specify an icon for iOS and Android") !!use appicon.png by default!!
#	signKey       = flag.String("signkey", "", "specify the path of the keystore to be used to sign Android apk files.")
#	signPass      = flag.String("signpass", "", "specify the password to decrypt the signkey.")

gogio -target $GOOS -arch $GOARCH -x -ldflags "$FLAGS" -appid $APPID -version $VERSION_INCREMENT -o "$OUTPUT" .
