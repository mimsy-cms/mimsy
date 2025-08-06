package util

import (
	"fmt"
	"net/http"
	"strconv"
)

func PathInt(r *http.Request, name string) (int64, error) {
	value := r.PathValue(name)
	if value == "" {
		return 0, fmt.Errorf("missing path parameter: %s", name)
	}
	return strconv.ParseInt(value, 10, 64)
}
