package minequery

import (
	"fmt"
	"testing"
)

func TestPing14(t *testing.T) {
	if err := ping14WithDefaultPinger(); err != nil {
		t.Errorf("default pinger test failed: %s", err)
	}
	if err := ping14WithNewPinger(); err != nil {
		t.Errorf("new pinger test failed: %s", err)
	}
}

func ping14WithDefaultPinger() error {
	res, err := Ping14(Hostname(), Port())
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

func ping14WithNewPinger() error {
	p := NewPinger()

	res, err := p.Ping14(Hostname(), Port())
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
