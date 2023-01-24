//go:build !linux
// +build !linux

package util

import (
	"os"
)

func OpenTAP(ifaceName string) (*os.File, error) {
	return nil, FeatUnsupported{Feat: "tap"}
}
