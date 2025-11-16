// ! Written by Matthew Wellington
package util

import "os"

//? This is the variables that could be modified with environment variables
var EnvDefaults = map[string]string{
	"RELEASE_VERSION": "1.5.0",
	"RELEASE_DATE":    "10/30/25",
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
