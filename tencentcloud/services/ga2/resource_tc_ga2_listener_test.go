package ga2_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcga2 "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
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

func ptrStringGa2Listener(s string) *string {
	return &s
}

func ptrUint64Ga2Listener(v uint64) *uint64 {
	return &v
}

func ptrBoolGa2Listener(v bool) *bool {
	return &v
}

// go test ./tencentcloud/services/ga2/ -run "TestGa2Listener_" -v -count=1 -gcflags="all=-l"

func TestGa2Listener_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2Listener().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateListenerWithContext", func(_ context.Context, _ *ga2v20250115.CreateListenerRequest) (*ga2v20250115.CreateListenerResponse, error) {
		resp := ga2v20250115.NewCreateListenerResponse()
		resp.Response = &ga2v20250115.CreateListenerResponseParams{
			TaskId:     ptrStringGa2Listener("task-12345"),
			ListenerId: ptrStringGa2Listener("lbl-abcdef"),
			RequestId:  ptrStringGa2Listener("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "WaitForGa2TaskFinish", func(_ context.Context, taskId string, _ interface{}) error {
		assert.Equal(t, "task-12345", taskId)
		return nil
	})

	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "DescribeGa2ListenerById", func(_ context.Context, gaId string, listenerId string) (*ga2v20250115.ListenerSet, error) {
		assert.Equal(t, "ga2-test1234", gaId)
		assert.Equal(t, "lbl-abcdef", listenerId)
		return &ga2v20250115.ListenerSet{
			GlobalAcceleratorId: ptrStringGa2Listener("ga2-test1234"),
			ListenerId:          ptrStringGa2Listener("lbl-abcdef"),
			Name:                ptrStringGa2Listener("test-listener"),
			Protocol:            ptrStringGa2Listener("TCP"),
			ListenerType:        ptrStringGa2Listener("INTELLIGENT"),
			PortRanges: &ga2v20250115.PortRanges{
				FromPort: ptrUint64Ga2Listener(80),
				ToPort:   ptrUint64Ga2Listener(80),
			},
			IdleTimeout:    ptrUint64Ga2Listener(900),
			ClientAffinity: ptrStringGa2Listener("SOURCE_IP"),
			GetRealIpType:  ptrStringGa2Listener("TOA"),
		}, nil
	})

	meta := newMockMetaGa2Listener()
	res := svcga2.ResourceTencentCloudGa2Listener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test1234",
		"name":                  "test-listener",
		"protocol":              "TCP",
		"listener_type":         "INTELLIGENT",
		"port_ranges": []interface{}{
			map[string]interface{}{
				"from_port": 80,
				"to_port":   80,
			},
		},
		"idle_timeout":     900,
		"client_affinity":  "SOURCE_IP",
		"get_real_ip_type": "TOA",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, "ga2-test1234#lbl-abcdef", d.Id())
	assert.Equal(t, "lbl-abcdef", d.Get("listener_id").(string))
}

func TestGa2Listener_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "DescribeGa2ListenerById", func(_ context.Context, gaId string, listenerId string) (*ga2v20250115.ListenerSet, error) {
		assert.Equal(t, "ga2-test1234", gaId)
		assert.Equal(t, "lbl-abcdef", listenerId)
		return &ga2v20250115.ListenerSet{
			GlobalAcceleratorId: ptrStringGa2Listener("ga2-test1234"),
			ListenerId:          ptrStringGa2Listener("lbl-abcdef"),
			Name:                ptrStringGa2Listener("test-listener"),
			Description:         ptrStringGa2Listener("test description"),
			Protocol:            ptrStringGa2Listener("TCP"),
			ListenerType:        ptrStringGa2Listener("INTELLIGENT"),
			PortRanges: &ga2v20250115.PortRanges{
				FromPort: ptrUint64Ga2Listener(80),
				ToPort:   ptrUint64Ga2Listener(80),
			},
			IdleTimeout:          ptrUint64Ga2Listener(900),
			ClientAffinity:       ptrStringGa2Listener("SOURCE_IP"),
			ClientAffinityTime:   ptrUint64Ga2Listener(300),
			GetRealIpType:        ptrStringGa2Listener("TOA"),
			RequestTimeout:       ptrUint64Ga2Listener(60),
			XForwardedForRealIp:  ptrBoolGa2Listener(true),
			CertificationType:    ptrStringGa2Listener("UNIDIRECTIONAL"),
			CipherPolicyId:       ptrStringGa2Listener("cipher-001"),
			ServerCertificates:   []*string{ptrStringGa2Listener("cert-001")},
			ClientCaCertificates: []*string{ptrStringGa2Listener("ca-001")},
		}, nil
	})

	meta := newMockMetaGa2Listener()
	res := svcga2.ResourceTencentCloudGa2Listener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test1234",
		"name":                  "test-listener",
	})
	d.SetId("ga2-test1234#lbl-abcdef")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, "lbl-abcdef", d.Get("listener_id").(string))
	assert.Equal(t, "test-listener", d.Get("name").(string))
	assert.Equal(t, "test description", d.Get("description").(string))
	assert.Equal(t, "TCP", d.Get("protocol").(string))
	assert.Equal(t, "INTELLIGENT", d.Get("listener_type").(string))
	assert.Equal(t, 900, d.Get("idle_timeout").(int))
	assert.Equal(t, "SOURCE_IP", d.Get("client_affinity").(string))
	assert.Equal(t, 300, d.Get("client_affinity_time").(int))
	assert.Equal(t, "TOA", d.Get("get_real_ip_type").(string))
	assert.Equal(t, 60, d.Get("request_timeout").(int))
	assert.Equal(t, true, d.Get("x_forwarded_for_real_ip").(bool))
	assert.Equal(t, "UNIDIRECTIONAL", d.Get("certification_type").(string))
	assert.Equal(t, "cipher-001", d.Get("cipher_policy_id").(string))
}

