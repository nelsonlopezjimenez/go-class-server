// ! Written by Matthew Wellington
package util

import "os"

//? This is the variables that could be modified with environment variables
var EnvDefaults = map[string]string{
	"RELEASE_VERSION":       "2.3.1",
	"RELEASE_DATE":          "2/11/26",
	"WEBSITES_GITEA_ADDR":   "http://192.168.1.47:3000/api/v1/repos/search?uid=6&limit=200",
	"UPDATE_IP":             "http://192.168.1.28:3000",
	"GIT_INSTALL_ADDR":      "http://192.168.1.47:3000/OfflineWebsites/",
	"WEBSITES_REPO_ADDR":    "http://192.168.1.47:3000/api/v1/repos/ClassroomResources/",
	"CIS_CLASS_SERVER_PORT": "22022",
	"VIDEO_VIEWER_PATH":     GetOSPaths()["Public"] + "/video-viewer",
}

//? LoadEnv checks to see if a value exists as an environment variable
//? if the variable exists then it is returned.
//? Otherwise the value from EnvDefaults is returned

func LoadEnv(key string) string {
	//? check to see if variable exists
	value, exist := os.LookupEnv(key)

	if exist {
		//? if it does, return it from .env
		return value
	} else {
		//? Otherwise return it from the default values (these should be production values)
		return EnvDefaults[key]
	}

}

func IsDevelopment() bool {
	VERSION_MODE, _ := os.LookupEnv("RELEASE_VERSION")

	return VERSION_MODE == "DEV"
}
