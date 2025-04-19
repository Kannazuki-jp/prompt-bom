package domain

import (
	"testing"
	"github.com/Masterminds/semver/v3"
)

func TestSemverParse(t *testing.T) {
	v, err := semver.NewVersion("1.2.3")
	if err != nil {
		t.Fatalf("valid semver should parse: %v", err)
	}
	if v.String() != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", v.String())
	}

	_, err = semver.NewVersion("not-semver")
	if err == nil {
		t.Error("invalid semver should fail to parse")
	}
}

func TestComputeSHA256(t *testing.T) {
	data := []byte("hello world")
	hash := ComputeSHA256(data)
	if len(hash) != 71 || hash[:7] != "sha256:" {
		t.Errorf("unexpected hash format: %s", hash)
	}
	// 再現性
	hash2 := ComputeSHA256(data)
	if hash != hash2 {
		t.Error("hash should be deterministic")
	}
}

func TestDetectCycle(t *testing.T) {
	// サイクルなし
	deps := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {},
	}
	if DetectCycle(deps) {
		t.Error("should not detect cycle")
	}
	// サイクルあり
	deps2 := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {"A"},
	}
	if !DetectCycle(deps2) {
		t.Error("should detect cycle")
	}
} 