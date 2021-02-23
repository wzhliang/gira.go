package config

import (
	"os"

	"github.com/naoina/toml"
	"github.com/rs/zerolog/log"
)

type Status struct {
	In_progress    uint
	Done           uint
	Ready_for_test uint
	Reopen         uint
}

// Config ...
type Config struct {
	Jira struct {
		User   string
		Passwd string
		Url    string
	}
	YUN   Status
	CLOUD Status
	ZHY   Status
	Gitee struct {
		User  string
		Token string
	}
}

// New ...
func New(tomlFile string) *Config {
	log.Info().Msgf("Loading configuration from :%s ...", tomlFile)
	f, err := os.Open(tomlFile)
	if err != nil {
		return nil
	}
	defer f.Close()
	var config Config
	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		return nil
	}

	return &config
}
