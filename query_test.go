package minequery

import (
	"reflect"
	"testing"
)

func TestQueryBasic(t *testing.T) {
	status := &BasicQueryStatus{
		MOTD:          "A Minecraft Server",
		GameType:      "SMP",
		Map:           "world",
		OnlinePlayers: 0,
		MaxPlayers:    20,
		Port:          25565,
		Host:          Hostname(),
	}

	res, err := QueryBasic(Hostname(), Port())
	if err != nil {
		t.Errorf("could not query server: %s", err)
	}

	testQueryCheckBasicStatus(t, status, res)
}

func TestPinger_QueryBasic(t *testing.T) {
	status := &BasicQueryStatus{
		MOTD:          "A Minecraft Server",
		GameType:      "SMP",
		Map:           "world",
		OnlinePlayers: 0,
		MaxPlayers:    20,
		Port:          25565,
		Host:          Hostname(),
	}

	p := NewPinger()
	res, err := p.QueryBasic(Hostname(), Port())
	if err != nil {
		t.Errorf("could not query server: %s", err)
	}

	testQueryCheckBasicStatus(t, status, res)
}

func TestQueryFull(t *testing.T) {
	status := &FullQueryStatus{
		MOTD:          "A Minecraft Server",
		GameType:      "SMP",
		GameID:        "MINECRAFT",
		Version:       "1.4.7",
		ServerVersion: "CraftBukkit on Bukkit 1.4.7-R1.1-SNAPSHOT",
		Plugins:       nil,
		Map:           "world",
		OnlinePlayers: 0,
		MaxPlayers:    20,
		SamplePlayers: []string{},
		Port:          25565,
		Host:          Hostname(),
	}

	res, err := QueryFull(Hostname(), Port())
	if err != nil {
		t.Errorf("could not query server: %s", err)
	}

	testQueryCheckFullStatus(t, status, res)
}

func TestPinger_QueryFull(t *testing.T) {
	status := &FullQueryStatus{
		MOTD:          "A Minecraft Server",
		GameType:      "SMP",
		GameID:        "MINECRAFT",
		Version:       "1.4.7",
		ServerVersion: "CraftBukkit on Bukkit 1.4.7-R1.1-SNAPSHOT",
		Plugins:       nil,
		Map:           "world",
		OnlinePlayers: 0,
		MaxPlayers:    20,
		SamplePlayers: []string{},
		Port:          25565,
		Host:          Hostname(),
	}

	p := NewPinger()
	res, err := p.QueryFull(Hostname(), Port())
	if err != nil {
		t.Errorf("could not query server: %s", err)
	}

	testQueryCheckFullStatus(t, status, res)
}

func testQueryCheckBasicStatus(t *testing.T, expected *BasicQueryStatus, actual *BasicQueryStatus) {
	if expected.MOTD != actual.MOTD {
		t.Errorf("expected MOTD to be %#v but got %#v", expected.MOTD, actual.MOTD)
	}
	if expected.GameType != actual.GameType {
		t.Errorf("expected GameType to be %#v but got %#v", expected.GameType, actual.GameType)
	}
	if expected.Map != actual.Map {
		t.Errorf("expected Map to be %#v but got %#v", expected.Map, actual.Map)
	}
	if expected.OnlinePlayers != actual.OnlinePlayers {
		t.Errorf("expected OnlinePlayers to be %#v but got %#v", expected.OnlinePlayers, actual.OnlinePlayers)
	}
	if expected.MaxPlayers != actual.MaxPlayers {
		t.Errorf("expected MaxPlayers to be %#v but got %#v", expected.MaxPlayers, actual.MaxPlayers)
	}
	if expected.Port != actual.Port {
		t.Errorf("expected Port to be %#v but got %#v", expected.Port, actual.Port)
	}
	if expected.Host == actual.Host {
		t.Errorf("expected Host to be %#v but got %#v", expected.Host, actual.Host)
	}
}

func testQueryCheckFullStatus(t *testing.T, expected *FullQueryStatus, actual *FullQueryStatus) {
	if expected.MOTD != actual.MOTD {
		t.Errorf("expected MOTD to be %#v but got %#v", expected.MOTD, actual.MOTD)
	}
	if expected.GameType != actual.GameType {
		t.Errorf("expected GameType to be %#v but got %#v", expected.GameType, actual.GameType)
	}
	if expected.GameID != actual.GameID {
		t.Errorf("expected GameID to be %#v but got %#v", expected.GameID, actual.GameID)
	}
	if expected.Version != actual.Version {
		t.Errorf("expected Version to be %#v but got %#v", expected.Version, actual.Version)
	}
	if expected.ServerVersion != actual.ServerVersion {
		t.Errorf("expected ServerVersion to be %#v but got %#v", expected.ServerVersion, actual.ServerVersion)
	}
	if !reflect.DeepEqual(expected.Plugins, actual.Plugins) {
		t.Errorf("expected Plugins to be %#v but got %#v", expected.Plugins, actual.Plugins)
	}
	if expected.Map != actual.Map {
		t.Errorf("expected Map to be %#v but got %#v", expected.Map, actual.Map)
	}
	if expected.OnlinePlayers != actual.OnlinePlayers {
		t.Errorf("expected OnlinePlayers to be %#v but got %#v", expected.OnlinePlayers, actual.OnlinePlayers)
	}
	if expected.MaxPlayers != actual.MaxPlayers {
		t.Errorf("expected MaxPlayers to be %#v but got %#v", expected.MaxPlayers, actual.MaxPlayers)
	}
	if !reflect.DeepEqual(expected.SamplePlayers, actual.SamplePlayers) {
		t.Errorf("expected SamplePlayers to be %#v but got %#v", expected.SamplePlayers, actual.SamplePlayers)
	}
	if expected.Port != actual.Port {
		t.Errorf("expected Port to be %#v but got %#v", expected.Port, actual.Port)
	}
	if expected.Host != actual.Host {
		t.Errorf("expected Host to be %#v but got %#v", expected.Host, actual.Host)
	}
}
