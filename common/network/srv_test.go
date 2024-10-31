package network

import (
	"context"
	"net"
	"net/http"
	"testing"
)

func TestSystemOutbound_DialContextWithSRV(t *testing.T) {
	// _httpsgstatic._tcp.dummy.launium.com is a private SRV record
	// that points to gstatic.com:443. This record is for testing purposes only.
	// Availability is not guaranteed. DO NOT USE IT IN PRODUCTION.
	transport := http.Transport{
		DisableKeepAlives: true,
		DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
			t.Log("Dialing with SRV")
			return SystemDialer.DialContextWithSRV(context.Background(), "tcp", "dummy.launium.com", "httpsgstatic")
		},
	}
	req, err := http.NewRequest(http.MethodGet, "https://gstatic.com/generate_204", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent { // 204
		t.Fatalf("got %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}
