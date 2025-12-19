package config

import (
	"encoding/json"
	"time"
)

type (
	Opts struct {
		// logger
		Logger struct {
			Level  string `long:"log.level"    env:"LOG_LEVEL"   description:"Log level" choice:"trace" choice:"debug" choice:"info" choice:"warning" choice:"error" default:"info"`                          // nolint:staticcheck // multiple choices are ok
			Format string `long:"log.format"   env:"LOG_FORMAT"  description:"Log format" choice:"logfmt" choice:"json" default:"logfmt"`                                                                     // nolint:staticcheck // multiple choices are ok
			Source string `long:"log.source"   env:"LOG_SOURCE"  description:"Show source for every log message (useful for debugging and bug reports)" choice:"" choice:"short" choice:"file" choice:"full"` // nolint:staticcheck // multiple choices are ok
			Color  string `long:"log.color"    env:"LOG_COLOR"   description:"Enable color for logs" choice:"" choice:"auto" choice:"yes" choice:"no"`                                                        // nolint:staticcheck // multiple choices are ok
			Time   bool   `long:"log.time"     env:"LOG_TIME"    description:"Show log time"`
		}

		MyUplink struct {
			Url string `long:"myuplink.url" env:"MYUPLINK_URL" description:"Url to myUplink API" default:"https://api.myuplink.com"`

			Auth struct {
				ClientID     string `long:"myuplink.auth.clientid"      env:"MYUPLINK_AUTH_CLIENTID"      description:"ClientID from myUplink" required:"yes"`
				ClientSecret string `long:"myuplink.auth.clientsecret"  env:"MYUPLINK_AUTH_CLIENTSECRET"  description:"ClientSecret from myUplink" json:"-"  required:"yes"`
			}

			Device struct {
				AllowedConnectionStates []string `long:"myuplink.device.allowed-connectionstates"  env:"MYUPLINK_DEVICE_ALLOWED_CONNECTIONSTATES"  env-delim:"," description:"Allowed device connection states" default:"Connected"`
				CalcTotalParameters     []string `long:"myuplink.device.calc-total-parameters"  env:"MYUPLINK_DEVICE_CALC_TOTAL_PARAMETRS"  env-delim:"," description:"Calculate total metrics for these parameters (eg. energey log parameters)"`
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
