package minequery

import (
	"net"
	"time"
)

// PingerOption is a configuring function that applies certain changes to Pinger.
type PingerOption func(*Pinger)

// WithTimeout sets Pinger Dialer timeout to the provided value.
func WithTimeout(timeout time.Duration) PingerOption {
	return func(p *Pinger) { p.Timeout = timeout }
}

// WithUseStrict sets Pinger UseStrict to the provided value.
func WithUseStrict(useStrict bool) PingerOption {
	return func(p *Pinger) { p.UseStrict = useStrict }
}

// WithProtocol16Version sets Pinger ProtocolVersion16 value.
func WithProtocol16Version(version byte) PingerOption {
	return func(p *Pinger) {
		p.ProtocolVersion16 = version
	}
}

// WithProtocol17Version sets Pinger ProtocolVersion17 value.
func WithProtocol17Version(version int32) PingerOption {
	return func(p *Pinger) {
		p.ProtocolVersion17 = version
	}
}

// defaultPinger is a default (zero-value) Pinger used in functions
// that don't have Pinger as receiver.
var defaultPinger Pinger

// Pinger contains options to ping Minecraft servers.
type Pinger struct {
	// Dialer used to establish and maintain connection with servers.
	Dialer net.Dialer

	// Timeout TODO
	Timeout time.Duration

	// UseStrict is a configuration value that defines if tolerable errors (in server ping responses)
	// that are by default silently ignored should be actually returned as errors.
	UseStrict bool

	// ProtocolVersion16 is protocol version to use when pinging with Ping16 function.
	// By default, Ping16ProtocolVersion162 (=74) will be used.
	// See ping_16.go for full list of built-in constants.
	ProtocolVersion16 byte

	// ProtocolVersion17 is protocol version to use when pinging with Ping17 function.
	// By default, Ping17ProtocolVersionUndefined (=-1) will be used.
	// See ping_17.go for full list of built-in constants.
	ProtocolVersion17 int32
}

// NewPinger constructs new Pinger instance optionally with additional options.
func NewPinger(options ...PingerOption) Pinger {
	var pinger Pinger
	for _, configure := range options {
		configure(&pinger)
	}
	return pinger
}
