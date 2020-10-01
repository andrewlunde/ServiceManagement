package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/buger/jsonparser"
)

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOARCH)

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	homeDirectory := user.HomeDir
	fmt.Printf("Home Directory: %s\n", homeDirectory)

	appData := os.Getenv("APPDATA")
	fmt.Printf("appData: %s\n", appData)

	var defaultsFile = "Unknown"
	//var defaultsExists = false

	var settingsFile = "Unknown"

	defaultsFile = appData + "/Code/storage.json"

	fmt.Println("defaultsFile: " + defaultsFile)

	_, err = os.Stat(defaultsFile)
	handleError(err)

	byteValue, err := ioutil.ReadFile(defaultsFile)
	if err == nil {
		configURIPath, err := jsonparser.GetString(byteValue, "windowsState", "lastActiveWindow", "workspaceIdentifier", "configURIPath")
		if err == nil {
			fmt.Println("configURIPath: " + configURIPath)
			settingsFile = strings.TrimLeft(configURIPath, "file:/")
		}
	}

	settingsFile = strings.Replace(settingsFile, "%3A", ":", -1)

	fmt.Println("settingsFile: " + settingsFile)

	_, err = os.Stat(settingsFile)
	handleError(err)

	fmt.Println("OK: Done.")

}
