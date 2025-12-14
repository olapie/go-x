SCRIPT_DIR=$(dirname "$0")
BUILD_DIR="$SCRIPT_DIR/build/ios"
IOS_FRAMEWORK="$BUILD_DIR"/mobilex.xcframework

MODULES="
go.olapie.com/mob
"

rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

#export GO111MODULE=off
export GOPROXY=direct
export GOSUMDB=off
gomobile bind -v  -target=ios -o "$IOS_FRAMEWORK" $MODULES