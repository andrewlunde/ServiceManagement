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
	"path/filepath"
	"runtime"

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

	var isWorkspace = false // Scan from the root for *.code-workspace files???
	// Scan for *.theia-workspace files in BAS ??
	var settingsFile = "Unknown"
	var settingsExists = false

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("On Mac:")
		if isWorkspace {
			settingsFile = "~/Code/User/"
		} else { //User(Global) Settings
			// settingsFile = "$HOME/Library/Application Support/Code/User/settings.json"
			settingsFile = homeDirectory + "/Library/Application Support/Code/User/settings.json"
		}

	case "linux":
		fmt.Println("On Linux:")
		if isWorkspace {
			settingsFile = "~/Code/User/"
		} else { //User(Global) Settings
			// Check to see if BAS
			settingsFile = homeDirectory + "/.theia/settings.json"
			if _, err := os.Stat(settingsFile); err == nil {
				// path/to/whatever exists
				fmt.Println("We are in BAS since " + settingsFile + " Exists!")

			} else {
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
		"name": "CAPMT_SMC:b44f32d4-6e31-4d95-b17f-6c6fcdb37e1f", 
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

	if settingsExists {
		// read file
		byteValue, err := ioutil.ReadFile(settingsFile)
		if err != nil {
			fmt.Print(err)
		} else {
			//err := jsonparser.GetString(data, "items", "[0]", "id")
			colorTheme, err := jsonparser.GetString(byteValue, "workbench.colorTheme")
			handleError(err)
			fmt.Println("colorTheme: " + colorTheme)

			var dataValue []byte
			var dataType jsonparser.ValueType
			var offset int = 0

			dataValue, dataType, offset, err = jsonparser.Get(byteValue, "sqltools.connections")
			if err != nil {
				fmt.Println("sqltools.connections" + " Key path not found")
				// We can go ahead and add it.
			}

			fmt.Println("dataValue: " + string(dataValue))
			fmt.Println("offset: ", offset)

			if dataType == jsonparser.NotExist {
				fmt.Println("sqltools.connections" + " is NotExist")

			} else if dataType == jsonparser.Array {
				fmt.Println("sqltools.connections" + " is Array")

				jsonparser.ArrayEach(dataValue, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					name, _ := jsonparser.GetString(value, "name")
					fmt.Println("name: " + name)
				})
				// https://github.com/buger/jsonparser#set

				fmt.Println("newConn: " + newConn)

			} else if dataType == jsonparser.Object {
				fmt.Println("sqltools.connections" + " is Object")

			} else if dataType == jsonparser.Null {
				fmt.Println("sqltools.connections" + " is Null")

			} else {
				fmt.Println("sqltools.connections" + " is unexpected")

			}

			/*
				var result map[string]interface{}
				json.Unmarshal([]byte(byteValue), &result)

				sqltoolsConnections := result["sqltools.connections"].(map[string]interface{})

				for key, value := range sqltoolsConnections {
					fmt.Println(key, value.(string))
				}
			*/
			//connections, dtype, offset, err := jsonparser.Get(byteValue, "sqltools.connections")
			//handleError(err)

			// fmt.Println("result: ", json.MarshalIndent(result, "", "    "))
			//fmt.Println("dtype: " + jsonparser.ValueType.String(dtype))
			//fmt.Println("offset: %d", offset)

		}
	}

	fmt.Println(filepath.Join("a", "b", "c"))
	fmt.Println(filepath.Join("a", "b/c"))
	fmt.Println(filepath.Join("a/b", "c"))
	fmt.Println(filepath.Join("a/b", "/c"))
}
