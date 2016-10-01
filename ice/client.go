package ice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strings"
)

// NewClient creates an iCE client for the given endpoint URL.
func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// Client provides an interface to the iCE Server REST API.
type Client struct {
	endpoint string
}

// MyIP returns the IP of the running machine as seen from the iCE server.
func (i *Client) MyIP(ctx context.Context) (net.IP, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/my_ip", i.endpoint))
	if err != nil {
		return net.IP{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return net.IP{}, fmt.Errorf("Error: got HTTP response %s", resp.Status)
	}

	respContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return net.IP{}, err
	}
	ipStr := strings.TrimSpace(string(respContents))

	return net.ParseIP(ipStr), nil
}

type storeInstanceResponseError struct {
	Message string `json:"message"`
}
type storeInstanceResponse struct {
	ID string `json:"_id"`

	// in case of a failure
	Error  storeInstanceResponseError `json:"_error"`
	Issues map[string]string          `json:"_issues"`
}

func storeInstanceErrorMessage(resp storeInstanceResponse) string {
	errMsg := resp.Error.Message

	if len(resp.Issues) > 0 {
		parts := []string{}
		for field, issue := range resp.Issues {
			parts = append(parts, fmt.Sprintf("`%s`: %s", field, issue))
		}
		// Sort the parts of the error message to avoid surprises/flakes with the
		// test assertions.
		sort.Sort(sort.StringSlice(parts))
		errMsg += " (" + strings.Join(parts, ", ") + ")"
	}

	return errMsg
}

// StoreInstance submits an iCE instance to the iCE server.
func (i *Client) StoreInstance(ctx context.Context, inst Instance) (string, error) {
	bodyBuffer := bytes.NewBuffer([]byte{})

	if err := json.NewEncoder(bodyBuffer).Encode(inst); err != nil {
		return "", err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/v2/instances", i.endpoint),
		"application/json", bodyBuffer,
	)
	if err != nil {
		return "", err
	}

	var respParsed storeInstanceResponse
	err = json.NewDecoder(resp.Body).Decode(&respParsed)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("Failed to parse response: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := storeInstanceErrorMessage(respParsed)
		if errMsg == "" {
			errMsg = fmt.Sprintf("Erroror: got HTTP response %s", resp.Status)
		}
		return "", errors.New(errMsg)
	}

	instID := respParsed.ID
	if instID == "" {
		return "", errors.New("Error: response does not include the `_id` field")
	}

	return instID, nil
}
