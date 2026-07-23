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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
)

type mockMetaGa2ForwardingPolicy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaGa2ForwardingPolicy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaGa2ForwardingPolicy{}

func newMockMetaGa2ForwardingPolicy() *mockMetaGa2ForwardingPolicy {
	return &mockMetaGa2ForwardingPolicy{client: &connectivity.TencentCloudClient{}}
}

func ptrStringGa2(s string) *string {
	return &s
}

func ptrBoolGa2(v bool) *bool {
	return &v
}

func ptrUint64Ga2(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/ga2/ -run "TestGa2ForwardingPolicy" -v -count=1 -gcflags="all=-l"

// TestGa2ForwardingPolicy_Create tests the Create operation
func TestGa2ForwardingPolicy_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2ForwardingPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateForwardingPolicyWithContext", func(ctx context.Context, request *ga2v20250115.CreateForwardingPolicyRequest) (*ga2v20250115.CreateForwardingPolicyResponse, error) {
		assert.Equal(t, "ga-test123", *request.GlobalAcceleratorId)
		assert.Equal(t, "lsr-test456", *request.ListenerId)
		assert.Equal(t, "example.com", *request.Host)

		resp := ga2v20250115.NewCreateForwardingPolicyResponse()
		resp.Response = &ga2v20250115.CreateForwardingPolicyResponseParams{
			TaskId:             ptrStringGa2("task-001"),
			ForwardingPolicyId: ptrStringGa2("fp-test789"),
			RequestId:          ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeForwardingPolicyWithContext for the Read call after Create
	patches.ApplyMethodFunc(ga2Client, "DescribeForwardingPolicyWithContext", func(ctx context.Context, request *ga2v20250115.DescribeForwardingPolicyRequest) (*ga2v20250115.DescribeForwardingPolicyResponse, error) {
		resp := ga2v20250115.NewDescribeForwardingPolicyResponse()
		resp.Response = &ga2v20250115.DescribeForwardingPolicyResponseParams{
			ForwardingPolicySet: []*ga2v20250115.ForwardingPolicySet{
				{
					GlobalAcceleratorId: ptrStringGa2("ga-test123"),
					ListenerId:          ptrStringGa2("lsr-test456"),
					ForwardingPolicyId:  ptrStringGa2("fp-test789"),
					Host:                ptrStringGa2("example.com"),
					DefaultHostFlag:     ptrBoolGa2(false),
				},
			},
			TotalCount: ptrUint64Ga2(1),
			RequestId:  ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResultWithContext for async task polling
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2ForwardingPolicy()
	res := ga2.ResourceTencentCloudGa2ForwardingPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test123",
		"listener_id":           "lsr-test456",
		"host":                  "example.com",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-test123#lsr-test456#fp-test789", d.Id())
	assert.Equal(t, "example.com", d.Get("host").(string))
	assert.Equal(t, false, d.Get("default_host_flag").(bool))
}

// TestGa2ForwardingPolicy_Read tests the Read operation
func TestGa2ForwardingPolicy_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2ForwardingPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeForwardingPolicyWithContext", func(ctx context.Context, request *ga2v20250115.DescribeForwardingPolicyRequest) (*ga2v20250115.DescribeForwardingPolicyResponse, error) {
		resp := ga2v20250115.NewDescribeForwardingPolicyResponse()
		resp.Response = &ga2v20250115.DescribeForwardingPolicyResponseParams{
			ForwardingPolicySet: []*ga2v20250115.ForwardingPolicySet{
				{
					GlobalAcceleratorId: ptrStringGa2("ga-read123"),
					ListenerId:          ptrStringGa2("lsr-read456"),
					ForwardingPolicyId:  ptrStringGa2("fp-read789"),
					Host:                ptrStringGa2("read.example.com"),
					DefaultHostFlag:     ptrBoolGa2(true),
				},
			},
			TotalCount: ptrUint64Ga2(1),
			RequestId:  ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2ForwardingPolicy()
	res := ga2.ResourceTencentCloudGa2ForwardingPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "",
		"listener_id":           "",
		"host":                  "",
	})
	d.SetId("ga-read123#lsr-read456#fp-read789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-read123#lsr-read456#fp-read789", d.Id())
	assert.Equal(t, "read.example.com", d.Get("host").(string))
	assert.Equal(t, true, d.Get("default_host_flag").(bool))
}

// TestGa2ForwardingPolicy_Read_NotFound tests Read when policy is not found
func TestGa2ForwardingPolicy_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2ForwardingPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeForwardingPolicyWithContext", func(ctx context.Context, request *ga2v20250115.DescribeForwardingPolicyRequest) (*ga2v20250115.DescribeForwardingPolicyResponse, error) {
		resp := ga2v20250115.NewDescribeForwardingPolicyResponse()
		resp.Response = &ga2v20250115.DescribeForwardingPolicyResponseParams{
			ForwardingPolicySet: []*ga2v20250115.ForwardingPolicySet{},
			TotalCount:          ptrUint64Ga2(0),
			RequestId:           ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2ForwardingPolicy()
	res := ga2.ResourceTencentCloudGa2ForwardingPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "",
		"listener_id":           "",
		"host":                  "",
	})
	d.SetId("ga-notfound#lsr-notfound#fp-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestGa2ForwardingPolicy_Update tests the Update operation
func TestGa2ForwardingPolicy_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2ForwardingPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "ModifyForwardingPolicyWithContext", func(ctx context.Context, request *ga2v20250115.ModifyForwardingPolicyRequest) (*ga2v20250115.ModifyForwardingPolicyResponse, error) {
		assert.Equal(t, "ga-upd123", *request.GlobalAcceleratorId)
		assert.Equal(t, "lsr-upd456", *request.ListenerId)
		assert.Equal(t, "fp-upd789", *request.ForwardingPolicyId)
		assert.Equal(t, "updated.example.com", *request.Host)

		resp := ga2v20250115.NewModifyForwardingPolicyResponse()
		resp.Response = &ga2v20250115.ModifyForwardingPolicyResponseParams{
			TaskId:    ptrStringGa2("task-002"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeForwardingPolicyWithContext for the Read call after Update
	patches.ApplyMethodFunc(ga2Client, "DescribeForwardingPolicyWithContext", func(ctx context.Context, request *ga2v20250115.DescribeForwardingPolicyRequest) (*ga2v20250115.DescribeForwardingPolicyResponse, error) {
		resp := ga2v20250115.NewDescribeForwardingPolicyResponse()
		resp.Response = &ga2v20250115.DescribeForwardingPolicyResponseParams{
			ForwardingPolicySet: []*ga2v20250115.ForwardingPolicySet{
				{
					GlobalAcceleratorId: ptrStringGa2("ga-upd123"),
					ListenerId:          ptrStringGa2("lsr-upd456"),
					ForwardingPolicyId:  ptrStringGa2("fp-upd789"),
					Host:                ptrStringGa2("updated.example.com"),
					DefaultHostFlag:     ptrBoolGa2(false),
				},
			},
			TotalCount: ptrUint64Ga2(1),
			RequestId:  ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResultWithContext for async task polling
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2ForwardingPolicy()
	res := ga2.ResourceTencentCloudGa2ForwardingPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-upd123",
		"listener_id":           "lsr-upd456",
		"host":                  "updated.example.com",
	})
	d.SetId("ga-upd123#lsr-upd456#fp-upd789")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "updated.example.com", d.Get("host").(string))
}

// TestGa2ForwardingPolicy_Delete tests the Delete operation
func TestGa2ForwardingPolicy_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2ForwardingPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DeleteForwardingPolicyWithContext", func(ctx context.Context, request *ga2v20250115.DeleteForwardingPolicyRequest) (*ga2v20250115.DeleteForwardingPolicyResponse, error) {
		assert.Equal(t, "ga-del123", *request.GlobalAcceleratorId)
		assert.Equal(t, "lsr-del456", *request.ListenerId)
		assert.Equal(t, "fp-del789", *request.ForwardingPolicyId)

		resp := ga2v20250115.NewDeleteForwardingPolicyResponse()
		resp.Response = &ga2v20250115.DeleteForwardingPolicyResponseParams{
			TaskId:    ptrStringGa2("task-003"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResultWithContext for async task polling
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2ForwardingPolicy()
	res := ga2.ResourceTencentCloudGa2ForwardingPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-del123",
		"listener_id":           "lsr-del456",
		"host":                  "del.example.com",
	})
	d.SetId("ga-del123#lsr-del456#fp-del789")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestGa2ForwardingPolicy_Import tests import by reading a resource set via composite ID
func TestGa2ForwardingPolicy_Import(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2ForwardingPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeForwardingPolicyWithContext", func(ctx context.Context, request *ga2v20250115.DescribeForwardingPolicyRequest) (*ga2v20250115.DescribeForwardingPolicyResponse, error) {
		resp := ga2v20250115.NewDescribeForwardingPolicyResponse()
		resp.Response = &ga2v20250115.DescribeForwardingPolicyResponseParams{
			ForwardingPolicySet: []*ga2v20250115.ForwardingPolicySet{
				{
					GlobalAcceleratorId: ptrStringGa2("ga-imp123"),
					ListenerId:          ptrStringGa2("lsr-imp456"),
					ForwardingPolicyId:  ptrStringGa2("fp-imp789"),
					Host:                ptrStringGa2("import.example.com"),
					DefaultHostFlag:     ptrBoolGa2(false),
				},
			},
			TotalCount: ptrUint64Ga2(1),
			RequestId:  ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2ForwardingPolicy()
	res := ga2.ResourceTencentCloudGa2ForwardingPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "",
		"listener_id":           "",
		"host":                  "",
	})
	d.SetId("ga-imp123#lsr-imp456#fp-imp789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-imp123#lsr-imp456#fp-imp789", d.Id())
	assert.Equal(t, "ga-imp123", d.Get("global_accelerator_id").(string))
	assert.Equal(t, "lsr-imp456", d.Get("listener_id").(string))
	assert.Equal(t, "fp-imp789", d.Get("forwarding_policy_id").(string))
}
