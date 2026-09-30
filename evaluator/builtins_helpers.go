package evaluator

import (
	"os"
	"strings"
	"time"
)

func containsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
}

func splitString(s, sep string) []string {
	return strings.Split(s, sep)
}

func joinStrings(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

func toUpper(s string) string {
	return strings.ToUpper(s)
}

func toLower(s string) string {
	return strings.ToLower(s)
}

func getCurrentUnixTimestamp() int64 {
	return time.Now().Unix()
}

func sleepMilliseconds(ms int64) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func readFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFileContent(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
