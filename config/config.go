package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Data struct {
	// The port to listen on. e.g '8080'
	Port string

	// Boards is a map of board to board URIs. The board should be in the format '/<characters>/'
	// The board names are the enabled boards.
	Boards map[string]Board `json:"boards"`
	// DefaultTime is the default Unix seconds - seconds since the Epoch
	// This is used if no timestamp is provided.
	DefaultTime string `json:"defaultTime"`

	// Displayed version
	Version string `json:"version"`

	// DSN is the data source name for the SQL DB
	DSN string `json:"DSN"`
}

type Board struct {
	// URI is the URI for the API without a trailing slash.
	// It can reference this server or another. e.g https://hostname.tld/c or http://localhost:8080/c
	URI string `json:"URI"`
	// ImageURI is the URI for the image server. Only include a trailing slash for a local path.
	// It can reference a relative path or another file server. e.g https://external.tld/images or /local/images/
	ImageURI string `json:"imageURI"`
}

func LoadConfig(path string) Data {
	v := viper.New()
	v.SetEnvPrefix("IGIARI")

	v.SetDefault("port", "8080")
	v.SetDefault("boards", map[string]Board{})
	v.SetDefault("dsn", "root:mariadbrootpassword@tcp(127.0.0.1:3306)/asagi?charset=utf8&parseTime=True&loc=Local")
	v.SetDefault("defaultTime", "1343080585")
	v.SetDefault("version", "-1")

	v.AutomaticEnv()

	v.SetConfigName("config")        // name of config file (without extension)
	v.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath("/etc/glaive/")  // path to look for the config file in
	v.AddConfigPath("$HOME/.glaive") // call multiple times to add many search paths
	v.AddConfigPath(".")             // optionally look for config in the working directory
	v.AddConfigPath(path)
	err := v.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, continuing.")
		} else {
			log.Fatalln("Config could not be loaded: " + err.Error())
		}
	}

	data := &Data{}
	err = v.Unmarshal(data)
	if err != nil {
		log.Fatalln("Failed to read in config: " + err.Error())
	}

	// TODO: This is a backwards compatibility hack. This can be removed after this ENV Var is no longer used
	if dsn := os.Getenv("IGIARI_SQL_DSN"); dsn != "" {
		data.DSN = dsn
	}

	// Cloud Run requirement
	if port := os.Getenv("PORT"); port != "" {
		data.Port = port
	}

	return *data
}
