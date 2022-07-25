package minequery

import (
	"testing"
)

func TestPingBeta18(t *testing.T) {
	testPingBeta18WithDefaultPinger(t)
	testPingBeta18WithNewPinger(t)
}

func testPingBeta18WithDefaultConfig(t *testing.T, status StatusBeta18) {
	if status.MOTD != "A Minecraft Server" {
		t.Errorf("expected MOTD = %#v, got %#v", "A Minecraft Server", status.MOTD)
	}
	if status.OnlinePlayers != 0 {
		t.Errorf("expected OnlinePlayers = %#v, got %#v", 0, status.MOTD)
	}
	if status.MaxPlayers != 20 {
		t.Errorf("expected MaxPlayers = %#v, got %#v", 20, status.MOTD)
	}
}

func testPingBeta18WithDefaultPinger(t *testing.T) {
	res, err := PingBeta18(Hostname(), Port())
	if err != nil {
		t.Errorf("default pinger test failed: %s", err)
		return
	}
	testPingBeta18WithDefaultConfig(t, res)
}

func testPingBeta18WithNewPinger(t *testing.T) {
	res, err := NewPinger().PingBeta18(Hostname(), Port())
	if err != nil {
		t.Errorf("new pinger test failed: %s", err)
		return
	}
	testPingBeta18WithDefaultConfig(t, res)
}
