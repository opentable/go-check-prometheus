package main

import "testing"

func TestCheckRequiredOptions(t *testing.T) {
	err := checkRequiredOptions()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	host = "example.com"
	err = checkRequiredOptions()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	host = ""
	query = "some query"
	err = checkRequiredOptions()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	query = ""
	warning = "10"
	err = checkRequiredOptions()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	warning = ""
	critical = "100"
	err = checkRequiredOptions()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	critical = ""
	err = checkRequiredOptions()
	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}
}
