package globals

import "testing"

func TestGlobalTestfireNetworkConfiguration(t *testing.T) {
	if GameServerID != "1017E300" {
		t.Fatalf("unexpected Testfire game server ID: %s", GameServerID)
	}

	if AccessKey != "da693ee5" {
		t.Fatalf("unexpected Testfire access key: %s", AccessKey)
	}

	if NEXMajor != 3 || NEXMinor != 8 || NEXPatch != 3 {
		t.Fatalf("unexpected Testfire NEX version: %d.%d.%d", NEXMajor, NEXMinor, NEXPatch)
	}
}
