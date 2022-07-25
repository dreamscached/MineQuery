package minequery

import (
	"testing"
)

func TestPing17(t *testing.T) {
	testPing17WithDefaultPinger(t)
	testPing17WithNewPinger(t)
}

func testPing17WithDefaultConfig(t *testing.T, res Status17) {
	motd := res.DescriptionText()
	if motd != "A Minecraft Server" {
		t.Errorf("expected DescriptionText() = %#v, got %#v", "A Minecraft Server", motd)
	}
	if res.OnlinePlayers != 0 {
		t.Errorf("expected OnlinePlayers = %#v, got %#v", 0, res.OnlinePlayers)
	}
	if res.MaxPlayers != 20 {
		t.Errorf("expected MaxPlayers = %#v, got %#v", 20, res.MaxPlayers)
	}
	if res.ProtocolVersion != int(Ping17ProtocolVersion172) {
		t.Errorf("expected ProtocolVersion = %#v, got %#v", Ping17ProtocolVersion172, res.ProtocolVersion)
	}
	if res.VersionName != "1.7.2" {
		t.Errorf("expected VersionName = %#v, got %#v", "1.7.2", res.VersionName)
	}
}

func testPing17WithDefaultPinger(t *testing.T) {
	res, err := Ping17(Hostname(), Port())
	if err != nil {
		t.Errorf("default pinger test failed: %s", err)
		return
	}
	testPing17WithDefaultConfig(t, res)
}

func testPing17WithNewPinger(t *testing.T) {
	res, err := NewPinger().Ping17(Hostname(), Port())
	if err != nil {
		t.Errorf("default pinger test failed: %s", err)
		return
	}
	testPing17WithDefaultConfig(t, res)
}
