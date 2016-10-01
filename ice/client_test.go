package ice

import (
	"net/http"
	"strings"
	"testing"

	ice_testing "github.com/glestaris/ice-agent/testing"
)

func TestMyIP(t *testing.T) {
	fakeServer := ice_testing.NewFakeIceServer()
	fakeServer.Start()
	defer fakeServer.Stop()
	fakeServer.Handle(
		http.MethodGet, "/v2/my_ip",
		func(resp http.ResponseWriter, req *http.Request) {
			resp.WriteHeader(http.StatusOK)
			if _, err := resp.Write([]byte("10.20.30.40\n")); err != nil {
				t.Fatal(err)
			}
		},
	)

	c := NewClient(fakeServer.Endpoint())
	myIP, err := c.MyIP(nil)
	if err != nil {
		t.Fatal(err)
	}

	if myIP.String() != "10.20.30.40" {
		t.Fatalf("Expected my IP to be `10.20.30.40`, got `%s`", myIP.String())
	}
}

var myIPErrorTests = []struct {
	httpResponseStatus     int
	httpResponseBody       string
	expectedErrorSubstring string
}{
	{http.StatusNotFound, "", http.StatusText(http.StatusNotFound)},
}

func TestMyIPErrors(t *testing.T) {
	for _, test := range myIPErrorTests {
		fakeServer := ice_testing.NewFakeIceServer()
		fakeServer.Start()
		defer fakeServer.Stop()
		fakeServer.Handle(
			http.MethodGet, "/v2/my_ip",
			func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(test.httpResponseStatus)
				if _, err := resp.Write([]byte(test.httpResponseBody)); err != nil {
					t.Fatal(err)
				}
			},
		)

		c := NewClient(fakeServer.Endpoint())
		_, err := c.MyIP(nil)
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), test.expectedErrorSubstring) {
			t.Errorf(
				"Expected error message to contain `%s`, got `%s`",
				test.expectedErrorSubstring,
				err.Error(),
			)
		}
	}
}

func TestStoreInstance(t *testing.T) {
	fakeServer := ice_testing.NewFakeIceServer()
	fakeServer.Start()
	defer fakeServer.Stop()
	fakeServer.Handle(
		http.MethodPost, "/v2/instances",
		func(resp http.ResponseWriter, req *http.Request) {
			resp.WriteHeader(http.StatusCreated)
			if _, err := resp.Write([]byte(`{
	"_id": "test-id"
}`)); err != nil {
				t.Fatal(err)
			}
		},
	)

	c := NewClient(fakeServer.Endpoint())
	instID, err := c.StoreInstance(nil, Instance{})
	if err != nil {
		t.Fatal(err)
	}

	if instID != "test-id" {
		t.Fatalf("Expected the instance ID to be `test-id`, got `%s`", instID)
	}
}

var storeInstnanceErrorTests = []struct {
	httpResponseStatus     int
	httpResponseBody       string
	expectedErrorSubstring string
}{
	{http.StatusNotFound, "", http.StatusText(http.StatusNotFound)},
	{http.StatusOK, "{[", "Failed to parse response"},
	{http.StatusOK, "{}", "Error: response does not include the `_id` field"},
	{
		http.StatusBadRequest, `{
	"_error": {
		"message": "failed to do stuff"
	}
}`,
		"failed to do stuff",
	},
	{
		http.StatusBadRequest, `{
	"_issues": {
		"key_a": "issue_a",
		"key_b": "issue_b"
	}
}`,
		"`key_a`: issue_a, `key_b`: issue_b",
	},
}

func TestStoreInstanceErrors(t *testing.T) {
	for _, test := range storeInstnanceErrorTests {
		fakeServer := ice_testing.NewFakeIceServer()
		fakeServer.Start()
		defer fakeServer.Stop()
		fakeServer.Handle(
			http.MethodPost, "/v2/instances",
			func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(test.httpResponseStatus)
				if _, err := resp.Write([]byte(test.httpResponseBody)); err != nil {
					t.Fatal(err)
				}
			},
		)

		c := NewClient(fakeServer.Endpoint())
		_, err := c.StoreInstance(nil, Instance{})
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), test.expectedErrorSubstring) {
			t.Errorf(
				"Expected error message to contain `%s`, got `%s`",
				test.expectedErrorSubstring,
				err.Error(),
			)
		}
	}
}
