package dlc_test

import (
	"context"
	"fmt"
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

// PolicyId format: v1|{SubjectType}|{SubjectId}|{PolicyType}|{Mode}|{Catalog}|{Database}|{Table}|{View}|{Function}|{Column}|{DataEngine}|{Operation}
const (
	dlcAupaPolicyIdTable = "v1|USER|100032676511|TABLE|COMMON|DataLakeCatalog|tf_example_db|tf_example_table||||SELECT"
	dlcAupaPolicyIdOther = "v1|USER|100032676511|TABLE|COMMON|DataLakeCatalog|other_db|other_table||||SELECT"
)

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
					PolicyId:   ptrStrDlcAupa(dlcAupaPolicyIdTable),
					Source:     ptrStrDlcAupa("USER"),
					Mode:       ptrStrDlcAupa("COMMON"),
				},
			},
		}
		return resp, nil
	})

	var describeRequest *dlc.DescribeUserInfoRequest
	patches.ApplyMethodFunc(dlcClient, "DescribeUserInfoWithContext", func(_ context.Context, request *dlc.DescribeUserInfoRequest) (*dlc.DescribeUserInfoResponse, error) {
		describeRequest = request
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
							PolicyId:   ptrStrDlcAupa(dlcAupaPolicyIdTable),
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
	assert.Equal(t, fmt.Sprintf("100032676511#%s", dlcAupaPolicyIdTable), d.Id())

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "100032676511", *capturedRequest.UserId)
	assert.Equal(t, "TencentAccount", *capturedRequest.AccountType)
	assert.Len(t, capturedRequest.PolicySet, 1)
	assert.Equal(t, "tf_example_db", *capturedRequest.PolicySet[0].Database)
	assert.Equal(t, "SELECT", *capturedRequest.PolicySet[0].Operation)

	// PolicyType `TABLE` in the PolicyId must be mapped to DescribeUserInfo `Type`=`DataAuth`.
	assert.NotNil(t, describeRequest)
	assert.Equal(t, "DataAuth", *describeRequest.Type)

	policySet := d.Get("policy_set").([]interface{})
	assert.Len(t, policySet, 1)
	policyMap := policySet[0].(map[string]interface{})
	assert.Equal(t, dlcAupaPolicyIdTable, policyMap["policy_id"])
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
							PolicyId:   ptrStrDlcAupa(dlcAupaPolicyIdTable),
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
							PolicyId:   ptrStrDlcAupa(dlcAupaPolicyIdOther),
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
	d.SetId(fmt.Sprintf("100032676511#%s", dlcAupaPolicyIdTable))

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("100032676511#%s", dlcAupaPolicyIdTable), d.Id())
	assert.Equal(t, "100032676511", d.Get("user_id"))
	assert.Equal(t, "TencentAccount", d.Get("account_type"))

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "100032676511", *capturedRequest.UserId)
	// PolicyType `TABLE` parsed from the PolicyId must be mapped to `DataAuth`.
	assert.Equal(t, "DataAuth", *capturedRequest.Type)
	assert.Equal(t, "TencentAccount", *capturedRequest.AccountType)

	// only the matched policy (by policy_id from the composite ID) should be
	// kept in state, even though the API returned 2 policies for the user.
	policySet := d.Get("policy_set").([]interface{})
	assert.Len(t, policySet, 1)
	policyMap := policySet[0].(map[string]interface{})
	assert.Equal(t, dlcAupaPolicyIdTable, policyMap["policy_id"])
	assert.Equal(t, "USER", policyMap["source"])
}

