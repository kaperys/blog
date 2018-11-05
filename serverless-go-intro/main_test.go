package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

const mockServerURL = "127.0.0.1:9294"

func TestHandlerSuccess(t *testing.T) {
	defer teardown(setup(t, "[ { \"emoji\": \"👌\", \"aliases\": [\"dummy-emoji\"] } ]"))

	res, _ := handler(events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{searchKey: "dummy-emoji"}})

	if v := res.Body; v != "👌" {
		t.Fatalf("TestHandlerSuccess failed: have %q, want %q", v, "👌")
	}
}

func TestHandlerFailEmojiDoesNotExist(t *testing.T) {
	defer teardown(setup(t, "[ { \"emoji\": \"👌\", \"aliases\": [\"dummy-emoji\"] } ]"))

	res, _ := handler(events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{searchKey: "emoji-that-does-not-exist"}})
	if res.StatusCode == http.StatusOK {
		t.Fatalf("TestHandlerFailEmojiDoesNotExist failed: expected an error")
	}

	x := "no results for \"emoji-that-does-not-exist\""
	if v := res.Body; v != x {
		t.Fatalf("TestHandlerFailEmojiDoesNotExist failed: have %q, want %q", v, x)
	}
}

func TestHandlerFailNoRouteToHost(t *testing.T) {
	res, _ := handler(events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{searchKey: "dummy-emoji"}})
	if res.StatusCode == http.StatusOK {
		t.Fatalf("TestHandlerFailNoRouteToHost failed: expected an error")
	}

	x := "could not retrieve emoji data: could not make the request: Get http://127.0.0.1:9294: dial tcp 127.0.0.1:9294: connect: connection refused"
	if v := res.Body; v != x {
		t.Fatalf("TestHandlerFailNoRouteToHost failed: have %q, want %q", v, x)
	}
}

func setup(t *testing.T, body string) *httptest.Server {
	os.Setenv("SOURCE_URL", "http://"+mockServerURL)

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, body) }))

	l, err := net.Listen("tcp", mockServerURL)
	if err != nil {
		t.Fatalf("could not start the mock server: %v", err)
	}

	ts.Listener = l
	ts.Start()

	return ts
}

func teardown(ts *httptest.Server) {
	ts.Close()
}
