package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
	//log "code.google.com/p/log4go"
)

var (
	Conf     = &Config{}
	confFile string
)

func init() {
	Conf.ServerName = "deploy-tools"
	flag.StringVar(&confFile, "c", "/etc/gm/conf/gm.toml", "config file path")
}

type Config struct {
	Debug         bool       `toml:"debug"`
	ServerName    string     `toml:"-"`
	Env           string     `toml:"env"`
	BaseDirectory string     `toml:"base_directory"`
	Projects      []*Project `toml:"project"`
}

type Project struct {
	Id            string `toml:"id"`
	Git           string `toml:"git"`
	Name          string `toml:"name"`
	Enabled       bool   `toml:"enabled"`
	InstallScript string `toml:"install_script"`
	UpdateScript  string `toml:"update_script"`
	RemoveScript  string `toml:"remove_script"`
	Description   string `toml:"description"`
}

func (p *Project) exists() (bool, error) {
	path := Conf.BaseDirectory + "/" + p.Name
	return pathExists(path)
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func InitConfig() (err error) {
	var (
		v []byte
	)

	if v, err = ioutil.ReadFile(confFile); err != nil {
		return
	}

	_, err = toml.Decode(string(v), Conf)

	Debug = Conf.Debug

	id := 1
	for _, p := range Conf.Projects {
		p.Id = fmt.Sprintf("%d", id)
		id++
	}
	return
}
