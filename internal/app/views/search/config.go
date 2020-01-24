package search

import (
	"github.com/sirupsen/logrus"
)

// Config is the config of search view.
type Config struct {
	log *logrus.Entry
}

// NewConfig create a new config.
func NewConfig() *Config {
	return &Config{
		log: logrus.WithFields(logrus.Fields{
			"pkg": "search",
		}),
	}
}
