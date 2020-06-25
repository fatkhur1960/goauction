package utils

import (
	"strings"
	"time"
)

// NOW generate current datetime
var NOW = time.Now().UTC()

// ReplacePackages --
func ReplacePackages(input string) string {
	paths := strings.Split(input, "/")
	return paths[len(paths)-1]
}
