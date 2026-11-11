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

type mockMetaGa2AclRuleSet struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaGa2AclRuleSet) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaGa2AclRuleSet{}

func newMockMetaGa2AclRuleSet() *mockMetaGa2AclRuleSet {
	return &mockMetaGa2AclRuleSet{client: &connectivity.TencentCloudClient{}}
}

func ptrStringGa2AclRuleSet(s string) *string { return &s }
func ptrUint64Ga2AclRuleSet(v uint64) *uint64 { return &v }

// go test ./tencentcloud/services/ga2/ -run "TestGa2GlobalAcceleratorAclRuleSet" -v -count=1 -gcflags="all=-l"

// TestGa2GlobalAcceleratorAclRuleSet_Create tests the Create operation with multiple acl entries.
func TestGa2GlobalAcceleratorAclRuleSet_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateGlobalAcceleratorAclRuleWithContext", func(ctx context.Context, request *ga2v20250115.CreateGlobalAcceleratorAclRuleRequest) (*ga2v20250115.CreateGlobalAcceleratorAclRuleResponse, error) {
		assert.Equal(t, "ga-test123", *request.GlobalAcceleratorId)
		assert.Equal(t, "sp-test456", *request.GlobalAcceleratorAclPolicyId)
		assert.Equal(t, 2, len(request.AclEntries))

		resp := ga2v20250115.NewCreateGlobalAcceleratorAclRuleResponse()
		resp.Response = &ga2v20250115.CreateGlobalAcceleratorAclRuleResponseParams{
			TaskId:                      ptrStringGa2AclRuleSet("task-001"),
			GlobalAcceleratorAclRuleIds: []*string{ptrStringGa2AclRuleSet("sr-001"), ptrStringGa2AclRuleSet("sr-002")},
			RequestId:                   ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeGlobalAcceleratorAclRulesWithContext for the Read call after Create.
	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclRulesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclRulesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclRulesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclRulesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclRulesResponseParams{
			GlobalAcceleratorAclRuleSet: []*ga2v20250115.GlobalAcceleratorAclRuleSet{
				{
					GlobalAcceleratorAclRuleId: ptrStringGa2AclRuleSet("sr-001"),
					Protocol:                   ptrStringGa2AclRuleSet("TCP"),
					Port:                       ptrStringGa2AclRuleSet("80"),
					SourceCidrBlock:            ptrStringGa2AclRuleSet("10.0.0.0/24"),
					Policy:                     ptrStringGa2AclRuleSet("ACCEPT"),
					Description:                ptrStringGa2AclRuleSet("rule 1"),
				},
				{
					GlobalAcceleratorAclRuleId: ptrStringGa2AclRuleSet("sr-002"),
					Protocol:                   ptrStringGa2AclRuleSet("UDP"),
					Port:                       ptrStringGa2AclRuleSet("443"),
					SourceCidrBlock:            ptrStringGa2AclRuleSet("10.0.1.0/24"),
					Policy:                     ptrStringGa2AclRuleSet("DROP"),
					Description:                ptrStringGa2AclRuleSet("rule 2"),
				},
			},
			TotalCount: ptrUint64Ga2AclRuleSet(2),
			RequestId:  ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResultWithContext for async task polling (WaitForGa2TaskFinish returns nil).
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2AclRuleSet("SUCCESS"),
			RequestId: ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "ga-test123",
		"global_accelerator_acl_policy_id": "sp-test456",
		"acl_entries": []interface{}{
			map[string]interface{}{
				"protocol":          "TCP",
				"port":              "80",
				"source_cidr_block": "10.0.0.0/24",
				"policy":            "ACCEPT",
				"description":       "rule 1",
			},
			map[string]interface{}{
				"protocol":          "UDP",
				"port":              "443",
				"source_cidr_block": "10.0.1.0/24",
				"policy":            "DROP",
				"description":       "rule 2",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-test123#sp-test456", d.Id())
	assert.Equal(t, "ga-test123", d.Get("global_accelerator_id").(string))
	assert.Equal(t, "sp-test456", d.Get("global_accelerator_acl_policy_id").(string))
	assert.Equal(t, "task-001", d.Get("task_id").(string))

	entries := d.Get("acl_entries").([]interface{})
	assert.Equal(t, 2, len(entries))
}

// TestGa2GlobalAcceleratorAclRuleSet_Create_Empty tests Create with empty acl_entries (no API call).
func TestGa2GlobalAcceleratorAclRuleSet_Create_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	// CreateGlobalAcceleratorAclRuleWithContext must NOT be called; if it is, the test fails.
	createCalled := false
	patches.ApplyMethodFunc(ga2Client, "CreateGlobalAcceleratorAclRuleWithContext", func(ctx context.Context, request *ga2v20250115.CreateGlobalAcceleratorAclRuleRequest) (*ga2v20250115.CreateGlobalAcceleratorAclRuleResponse, error) {
		createCalled = true
		return nil, nil
	})

	// Mock DescribeGlobalAcceleratorAclRulesWithContext for the Read call after Create (empty result).
	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclRulesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclRulesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclRulesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclRulesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclRulesResponseParams{
			GlobalAcceleratorAclRuleSet: []*ga2v20250115.GlobalAcceleratorAclRuleSet{},
			TotalCount:                  ptrUint64Ga2AclRuleSet(0),
			RequestId:                   ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "ga-test123",
		"global_accelerator_acl_policy_id": "sp-test456",
		"acl_entries":                      []interface{}{},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-test123#sp-test456", d.Id())
	assert.False(t, createCalled, "CreateGlobalAcceleratorAclRule must not be called for empty acl_entries")

	entries := d.Get("acl_entries").([]interface{})
	assert.Equal(t, 0, len(entries))
}

// TestGa2GlobalAcceleratorAclRuleSet_Read tests the Read operation.
func TestGa2GlobalAcceleratorAclRuleSet_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclRulesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclRulesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclRulesResponse, error) {
		assert.Equal(t, "sp-read456", *request.GlobalAcceleratorAclPolicyId)
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclRulesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclRulesResponseParams{
			GlobalAcceleratorAclRuleSet: []*ga2v20250115.GlobalAcceleratorAclRuleSet{
				{
					GlobalAcceleratorAclRuleId: ptrStringGa2AclRuleSet("sr-002"),
					Protocol:                   ptrStringGa2AclRuleSet("UDP"),
					Port:                       ptrStringGa2AclRuleSet("443"),
					SourceCidrBlock:            ptrStringGa2AclRuleSet("10.0.1.0/24"),
					Policy:                     ptrStringGa2AclRuleSet("DROP"),
					Description:                ptrStringGa2AclRuleSet("rule 2"),
				},
				{
					GlobalAcceleratorAclRuleId: ptrStringGa2AclRuleSet("sr-001"),
					Protocol:                   ptrStringGa2AclRuleSet("TCP"),
					Port:                       ptrStringGa2AclRuleSet("80"),
					SourceCidrBlock:            ptrStringGa2AclRuleSet("10.0.0.0/24"),
					Policy:                     ptrStringGa2AclRuleSet("ACCEPT"),
					Description:                ptrStringGa2AclRuleSet("rule 1"),
				},
			},
			TotalCount: ptrUint64Ga2AclRuleSet(2),
			RequestId:  ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "",
		"global_accelerator_acl_policy_id": "",
		"acl_entries":                      []interface{}{},
	})
	d.SetId("ga-read123#sp-read456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-read123#sp-read456", d.Id())
	assert.Equal(t, "ga-read123", d.Get("global_accelerator_id").(string))
	assert.Equal(t, "sp-read456", d.Get("global_accelerator_acl_policy_id").(string))

	// Entries should be sorted by rule id (sr-001 before sr-002).
	entries := d.Get("acl_entries").([]interface{})
	assert.Equal(t, 2, len(entries))
	first := entries[0].(map[string]interface{})
	assert.Equal(t, "sr-001", first["global_accelerator_acl_rule_id"])
}

// TestGa2GlobalAcceleratorAclRuleSet_Read_Empty tests Read when the policy has no rules.
func TestGa2GlobalAcceleratorAclRuleSet_Read_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclRulesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclRulesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclRulesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclRulesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclRulesResponseParams{
			GlobalAcceleratorAclRuleSet: []*ga2v20250115.GlobalAcceleratorAclRuleSet{},
			TotalCount:                  ptrUint64Ga2AclRuleSet(0),
			RequestId:                   ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "",
		"global_accelerator_acl_policy_id": "",
		"acl_entries":                      []interface{}{},
	})
	d.SetId("ga-read123#sp-read456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// Resource must remain in state (not removed) even with an empty rule set.
	assert.Equal(t, "ga-read123#sp-read456", d.Id())

	entries := d.Get("acl_entries").([]interface{})
	assert.Equal(t, 0, len(entries))
}

