package cynosdb_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cynosdbSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"
)

type mockMetaClusterReadOnlyInstanceGroupAccesOperation struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaClusterReadOnlyInstanceGroupAccesOperation) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaClusterReadOnlyInstanceGroupAccesOperation{}

func newMockMetaClusterReadOnlyInstanceGroupAccesOperation() *mockMetaClusterReadOnlyInstanceGroupAccesOperation {
	return &mockMetaClusterReadOnlyInstanceGroupAccesOperation{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCROIGAO(s string) *string {
	return &s
}

func ptrInt64CROIGAO(i int64) *int64 {
	return &i
}

func TestCynosdbClusterReadOnlyInstanceGroupAccesOperation_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cynosdbClient := &cynosdbSDK.Client{}
	meta := newMockMetaClusterReadOnlyInstanceGroupAccesOperation()
	patches.ApplyMethodReturn(meta.client, "UseCynosdbClient", cynosdbClient)

	patches.ApplyMethodFunc(cynosdbClient, "OpenClusterReadOnlyInstanceGroupAccessWithContext", func(_ context.Context, request *cynosdbSDK.OpenClusterReadOnlyInstanceGroupAccessRequest) (*cynosdbSDK.OpenClusterReadOnlyInstanceGroupAccessResponse, error) {
		assert.NotNil(t, request.ClusterId)
		assert.Equal(t, "cynosdbmysql-12345678", *request.ClusterId)
		assert.NotNil(t, request.Port)
		assert.Equal(t, "3306", *request.Port)
		assert.Equal(t, 2, len(request.SecurityGroupIds))
		assert.Equal(t, "sg-aaa", *request.SecurityGroupIds[0])
		assert.Equal(t, "sg-bbb", *request.SecurityGroupIds[1])

		resp := cynosdbSDK.NewOpenClusterReadOnlyInstanceGroupAccessResponse()
		resp.Response = &cynosdbSDK.OpenClusterReadOnlyInstanceGroupAccessResponseParams{
			FlowId:    ptrInt64CROIGAO(123456),
			RequestId: ptrStringCROIGAO("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cynosdbClient, "DescribeFlow", func(request *cynosdbSDK.DescribeFlowRequest) (*cynosdbSDK.DescribeFlowResponse, error) {
		assert.NotNil(t, request.FlowId)
		assert.Equal(t, int64(123456), *request.FlowId)

		resp := cynosdbSDK.NewDescribeFlowResponse()
		resp.Response = &cynosdbSDK.DescribeFlowResponseParams{
			Status:    ptrInt64CROIGAO(0),
			RequestId: ptrStringCROIGAO("fake-describe-flow-request-id"),
		}
		return resp, nil
	})

	res := cynosdb.ResourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id":         "cynosdbmysql-12345678",
		"port":               "3306",
		"security_group_ids": []interface{}{"sg-aaa", "sg-bbb"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, 123456, d.Get("flow_id"))
}

func TestCynosdbClusterReadOnlyInstanceGroupAccesOperation_Read(t *testing.T) {
	meta := newMockMetaClusterReadOnlyInstanceGroupAccesOperation()
	res := cynosdb.ResourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cynosdbmysql-12345678",
	})
	d.SetId("some-token-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "some-token-id", d.Id())
}

func TestCynosdbClusterReadOnlyInstanceGroupAccesOperation_Delete(t *testing.T) {
	meta := newMockMetaClusterReadOnlyInstanceGroupAccesOperation()
	res := cynosdb.ResourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cynosdbmysql-12345678",
	})
	d.SetId("some-token-id")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
