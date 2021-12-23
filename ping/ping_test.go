package ping

import (
	"testing"
)

func TestPing(t *testing.T) {
	res, err := Ping("127.0.0.1", 25565)
	if err != nil {
		t.Fatalf("Failed to ping server: %s.", err)
	}

	if res.Version.Name != "1.7.2" {
		t.Errorf("Expected version name 1.7.2, got %s", res.Version.Name)
	}

	if res.Version.Protocol != 4 {
		t.Errorf("Expected protocol version 4, got %d", res.Version.Protocol)
	}

	if res.Description.(map[string]interface{})["text"].(string) != "A Minecraft Server" {
		t.Errorf("Expected description A Minecraft Server, got %s", res.Description.(map[string]interface{})["text"].(string))
	}

	if res.Players.Max != 20 {
		t.Errorf("Expected max players of 20, got %d", res.Players.Max)
	}

	if res.Players.Online != 0 {
		t.Errorf("Expected online players of 0, got %d", res.Players.Online)
	}
}

func TestPingLegacy(t *testing.T) {
	res, err := PingLegacy("127.0.0.1", 25566)
	if err != nil {
		t.Fatalf("Failed to ping server: %s.", err)
	}

	if err != nil {
		t.Fatalf("Failed to ping server: %s.", err)
	}

	if res.Version != "1.6.2" {
		t.Errorf("Expected version name 1.6.2, got %s", res.Version)
	}

	if res.ProtocolVersion != 74 {
		t.Errorf("Expected protocol version 74, got %d", res.ProtocolVersion)
	}

	if res.MessageOfTheDay != "A Minecraft Server" {
		t.Errorf("Expected message of the day A Minecraft Server, got %s", res.MessageOfTheDay)
	}

	if res.MaxPlayers != 20 {
		t.Errorf("Expected max players of 20, got %d", res.MaxPlayers)
	}

	if res.PlayerCount != 0 {
		t.Errorf("Expected online players of 0, got %d", res.PlayerCount)
	}
}

func TestPingAncient(t *testing.T) {
	res, err := PingAncient("127.0.0.1", 25568)
	if err != nil {
		t.Fatalf("Failed to ping server: %s.", err)
	}

	if res.MessageOfTheDay != "A Minecraft Server" {
		t.Errorf("Expected message of the day A Minecraft Server, got %s", res.MessageOfTheDay)
	}

	if res.MaxPlayers != 20 {
		t.Errorf("Expected max players of 20, got %d", res.MaxPlayers)
	}

	if res.PlayerCount != 0 {
		t.Errorf("Expected online players of 0, got %d", res.PlayerCount)
	}
}