// TestGa2GlobalAcceleratorAclRuleSet_Update_Add tests the Update operation when adding new rules.
func TestGa2GlobalAcceleratorAclRuleSet_Update_Add(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateGlobalAcceleratorAclRuleWithContext", func(ctx context.Context, request *ga2v20250115.CreateGlobalAcceleratorAclRuleRequest) (*ga2v20250115.CreateGlobalAcceleratorAclRuleResponse, error) {
		assert.Equal(t, "ga-upd123", *request.GlobalAcceleratorId)
		assert.Equal(t, "sp-upd456", *request.GlobalAcceleratorAclPolicyId)
		// One new entry being created.
		assert.Equal(t, 1, len(request.AclEntries))

		resp := ga2v20250115.NewCreateGlobalAcceleratorAclRuleResponse()
		resp.Response = &ga2v20250115.CreateGlobalAcceleratorAclRuleResponseParams{
			TaskId:                      ptrStringGa2AclRuleSet("task-add"),
			GlobalAcceleratorAclRuleIds: []*string{ptrStringGa2AclRuleSet("sr-new")},
			RequestId:                   ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeGlobalAcceleratorAclRulesWithContext for the Read call after Update.
	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclRulesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclRulesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclRulesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclRulesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclRulesResponseParams{
			GlobalAcceleratorAclRuleSet: []*ga2v20250115.GlobalAcceleratorAclRuleSet{
				{
					GlobalAcceleratorAclRuleId: ptrStringGa2AclRuleSet("sr-new"),
					Protocol:                   ptrStringGa2AclRuleSet("TCP"),
					Port:                       ptrStringGa2AclRuleSet("8080"),
					SourceCidrBlock:            ptrStringGa2AclRuleSet("10.0.2.0/24"),
					Policy:                     ptrStringGa2AclRuleSet("ACCEPT"),
					Description:                ptrStringGa2AclRuleSet("added rule"),
				},
			},
			TotalCount: ptrUint64Ga2AclRuleSet(1),
			RequestId:  ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2AclRuleSet("SUCCESS"),
			RequestId: ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "ga-upd123",
		"global_accelerator_acl_policy_id": "sp-upd456",
		"acl_entries": []interface{}{
			map[string]interface{}{
				"protocol":          "TCP",
				"port":              "8080",
				"source_cidr_block": "10.0.2.0/24",
				"policy":            "ACCEPT",
				"description":       "added rule",
			},
		},
	})
	d.SetId("ga-upd123#sp-upd456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-add", d.Get("task_id").(string))
}

