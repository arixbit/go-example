package model

import "testing"

func TestExampleTableName(t *testing.T) {
	if got := (Example{}).TableName(); got != "examples" {
		t.Fatalf("TableName() = %q, want %q", got, "examples")
	}
}
