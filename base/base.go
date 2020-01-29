package base

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func Check(a interface{}) {
	if a == nil {
		_, f, l, _ := runtime.Caller(1)
		strs := strings.Split(f, "/")
		size := len(strs)
		if size > 4 {
			size = 4
		}
		f = filepath.Join(strs[len(strs)-size:]...)
		logrus.Printf("Object can't be nil, cause by: %s(%d)", f, l)
	}
}
