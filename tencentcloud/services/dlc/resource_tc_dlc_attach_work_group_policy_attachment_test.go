package dlc_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dlc_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"
)

type mockMetaAttachWorkGroupPolicyAttachment struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaAttachWorkGroupPolicyAttachment) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaAttachWorkGroupPolicyAttachment{}

func newMockMetaAttachWorkGroupPolicyAttachment() *mockMetaAttachWorkGroupPolicyAttachment {
	return &mockMetaAttachWorkGroupPolicyAttachment{client: &connectivity.TencentCloudClient{}}
}

func ptrStrAWGPA(s string) *string {
	return &s
}

func ptrInt64AWGPA(i int64) *int64 {
	return &i
}

func ptrBoolAWGPA(b bool) *bool {
	return &b
}

func TestDlcAttachWorkGroupPolicyAttachment_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaAttachWorkGroupPolicyAttachment().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "AttachWorkGroupPolicy", func(request *dlc_sdk.AttachWorkGroupPolicyRequest) (*dlc_sdk.AttachWorkGroupPolicyResponse, error) {
		assert.NotNil(t, request.WorkGroupId)
		assert.Equal(t, int64(23184), *request.WorkGroupId)
		assert.NotNil(t, request.PolicySet)
		assert.Equal(t, 1, len(request.PolicySet))

		resp := dlc_sdk.NewAttachWorkGroupPolicyResponse()
		resp.Response = &dlc_sdk.AttachWorkGroupPolicyResponseParams{
			PolicySet: []*dlc_sdk.Policy{
				{
					PolicyId: ptrStrAWGPA("policy-xxxx"),
				},
			},
			RequestId: ptrStrAWGPA("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeWorkGroupInfo", func(request *dlc_sdk.DescribeWorkGroupInfoRequest) (*dlc_sdk.DescribeWorkGroupInfoResponse, error) {
		resp := dlc_sdk.NewDescribeWorkGroupInfoResponse()
		resp.Response = &dlc_sdk.DescribeWorkGroupInfoResponseParams{
			WorkGroupInfo: &dlc_sdk.WorkGroupDetailInfo{
				WorkGroupId: ptrInt64AWGPA(23184),
				DataPolicyInfo: &dlc_sdk.Policys{
					PolicySet: []*dlc_sdk.Policy{
						{
							PolicyId: ptrStrAWGPA("policy-xxxx"),
							Database: ptrStrAWGPA("tf_example_db"),
							Catalog:  ptrStrAWGPA("DataLakeCatalog"),
							Table:    ptrStrAWGPA("tf_example_table"),
						},
					},
				},
			},
			RequestId: ptrStrAWGPA("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAttachWorkGroupPolicyAttachment()
	res := dlc.ResourceTencentCloudDlcAttachWorkGroupPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"work_group_id": 23184,
		"policy_set": []interface{}{
			map[string]interface{}{
				"database":    "tf_example_db",
				"catalog":     "DataLakeCatalog",
				"table":       "tf_example_table",
				"operation":   "ASSAYER",
				"policy_type": "DATABASE",
				"source":      "USER",
				"mode":        "COMMON",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "23184#policy-xxxx", d.Id())
}

func TestDlcAttachWorkGroupPolicyAttachment_ReadBindingExists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaAttachWorkGroupPolicyAttachment().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeWorkGroupInfo", func(request *dlc_sdk.DescribeWorkGroupInfoRequest) (*dlc_sdk.DescribeWorkGroupInfoResponse, error) {
		assert.NotNil(t, request.WorkGroupId)
		assert.Equal(t, int64(23184), *request.WorkGroupId)
		assert.NotNil(t, request.Type)
		assert.Equal(t, "DataAuth", *request.Type)

		resp := dlc_sdk.NewDescribeWorkGroupInfoResponse()
		resp.Response = &dlc_sdk.DescribeWorkGroupInfoResponseParams{
			WorkGroupInfo: &dlc_sdk.WorkGroupDetailInfo{
				WorkGroupId: ptrInt64AWGPA(23184),
				DataPolicyInfo: &dlc_sdk.Policys{
					PolicySet: []*dlc_sdk.Policy{
						{
							PolicyId:  ptrStrAWGPA("policy-xxxx"),
							Database:  ptrStrAWGPA("tf_example_db"),
							Catalog:   ptrStrAWGPA("DataLakeCatalog"),
							Table:     ptrStrAWGPA("tf_example_table"),
							Operation: ptrStrAWGPA("ASSAYER"),
							ReAuth:    ptrBoolAWGPA(false),
						},
					},
				},
			},
			RequestId: ptrStrAWGPA("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAttachWorkGroupPolicyAttachment()
	res := dlc.ResourceTencentCloudDlcAttachWorkGroupPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"work_group_id": 23184,
		"policy_set": []interface{}{
			map[string]interface{}{
				"database":  "tf_example_db",
				"catalog":   "DataLakeCatalog",
				"table":     "tf_example_table",
				"operation": "ASSAYER",
			},
		},
	})
	d.SetId("23184#policy-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "23184#policy-xxxx", d.Id())
	assert.Equal(t, 23184, d.Get("work_group_id"))
}

func TestDlcAttachWorkGroupPolicyAttachment_ReadDriftClear(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaAttachWorkGroupPolicyAttachment().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeWorkGroupInfo", func(request *dlc_sdk.DescribeWorkGroupInfoRequest) (*dlc_sdk.DescribeWorkGroupInfoResponse, error) {
		resp := dlc_sdk.NewDescribeWorkGroupInfoResponse()
		resp.Response = &dlc_sdk.DescribeWorkGroupInfoResponseParams{
			WorkGroupInfo: &dlc_sdk.WorkGroupDetailInfo{
				WorkGroupId: ptrInt64AWGPA(23184),
				DataPolicyInfo: &dlc_sdk.Policys{
					PolicySet: []*dlc_sdk.Policy{
						{
							PolicyId: ptrStrAWGPA("policy-other"),
						},
					},
				},
			},
			RequestId: ptrStrAWGPA("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAttachWorkGroupPolicyAttachment()
	res := dlc.ResourceTencentCloudDlcAttachWorkGroupPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"work_group_id": 23184,
	})
	d.SetId("23184#policy-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestDlcAttachWorkGroupPolicyAttachment_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaAttachWorkGroupPolicyAttachment().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DetachWorkGroupPolicy", func(request *dlc_sdk.DetachWorkGroupPolicyRequest) (*dlc_sdk.DetachWorkGroupPolicyResponse, error) {
		assert.NotNil(t, request.WorkGroupId)
		assert.Equal(t, int64(23184), *request.WorkGroupId)
		assert.NotNil(t, request.PolicyIds)
		assert.Equal(t, 1, len(request.PolicyIds))
		assert.Equal(t, "policy-xxxx", *request.PolicyIds[0])

		resp := dlc_sdk.NewDetachWorkGroupPolicyResponse()
		resp.Response = &dlc_sdk.DetachWorkGroupPolicyResponseParams{
			RequestId: ptrStrAWGPA("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAttachWorkGroupPolicyAttachment()
	res := dlc.ResourceTencentCloudDlcAttachWorkGroupPolicyAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"work_group_id": 23184,
	})
	d.SetId("23184#policy-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
