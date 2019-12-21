package config

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	data := []struct {
		file   string
		hasErr bool
		errStr string
	}{
		{"file1", true, "unknown field"},
		{"file2", true, "must be numeric"},
		// {"file3", true, "must be boolean"},
	}
	for i, d := range data {
		err := Load("./testData/" + d.file + ".ini")
		if d.hasErr {
			if err == nil {
				t.Fatalf("test load config failed:%d", i)
			} else if !strings.Contains(err.Error(), d.errStr) {
				t.Fatalf("test load config failed:%d, errmsg is wrong", i)
			}
		}
	}
}
