package config

// TCPPort returns tcp_port
func TCPPort() int {
	return srvCfg.TCPPort
}

// HTTPPort returns http_port
func HTTPPort() int {
	return srvCfg.HTTPPort
}

// UDPPort returns udp_port
func UDPPort() int {
	return srvCfg.UDPPort
}

// LogLevel returns log_level
func LogLevel() string {
	return srvCfg.LogLevel
}

// LogPath returns log_path
func LogPath() string {
	return srvCfg.LogPath
}
