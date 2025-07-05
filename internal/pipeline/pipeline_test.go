package pipeline

import "testing"

func TestAlwaysPass(t *testing.T) {
	t.Log("✅ Test ran successfully")
}

func TestAlwaysFail(t *testing.T) {
	t.Error("❌ This test is designed to fail")
}
