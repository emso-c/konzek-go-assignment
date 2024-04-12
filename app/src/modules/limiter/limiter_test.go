package limiter

import (
	"os"
	"testing"
	"time"
)

func setup() {
	os.Setenv("HTTP_RATE_LIMIT", "2")
	os.Setenv("HTTP_RATE_LIMIT_WINDOW", "1")
}

func teardown() {
	os.Unsetenv("HTTP_RATE_LIMIT")
	os.Unsetenv("HTTP_RATE_LIMIT_WINDOW")
}

func TestExceedsLimit(t *testing.T) {
	setup()
	defer teardown()

	os.Setenv("HTTP_RATE_LIMIT", "2")
	l := NewLimiter()
	l.Initialize()

	// Test with a client that has not exceeded the limit
	l.Increment("client1")
	if l.ExceedsLimit("client1") {
		t.Errorf("ExceedsLimit returned true for client1, want false")
	}

	// Test with a client that has exceeded the limit
	l.Increment("client2")
	l.Increment("client2")
	l.Increment("client2")
	if !l.ExceedsLimit("client2") {
		t.Errorf("ExceedsLimit returned false for client2, want true")
	}

	// Test with a client that has not made any requests
	if l.ExceedsLimit("client3") {
		t.Errorf("ExceedsLimit returned true for client3, want false")
	}

	// Should throw Fatal log message if HTTP_RATE_LIMIT is not set
	os.Unsetenv("HTTP_RATE_LIMIT")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ExceedsLimit did not panic when HTTP_RATE_LIMIT is not set")
		}
	}()

	l.ExceedsLimit("client1")
}

func TestIncrement(t *testing.T) {
	setup()
	defer teardown()

	l := NewLimiter()
	l.Initialize()

	// Test incrementing the request count for a client
	l.Increment("client1")
	if l.UsageMap["client1"] != 1 {
		t.Errorf("Increment did not increase the request count for client1")
	}

	// Test incrementing the request count for a different client
	l.Increment("client2")
	if l.UsageMap["client2"] != 1 {
		t.Errorf("Increment did not increase the request count for client2")
	}
}

func TestInitialize(t *testing.T) {
	setup()
	defer teardown()

	os.Setenv("HTTP_RATE_LIMIT_WINDOW", "1")
	l := NewLimiter()
	l.Initialize()

	// Test that the usage map is reset after the time window
	l.Increment("client1")
	time.Sleep(2 * time.Second)
	if l.UsageMap["client1"] != 0 {
		t.Errorf("Initialize did not reset the usage map after the time window")
	}
}

func TestGetLimiter(t *testing.T) {
	setup()
	defer teardown()

	l1 := GetLimiter()
	l2 := GetLimiter()

	if l1 != l2 {
		t.Errorf("GetLimiter did not return the shared singleton instance of the Limiter")
	}
}
