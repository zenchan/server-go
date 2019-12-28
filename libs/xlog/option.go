package xlog

// Option logger option
type Option func(l *loggingT)

// WithStdout logger output with stdout
func WithStdout(b bool) Option {
	return func(l *loggingT) {
		if b {
			l.logStdout = 1
		} else {
			l.logStdout = 0
		}
	}
}
