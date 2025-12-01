package util

import (
	"strings"
)

func CheckExt(url string) bool {
	var acceptedFileExt = []string{".css", ".js", ".woff2", ".svg", ".jpg", ".jpeg", ".gif", ".manifest", ".webmanifest", ".webp", ".mp4", ".png", ".ico", ".asp", ".php"}

	for _, ext := range acceptedFileExt {
		if strings.HasSuffix(url, ext) {
			return strings.HasSuffix(url, ext)
		}

	}
	return false
}
