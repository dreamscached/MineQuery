package minequery

import (
	"fmt"
	"testing"
)

func TestPing16(t *testing.T) {
	if err := ping16WithDefaultPinger(); err != nil {
		t.Errorf("default pinger test failed: %s", err)
	}
	if err := ping16WithNewPinger(); err != nil {
		t.Errorf("new pinger test failed: %s", err)
	}
}

func ping16WithDefaultPinger() error {
	res, err := Ping16(Hostname(), Port())
	if err != nil {
		return err
	}

	if res.MOTD != "A Minecraft Server" {
		return fmt.Errorf("expected MOTD = %#v, got %#v", "A Minecraft Server", res.MOTD)
	}
	if res.OnlinePlayers != 0 {
		return fmt.Errorf("expected OnlinePlayers = %#v, got %#v", 0, res.MOTD)
	}
	if res.MaxPlayers != 20 {
		return fmt.Errorf("expected MaxPlayers = %#v, got %#v", 20, res.MOTD)
	}

	return nil
}

func ping16WithNewPinger() error {
	p := NewPinger()

	res, err := p.Ping16(Hostname(), Port())
	if err != nil {
		return err
	}

	if res.ProtocolVersion != int(Ping16ProtocolVersion161) {
		return fmt.Errorf("expected ProtocolVersion = %#v, got %#v", Ping16ProtocolVersion161, res.ProtocolVersion)
	}
	if res.ServerVersion != "1.6.1" {
		return fmt.Errorf("expected ProtocolVersion = %#v, got %#v", Ping16ProtocolVersion161, res.ProtocolVersion)
	}
	if res.MOTD != "A Minecraft Server" {
		return fmt.Errorf("expected MOTD = %#v, got %#v", "A Minecraft Server", res.MOTD)
	}
	if res.OnlinePlayers != 0 {
		return fmt.Errorf("expected OnlinePlayers = %#v, got %#v", 0, res.MOTD)
	}
	if res.MaxPlayers != 20 {
		return fmt.Errorf("expected MaxPlayers = %#v, got %#v", 20, res.MOTD)
	}

	return nil
}
