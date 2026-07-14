package dlc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcdlc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"
)

type mockMetaDlcAttachUserPolicyAttachment struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDlcAttachUserPolicyAttachment) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDlcAttachUserPolicyAttachment{}

func newMockMetaDlcAttachUserPolicyAttachment() *mockMetaDlcAttachUserPolicyAttachment {
	return &mockMetaDlcAttachUserPolicyAttachment{client: &connectivity.TencentCloudClient{}}
}

func ptrStrDlcAupa(s string) *string {
	return &s
}

func ptrBoolDlcAupa(b bool) *bool {
	return &b
}

func ptrInt64DlcAupa(i int64) *int64 {
	return &i
}

func TestDlcAttachUserPolicyAttachment_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.AttachUserPolicyRequest
	patches.ApplyMethodFunc(dlcClient, "AttachUserPolicyWithContext", func(_ context.Context, request *dlc.AttachUserPolicyRequest) (*dlc.AttachUserPolicyResponse, error) {
		capturedRequest = request
		resp := dlc.NewAttachUserPolicyResponse()
		resp.Response = &dlc.AttachUserPolicyResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-create"),
			PolicySet: []*dlc.Policy{
				{
					Database:   ptrStrDlcAupa("tf_example_db"),
					Catalog:    ptrStrDlcAupa("DataLakeCatalog"),
					Table:      ptrStrDlcAupa("tf_example_table"),
					Operation:  ptrStrDlcAupa("SELECT"),
					PolicyType: ptrStrDlcAupa("TABLE"),
					PolicyId:   ptrStrDlcAupa("policy-id-001"),
					Source:     ptrStrDlcAupa("USER"),
					Mode:       ptrStrDlcAupa("COMMON"),
				},
			},
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeUserInfoWithContext", func(_ context.Context, request *dlc.DescribeUserInfoRequest) (*dlc.DescribeUserInfoResponse, error) {
		resp := dlc.NewDescribeUserInfoResponse()
		totalCount := int64(1)
		resp.Response = &dlc.DescribeUserInfoResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-read"),
			UserInfo: &dlc.UserDetailInfo{
				UserId:      ptrStrDlcAupa("100032676511"),
				Type:        ptrStrDlcAupa("DataAuth"),
				AccountType: ptrStrDlcAupa("TencentAccount"),
				DataPolicyInfo: &dlc.Policys{
					TotalCount: &totalCount,
					PolicySet: []*dlc.Policy{
						{
							Database:   ptrStrDlcAupa("tf_example_db"),
							Catalog:    ptrStrDlcAupa("DataLakeCatalog"),
							Table:      ptrStrDlcAupa("tf_example_table"),
							Operation:  ptrStrDlcAupa("SELECT"),
							PolicyType: ptrStrDlcAupa("TABLE"),
							PolicyId:   ptrStrDlcAupa("policy-id-001"),
							Source:     ptrStrDlcAupa("USER"),
							Mode:       ptrStrDlcAupa("COMMON"),
						},
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id":      "100032676511",
		"account_type": "TencentAccount",
		"policy_set": []interface{}{
			map[string]interface{}{
				"database":    "tf_example_db",
				"catalog":     "DataLakeCatalog",
				"table":       "tf_example_table",
				"operation":   "SELECT",
				"policy_type": "TABLE",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "100032676511#policy-id-001", d.Id())

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "100032676511", *capturedRequest.UserId)
	assert.Equal(t, "TencentAccount", *capturedRequest.AccountType)
	assert.Len(t, capturedRequest.PolicySet, 1)
	assert.Equal(t, "tf_example_db", *capturedRequest.PolicySet[0].Database)
	assert.Equal(t, "SELECT", *capturedRequest.PolicySet[0].Operation)

	policySet := d.Get("policy_set").([]interface{})
	assert.Len(t, policySet, 1)
	policyMap := policySet[0].(map[string]interface{})
	assert.Equal(t, "policy-id-001", policyMap["policy_id"])
}

func TestDlcAttachUserPolicyAttachment_Create_MultiplePolicies_Error(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

	called := false
	patches.ApplyMethodFunc(dlcClient, "AttachUserPolicyWithContext", func(_ context.Context, request *dlc.AttachUserPolicyRequest) (*dlc.AttachUserPolicyResponse, error) {
		called = true
		resp := dlc.NewAttachUserPolicyResponse()
		resp.Response = &dlc.AttachUserPolicyResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-create"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
	// schema MaxItems=1 is enforced by Terraform core during plan/apply, not by
	// TestResourceDataRaw; the resource's own defensive check (len != 1) is what
	// this test exercises by passing zero policy items.
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id":      "100032676511",
		"account_type": "TencentAccount",
		"policy_set":   []interface{}{},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.False(t, called)
}

func TestDlcAttachUserPolicyAttachment_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.DescribeUserInfoRequest
	patches.ApplyMethodFunc(dlcClient, "DescribeUserInfoWithContext", func(_ context.Context, request *dlc.DescribeUserInfoRequest) (*dlc.DescribeUserInfoResponse, error) {
		capturedRequest = request
		totalCount := int64(1)
		resp := dlc.NewDescribeUserInfoResponse()
		resp.Response = &dlc.DescribeUserInfoResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-read"),
			UserInfo: &dlc.UserDetailInfo{
				UserId:      ptrStrDlcAupa("100032676511"),
				Type:        ptrStrDlcAupa("DataAuth"),
				AccountType: ptrStrDlcAupa("TencentAccount"),
				DataPolicyInfo: &dlc.Policys{
					TotalCount: &totalCount,
					PolicySet: []*dlc.Policy{
						{
							Database:   ptrStrDlcAupa("tf_example_db"),
							Catalog:    ptrStrDlcAupa("DataLakeCatalog"),
							Table:      ptrStrDlcAupa("tf_example_table"),
							Operation:  ptrStrDlcAupa("SELECT"),
							PolicyType: ptrStrDlcAupa("TABLE"),
							PolicyId:   ptrStrDlcAupa("policy-id-001"),
							ReAuth:     ptrBoolDlcAupa(false),
							Source:     ptrStrDlcAupa("USER"),
							Mode:       ptrStrDlcAupa("COMMON"),
						},
						{
							Database:   ptrStrDlcAupa("other_db"),
							Catalog:    ptrStrDlcAupa("DataLakeCatalog"),
							Table:      ptrStrDlcAupa("other_table"),
							Operation:  ptrStrDlcAupa("SELECT"),
							PolicyType: ptrStrDlcAupa("TABLE"),
							PolicyId:   ptrStrDlcAupa("policy-id-999"),
							Source:     ptrStrDlcAupa("USER"),
							Mode:       ptrStrDlcAupa("COMMON"),
						},
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id":      "100032676511",
		"account_type": "TencentAccount",
		"policy_set": []interface{}{
			map[string]interface{}{
				"database":    "tf_example_db",
				"catalog":     "DataLakeCatalog",
				"table":       "tf_example_table",
				"operation":   "SELECT",
				"policy_type": "TABLE",
			},
		},
	})
	d.SetId("100032676511#policy-id-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "100032676511#policy-id-001", d.Id())
	assert.Equal(t, "100032676511", d.Get("user_id"))
	assert.Equal(t, "TencentAccount", d.Get("account_type"))

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "100032676511", *capturedRequest.UserId)
	assert.Equal(t, "DataAuth", *capturedRequest.Type)
	assert.Equal(t, "TencentAccount", *capturedRequest.AccountType)

	// only the matched policy (by policy_id from the composite ID) should be
	// kept in state, even though the API returned 2 policies for the user.
	policySet := d.Get("policy_set").([]interface{})
	assert.Len(t, policySet, 1)
	policyMap := policySet[0].(map[string]interface{})
	assert.Equal(t, "policy-id-001", policyMap["policy_id"])
	assert.Equal(t, "USER", policyMap["source"])
}

func TestDlcAttachUserPolicyAttachment_Read_PolicyNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeUserInfoWithContext", func(_ context.Context, request *dlc.DescribeUserInfoRequest) (*dlc.DescribeUserInfoResponse, error) {
		totalCount := int64(1)
		resp := dlc.NewDescribeUserInfoResponse()
		resp.Response = &dlc.DescribeUserInfoResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-read"),
			UserInfo: &dlc.UserDetailInfo{
				UserId: ptrStrDlcAupa("100032676511"),
				DataPolicyInfo: &dlc.Policys{
					TotalCount: &totalCount,
					PolicySet: []*dlc.Policy{
						{
							Database: ptrStrDlcAupa("other_db"),
							PolicyId: ptrStrDlcAupa("policy-id-999"),
						},
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id":      "100032676511",
		"account_type": "TencentAccount",
	})
	d.SetId("100032676511#policy-id-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestDlcAttachUserPolicyAttachment_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeUserInfoWithContext", func(_ context.Context, request *dlc.DescribeUserInfoRequest) (*dlc.DescribeUserInfoResponse, error) {
		resp := dlc.NewDescribeUserInfoResponse()
		resp.Response = &dlc.DescribeUserInfoResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-read-empty"),
			UserInfo:  &dlc.UserDetailInfo{},
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id":      "100032676511",
		"account_type": "TencentAccount",
	})
	d.SetId("100032676511#policy-id-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestDlcAttachUserPolicyAttachment_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.DetachUserPolicyRequest
	patches.ApplyMethodFunc(dlcClient, "DetachUserPolicyWithContext", func(_ context.Context, request *dlc.DetachUserPolicyRequest) (*dlc.DetachUserPolicyResponse, error) {
		capturedRequest = request
		resp := dlc.NewDetachUserPolicyResponse()
		resp.Response = &dlc.DetachUserPolicyResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id":      "100032676511",
		"account_type": "TencentAccount",
		"policy_set": []interface{}{
			map[string]interface{}{
				"database":    "tf_example_db",
				"catalog":     "DataLakeCatalog",
				"table":       "tf_example_table",
				"operation":   "SELECT",
				"policy_type": "TABLE",
			},
		},
	})
	d.SetId("100032676511#policy-id-001")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "100032676511", *capturedRequest.UserId)
	assert.Equal(t, "TencentAccount", *capturedRequest.AccountType)
	assert.Len(t, capturedRequest.PolicyIds, 1)
	assert.Equal(t, "policy-id-001", *capturedRequest.PolicyIds[0])
}

func TestDlcAttachUserPolicyAttachment_Update_Immutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

	meta := newMockMetaDlcAttachUserPolicyAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id":      "100032676511",
		"account_type": "TencentAccount",
		"policy_set": []interface{}{
			map[string]interface{}{
				"database":  "tf_example_db",
				"catalog":   "DataLakeCatalog",
				"table":     "tf_example_table",
				"operation": "SELECT",
			},
		},
	})
	d.SetId("100032676511#policy-id-001")

	// simulate a change on user_id (ForceNew top-level arg) -> update should error
	if err := d.Set("user_id", "100032676512"); err != nil {
		t.Fatalf("failed to set user_id: %v", err)
	}

	err := res.Update(d, meta)
	assert.Error(t, err)
}

func TestDlcAttachUserPolicyAttachment_Create_NilPolicyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "AttachUserPolicyWithContext", func(_ context.Context, request *dlc.AttachUserPolicyRequest) (*dlc.AttachUserPolicyResponse, error) {
		resp := dlc.NewAttachUserPolicyResponse()
		resp.Response = &dlc.AttachUserPolicyResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-create-empty"),
			PolicySet: []*dlc.Policy{
				{
					Database: ptrStrDlcAupa("tf_example_db"),
					PolicyId: nil,
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id":      "100032676511",
		"account_type": "TencentAccount",
		"policy_set": []interface{}{
			map[string]interface{}{
				"database":  "tf_example_db",
				"catalog":   "DataLakeCatalog",
				"table":     "tf_example_table",
				"operation": "SELECT",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// keep references to helpers used to avoid unused warnings in some build configs
var _ = ptrInt64DlcAupa
