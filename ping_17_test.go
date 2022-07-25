package minequery

import (
	"testing"
)

type TestPing17Config struct {
	MOTD            string
	OnlinePlayers   int
	MaxPlayers      int
	ProtocolVersion int32
	VersionName     string
}

func (c TestPing17Config) Check(t *testing.T, status Status17) {
	motd := status.DescriptionText()
	if motd != c.MOTD {
		t.Errorf("expected DescriptionText() = %#v, got %#v", c.MOTD, motd)
	}
	if status.OnlinePlayers != c.OnlinePlayers {
		t.Errorf("expected OnlinePlayers = %#v, got %#v", c.OnlinePlayers, status.OnlinePlayers)
	}
	if status.MaxPlayers != c.MaxPlayers {
		t.Errorf("expected MaxPlayers = %#v, got %#v", c.MaxPlayers, status.MaxPlayers)
	}
	if int32(status.ProtocolVersion) != c.ProtocolVersion {
		t.Errorf("expected ProtocolVersion = %#v, got %#v", c.ProtocolVersion, status.ProtocolVersion)
	}
	if status.VersionName != c.VersionName {
		t.Errorf("expected VersionName = %#v, got %#v", c.VersionName, status.VersionName)
	}
}

var (
	testPing17VanillaConfig = TestPing17Config{
		MOTD:            "A Minecraft Server",
		OnlinePlayers:   0,
		MaxPlayers:      20,
		ProtocolVersion: Ping17ProtocolVersion172,
		VersionName:     "1.7.2",
	}
	testPing17CraftBukkitConfig = TestPing17Config{
		MOTD:            "A Minecraft Server",
		OnlinePlayers:   0,
		MaxPlayers:      20,
		ProtocolVersion: Ping17ProtocolVersion172,
		VersionName:     "CraftBukkit 1.7.2",
	}
	testPing17SpigotConfig = TestPing17Config{
		MOTD:            "A Minecraft Server",
		OnlinePlayers:   0,
		MaxPlayers:      20,
		ProtocolVersion: Ping17ProtocolVersion172,
		VersionName:     "Spigot 1.7.2",
	}
)

func TestPing17(t *testing.T) {
	config := testPing17GetConfig(t)
	testPing17WithDefaultPinger(t, config)
	testPing17WithNewPinger(t, config)
}

func testPing17GetConfig(t *testing.T) TestPing17Config {
	switch Type() {
	case serverTypeVanilla:
		return testPing17VanillaConfig
	case serverTypeCraftBukkit:
		return testPing17CraftBukkitConfig
	case serverTypeSpigot:
		return testPing17SpigotConfig
	}

	t.Fatalf("unknown server type %#v", Type())
	return TestPing17Config{} // Required by compiler.
}

func testPing17WithDefaultPinger(t *testing.T, config TestPing17Config) {
	res, err := Ping17(Hostname(), Port())
	if err != nil {
		t.Errorf("could not ping %s:%d with default pinger: %s", Hostname(), Port(), err)
		return
	}
	config.Check(t, res)
}

func testPing17WithNewPinger(t *testing.T, config TestPing17Config) {
	res, err := NewPinger().Ping17(Hostname(), Port())
	if err != nil {
		t.Errorf("could not ping %s:%d with new pinger: %s", Hostname(), Port(), err)
		return
	}
	config.Check(t, res)
}