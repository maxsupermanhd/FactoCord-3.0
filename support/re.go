package support

import "regexp"

var ModFileRegexp = regexp.MustCompile(`([A-Za-z0-9_\- ]+)_(\d+\.\d+\.\d+)(\.zip)?`)
