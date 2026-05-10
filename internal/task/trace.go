package task

import "encoding/json"

// TraceIDFromPayload extracts trace_id from a JSON task payload.
func TraceIDFromPayload(payload []byte) string {
	var envelope struct {
		TraceID string `json:"trace_id"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return ""
	}
	return envelope.TraceID
}
