package config

import (
	"encoding/json"
	"time"
)

type (
	Opts struct {
		// logger
		Logger struct {
			Debug       bool `long:"log.debug"    env:"LOG_DEBUG"  description:"debug mode"`
			Development bool `long:"log.devel"    env:"LOG_DEVEL"  description:"development mode"`
			Json        bool `long:"log.json"     env:"LOG_JSON"   description:"Switch log output to json format"`
		}

		MyUplink struct {
			Url string `long:"myuplink.url" env:"MYUPLINK_URL" description:"Url to myUplink API" default:"https://api.myuplink.com"`

			Auth struct {
				ClientID     string `long:"myuplink.auth.clientid"      env:"MYUPLINK_AUTH_CLIENTID"      description:"ClientID from myUplink" required:"yes"`
				ClientSecret string `long:"myuplink.auth.clientsecret"  env:"MYUPLINK_AUTH_CLIENTSECRET"  description:"ClientSecret from myUplink" json:"-"  required:"yes"`
			}
		}

		// general options
		Server struct {
			// general options
			Bind         string        `long:"server.bind"              env:"SERVER_BIND"           description:"Server address"        default:":8080"`
			ReadTimeout  time.Duration `long:"server.timeout.read"      env:"SERVER_TIMEOUT_READ"   description:"Server read timeout"   default:"5s"`
			WriteTimeout time.Duration `long:"server.timeout.write"     env:"SERVER_TIMEOUT_WRITE"  description:"Server write timeout"  default:"10s"`
		}
	}
)

func (o *Opts) GetJson() []byte {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return jsonBytes
}