// TestGa2GlobalAcceleratorAclRuleSet_Delete tests the Delete operation.
func TestGa2GlobalAcceleratorAclRuleSet_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DeleteGlobalAcceleratorAclRuleWithContext", func(ctx context.Context, request *ga2v20250115.DeleteGlobalAcceleratorAclRuleRequest) (*ga2v20250115.DeleteGlobalAcceleratorAclRuleResponse, error) {
		assert.Equal(t, "ga-del123", *request.GlobalAcceleratorId)
		assert.Equal(t, "sp-del456", *request.GlobalAcceleratorAclPolicyId)
		assert.Equal(t, 1, len(request.GlobalAcceleratorAclRuleIds))
		assert.Equal(t, "sr-del001", *request.GlobalAcceleratorAclRuleIds[0])

		resp := ga2v20250115.NewDeleteGlobalAcceleratorAclRuleResponse()
		resp.Response = &ga2v20250115.DeleteGlobalAcceleratorAclRuleResponseParams{
			TaskId:    ptrStringGa2AclRuleSet("task-del"),
			RequestId: ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2AclRuleSet("SUCCESS"),
			RequestId: ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "ga-del123",
		"global_accelerator_acl_policy_id": "sp-del456",
		"acl_entries": []interface{}{
			map[string]interface{}{
				"protocol":                       "TCP",
				"port":                           "80",
				"source_cidr_block":              "10.0.0.0/24",
				"policy":                         "ACCEPT",
				"description":                    "to delete",
				"global_accelerator_acl_rule_id": "sr-del001",
			},
		},
	})
	d.SetId("ga-del123#sp-del456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestGa2GlobalAcceleratorAclRuleSet_Delete_Empty tests Delete with empty acl_entries (no API call).
func TestGa2GlobalAcceleratorAclRuleSet_Delete_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	deleteCalled := false
	patches.ApplyMethodFunc(ga2Client, "DeleteGlobalAcceleratorAclRuleWithContext", func(ctx context.Context, request *ga2v20250115.DeleteGlobalAcceleratorAclRuleRequest) (*ga2v20250115.DeleteGlobalAcceleratorAclRuleResponse, error) {
		deleteCalled = true
		return nil, nil
	})

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "ga-del123",
		"global_accelerator_acl_policy_id": "sp-del456",
		"acl_entries":                      []interface{}{},
	})
	d.SetId("ga-del123#sp-del456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.False(t, deleteCalled, "DeleteGlobalAcceleratorAclRule must not be called for empty acl_entries")
}

