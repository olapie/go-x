package xapp

import (
	"fmt"
	"os"
)

/*
*

	LD_FLAGS="-X 'go.olapie.com/x/xapp.BinaryVersion=1.0.0' \

-X 'go.olapie.com/x/xapp.BinaryBuildTime=$(date -u +"%Y-%m-%dT%H:%MZ")' \
-X 'go.olapie.com/x/xapp.BinarySourceCommitID=$(git rev-parse HEAD)' \
-X 'go.olapie.com/x/xapp.BinarySourceShortCommitID=$(git rev-parse --short HEAD)' \
-s -w"
*/
var (
	BinaryVersion             string
	BinaryBuildTime           string
	BinarySourceShortCommitID string
)

func CheckVersionArgument(binaryName string) {
	if len(os.Args) != 2 {
		return
	}
	switch os.Args[1] {
	case "-v", "-version", "--version", "version":
		fmt.Println(binaryName, "version", BinaryVersion+"-"+BinarySourceShortCommitID, BinaryBuildTime)
		os.Exit(0)
	}
}
