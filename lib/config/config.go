package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	defaultConfigFileName = ".vshop.json"
)

type Config struct {
	InProduction bool   `json:"inProduction,omitempty"`
	AppRoot      string `json:"appRoot,omitempty"`
	TemplateRoot string `json:"templateRoot,omitempty"`
	S3Root       string `json:"s3Root,omitempty"`
	CheckoutURL  string `json:"checkoutURL,omitempty"`
	Db           struct {
		User                 string `json:"user,omitempty"`
		Passwd               string `json:"passwd,omitempty"`
		Net                  string `json:"net,omitempty"`
		Addr                 string `json:"addr,omitempty"`
		DBName               string `json:"dbName,omitempty"`
		ParseTime            bool   `json:"parseTime,omitempty"`
		Loc                  string `json:"loc,omitempty"`
		AllowNativePasswords bool   `json:"allowNativePasswords,omitempty"`
	} `json:"db,omitempty"`
	Security struct {
		CSRFKey string `json:"csrfKey,omitempty"`
	} `json:"security,omitempty"`
	Session struct {
		Name        string `json:"name,omitempty"`
		Path        string `json:"path,omitempty"`
		Domain      string `json:"domain,omitempty"`
		MaxAgeHours int    `json:"maxAgeHours,omitempty"`
	} `json:"session,omitempty"`
}

var c = Config{}

func Configuration(configFileName ...string) (*Config, error) {

	if (c == Config{}) {

		var cfname string
		switch len(configFileName) {
		case 0:
			dirname, err := os.UserHomeDir()
			if err != nil {
				return &c, err
			}
			cfname = fmt.Sprintf("%s/%s", dirname, defaultConfigFileName)
		case 1:
			cfname = configFileName[0]
		default:
			return &c, fmt.Errorf("incorrect arguments for configuration file name")
		}

		configFile, err := os.Open(cfname)
		if err != nil {
			return &c, err
		}
		defer configFile.Close()

		decoder := json.NewDecoder(configFile)
		err = decoder.Decode(&c)
		if err != nil {
			return &c, err
		}
	}

	return &c, nil
}
