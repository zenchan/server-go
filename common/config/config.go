package config

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/vaughan0/go-ini"
)

// Config server configuration
type Config struct {
	TCPPort  int `json:"tcp_port"`
	HTTPPort int `json:"http_port"`
	UDPPort  int `json:"udp_port"`
	RPCPort  int `json:"rpc_port"`

	LogLevel string `json:"log_level"`
	LogPath  string `json:"log_path"`
}

var (
	// global server config
	srvCfg *Config
)

// Load load configuration
func Load(file string) (err error) {
	f, err := ini.LoadFile(file)
	if err != nil {
		return
	}

	cfg := &Config{}
	pt := reflect.TypeOf(cfg)
	pv := reflect.ValueOf(cfg)
	vt := pt.Elem()
	vv := pv.Elem()

	for _, sec := range f {
		for k, v := range sec {
			if v == "" {
				continue
			}
			if err = fillStructs(vv, vt, k, v); err != nil {
				return
			}
		}
	}

	fillDefaultValue(cfg)
	srvCfg = cfg
	return
}

func fillStructs(vv reflect.Value, vt reflect.Type, k, v string) (err error) {
	found := false
	for i := 0; i < vv.NumField(); i++ {
		tag := vt.Field(i).Tag.Get("json")
		kind := vv.Field(i).Kind()

		if kind == reflect.Struct {
			fillStructs(vv.Field(i), vt.Field(i).Type, k, v)
			continue
		}

		if tag == k {
			found = true
			switch kind {
			case reflect.String:
				vv.Field(i).SetString(v)
			case reflect.Int:
				var n int
				n, err = strconv.Atoi(v)
				if err != nil {
					err = errors.New(`"` + k + `" must be numeric`)
					return
				}
				vv.Field(i).SetInt(int64(n))
			case reflect.Slice:
				arr := strings.Split(v, ",")
				hdr := (*reflect.SliceHeader)(unsafe.Pointer(vv.Field(i).UnsafeAddr()))
				valHdr := (*reflect.SliceHeader)(unsafe.Pointer(&arr))
				hdr.Len = valHdr.Len
				hdr.Data = valHdr.Data
				hdr.Cap = valHdr.Cap
			case reflect.Bool:
				var b bool
				b, err = strconv.ParseBool(v)
				if err != nil {
					err = errors.New(`"` + k + `" must be boolean`)
					return
				}
				vv.Field(i).SetBool(b)
			default:
				err = errors.New("unsupported config type: " + k)
				return
			}
			return
		}
	}

	if !found {
		err = errors.New(`unknown field "` + k + `" was found`)
	}
	return
}

func fillDefaultValue(cfg *Config) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.LogPath == "" {
		cfg.LogPath = "./logs"
	}
}
