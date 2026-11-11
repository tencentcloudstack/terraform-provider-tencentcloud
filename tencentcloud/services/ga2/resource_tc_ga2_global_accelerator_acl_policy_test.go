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

type mockMetaGa2AclPolicy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaGa2AclPolicy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaGa2AclPolicy{}

func newMockMetaGa2AclPolicy() *mockMetaGa2AclPolicy {
	return &mockMetaGa2AclPolicy{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/ga2/ -run "TestGa2GlobalAcceleratorAclPolicy" -v -count=1 -gcflags="all=-l"

// TestGa2GlobalAcceleratorAclPolicy_Create tests the Create operation
func TestGa2GlobalAcceleratorAclPolicy_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateGlobalAcceleratorAclPolicyWithContext", func(ctx context.Context, request *ga2v20250115.CreateGlobalAcceleratorAclPolicyRequest) (*ga2v20250115.CreateGlobalAcceleratorAclPolicyResponse, error) {
		assert.Equal(t, "ga-test123", *request.GlobalAcceleratorId)
		assert.Equal(t, "ACCEPT", *request.DefaultAction)

		resp := ga2v20250115.NewCreateGlobalAcceleratorAclPolicyResponse()
		resp.Response = &ga2v20250115.CreateGlobalAcceleratorAclPolicyResponseParams{
			TaskId:                       ptrStringGa2("task-001"),
			GlobalAcceleratorAclPolicyId: ptrStringGa2("gapolicy-test789"),
			RequestId:                    ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeGlobalAcceleratorAclPoliciesWithContext for the Read call after Create
	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclPoliciesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclPoliciesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponse, error) {
		assert.Equal(t, "ga-test123", *request.GlobalAcceleratorId)
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclPoliciesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponseParams{
			GlobalAcceleratorAclPolicySet: []*ga2v20250115.GlobalAcceleratorAclPolicies{
				{
					GlobalAcceleratorAclPolicyId: ptrStringGa2("gapolicy-test789"),
					DefaultAction:                ptrStringGa2("ACCEPT"),
					Status:                       ptrStringGa2("OPEN"),
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

	meta := newMockMetaGa2AclPolicy()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test123",
		"default_action":        "ACCEPT",
		"status":                "OPEN",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-test123#gapolicy-test789", d.Id())
	assert.Equal(t, "OPEN", d.Get("status").(string))
	assert.Equal(t, "gapolicy-test789", d.Get("global_accelerator_acl_policy_id").(string))
	assert.Equal(t, "task-001", d.Get("task_id").(string))
}

// TestGa2GlobalAcceleratorAclPolicy_Read tests the Read operation
func TestGa2GlobalAcceleratorAclPolicy_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclPoliciesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclPoliciesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclPoliciesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponseParams{
			GlobalAcceleratorAclPolicySet: []*ga2v20250115.GlobalAcceleratorAclPolicies{
				{
					GlobalAcceleratorAclPolicyId: ptrStringGa2("gapolicy-read789"),
					DefaultAction:                ptrStringGa2("DROP"),
					Status:                       ptrStringGa2("CLOSE"),
				},
			},
			TotalCount: ptrUint64Ga2(1),
			RequestId:  ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclPolicy()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "",
		"default_action":        "",
		"status":                "",
	})
	d.SetId("ga-read123#gapolicy-read789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-read123#gapolicy-read789", d.Id())
	assert.Equal(t, "ga-read123", d.Get("global_accelerator_id").(string))
	assert.Equal(t, "DROP", d.Get("default_action").(string))
	assert.Equal(t, "CLOSE", d.Get("status").(string))
	assert.Equal(t, "gapolicy-read789", d.Get("global_accelerator_acl_policy_id").(string))
}

// TestGa2GlobalAcceleratorAclPolicy_Read_NotFound tests Read when the policy is not found
func TestGa2GlobalAcceleratorAclPolicy_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclPoliciesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclPoliciesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclPoliciesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponseParams{
			GlobalAcceleratorAclPolicySet: []*ga2v20250115.GlobalAcceleratorAclPolicies{},
			TotalCount:                    ptrUint64Ga2(0),
			RequestId:                     ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclPolicy()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "",
		"default_action":        "",
		"status":                "",
	})
	d.SetId("ga-notfound#gapolicy-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestGa2GlobalAcceleratorAclPolicy_Update tests the Update operation on a status change
func TestGa2GlobalAcceleratorAclPolicy_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "ModifyGlobalAcceleratorAclPolicyWithContext", func(ctx context.Context, request *ga2v20250115.ModifyGlobalAcceleratorAclPolicyRequest) (*ga2v20250115.ModifyGlobalAcceleratorAclPolicyResponse, error) {
		assert.Equal(t, "ga-upd123", *request.GlobalAcceleratorId)
		assert.Equal(t, "gapolicy-upd789", *request.GlobalAcceleratorAclPolicyId)
		assert.Equal(t, "CLOSE", *request.Status)

		resp := ga2v20250115.NewModifyGlobalAcceleratorAclPolicyResponse()
		resp.Response = &ga2v20250115.ModifyGlobalAcceleratorAclPolicyResponseParams{
			TaskId:    ptrStringGa2("task-002"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeGlobalAcceleratorAclPoliciesWithContext for the Read call after Update
	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclPoliciesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclPoliciesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclPoliciesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponseParams{
			GlobalAcceleratorAclPolicySet: []*ga2v20250115.GlobalAcceleratorAclPolicies{
				{
					GlobalAcceleratorAclPolicyId: ptrStringGa2("gapolicy-upd789"),
					DefaultAction:                ptrStringGa2("ACCEPT"),
					Status:                       ptrStringGa2("CLOSE"),
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

	meta := newMockMetaGa2AclPolicy()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-upd123",
		"default_action":        "ACCEPT",
		"status":                "CLOSE",
	})
	d.SetId("ga-upd123#gapolicy-upd789")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "CLOSE", d.Get("status").(string))
	assert.Equal(t, "task-002", d.Get("task_id").(string))
}

// TestGa2GlobalAcceleratorAclPolicy_Delete tests the Delete operation
func TestGa2GlobalAcceleratorAclPolicy_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DeleteGlobalAcceleratorAclPolicyWithContext", func(ctx context.Context, request *ga2v20250115.DeleteGlobalAcceleratorAclPolicyRequest) (*ga2v20250115.DeleteGlobalAcceleratorAclPolicyResponse, error) {
		assert.Equal(t, "ga-del123", *request.GlobalAcceleratorId)
		assert.Equal(t, "gapolicy-del789", *request.GlobalAcceleratorAclPolicyId)

		resp := ga2v20250115.NewDeleteGlobalAcceleratorAclPolicyResponse()
		resp.Response = &ga2v20250115.DeleteGlobalAcceleratorAclPolicyResponseParams{
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

	meta := newMockMetaGa2AclPolicy()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-del123",
		"default_action":        "ACCEPT",
		"status":                "OPEN",
	})
	d.SetId("ga-del123#gapolicy-del789")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-003", d.Get("task_id").(string))
}

// TestGa2GlobalAcceleratorAclPolicy_Import tests import by reading a resource set via composite ID
func TestGa2GlobalAcceleratorAclPolicy_Import(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclPolicy().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclPoliciesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclPoliciesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclPoliciesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclPoliciesResponseParams{
			GlobalAcceleratorAclPolicySet: []*ga2v20250115.GlobalAcceleratorAclPolicies{
				{
					GlobalAcceleratorAclPolicyId: ptrStringGa2("gapolicy-imp789"),
					DefaultAction:                ptrStringGa2("DROP"),
					Status:                       ptrStringGa2("OPEN"),
				},
			},
			TotalCount: ptrUint64Ga2(1),
			RequestId:  ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclPolicy()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "",
		"default_action":        "",
		"status":                "",
	})
	d.SetId("ga-imp123#gapolicy-imp789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-imp123#gapolicy-imp789", d.Id())
	assert.Equal(t, "ga-imp123", d.Get("global_accelerator_id").(string))
	assert.Equal(t, "DROP", d.Get("default_action").(string))
	assert.Equal(t, "gapolicy-imp789", d.Get("global_accelerator_acl_policy_id").(string))
}
