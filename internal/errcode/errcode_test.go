package errcode

import "testing"

func TestErrorFields(t *testing.T) {
	if InvalidParams.Code() != 1001 {
		t.Fatalf("Code() = %d, want 1001", InvalidParams.Code())
	}
	if InvalidParams.Reason() != "INVALID_PARAMS" {
		t.Fatalf("Reason() = %q, want INVALID_PARAMS", InvalidParams.Reason())
	}
	if InvalidParams.Error() != InvalidParams.Reason() {
		t.Fatalf("Error() = %q, want %q", InvalidParams.Error(), InvalidParams.Reason())
	}
}
