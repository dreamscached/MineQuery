package minequery

import (
	"testing"
)

func TestPing16(t *testing.T) {
	ping16WithDefaultPinger(t)
	ping16WithNewPinger(t)
}

func testPing16WithDefaultConfig(t *testing.T, status *Status16) {
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

func ping16WithDefaultPinger(t *testing.T) {
	res, err := Ping16(Hostname(), Port())
	if err != nil {
		t.Errorf("default pinger test failed: %s", err)
		return
	}
	testPing16WithDefaultConfig(t, res)
}

func ping16WithNewPinger(t *testing.T) {
	res, err := NewPinger().Ping16(Hostname(), Port())
	if err != nil {
		t.Errorf("new pinger test failed: %s", err)
		return
	}
	testPing16WithDefaultConfig(t, res)
}
