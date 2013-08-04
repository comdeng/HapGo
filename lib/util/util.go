package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"
)

func UniqId() string {
	now := time.Now()
	nanoTime := int(now.Unix()) ^ now.Nanosecond()
	randNum := int(rand.ExpFloat64() * 1000000)
	h := md5.New()
	io.WriteString(h, strconv.Itoa(nanoTime)+strconv.Itoa(randNum))
	bytes := h.Sum(nil)
	return fmt.Sprintf("%x", bytes)
}
