package tpulsar_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tdmqsdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tpulsar"
)

type mockMetaTdmqProInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTdmqProInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTdmqProInstance{}

func newMockMetaTdmqProInstance() *mockMetaTdmqProInstance {
	return &mockMetaTdmqProInstance{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTdmqPro(s string) *string { return &s }
func ptrInt64TdmqPro(i int64) *int64    { return &i }
func ptrUint64TdmqPro(i uint64) *uint64 { return &i }

func TestTdmqProInstance_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmqClient := &tdmqsdk.Client{}
	patches.ApplyMethodReturn(newMockMetaTdmqProInstance().client, "UseTdmqClient", tdmqClient)
	patches.ApplyMethodFunc(tdmqClient, "CreateProCluster", func(request *tdmqsdk.CreateProClusterRequest) (*tdmqsdk.CreateProClusterResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid parameter")
	})

	meta := newMockMetaTdmqProInstance()
	res := tpulsar.ResourceTencentCloudTdmqProInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_ids":     []interface{}{200002},
		"product_name": "PULSAR.P1.MINI2",
		"storage_size": 200,
		"vpc":          []interface{}{map[string]interface{}{"vpc_id": "vpc-123", "subnet_id": "subnet-123"}},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

func TestTdmqProInstance_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmqClient := &tdmqsdk.Client{}
	patches.ApplyMethodReturn(newMockMetaTdmqProInstance().client, "UseTdmqClient", tdmqClient)
	patches.ApplyMethodFunc(tdmqClient, "DescribeClusters", func(request *tdmqsdk.DescribeClustersRequest) (*tdmqsdk.DescribeClustersResponse, error) {
		resp := tdmqsdk.NewDescribeClustersResponse()
		resp.Response = &tdmqsdk.DescribeClustersResponseParams{
			ClusterSet: []*tdmqsdk.Cluster{},
			TotalCount: ptrInt64TdmqPro(0),
			RequestId:  ptrStringTdmqPro("req-empty"),
		}
		return resp, nil
	})

	meta := newMockMetaTdmqProInstance()
	res := tpulsar.ResourceTencentCloudTdmqProInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("pulsar-nonexistent")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTdmqProInstance_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmqClient := &tdmqsdk.Client{}
	patches.ApplyMethodReturn(newMockMetaTdmqProInstance().client, "UseTdmqClient", tdmqClient)
	patches.ApplyMethodFunc(tdmqClient, "DeleteCluster", func(request *tdmqsdk.DeleteClusterRequest) (*tdmqsdk.DeleteClusterResponse, error) {
		resp := tdmqsdk.NewDeleteClusterResponse()
		resp.Response = &tdmqsdk.DeleteClusterResponseParams{
			ClusterId: ptrStringTdmqPro("pulsar-test123"),
			RequestId: ptrStringTdmqPro("req-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaTdmqProInstance()
	res := tpulsar.ResourceTencentCloudTdmqProInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("pulsar-test123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestTdmqProInstance_Schema(t *testing.T) {
	res := tpulsar.ResourceTencentCloudTdmqProInstance()
	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.Contains(t, res.Schema, "zone_ids")
	assert.True(t, res.Schema["zone_ids"].Required)
	assert.True(t, res.Schema["zone_ids"].ForceNew)
	assert.Contains(t, res.Schema, "product_name")
	assert.True(t, res.Schema["product_name"].Required)
	assert.Contains(t, res.Schema, "storage_size")
	assert.True(t, res.Schema["storage_size"].Required)
	assert.Contains(t, res.Schema, "vpc")
	assert.True(t, res.Schema["vpc"].Required)
	assert.Contains(t, res.Schema, "cluster_name")
	assert.True(t, res.Schema["cluster_name"].Optional)
	assert.Contains(t, res.Schema, "cluster_id")
	assert.True(t, res.Schema["cluster_id"].Computed)
	assert.Contains(t, res.Schema, "status")
	assert.True(t, res.Schema["status"].Computed)
}
