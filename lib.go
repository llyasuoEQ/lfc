package lfc

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"
)

func getMd5Bytes(bytes []byte) string {
	hashOb := md5.New()
	hashOb.Write(bytes)
	return hex.EncodeToString(hashOb.Sum(nil))
}

func getStrListHash(strs []string) string {
	str := strings.Join(strs, ",")
	return getMd5Bytes([]byte(str))
}

func containInStringSlice(source []string, find string) (result bool) {
	sort.Strings(source)
	index := sort.SearchStrings(source, find)
	if source[index] == find {
		result = true
	}
	return
}
