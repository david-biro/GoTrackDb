package config

import (
	"bufio"
	"os"
	"strings"
)

var Cfg Config

type Config struct {
	DataXMLPathLinux string
	DataXMLPathOther string
	Pattern          string
	Port             string
	Dbhost           string
	Dbport           string
	Dbuser           string
	Dbpass           string
	Dbname           string
}

func ReadConfig(filename string) (Config, error) {
	config := Config{}

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := strings.Trim(parts[1], `"`)

		switch key {
		case "dataxmlpathlinux":
			config.DataXMLPathLinux = value
		case "dataxmlpathother":
			config.DataXMLPathOther = value
		case "pattern":
			config.Pattern = value
		case "port":
			config.Port = value
		case "dbhost":
			config.Dbhost = value
		case "dbport":
			config.Dbport = value
		case "dbuser":
			config.Dbuser = value
		case "dbpass":
			config.Dbpass = value
		case "dbname":
			config.Dbname = value
		}
	}

	if err := scanner.Err(); err != nil {
		return config, err
	}

	return config, nil
}
