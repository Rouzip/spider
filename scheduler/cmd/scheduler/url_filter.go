package main

import "strings"

const domain = "https://wiki.archlinux.org"

func isLegitURL(u string) bool {
	return strings.HasPrefix(u, domain)
}