// TestGa2GlobalAcceleratorAclRuleSet_ParseId_Invalid tests parseGa2GlobalAcceleratorAclRuleSetId error cases
// indirectly via Read, which returns the parse error.
func TestGa2GlobalAcceleratorAclRuleSet_ParseId_Invalid(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()

	// Too few parts.
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "",
		"global_accelerator_acl_policy_id": "",
		"acl_entries":                      []interface{}{},
	})
	d.SetId("ga-only")
	err := res.Read(d, meta)
	assert.Error(t, err)

	// Empty component.
	d2 := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "",
		"global_accelerator_acl_policy_id": "",
		"acl_entries":                      []interface{}{},
	})
	d2.SetId("ga-123#")
	err = res.Read(d2, meta)
	assert.Error(t, err)

	// Too many parts.
	d3 := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "",
		"global_accelerator_acl_policy_id": "",
		"acl_entries":                      []interface{}{},
	})
	d3.SetId("ga-123#sp-456#sr-789")
	err = res.Read(d3, meta)
	assert.Error(t, err)
}

// TestGa2GlobalAcceleratorAclRuleSet_Update_Remove tests the Update operation when removing rules.
// It patches (*schema.ResourceData).GetChange so the diff sees an old rule that is absent in the new
// list, exercising the batch-delete path.
func TestGa2GlobalAcceleratorAclRuleSet_Update_Remove(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	// old has one rule (sr-001); new is empty -> sr-001 is removed.
	patches.ApplyMethodFunc(&schema.ResourceData{}, "GetChange", func(key string) (interface{}, interface{}) {
		if key == "acl_entries" {
			old := []interface{}{
				map[string]interface{}{
					"protocol":                       "TCP",
					"port":                           "80",
					"source_cidr_block":              "10.0.0.0/24",
					"policy":                         "ACCEPT",
					"description":                    "rule 1",
					"global_accelerator_acl_rule_id": "sr-001",
				},
			}
			newVal := []interface{}{}
			return old, newVal
		}
		return nil, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DeleteGlobalAcceleratorAclRuleWithContext", func(ctx context.Context, request *ga2v20250115.DeleteGlobalAcceleratorAclRuleRequest) (*ga2v20250115.DeleteGlobalAcceleratorAclRuleResponse, error) {
		assert.Equal(t, "ga-upd123", *request.GlobalAcceleratorId)
		assert.Equal(t, "sp-upd456", *request.GlobalAcceleratorAclPolicyId)
		assert.Equal(t, 1, len(request.GlobalAcceleratorAclRuleIds))
		assert.Equal(t, "sr-001", *request.GlobalAcceleratorAclRuleIds[0])

		resp := ga2v20250115.NewDeleteGlobalAcceleratorAclRuleResponse()
		resp.Response = &ga2v20250115.DeleteGlobalAcceleratorAclRuleResponseParams{
			TaskId:    ptrStringGa2AclRuleSet("task-rem"),
			RequestId: ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclRulesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclRulesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclRulesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclRulesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclRulesResponseParams{
			GlobalAcceleratorAclRuleSet: []*ga2v20250115.GlobalAcceleratorAclRuleSet{},
			TotalCount:                  ptrUint64Ga2AclRuleSet(0),
			RequestId:                   ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2AclRuleSet("SUCCESS"),
			RequestId: ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "ga-upd123",
		"global_accelerator_acl_policy_id": "sp-upd456",
		"acl_entries":                      []interface{}{},
	})
	d.SetId("ga-upd123#sp-upd456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-rem", d.Get("task_id").(string))
}

// TestGa2GlobalAcceleratorAclRuleSet_Update_Modify tests the Update operation when modifying an existing rule.
// It patches (*schema.ResourceData).GetChange so the diff sees a changed description for an existing rule id,
// exercising the per-rule Modify path.
func TestGa2GlobalAcceleratorAclRuleSet_Update_Modify(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2AclRuleSet().client, "UseGa2V20250115Client", ga2Client)

	// old and new share the same rule id (sr-001) but the description differs -> modify path.
	patches.ApplyMethodFunc(&schema.ResourceData{}, "GetChange", func(key string) (interface{}, interface{}) {
		if key == "acl_entries" {
			old := []interface{}{
				map[string]interface{}{
					"protocol":                       "TCP",
					"port":                           "80",
					"source_cidr_block":              "10.0.0.0/24",
					"policy":                         "ACCEPT",
					"description":                    "rule 1",
					"global_accelerator_acl_rule_id": "sr-001",
				},
			}
			newVal := []interface{}{
				map[string]interface{}{
					"protocol":                       "TCP",
					"port":                           "80",
					"source_cidr_block":              "10.0.0.0/24",
					"policy":                         "ACCEPT",
					"description":                    "rule 1 updated",
					"global_accelerator_acl_rule_id": "sr-001",
				},
			}
			return old, newVal
		}
		return nil, nil
	})

	patches.ApplyMethodFunc(ga2Client, "ModifyGlobalAcceleratorAclRuleWithContext", func(ctx context.Context, request *ga2v20250115.ModifyGlobalAcceleratorAclRuleRequest) (*ga2v20250115.ModifyGlobalAcceleratorAclRuleResponse, error) {
		assert.Equal(t, "ga-upd123", *request.GlobalAcceleratorId)
		assert.Equal(t, "sp-upd456", *request.GlobalAcceleratorAclPolicyId)
		assert.Equal(t, "sr-001", *request.GlobalAcceleratorAclRuleId)
		assert.Equal(t, "rule 1 updated", *request.Description)

		resp := ga2v20250115.NewModifyGlobalAcceleratorAclRuleResponse()
		resp.Response = &ga2v20250115.ModifyGlobalAcceleratorAclRuleResponseParams{
			TaskId:    ptrStringGa2AclRuleSet("task-mod"),
			RequestId: ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorAclRulesWithContext", func(ctx context.Context, request *ga2v20250115.DescribeGlobalAcceleratorAclRulesRequest) (*ga2v20250115.DescribeGlobalAcceleratorAclRulesResponse, error) {
		resp := ga2v20250115.NewDescribeGlobalAcceleratorAclRulesResponse()
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorAclRulesResponseParams{
			GlobalAcceleratorAclRuleSet: []*ga2v20250115.GlobalAcceleratorAclRuleSet{
				{
					GlobalAcceleratorAclRuleId: ptrStringGa2AclRuleSet("sr-001"),
					Protocol:                   ptrStringGa2AclRuleSet("TCP"),
					Port:                       ptrStringGa2AclRuleSet("80"),
					SourceCidrBlock:            ptrStringGa2AclRuleSet("10.0.0.0/24"),
					Policy:                     ptrStringGa2AclRuleSet("ACCEPT"),
					Description:                ptrStringGa2AclRuleSet("rule 1 updated"),
				},
			},
			TotalCount: ptrUint64Ga2AclRuleSet(1),
			RequestId:  ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(ctx context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2AclRuleSet("SUCCESS"),
			RequestId: ptrStringGa2AclRuleSet("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2AclRuleSet()
	res := ga2.ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id":            "ga-upd123",
		"global_accelerator_acl_policy_id": "sp-upd456",
		"acl_entries": []interface{}{
			map[string]interface{}{
				"protocol":                       "TCP",
				"port":                           "80",
				"source_cidr_block":              "10.0.0.0/24",
				"policy":                         "ACCEPT",
				"description":                    "rule 1 updated",
				"global_accelerator_acl_rule_id": "sr-001",
			},
		},
	})
	d.SetId("ga-upd123#sp-upd456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-mod", d.Get("task_id").(string))
}
