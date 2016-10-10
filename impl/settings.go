package impl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	CONF_FILE    = "/etc/puppetlabs/puppet/monitor.conf"
	def_port     = 3840
	def_ip       = "127.0.0.1"
	def_pid_file = "/var/run/puppet_monitoring.pid"
	def_rpc      = 3841
	def_ctime    = 35
)

type Settings struct {
	Port        int    `json:"port,omitempty"`
	Ip          string `json:"ip,omitempty"`
	PidFile     string `json:"pid,omitempty"`
	RpcPort     int    `json:"rpc,omitempty"`
	RpcComputed string `json:"-"`
	ControlTime int    `json:"ctime,omitempty"`
}

func (s Settings) LoadSettings() Settings {
	var _, err = os.Stat(CONF_FILE)
	if err == nil { // if file exists, loading from it
		var set_data, _ = ioutil.ReadFile(CONF_FILE)
		var settings = Settings{}
		err = json.Unmarshal(set_data, &settings)
		if err != nil {
			panic(err)
		}
		if settings.Ip == "" {
			settings.Ip = def_ip
		}
		if settings.Port == 0 {
			settings.Port = def_port
		}
		if settings.PidFile == "" {
			settings.PidFile = def_pid_file
		}
		if settings.RpcPort == 0 {
			settings.RpcPort = def_rpc
		}
		if settings.ControlTime == 0 {
			settings.ControlTime = def_ctime
		}

		settings.RpcComputed = settings.Ip + ":" + strconv.Itoa(settings.RpcPort)

		return settings
	}
	return Settings{Ip: def_ip, Port: def_port, PidFile: def_pid_file, RpcPort: def_rpc, ControlTime: def_ctime, RpcComputed: def_ip + ":" + strconv.Itoa(def_rpc)}
}