// TestDlcAttachUserPolicyAttachment_Read_TypeMapping verifies that the `PolicyType`
// segment (the 4th `|`-separated field) parsed from the composite ID's policy_id is
// correctly mapped to the DescribeUserInfo `Type` field, and that the resource reads
// the policy back from the corresponding response field (CatalogPolicyInfo,
// EnginePolicyInfo, RowFilterInfo, ModelPolicyInfo).
func TestDlcAttachUserPolicyAttachment_Read_TypeMapping(t *testing.T) {
	cases := []struct {
		name          string
		policyType    string
		expectedType  string
		buildUserInfo func(policyId string) *dlc.UserDetailInfo
	}{
		{
			name:         "DATASOURCE_maps_to_CatalogAuth",
			policyType:   "DATASOURCE",
			expectedType: "CatalogAuth",
			buildUserInfo: func(policyId string) *dlc.UserDetailInfo {
				return &dlc.UserDetailInfo{
					UserId: ptrStrDlcAupa("100032676511"),
					CatalogPolicyInfo: &dlc.Policys{
						TotalCount: ptrInt64DlcAupa(1),
						PolicySet: []*dlc.Policy{
							{PolicyId: ptrStrDlcAupa(policyId), Catalog: ptrStrDlcAupa("COSDataCatalog")},
						},
					},
				}
			},
		},
		{
			name:         "ENGINE_maps_to_EngineAuth",
			policyType:   "ENGINE",
			expectedType: "EngineAuth",
			buildUserInfo: func(policyId string) *dlc.UserDetailInfo {
				return &dlc.UserDetailInfo{
					UserId: ptrStrDlcAupa("100032676511"),
					EnginePolicyInfo: &dlc.Policys{
						TotalCount: ptrInt64DlcAupa(1),
						PolicySet: []*dlc.Policy{
							{PolicyId: ptrStrDlcAupa(policyId), DataEngine: ptrStrDlcAupa("engine1")},
						},
					},
				}
			},
		},
		{
			name:         "ROWFILTER_maps_to_RowFilter",
			policyType:   "ROWFILTER",
			expectedType: "RowFilter",
			buildUserInfo: func(policyId string) *dlc.UserDetailInfo {
				return &dlc.UserDetailInfo{
					UserId: ptrStrDlcAupa("100032676511"),
					RowFilterInfo: &dlc.Policys{
						TotalCount: ptrInt64DlcAupa(1),
						PolicySet: []*dlc.Policy{
							{PolicyId: ptrStrDlcAupa(policyId), Database: ptrStrDlcAupa("tf_example_db")},
						},
					},
				}
			},
		},
		{
			name:         "MODEL_maps_to_MODEL",
			policyType:   "MODEL",
			expectedType: "MODEL",
			buildUserInfo: func(policyId string) *dlc.UserDetailInfo {
				return &dlc.UserDetailInfo{
					UserId: ptrStrDlcAupa("100032676511"),
					ModelPolicyInfo: &dlc.Policys{
						TotalCount: ptrInt64DlcAupa(1),
						PolicySet: []*dlc.Policy{
							{PolicyId: ptrStrDlcAupa(policyId), Model: ptrStrDlcAupa("model1")},
						},
					},
				}
			},
		},
		{
			name:         "ADMIN_maps_to_DataAuth",
			policyType:   "ADMIN",
			expectedType: "DataAuth",
			buildUserInfo: func(policyId string) *dlc.UserDetailInfo {
				return &dlc.UserDetailInfo{
					UserId: ptrStrDlcAupa("100032676511"),
					DataPolicyInfo: &dlc.Policys{
						TotalCount: ptrInt64DlcAupa(1),
						PolicySet: []*dlc.Policy{
							{PolicyId: ptrStrDlcAupa(policyId), Database: ptrStrDlcAupa("*")},
						},
					},
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			policyId := fmt.Sprintf("v1|USER|100032676511|%s|COMMON||||||||", tc.policyType)

			patches := gomonkey.NewPatches()
			defer patches.Reset()

			dlcClient := &dlc.Client{}
			patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

			var capturedRequest *dlc.DescribeUserInfoRequest
			patches.ApplyMethodFunc(dlcClient, "DescribeUserInfoWithContext", func(_ context.Context, request *dlc.DescribeUserInfoRequest) (*dlc.DescribeUserInfoResponse, error) {
				capturedRequest = request
				resp := dlc.NewDescribeUserInfoResponse()
				resp.Response = &dlc.DescribeUserInfoResponseParams{
					RequestId: ptrStrDlcAupa("fake-request-id-read"),
					UserInfo:  tc.buildUserInfo(policyId),
				}
				return resp, nil
			})

			meta := newMockMetaDlcAttachUserPolicyAttachment()
			res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
				"user_id": "100032676511",
			})
			d.SetId(fmt.Sprintf("100032676511#%s", policyId))

			err := res.Read(d, meta)
			assert.NoError(t, err)
			assert.NotEqual(t, "", d.Id())

			assert.NotNil(t, capturedRequest)
			assert.Equal(t, tc.expectedType, *capturedRequest.Type)

			policySet := d.Get("policy_set").([]interface{})
			assert.Len(t, policySet, 1)
			policyMap := policySet[0].(map[string]interface{})
			assert.Equal(t, policyId, policyMap["policy_id"])
		})
	}
}

