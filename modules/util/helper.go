package util

// Rocky Connor 420711

import (
	"strings"

	"github.com/gin-gonic/gin"
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

// IsAuthorized
//
// This middleware function prevents the access of the students' server
// to anyone on the outside network. This prevents the possibility of data
// transfer from system to system as prohibited by DOC.
// Note: This can only be controlled via the source code. After compile,
// This cannot be circumvented from the executable itself.
func IsAuthorized(ctx *gin.Context) {
	switch ctx.RemoteIP() {
	case "::1", "127.0.0.1":
		ctx.Next()
	default:
		ctx.String(403, "I'm sorry. Your are not authorized to see this.")
		ctx.Abort()
	}
}
