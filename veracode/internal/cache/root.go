package cache

import (
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	cachePathRelative = ".veracode" + string(filepath.Separator) + "cache"
)

type Cache struct {
	path string
}

var theCache Cache

func init() {
	// fmt.Printf("Entering Cache.init()\n")
	theCache.path = CachePath()
	err := os.MkdirAll(theCache.path, 0755)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Created cache directory: %s\n", theCache)
}

func Path() string {
	return theCache.path
}

func Clear() {

	os.RemoveAll(theCache.path)

}

func CachePath() string {

	// Create a cache directory
	// user, err := user.Current()
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return homeDir + string(filepath.Separator) + cachePathRelative
}
