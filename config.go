package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Configuration stores dbs configuration parameters
type Configuration struct {
	Verbose      int    `json:"verbose"`       // verbosity level
	Interval     int    `json:"interval"`      // server interval
	HTTPTimeout  int    `json:"http_timeout"`  // http timeout interval
	AlertManager string `json:"alert_manager"` // alert manager URL
	LogFile      string `json:"log_file"`      // server log file
	Rules        []Rule `json:"rules"`         // pod rules
}

// global variables
var Config Configuration

// String returns string representation of dbs Config
func (c *Configuration) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("ERROR: fail to marshal configuration", err)
		return ""
	}
	return string(data)
}

func ParseConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("unable to read config file", configFile, err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Println("unable to parse config file", configFile, err)
		return err
	}
	if Config.Interval == 0 {
		Config.Interval = 60 // default value
	}
	return nil
}
