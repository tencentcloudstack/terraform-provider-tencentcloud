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

type mockMetaDlcAttachUserPolicyrAttachment struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDlcAttachUserPolicyrAttachment) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDlcAttachUserPolicyrAttachment{}

func newMockMetaDlcAttachUserPolicyrAttachment() *mockMetaDlcAttachUserPolicyrAttachment {
	return &mockMetaDlcAttachUserPolicyrAttachment{client: &connectivity.TencentCloudClient{}}
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

func TestDlcAttachUserPolicyrAttachment_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyrAttachment().client, "UseDlcClient", dlcClient)

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

	meta := newMockMetaDlcAttachUserPolicyrAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyrAttachment()
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
	assert.Equal(t, "100032676511#TencentAccount", d.Id())

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

func TestDlcAttachUserPolicyrAttachment_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyrAttachment().client, "UseDlcClient", dlcClient)

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
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyrAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyrAttachment()
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
	d.SetId("100032676511#TencentAccount")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "100032676511#TencentAccount", d.Id())
	assert.Equal(t, "100032676511", d.Get("user_id"))
	assert.Equal(t, "TencentAccount", d.Get("account_type"))

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "100032676511", *capturedRequest.UserId)
	assert.Equal(t, "DataAuth", *capturedRequest.Type)
	assert.Equal(t, "TencentAccount", *capturedRequest.AccountType)

	policySet := d.Get("policy_set").([]interface{})
	assert.Len(t, policySet, 1)
	policyMap := policySet[0].(map[string]interface{})
	assert.Equal(t, "policy-id-001", policyMap["policy_id"])
	assert.Equal(t, "USER", policyMap["source"])
}

func TestDlcAttachUserPolicyrAttachment_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyrAttachment().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeUserInfoWithContext", func(_ context.Context, request *dlc.DescribeUserInfoRequest) (*dlc.DescribeUserInfoResponse, error) {
		resp := dlc.NewDescribeUserInfoResponse()
		resp.Response = &dlc.DescribeUserInfoResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-read-empty"),
			UserInfo:  &dlc.UserDetailInfo{},
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyrAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyrAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id":      "100032676511",
		"account_type": "TencentAccount",
	})
	d.SetId("100032676511#TencentAccount")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestDlcAttachUserPolicyrAttachment_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyrAttachment().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.DetachUserPolicyRequest
	patches.ApplyMethodFunc(dlcClient, "DetachUserPolicyWithContext", func(_ context.Context, request *dlc.DetachUserPolicyRequest) (*dlc.DetachUserPolicyResponse, error) {
		capturedRequest = request
		resp := dlc.NewDetachUserPolicyResponse()
		resp.Response = &dlc.DetachUserPolicyResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcAttachUserPolicyrAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyrAttachment()
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
	d.SetId("100032676511#TencentAccount")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "100032676511", *capturedRequest.UserId)
	assert.Equal(t, "TencentAccount", *capturedRequest.AccountType)
	assert.Len(t, capturedRequest.PolicySet, 1)
	assert.Equal(t, "tf_example_db", *capturedRequest.PolicySet[0].Database)
	assert.Equal(t, "SELECT", *capturedRequest.PolicySet[0].Operation)
}

func TestDlcAttachUserPolicyrAttachment_Update_Immutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyrAttachment().client, "UseDlcClient", dlcClient)

	meta := newMockMetaDlcAttachUserPolicyrAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyrAttachment()
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
	d.SetId("100032676511#TencentAccount")

	// simulate a change on user_id (ForceNew top-level arg) -> update should error
	d.Set("user_id", "100032676512")

	err := res.Update(d, meta)
	assert.Error(t, err)
}

func TestDlcAttachUserPolicyrAttachment_Create_NilPolicyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyrAttachment().client, "UseDlcClient", dlcClient)

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

	meta := newMockMetaDlcAttachUserPolicyrAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyrAttachment()
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
