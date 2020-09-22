// Find out what Platfrom / Arch I am?
// User/Workspace file
// Open
// Read existing
// Replace/Append
// Close
// Prompt to restart?

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"strconv"
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

	var isWorkspace = true // Scan from the root for *.code-workspace files???
	var isBAS = false
	// Scan for *.theia-workspace files in BAS ??
	var defaultsFile = "Unknown"
	//var defaultsExists = false

	var settingsFile = "Unknown"
	var settingsExists = false

	var skipping = false
	var forceReplace = true

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("On Mac:")
		if isWorkspace {
			//settingsFile = "~/Code/User/"
			// The current code-workspace file can be found by looking here.
			// cat $HOME/Library/Application\ Support/Code/storage.json | grep -A 3 lastActiveWindow
			defaultsFile = homeDirectory + "/Library/Application Support/Code/storage.json"
			byteValue, err := ioutil.ReadFile(defaultsFile)
			handleError(err)

			configURIPath, err := jsonparser.GetString(byteValue, "windowsState", "lastActiveWindow", "workspaceIdentifier", "configURIPath")
			handleError(err)

			fmt.Println("configURIPath: " + configURIPath)

			settingsFile = "/" + strings.TrimLeft(configURIPath, "file:/")
			//settingsFile = homeDirectory + "/git/vsws/mta.code-workspace"
		} else { //User(Global) Settings
			// settingsFile = "$HOME/Library/Application Support/Code/User/settings.json"
			settingsFile = homeDirectory + "/Library/Application Support/Code/User/settings.json"
		}

	case "linux":
		fmt.Println("On Linux:")

		// Check to see if BAS
		settingsFile = homeDirectory + "/.theia/settings.json"
		if _, err := os.Stat(settingsFile); err == nil {
			// path/to/whatever exists
			fmt.Println("We are in BAS since " + settingsFile + " Exists!")
			isWorkspace = false
			isBAS = true
		}

		if isWorkspace {
			settingsFile = "~/Code/User/"
		} else { //User(Global) Settings
			if !isBAS {
				settingsFile = homeDirectory + "/.config/Code/User/settings.json"
			}
		}

	case "windows":
		fmt.Println("On Windoz:")
		if isWorkspace {
			settingsFile = "~/Code/User/"
		} else { //User(Global) Settings
			settingsFile = homeDirectory + "%APPDATA%\\Code\\User\\settings.json"
		}

	}

	fmt.Println("settingsFile: " + settingsFile)

	if _, err := os.Stat(settingsFile); err == nil {
		// path/to/whatever exists
		fmt.Println("settingsFile: " + settingsFile + " Exists!")
		settingsExists = true

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		fmt.Println("settingsFile: " + settingsFile + " Does NOT Exist!")

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		fmt.Println("settingsFile: " + settingsFile + " Existence Unknown!")

	}

	var newConn = `{
						"name": "CAPMT_SMC:subAcct",
						"group": "SMSI", 
						"driver": "SAPHana", 
						"dialect": "SAPHana",
						"server": "833726c5-cca3-4dce-a325-4385426009e7.hana.trial-us10.hanacloud.ondemand.com", 
						"port": 443, 
						"database": "D53EE042B6AD4E8093FF0A24F931586B", 
						"username": "D53EE042B6AD4E8093FF0A24F931586B_B5IBO9PWMQ841D52POXNE26XN_RT", 
						"password": "Mw9h7H.5r6CBidD2vtq.vxmzisxLAMx2_UJ9YrjZim2Yop-kUOcBII-g6VHYZMDpPzjT0PCQua.8i-V2f8MrjDqkGG6hRZAct2a2YIL7PFrlzeSDhO5qBOl6ni-VRF3t", 
						"connectionTimeout": 30, 
						"hanaOptions": {
							"encrypt": true, 
							"sslValidateCertificate": true, 
							"sslCryptoProvider": "openssl", 
							"sslTrustStore": "-----BEGIN CERTIFICATE-----MIIDrzCCApegAwIBAgIQCDvgVpBCRrGhdWrJWZHHSjANBgkqhkiG9w0BAQUFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0wNjExMTAwMDAwMDBaFw0zMTExMTAwMDAwMDBaMGExCzAJBgNVBAYTAlVTMRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5jb20xIDAeBgNVBAMTF0RpZ2lDZXJ0IEdsb2JhbCBSb290IENBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4jvhEXLeqKTTo1eqUKKPC3eQyaKl7hLOllsBCSDMAZOnTjC3U/dDxGkAV53ijSLdhwZAAIEJzs4bg7/fzTtxRuLWZscFs3YnFo97nh6Vfe63SKMI2tavegw5BmV/Sl0fvBf4q77uKNd0f3p4mVmFaG5cIzJLv07A6Fpt43C/dxC//AH2hdmoRBBYMql1GNXRor5H4idq9Joz+EkIYIvUX7Q6hL+hqkpMfT7PT19sdl6gSzeRntwi5m3OFBqOasv+zbMUZBfHWymeMr/y7vrTC0LUq7dBMtoM1O/4gdW7jVg/tRvoSSiicNoxBN33shbyTApOB6jtSj1etX+jkMOvJwIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUA95QNVbRTLtm8KPiGxvDl7I90VUwHwYDVR0jBBgwFoAUA95QNVbRTLtm8KPiGxvDl7I90VUwDQYJKoZIhvcNAQEFBQADggEBAMucN6pIExIK+t1EnE9SsPTfrgT1eXkIoyQY/EsrhMAtudXH/vTBH1jLuG2cenTnmCmrEbXjcKChzUyImZOMkXDiqw8cvpOp/2PV5Adg06O/nVsJ8dWO41P0jmP6P6fbtGbfYmbW0W5BjfIttep3Sp+dWOIrWcBAI+0tKIJFPnlUkiaY4IBIqDfv8NZ5YBberOgOzW6sRBc4L0na4UU+Krk2U886UAb3LujEV0lsYSEY1QSteDwsOoBrp+uvFRTp2InBuThs4pFsiv9kuXclVzDAGySj4dzp30d8tbQkCAUw7C29C79Fv1C5qfPrmAESrciIxpg0X40KPMbp1ZWVbd4=-----END CERTIFICATE-----"
							}
						}`

	newConnName, _ := jsonparser.GetString([]byte(newConn), "name")

	var foundIdx int = -1

	if settingsExists {
		// read file
		byteValue, err := ioutil.ReadFile(settingsFile)
		if err != nil {
			fmt.Print(err)
		} else {
			//err := jsonparser.GetString(data, "items", "[0]", "id")
			//colorTheme, err := jsonparser.GetString(byteValue, "workbench.colorTheme")
			//handleError(err)
			//fmt.Println("colorTheme: " + colorTheme)

			// var newValue []byte
			// var newType jsonparser.ValueType
			// var newOffset int = 0

			var dataValue []byte
			var dataType jsonparser.ValueType
			var dataOffset int = 0

			if isWorkspace {
				dataValue, dataType, dataOffset, err = jsonparser.Get(byteValue, "settings", "sqltools.connections")
			} else {
				dataValue, dataType, dataOffset, err = jsonparser.Get(byteValue, "sqltools.connections")
			}

			if err != nil {
				fmt.Println("sqltools.connections" + " Key path not found")
				// We can go ahead and add it.
			}

			// fmt.Println("dataValue: " + string(dataValue))
			fmt.Println("offset: ", dataOffset)

			if dataType == jsonparser.NotExist {
				fmt.Println("sqltools.connections" + " is NotExist")
				// IF this is the case then we can safely create a new sqltools.connections array and append it to settings

				var newSQLToolsConn string
				newSQLToolsConn = string(byteValue)
				newSQLToolsConn2 := strings.TrimRight(newSQLToolsConn, "}")
				newSQLToolsConn = newSQLToolsConn2
				newSQLToolsConn += ","
				newSQLToolsConn += `"sqltools.connections": [ `
				newSQLToolsConn += newConn + "] }"

				// write file
				err = ioutil.WriteFile(settingsFile, []byte(newSQLToolsConn), 0644)
				handleError(err)

			} else if dataType == jsonparser.Array {
				fmt.Println("sqltools.connections" + " is an Array")

				var scidx int = 0
				jsonparser.ArrayEach(dataValue, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					name, _ := jsonparser.GetString(value, "name")
					// fmt.Println("name: " + name)
					if newConnName != name {
						fmt.Println("keeping: " + name)
					} else {
						if forceReplace {
							fmt.Println("replacing: " + name)
						} else {
							fmt.Println("skipping: " + name)
						}
						foundIdx = scidx
						skipping = true
					}
					scidx = scidx + 1
				})
				// https://github.com/buger/jsonparser#set

				if !skipping {
					fmt.Println("Adding connection with name " + newConnName + ".")

					var newSQLToolsConn string

					newSQLToolsConn = string(dataValue)
					newSQLToolsConn2 := strings.TrimRight(newSQLToolsConn, "]")
					newSQLToolsConn = newSQLToolsConn2
					if scidx > 0 {
						newSQLToolsConn += ","
					}
					newSQLToolsConn += newConn + "]"

					var setValue []byte

					// fmt.Println("attempt set: ")

					if isWorkspace {
						setValue, err = jsonparser.Set(byteValue, []byte(newSQLToolsConn), "settings", "sqltools.connections")
					} else {
						setValue, err = jsonparser.Set(byteValue, []byte(newSQLToolsConn), "sqltools.connections")
					}
					handleError(err)

					//fmt.Println("after set: ")
					// jsonparser.ArrayEach(setValue, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					// 	name, _ := jsonparser.GetString(value, "name")
					// 	fmt.Println("name: " + name)
					// })

					// fmt.Println("newConn: " + newConn)
					// fmt.Println("setValue: " + string(setValue))

					// write file
					err = ioutil.WriteFile(settingsFile, setValue, 0644)
					handleError(err)
				} else {
					if forceReplace {
						fmt.Println("Connection with name " + newConnName + " already exists!  Forcing replacement.")
						idxStr := "[" + strconv.Itoa(foundIdx) + "]"
						// idxStr := strconv.Itoa(foundIdx)
						// fmt.Println("idxStr:" + idxStr)
						var setValue []byte
						if isWorkspace {
							// dataValue, dataType, dataOffset, err = jsonparser.Get(byteValue, "settings", "sqltools.connections", idxStr)
							setValue, err = jsonparser.Set(byteValue, []byte(newConn), "settings", "sqltools.connections", idxStr)
						} else {
							// dataValue, dataType, dataOffset, err = jsonparser.Get(byteValue, "sqltools.connections", idxStr)
							setValue, err = jsonparser.Set(byteValue, []byte(newConn), "settings", "sqltools.connections", idxStr)
						}
						handleError(err)

						//fmt.Println("setValue: " + string(setValue))
						//fmt.Println("offset: ", dataOffset)

						// fmt.Println("setValue: " + string(setValue))

						// write file
						err = ioutil.WriteFile(settingsFile, setValue, 0644)
						handleError(err)

					} else {
						fmt.Println("Connection with name " + newConnName + " already exists!  Delete it first and rerun.")
					}
				}

			} else if dataType == jsonparser.Object {
				fmt.Println("sqltools.connections" + " is Object")

			} else if dataType == jsonparser.Null {
				fmt.Println("sqltools.connections" + " is Null")

			} else {
				fmt.Println("sqltools.connections" + " is unexpected")

			}

		}
	}

}
