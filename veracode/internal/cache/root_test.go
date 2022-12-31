package cache

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCachePath_ContainsOnlyFileSeparator(t *testing.T) {
	expected := string(filepath.Separator)
	actual := CachePath()
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected %v to contain %v", actual, expected)
	}
	actual = strings.ReplaceAll(actual, string(filepath.Separator), "")
	unexpected := "/\\"
	if strings.ContainsAny(actual, unexpected) {
		t.Errorf("%v had one or more of %v after removing the correct separator %v", actual, unexpected, string(filepath.Separator))
	}
}
