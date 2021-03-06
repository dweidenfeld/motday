package config

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"
)

// Config Model
type Config struct {
	Flickr  FlickrConf   `json:"flickr"`
	Motives []MotiveConf `json:"motives"`
}

// FlickrConf Model
type FlickrConf struct {
	APIKey string `json:"apiKey"`
}

// MotiveConf Model
type MotiveConf struct {
	Title        string   `json:"title"`
	Descriptions []string `json:"descriptions"`
	Queries      []string `json:"queries"`
}

// Load loads the configuration from config.json
func Load(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)
	if nil != err {
		return nil, err
	}
	var config Config
	json.Unmarshal(b, &config)
	return &config, nil
}

// RandomMotive returns a random motive
func (c *Config) RandomMotive() *MotiveConf {
	return &c.Motives[random(0, len(c.Motives))]
}

// RandomQuery returns a random query
func (m *MotiveConf) RandomQuery() *string {
	return &m.Queries[random(0, len(m.Queries))]
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
