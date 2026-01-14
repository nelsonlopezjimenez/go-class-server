package util

import (
	"strings"
)

func CheckExt(url string) bool {
	var acceptedFileExt = []string{
		".css",
		".js",
		".woff2",
		".svg",
		".jpg",
		".jpeg",
		".gif",
		".manifest",
		".webmanifest",
		".webp",
		".mp4",
		".mp3",
		".png",
		".ico",
		".asp",
		".php",
		".json",
		".wasm",
		".ttf",
		".js.map",
		".otf",
		".txt",
		".md",
		".pdf",
		".htm",
		".wmv",
		".wav",
		".ogg",
	}

	for _, ext := range acceptedFileExt {
		if strings.HasSuffix(url, ext) {
			return strings.HasSuffix(url, ext)
		}

	}
	return false
}
