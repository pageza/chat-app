package main

import (
	"testing"
)

// TestMainFunction tests if the main function can run without panicking or exiting.
func TestMainFunction(t *testing.T) {
	// This function will recover from any panics and fail the test if one occurs.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did not complete successfully: %v", r)
		}
	}()

	// Call the main function.
	main()
}
