package minequery

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func (p Pinger) openTCPConn(host string, port int) (net.Conn, error) {
	conn, err := p.Dialer.Dial("tcp", toAddrString(host, port))
	if err != nil {
		return nil, err
	}
	if p.Timeout != 0 {
		if err = conn.SetDeadline(time.Now().Add(p.Timeout)); err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func shouldWrapIPv6(host string) bool {
	return len(host) >= 2 && !(host[0] == '[' && host[1] == ']') && strings.Count(host, ":") >= 2
}

func toAddrString(host string, port int) string {
	if shouldWrapIPv6(host) {
		return fmt.Sprintf(`[%s]:%d`, host, port)
	} else {
		return fmt.Sprintf(`%s:%d`, host, port)
	}
}
