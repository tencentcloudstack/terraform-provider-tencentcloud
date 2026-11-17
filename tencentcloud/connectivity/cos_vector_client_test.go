package connectivity

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func TestCosVectorClientConfiguration(t *testing.T) {
	shared := &TencentCloudClient{Region: "ap-guangzhou", CosDomain: "http://10.0.0.1", Credential: common.NewTokenCredential("local-id", "local-key", "local-token")}
	client, err := shared.UseTencentCosVectorClient()
	if err != nil {
		t.Fatal(err)
	}
	if client.BaseURL.VectorURL.String() != "https://vectors.ap-guangzhou.coslake.com" || client.BaseURL.BucketURL != nil {
		t.Fatal("vector and normal COS endpoints were mixed")
	}
	if client.Conf.RetryOpt.Count != 1 || !strings.HasPrefix(client.UserAgent, "Terraform-") {
		t.Fatal("unexpected retry count or User-Agent")
	}
	if _, err := (&TencentCloudClient{Region: "ap-guangzhou"}).UseTencentCosVectorClient(); err == nil {
		t.Fatal("missing credentials accepted")
	}
	shared.Region = "../other-host"
	if _, err := shared.UseTencentCosVectorClient(); err == nil {
		t.Fatal("invalid region accepted")
	}
}

func TestCosVectorClientReuseAndIsolation(t *testing.T) {
	shared := &TencentCloudClient{Region: "ap-guangzhou", Credential: common.NewTokenCredential("local-id", "local-key", "local-token")}
	first, err := shared.UseTencentCosVectorClient()
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := shared.UseTencentCosVectorClient()
			if err != nil || got != first {
				t.Error("parallel access did not reuse the client")
			}
		}()
	}
	wg.Wait()
	oldAuth := shared.cosVectorConn.auth
	shared.Credential = common.NewTokenCredential("new-id", "new-key", "new-token")
	rotated, err := shared.UseTencentCosVectorClient()
	if err != nil || rotated == first {
		t.Fatal("credential rotation reused stale client")
	}
	id, _, token, _ := oldAuth.GetCredential()
	if id != "local-id" || token != "local-token" {
		t.Fatal("rotation mutated old client credentials")
	}
	shared.Region = "ap-shanghai"
	regional, err := shared.UseTencentCosVectorClient()
	if err != nil || regional == rotated || regional.BaseURL.VectorURL.Host != "vectors.ap-shanghai.coslake.com" {
		t.Fatal("region change reused wrong client")
	}
	other := &TencentCloudClient{Region: shared.Region, Credential: shared.Credential}
	isolated, err := other.UseTencentCosVectorClient()
	if err != nil || isolated == regional {
		t.Fatal("different provider clients share a vector client")
	}
}

type cosVectorTestRoundTripper func(*http.Request) (*http.Response, error)

func (f cosVectorTestRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type cosVectorBrokenBody struct{}

func (cosVectorBrokenBody) Read(_ []byte) (int, error) { return 0, errors.New("broken response") }
func (cosVectorBrokenBody) Close() error               { return nil }

type cosVectorReadErrorBody struct{ cause error }

func (b cosVectorReadErrorBody) Read(_ []byte) (int, error) { return 0, b.cause }
func (cosVectorReadErrorBody) Close() error                 { return nil }

func TestCosVectorBodyReadFailureClassification(t *testing.T) {
	for _, tc := range []struct {
		name  string
		cause error
		code  string
	}{
		{"interrupted", io.ErrUnexpectedEOF, "ConnectionFailure"},
		{"canceled", context.Canceled, "Canceled"},
		{"timeout", context.DeadlineExceeded, "Timeout"},
		{"invalid", errors.New("private detail"), "InvalidResponse"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := &CosVectorResponseTransport{Transport: cosVectorTestRoundTripper(func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: cosVectorReadErrorBody{cause: tc.cause}}, nil
			})}
			req, err := http.NewRequest(http.MethodPost, "https://vectors.ap-guangzhou.coslake.com/CreateVectorBucket", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := transport.RoundTrip(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Error(err)
				}
			}()
			if resp.StatusCode != 200 || resp.Header.Get(CosVectorInvalidResponseHeader) != tc.code {
				t.Fatal("body failure lost success status or safe error category")
			}
		})
	}
}

func TestCosVectorResponseValidation(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		path    string
		body    io.ReadCloser
		invalid bool
	}{
		{"valid404", 404, "/GetVectorBucket", io.NopCloser(strings.NewReader(`{"message":"missing"}`)), false},
		{"malformed404", 404, "/GetVectorBucket", io.NopCloser(strings.NewReader(`{"message":`)), true},
		{"empty404", 404, "/GetVectorBucket", io.NopCloser(strings.NewReader("")), true},
		{"broken404", 404, "/GetVectorBucket", cosVectorBrokenBody{}, true},
		{"brokenCreate", 200, "/CreateVectorBucket", cosVectorBrokenBody{}, true},
		{"oversize", 200, "/GetVectorBucket", io.NopCloser(strings.NewReader(strings.Repeat("a", (1<<20)+2))), true},
		{"emptyDelete", 200, "/DeleteVectorBucket", io.NopCloser(strings.NewReader("")), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := &CosVectorResponseTransport{Transport: cosVectorTestRoundTripper(func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: tc.status, Header: make(http.Header), Body: tc.body}, nil
			})}
			req, err := http.NewRequest(http.MethodPost, "https://vectors.ap-guangzhou.coslake.com"+tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := transport.RoundTrip(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Error(err)
				}
			}()
			if (resp.Header.Get(CosVectorInvalidResponseHeader) != "") != tc.invalid || resp.StatusCode != tc.status {
				t.Fatal("invalid response classification or lost HTTP status")
			}
		})
	}
}

func TestCosVectorPublicIP(t *testing.T) {
	for _, ip := range []string{"127.0.0.1", "10.0.0.1", "172.16.0.1", "192.168.1.1", "9.1.1.1", "11.1.1.1", "21.1.1.1", "30.1.1.1", "100.64.0.1", "169.254.169.254", "::1", "fc00::1"} {
		if cosVectorPublicIP(net.ParseIP(ip)) {
			t.Fatalf("private or restricted IP allowed: %s", ip)
		}
	}
	if !cosVectorPublicIP(net.ParseIP("8.8.8.8")) {
		t.Fatal("public IP rejected")
	}
}
