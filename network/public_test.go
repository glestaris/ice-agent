package network

import (
	"strings"
	"testing"
)

func TestReverseDNS(t *testing.T) {
	ipAddr := "8.8.8.8"
	res, err := ReverseDNS(nil, ipAddr)
	if err != nil {
		t.Fatal(err)
	}

	expectedRes := "google-public-dns-a.google.com"
	if res != expectedRes {
		t.Fatalf(
			"Expected reverse DNS of `%s` to be `%s`, got `%s`",
			ipAddr, expectedRes, res,
		)
	}
}

func TestReverseDNSReturnsErrorWhenCalledWithAnInvalidAddress(t *testing.T) {
	ipAddr := "800.800.800.800"
	_, err := ReverseDNS(nil, ipAddr)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	expectedError := "unrecognized address"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf(
			"Expected error message to contain `%s`, got `%s`",
			expectedError, err.Error(),
		)
	}
}
