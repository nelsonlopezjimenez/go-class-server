package controller

import (
	"localhost/CIS/modules/util"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func CheckForSiteExt(ctx *gin.Context) {
	// ctc.Param returns the wildcard value in the url path
	param := ctx.Request.URL.Path
	if !strings.HasSuffix(param, ".html") && !strings.HasSuffix(param, "/") {
		if !util.CheckExt(param) {
			param = param + ".html"
			ctx.Redirect(301, param)
		}
	}

	if strings.HasSuffix(param, ".asp") {
		param = strings.Replace(param, ".asp", ".html", 1)
		ctx.Redirect(301, param)
	}
	ctx.Next()

}