func TestGa2Listener_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "DescribeGa2ListenerById", func(_ context.Context, gaId string, listenerId string) (*ga2v20250115.ListenerSet, error) {
		return nil, nil
	})

	meta := newMockMetaGa2Listener()
	res := svcga2.ResourceTencentCloudGa2Listener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test1234",
		"name":                  "test-listener",
	})
	d.SetId("ga2-test1234#lbl-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestGa2Listener_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2Listener().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "ModifyListenerWithContext", func(_ context.Context, req *ga2v20250115.ModifyListenerRequest) (*ga2v20250115.ModifyListenerResponse, error) {
		assert.Equal(t, "ga2-test1234", *req.GlobalAcceleratorId)
		assert.Equal(t, "lbl-abcdef", *req.ListenerId)
		resp := ga2v20250115.NewModifyListenerResponse()
		resp.Response = &ga2v20250115.ModifyListenerResponseParams{
			TaskId:    ptrStringGa2Listener("task-modify-001"),
			RequestId: ptrStringGa2Listener("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "WaitForGa2TaskFinish", func(_ context.Context, taskId string, _ interface{}) error {
		assert.Equal(t, "task-modify-001", taskId)
		return nil
	})

	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "DescribeGa2ListenerById", func(_ context.Context, gaId string, listenerId string) (*ga2v20250115.ListenerSet, error) {
		return &ga2v20250115.ListenerSet{
			GlobalAcceleratorId: ptrStringGa2Listener("ga2-test1234"),
			ListenerId:          ptrStringGa2Listener("lbl-abcdef"),
			Name:                ptrStringGa2Listener("updated-listener"),
			Description:         ptrStringGa2Listener("updated description"),
			Protocol:            ptrStringGa2Listener("TCP"),
			ListenerType:        ptrStringGa2Listener("INTELLIGENT"),
			PortRanges: &ga2v20250115.PortRanges{
				FromPort: ptrUint64Ga2Listener(80),
				ToPort:   ptrUint64Ga2Listener(80),
			},
			IdleTimeout:    ptrUint64Ga2Listener(1200),
			ClientAffinity: ptrStringGa2Listener("SOURCE_IP"),
			GetRealIpType:  ptrStringGa2Listener("TOA"),
		}, nil
	})

	meta := newMockMetaGa2Listener()
	res := svcga2.ResourceTencentCloudGa2Listener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test1234",
		"name":                  "updated-listener",
		"description":           "updated description",
		"protocol":              "TCP",
		"listener_type":         "INTELLIGENT",
		"port_ranges": []interface{}{
			map[string]interface{}{
				"from_port": 80,
				"to_port":   80,
			},
		},
		"idle_timeout":     1200,
		"client_affinity":  "SOURCE_IP",
		"get_real_ip_type": "TOA",
	})
	d.SetId("ga2-test1234#lbl-abcdef")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "updated-listener", d.Get("name").(string))
}

func TestGa2Listener_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2Listener().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DeleteListenerWithContext", func(_ context.Context, req *ga2v20250115.DeleteListenerRequest) (*ga2v20250115.DeleteListenerResponse, error) {
		assert.Equal(t, "ga2-test1234", *req.GlobalAcceleratorId)
		assert.Equal(t, "lbl-abcdef", *req.ListenerId)
		resp := ga2v20250115.NewDeleteListenerResponse()
		resp.Response = &ga2v20250115.DeleteListenerResponseParams{
			TaskId:    ptrStringGa2Listener("task-delete-001"),
			RequestId: ptrStringGa2Listener("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "WaitForGa2TaskFinish", func(_ context.Context, taskId string, _ interface{}) error {
		assert.Equal(t, "task-delete-001", taskId)
		return nil
	})

	meta := newMockMetaGa2Listener()
	res := svcga2.ResourceTencentCloudGa2Listener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test1234",
		"name":                  "test-listener",
	})
	d.SetId("ga2-test1234#lbl-abcdef")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestGa2Listener_Schema(t *testing.T) {
	res := svcga2.ResourceTencentCloudGa2Listener()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "listener_id")
	listenerIdSchema := res.Schema["listener_id"]
	assert.Equal(t, schema.TypeString, listenerIdSchema.Type)
	assert.True(t, listenerIdSchema.Computed)
	assert.False(t, listenerIdSchema.Optional)
	assert.False(t, listenerIdSchema.Required)

	assert.Contains(t, res.Schema, "global_accelerator_id")
	gaIdSchema := res.Schema["global_accelerator_id"]
	assert.True(t, gaIdSchema.ForceNew)
	assert.True(t, gaIdSchema.Required)

	assert.Contains(t, res.Schema, "port_ranges")
	portRangesSchema := res.Schema["port_ranges"]
	assert.True(t, portRangesSchema.ForceNew)
	assert.True(t, portRangesSchema.Required)
}