func TestDlcAttachUserPolicyAttachment_Read_InvalidPolicyIdFormat(t *testing.T) {
	meta := newMockMetaDlcAttachUserPolicyAttachment()
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_id": "100032676511",
	})
	// policy_id segment does not contain the required `|`-separated PolicyType field.
	d.SetId("100032676511#not-a-valid-policy-id")

	err := res.Read(d, meta)
	assert.Error(t, err)
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
							PolicyId: ptrStrDlcAupa(dlcAupaPolicyIdOther),
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
	d.SetId(fmt.Sprintf("100032676511#%s", dlcAupaPolicyIdTable))

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
	d.SetId(fmt.Sprintf("100032676511#%s", dlcAupaPolicyIdTable))

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
	d.SetId(fmt.Sprintf("100032676511#%s", dlcAupaPolicyIdTable))

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "100032676511", *capturedRequest.UserId)
	assert.Equal(t, "TencentAccount", *capturedRequest.AccountType)
	assert.Len(t, capturedRequest.PolicyIds, 1)
	assert.Equal(t, dlcAupaPolicyIdTable, *capturedRequest.PolicyIds[0])
}

// TestDlcAttachUserPolicyAttachment_Update_Immutable verifies that all of the
// resource's schema fields are ForceNew (user_id, policy_set, account_type),
// so the resource has no Update function and any change forces recreation.
func TestDlcAttachUserPolicyAttachment_Update_Immutable(t *testing.T) {
	res := svcdlc.ResourceTencentCloudDlcAttachUserPolicyAttachment()

	assert.Nil(t, res.Update, "resource should not define an Update function; all fields are ForceNew")

	for name, s := range res.Schema {
		assert.True(t, s.ForceNew, "field %q should be ForceNew", name)
	}
}

// TestDlcAttachUserPolicyAttachment_Read_OperationOrderInsensitive verifies that
// a drift is not reported when the API returns the `operation` value with the
// same comma-separated tokens as configured but in a different order
// (e.g. configured "MONITOR,USE" vs. API-returned "USE,MONITOR").
func TestDlcAttachUserPolicyAttachment_Read_OperationOrderInsensitive(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcAttachUserPolicyAttachment().client, "UseDlcClient", dlcClient)

	policyId := "v1|USER|100029759702|ENGINE|SENIOR|||||||test|MONITOR%2CUSE"

	patches.ApplyMethodFunc(dlcClient, "DescribeUserInfoWithContext", func(_ context.Context, request *dlc.DescribeUserInfoRequest) (*dlc.DescribeUserInfoResponse, error) {
		resp := dlc.NewDescribeUserInfoResponse()
		resp.Response = &dlc.DescribeUserInfoResponseParams{
			RequestId: ptrStrDlcAupa("fake-request-id-read"),
			UserInfo: &dlc.UserDetailInfo{
				UserId: ptrStrDlcAupa("100029759702"),
				EnginePolicyInfo: &dlc.Policys{
					TotalCount: ptrInt64DlcAupa(1),
					PolicySet: []*dlc.Policy{
						{
							PolicyId:   ptrStrDlcAupa(policyId),
							DataEngine: ptrStrDlcAupa("test"),
							// API returns tokens in a different order than what was configured.
							Operation: ptrStrDlcAupa("USE,MONITOR"),
							Mode:      ptrStrDlcAupa("SENIOR"),
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
		"user_id": "100029759702",
		"policy_set": []interface{}{
			map[string]interface{}{
				"data_engine": "test",
				"operation":   "MONITOR,USE",
				"policy_type": "ENGINE",
			},
		},
	})
	d.SetId(fmt.Sprintf("100029759702#%s", policyId))

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// After Read, state now holds the API's order ("USE,MONITOR"), while the
	// configured order was "MONITOR,USE". The DiffSuppressFunc on `operation`
	// must treat these as equal so no drift is reported.
	suppress := res.Schema["policy_set"].Elem.(*schema.Resource).Schema["operation"].DiffSuppressFunc
	assert.NotNil(t, suppress)
	assert.True(t, suppress("policy_set.0.operation", "MONITOR,USE", "USE,MONITOR", d))
	assert.False(t, suppress("policy_set.0.operation", "MONITOR,USE", "USE,MONITOR,EXTRA", d))
	assert.False(t, suppress("policy_set.0.operation", "MONITOR,USE", "SELECT,INSERT", d))
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
