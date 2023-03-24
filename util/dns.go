package util

import (
	"strings"
)

func ResolveIp(name string) string {
	index := strings.Index(name, ".in-addr.arpa.")
	if index == -1 {
		index = strings.Index(name, ".ip6.arpa.")
	}
	arr := strings.Split(name[0:index], ".")
	ArrayReverse(arr)
	return strings.Join(arr, ".")
}
