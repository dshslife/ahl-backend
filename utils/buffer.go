package utils

import (
	"io"
	"strings"
)

func ReadAll(reader io.ReadCloser) (string, int64, error) {
	buf := new(strings.Builder)
	n, e := io.Copy(buf, reader)
	return buf.String(), n, e
}
