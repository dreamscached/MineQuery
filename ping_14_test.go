package minequery

import (
	"testing"
)

func TestPing14(t *testing.T) {
	testPing14WithDefaultPinger(t)
	testPing14WithNewPinger(t)
}

func testPing14WithDefaultConfig(t *testing.T, status *Status14) {
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

func testPing14WithDefaultPinger(t *testing.T) {
	res, err := Ping14(Hostname(), Port())
	if err != nil {
		t.Errorf("default pinger test failed: %s", err)
		return
	}
	testPing14WithDefaultConfig(t, res)
}

func testPing14WithNewPinger(t *testing.T) {
	res, err := NewPinger().Ping14(Hostname(), Port())
	if err != nil {
		t.Errorf("new pinger test failed: %s", err)
		return
	}
	testPing14WithDefaultConfig(t, res)
}
