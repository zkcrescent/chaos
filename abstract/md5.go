package abstract

import (
	"crypto/md5"
	"fmt"
)

// MD5, only for small size of []byte
func MD5(bts []byte) string {
	h := md5.New()
	h.Write(bts)
	return fmt.Sprintf("%x", h.Sum(nil))
}
