package ga2_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
)

type mockMetaGa2Listener struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaGa2Listener) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaGa2Listener{}

func newMockMetaGa2Listener() *mockMetaGa2Listener {
	return &mockMetaGa2Listener{client: &connectivity.TencentCloudClient{}}
}

func ptrStringListener(s string) *string { return &s }
func ptrBoolListener(v bool) *bool       { return &v }
func ptrUint64Listener(v uint64) *uint64 { return &v }

// go test ./tencentcloud/services/ga2/ -run "TestGa2Listener" -v -count=1 -gcflags="all=-l"

// TestGa2Listener_Create_HttpVersion tests the Create operation with http_version set for HTTPS listener
func TestGa2Listener_Create_HttpVersion(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2Listener().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateListenerWithContext", func(ctx context.Context, request *ga2v20250115.CreateListenerRequest) (*ga2v20250115.CreateListenerResponse, error) {
		assert.Equal(t, "ga-test123", *request.GlobalAcceleratorId)
		assert.Equal(t, "HTTPS", *request.Protocol)
		assert.Equal(t, "HTTP/2", *request.HttpVersion)
		assert.Equal(t, "MUTUAL", *request.CertificationType)
		assert.Equal(t, "tls_policy_1.2_strict-1.3", *request.CipherPolicyId)

		resp := ga2v20250115.NewCreateListenerResponse()
		resp.Response = &ga2v20250115.CreateListenerResponseParams{
			TaskId:     ptrStringListener("task-001"),
			ListenerId: ptrStringListener("lsr-test456"),
			RequestId:  ptrStringListener("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeListenersWithContext for the Read call after Create
	patches.ApplyMethodFunc(ga2Client, "DescribeListenersWithContext", func(ctx context.Context, request *ga2v20250115.DescribeListenersRequest) (*ga2v20250115.DescribeListenersResponse, error) {
		resp := ga2v20250115.NewDescribeListenersResponse()
		resp.Response = &ga2v20250115.DescribeListenersResponseParams{
			ListenerSet: []*ga2v20250115.ListenerSet{
				{
					GlobalAcceleratorId: ptrStringListener("ga-test123"),
					ListenerId:          ptrStringListener("lsr-test456"),
					Name:                ptrStringListener("tf-example-https"),
					Protocol:            ptrStringListener("HTTPS"),
					PortRanges: &ga2v20250115.PortRanges{
						FromPort: ptrUint64Listener(9090),
						ToPort:   ptrUint64Listener(9090),
					},
					CertificationType:    ptrStringListener("MUTUAL"),
					CipherPolicyId:       ptrStringListener("tls_policy_1.2_strict-1.3"),
					ServerCertificates:   []*string{ptrStringListener("Yj6CmODs")},
					ClientCaCertificates: []*string{ptrStringListener("W6aH2tOc")},
					HttpVersion:          ptrStringListener("HTTP/2"),
					IdleTimeout:          ptrUint64Listener(38),
					RequestTimeout:       ptrUint64Listener(60),
					XForwardedForRealIp:  ptrBoolListener(true),
					ListenerType:         ptrStringListener("Standard"),
					Status:               ptrStringListener("RUNNING"),
				},
			},
			TotalCount: ptrUint64Listener(1),
			RequestId:  ptrStringListener("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResultWithContext for async task polling
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringListener("SUCCESS"),
			RequestId: ptrStringListener("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2Listener()
	res := ga2.ResourceTencentCloudGa2Listener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test123",
		"protocol":              "HTTPS",
		"port_ranges": []interface{}{
			map[string]interface{}{
				"from_port": 9090,
				"to_port":   9090,
			},
		},
		"certification_type":     "MUTUAL",
		"cipher_policy_id":       "tls_policy_1.2_strict-1.3",
		"server_certificates":    []interface{}{"Yj6CmODs"},
		"client_ca_certificates": []interface{}{"W6aH2tOc"},
		"http_version":           "HTTP/2",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-test123#lsr-test456", d.Id())
	assert.Equal(t, "HTTP/2", d.Get("http_version").(string))
}

// TestGa2Listener_Create_HttpVersionOmitted tests the Create operation without http_version set
func TestGa2Listener_Create_HttpVersionOmitted(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2Listener().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateListenerWithContext", func(ctx context.Context, request *ga2v20250115.CreateListenerRequest) (*ga2v20250115.CreateListenerResponse, error) {
		assert.Equal(t, "ga-test123", *request.GlobalAcceleratorId)
		assert.Equal(t, "TCP", *request.Protocol)
		// http_version should NOT be set on the request
		assert.Nil(t, request.HttpVersion)

		resp := ga2v20250115.NewCreateListenerResponse()
		resp.Response = &ga2v20250115.CreateListenerResponseParams{
			TaskId:     ptrStringListener("task-002"),
			ListenerId: ptrStringListener("lsr-test789"),
			RequestId:  ptrStringListener("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeListenersWithContext for the Read call after Create
	patches.ApplyMethodFunc(ga2Client, "DescribeListenersWithContext", func(ctx context.Context, request *ga2v20250115.DescribeListenersRequest) (*ga2v20250115.DescribeListenersResponse, error) {
		resp := ga2v20250115.NewDescribeListenersResponse()
		resp.Response = &ga2v20250115.DescribeListenersResponseParams{
			ListenerSet: []*ga2v20250115.ListenerSet{
				{
					GlobalAcceleratorId: ptrStringListener("ga-test123"),
					ListenerId:          ptrStringListener("lsr-test789"),
					Name:                ptrStringListener("tf-example-tcp"),
					Protocol:            ptrStringListener("TCP"),
					PortRanges: &ga2v20250115.PortRanges{
						FromPort: ptrUint64Listener(80),
						ToPort:   ptrUint64Listener(80),
					},
					IdleTimeout:  ptrUint64Listener(900),
					ListenerType: ptrStringListener("Standard"),
					Status:       ptrStringListener("RUNNING"),
				},
			},
			TotalCount: ptrUint64Listener(1),
			RequestId:  ptrStringListener("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResultWithContext for async task polling
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringListener("SUCCESS"),
			RequestId: ptrStringListener("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2Listener()
	res := ga2.ResourceTencentCloudGa2Listener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test123",
		"protocol":              "TCP",
		"port_ranges": []interface{}{
			map[string]interface{}{
				"from_port": 80,
				"to_port":   80,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-test123#lsr-test789", d.Id())
}

// TestGa2Listener_Read tests the Read operation with http_version populated
func TestGa2Listener_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2Listener().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeListenersWithContext", func(ctx context.Context, request *ga2v20250115.DescribeListenersRequest) (*ga2v20250115.DescribeListenersResponse, error) {
		resp := ga2v20250115.NewDescribeListenersResponse()
		resp.Response = &ga2v20250115.DescribeListenersResponseParams{
			ListenerSet: []*ga2v20250115.ListenerSet{
				{
					GlobalAcceleratorId: ptrStringListener("ga-read123"),
					ListenerId:          ptrStringListener("lsr-read456"),
					Name:                ptrStringListener("read-example"),
					Protocol:            ptrStringListener("HTTPS"),
					PortRanges: &ga2v20250115.PortRanges{
						FromPort: ptrUint64Listener(443),
						ToPort:   ptrUint64Listener(443),
					},
					CertificationType:   ptrStringListener("UNIDIRECTIONAL"),
					CipherPolicyId:      ptrStringListener("tls_policy_1.2"),
					ServerCertificates:  []*string{ptrStringListener("cert-001")},
					HttpVersion:         ptrStringListener("HTTP/2"),
					IdleTimeout:         ptrUint64Listener(30),
					RequestTimeout:      ptrUint64Listener(60),
					XForwardedForRealIp: ptrBoolListener(true),
					ListenerType:        ptrStringListener("Standard"),
					Status:              ptrStringListener("RUNNING"),
				},
			},
			TotalCount: ptrUint64Listener(1),
			RequestId:  ptrStringListener("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2Listener()
	res := ga2.ResourceTencentCloudGa2Listener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "",
		"protocol":              "",
		"port_ranges": []interface{}{
			map[string]interface{}{
				"from_port": 0,
				"to_port":   0,
			},
		},
	})
	d.SetId("ga-read123#lsr-read456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-read123#lsr-read456", d.Id())
	assert.Equal(t, "HTTP/2", d.Get("http_version").(string))
	assert.Equal(t, "HTTPS", d.Get("protocol").(string))
}

// TestGa2Listener_Read_NotFound tests Read when listener is not found
func TestGa2Listener_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2Listener().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeListenersWithContext", func(ctx context.Context, request *ga2v20250115.DescribeListenersRequest) (*ga2v20250115.DescribeListenersResponse, error) {
		resp := ga2v20250115.NewDescribeListenersResponse()
		resp.Response = &ga2v20250115.DescribeListenersResponseParams{
			ListenerSet: []*ga2v20250115.ListenerSet{},
			TotalCount:  ptrUint64Listener(0),
			RequestId:   ptrStringListener("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2Listener()
	res := ga2.ResourceTencentCloudGa2Listener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "",
		"protocol":              "",
		"port_ranges": []interface{}{
			map[string]interface{}{
				"from_port": 0,
				"to_port":   0,
			},
		},
	})
	d.SetId("ga-notfound#lsr-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestAccTencentCloudGa2ListenerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGa2Listener,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_listener.example", "id"),
				),
			},
			{
				Config: testAccGa2ListenerUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_listener.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_ga2_listener.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGa2Listener = `
resource "tencentcloud_ga2_listener" "example" {
  global_accelerator_id = "ga-4mredmiu"
  name                  = "tf-example"
  protocol              = "TCP"

  port_ranges {
    from_port = 80
    to_port   = 80
  }

  description     = "tf example listener"
  client_affinity = "Open"
  idle_timeout    = 900
}
`

const testAccGa2ListenerUpdate = `
resource "tencentcloud_ga2_listener" "example" {
  global_accelerator_id = "ga-4mredmiu"
  name                  = "tf-example"
  protocol              = "TCP"

  port_ranges {
    from_port = 80
    to_port   = 80
  }

  description     = "tf example listener"
  client_affinity = "Open"
  idle_timeout    = 900
}
`
