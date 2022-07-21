package minequery

import (
	"fmt"
	"testing"
)

func TestPing17(t *testing.T) {

}

func ping17WithDefaultPinger() error {
	res, err := Ping17(Hostname(), Port())
	if err != nil {
		return err
	}

	motd := res.DescriptionText()
	if motd != "A Minecraft Server" {
		return fmt.Errorf("expected DescriptionText() = %#v, got %#v", "A Minecraft Server", motd)
	}
	if res.OnlinePlayers != 0 {
		return fmt.Errorf("expected OnlinePlayers = %#v, got %#v", 0, res.OnlinePlayers)
	}
	if res.MaxPlayers != 20 {
		return fmt.Errorf("expected MaxPlayers = %#v, got %#v", 20, res.MaxPlayers)
	}
	if res.ProtocolVersion != int(Ping17ProtocolVersion172) {
		return fmt.Errorf("expected ProtocolVersion = %#v, got %#v", Ping17ProtocolVersion172, res.ProtocolVersion)
	}
	if res.VersionName != "1.7.2" {
		return fmt.Errorf("expected VersionName = %#v, got %#v", "1.7.2", res.VersionName)
	}

	return nil
}

func ping17WithNewPinger() error {
	p := NewPinger()

	res, err := p.Ping17(Hostname(), Port())
	if err != nil {
		return err
	}

	motd := res.DescriptionText()
	if motd != "A Minecraft Server" {
		return fmt.Errorf("expected DescriptionText() = %#v, got %#v", "A Minecraft Server", motd)
	}
	if res.OnlinePlayers != 0 {
		return fmt.Errorf("expected OnlinePlayers = %#v, got %#v", 0, res.OnlinePlayers)
	}
	if res.MaxPlayers != 20 {
		return fmt.Errorf("expected MaxPlayers = %#v, got %#v", 20, res.MaxPlayers)
	}
	if res.ProtocolVersion != int(Ping17ProtocolVersion172) {
		return fmt.Errorf("expected ProtocolVersion = %#v, got %#v", Ping17ProtocolVersion172, res.ProtocolVersion)
	}
	if res.VersionName != "1.7.2" {
		return fmt.Errorf("expected VersionName = %#v, got %#v", "1.7.2", res.VersionName)
	}

	return nil
}
