package client

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sabafly/discord-bot-template/bot/db"
	"gopkg.in/yaml.v2"
)

func LoadConfig(config_filepath string) (*Config, error) {
	file, err := os.Open(config_filepath)
	if os.IsNotExist(err) {
		if file, err = os.Create(config_filepath); err != nil {
			return nil, err
		}
		switch filepath.Ext(config_filepath) {
		case ".json":
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "\t")
			err = encoder.Encode(defaultConfig)
		case ".yml", ".yaml":
			err = yaml.NewEncoder(file).Encode(file)
		case ".toml":
			err = toml.NewEncoder(file).Encode(defaultConfig)
		case ".xml":
			encoder := xml.NewEncoder(file)
			encoder.Indent("", "\t")
			err = encoder.Encode(defaultConfig)
		case ".gob":
			err = gob.NewEncoder(file).Encode(defaultConfig)
		default:
			panic("unknown config file type " + filepath.Ext(config_filepath))
		}
		if err != nil {
			return nil, err
		}
		return nil, errors.New("config file not found, created new one")
	} else if err != nil {
		return nil, err
	}

	var cfg Config
	switch filepath.Ext(config_filepath) {
	case ".json":
		err = json.NewDecoder(file).Decode(&cfg)
	case ".yml", ".yaml":
		err = yaml.NewDecoder(file).Decode(&cfg)
	case ".tml", ".toml":
		_, err = toml.NewDecoder(file).Decode(&cfg)
	case ".xml":
		err = xml.NewDecoder(file).Decode(&cfg)
	case ".gob":
		err = gob.NewDecoder(file).Decode(&cfg)
	default:
		panic("unknown config file type" + filepath.Ext(config_filepath))
	}
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

var defaultConfig = Config{
	DBConfig: db.DBConfig{
		Host: "localhost",
		Port: "6432",
		DB:   0,
	},
}

type Config struct {
	DBConfig db.DBConfig `json:"db_config"`
}
