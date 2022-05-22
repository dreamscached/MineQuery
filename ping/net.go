package ping

import (
	"fmt"
	"net"
	"time"
)

func newTCPConn(host string, port uint16, deadline time.Time) (net.Conn, error) {
	dialer := &net.Dialer{Deadline: deadline}

	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	if err = conn.SetDeadline(deadline); err != nil {
		return nil, err
	}

	return conn, nil
}
