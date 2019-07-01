package config

import (
	"os"
	"fmt"
	"flag"
	"encoding/json"
)

type Configuration struct {
	PORT           string
	DATASTORE      string
	MONGO_SERVER   string
	MONGO_DB       string
	SECRET_KEY     string
}

var (
	configuration *Configuration = nil
	configFile    *string        = nil
)

func init() {
	configFile = flag.String("file", "", "a String")
}

func TestSetupConfiguration(configFile string) *Configuration {
	if configuration == nil {
		configuration = loadConfiguration(configFile)
		return configuration
	}
	return nil
}

func LoadAppConfiguration() {
	if configuration == nil {
		flag.Parse()
		if len(*configFile) == 0 {
			StopService("Mandatory arguments not provided for executing the App")
		}
		configuration = loadConfiguration(*configFile)
	}
}

func loadConfiguration(filename string) (*Configuration) {
	if configuration == nil {
		configFile, err := os.Open(filename)
		defer configFile.Close()
		if err != nil {
			StopService(err.Error())
		}
		jsonParser := json.NewDecoder(configFile)
		err1 := jsonParser.Decode(&configuration)
		if err1 != nil {
			fmt.Println("Failed to parse configuration file")
			StopService(err1.Error())
		}
	}
	return configuration
}

func GetAppConfiguration() *Configuration {
	if configuration == nil {
		fmt.Println("Unable to get the app configuration. Loading freshly. \t")
		LoadAppConfiguration()
	}
	return configuration
}

func StopService(log string) {
	fmt.Println(log)
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(os.Kill)
}

