//go:build integration

package integration_test

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const dependencyTimeout = 10 * time.Second

func requireEnv(t *testing.T, key string) string {
	t.Helper()

	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		t.Fatalf("%s is required for integration tests", key)
	}
	return value
}

func requireIntEnv(t *testing.T, key string) int {
	t.Helper()

	raw := requireEnv(t, key)
	value, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("parse %s: %v", key, err)
	}
	return value
}
