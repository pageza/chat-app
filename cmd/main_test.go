package main

import (
	"testing"
	"time"
)

// TestMainFunction tests if the main function can run without panicking or exiting.
func TestMainFunction(t *testing.T) {
	// This function will recover from any panics and fail the test if one occurs.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did not complete successfully: %v", r)
		}
	}()

	// Create a channel to signal completion of main function
	done := make(chan bool)

	// Run the main function in a goroutine
	go func() {
		main()
		close(done)
	}()

	// Wait for main to complete or time out
	select {
	case <-done:
		// main completed successfully
	case <-time.After(time.Second * 2): // 2 seconds timeout
		t.Errorf("The main function timed out")
	}
}
