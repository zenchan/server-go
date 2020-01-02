package xrpc

// Option rpc server option
type Option func(s *rpcServer)

// SetRequestCallback set request/reply callback function
func SetRequestCallback(cb RequestCallback) Option {
	return func(s *rpcServer) {
		s.requestCb = cb
	}
}

// SetRouteCallback set route callback function
func SetRouteCallback(cb RouteCallback) Option {
	return func(s *rpcServer) {
		s.routeCb = cb
	}
}

// SetErrorCallback set error message callback function
func SetErrorCallback(cb func(string)) Option {
	return func(s *rpcServer) {
		s.errCb = cb
	}
}
