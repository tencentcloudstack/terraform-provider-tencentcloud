package dbdc_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

type mockMetaDbdcClusterNodeResourcesDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDbdcClusterNodeResourcesDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDbdcClusterNodeResourcesDS{}

func newMockMetaDbdcClusterNodeResourcesDS() *mockMetaDbdcClusterNodeResourcesDS {
	return &mockMetaDbdcClusterNodeResourcesDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrCNR(s string) *string {
	return &s
}

func ptrFloat64CNR(f float64) *float64 {
	return &f
}

func ptrUint64CNR(u uint64) *uint64 {
	return &u
}

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomClusterNodeResourcesDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomClusterNodeResourcesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcClusterNodeResourcesDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterNodeResources", func(request *dbdcv20201029.DescribeDBCustomClusterNodeResourcesRequest) (*dbdcv20201029.DescribeDBCustomClusterNodeResourcesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterNodeResourcesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterNodeResourcesResponseParams{
			NodeSet: []*dbdcv20201029.DBCustomClusterNodeResource{
				{
					NodeId: ptrStrCNR("node-abc123"),
					Capacity: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(64),
						Memory: ptrFloat64CNR(512),
						Pods:   ptrUint64CNR(100),
					},
					Allocatable: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(60),
						Memory: ptrFloat64CNR(480),
						Pods:   ptrUint64CNR(90),
					},
					Requests: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(10),
						Memory: ptrFloat64CNR(80),
						Pods:   ptrUint64CNR(20),
					},
					Limits: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(20),
						Memory: ptrFloat64CNR(160),
						Pods:   ptrUint64CNR(40),
					},
					Available: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(50),
						Memory: ptrFloat64CNR(400),
						Pods:   ptrUint64CNR(70),
					},
				},
				{
					NodeId: ptrStrCNR("node-def456"),
					Capacity: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(128),
						Memory: ptrFloat64CNR(1024),
						Pods:   ptrUint64CNR(200),
					},
					Allocatable: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(120),
						Memory: ptrFloat64CNR(960),
						Pods:   ptrUint64CNR(180),
					},
					Requests: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(30),
						Memory: ptrFloat64CNR(240),
						Pods:   ptrUint64CNR(60),
					},
					Limits: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(60),
						Memory: ptrFloat64CNR(480),
						Pods:   ptrUint64CNR(120),
					},
					Available: &dbdcv20201029.MetaResource{
						Cpu:    ptrFloat64CNR(90),
						Memory: ptrFloat64CNR(720),
						Pods:   ptrUint64CNR(120),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcClusterNodeResourcesDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodeResources()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "dbcc-nmtmsew8",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	nodeSet := d.Get("node_set").([]interface{})
	assert.Len(t, nodeSet, 2)

	node0 := nodeSet[0].(map[string]interface{})
	assert.Equal(t, "node-abc123", node0["node_id"].(string))

	capacity0 := node0["capacity"].([]interface{})
	assert.Len(t, capacity0, 1)
	capacityMap0 := capacity0[0].(map[string]interface{})
	assert.Equal(t, float64(64), capacityMap0["cpu"].(float64))
	assert.Equal(t, float64(512), capacityMap0["memory"].(float64))
	assert.Equal(t, 100, capacityMap0["pods"].(int))

	allocatable0 := node0["allocatable"].([]interface{})
	assert.Len(t, allocatable0, 1)
	allocatableMap0 := allocatable0[0].(map[string]interface{})
	assert.Equal(t, float64(60), allocatableMap0["cpu"].(float64))
	assert.Equal(t, float64(480), allocatableMap0["memory"].(float64))
	assert.Equal(t, 90, allocatableMap0["pods"].(int))

	requests0 := node0["requests"].([]interface{})
	assert.Len(t, requests0, 1)
	requestsMap0 := requests0[0].(map[string]interface{})
	assert.Equal(t, float64(10), requestsMap0["cpu"].(float64))
	assert.Equal(t, float64(80), requestsMap0["memory"].(float64))
	assert.Equal(t, 20, requestsMap0["pods"].(int))

	limits0 := node0["limits"].([]interface{})
	assert.Len(t, limits0, 1)
	limitsMap0 := limits0[0].(map[string]interface{})
	assert.Equal(t, float64(20), limitsMap0["cpu"].(float64))
	assert.Equal(t, float64(160), limitsMap0["memory"].(float64))
	assert.Equal(t, 40, limitsMap0["pods"].(int))

	available0 := node0["available"].([]interface{})
	assert.Len(t, available0, 1)
	availableMap0 := available0[0].(map[string]interface{})
	assert.Equal(t, float64(50), availableMap0["cpu"].(float64))
	assert.Equal(t, float64(400), availableMap0["memory"].(float64))
	assert.Equal(t, 70, availableMap0["pods"].(int))

	node1 := nodeSet[1].(map[string]interface{})
	assert.Equal(t, "node-def456", node1["node_id"].(string))

	capacity1 := node1["capacity"].([]interface{})
	assert.Len(t, capacity1, 1)
	capacityMap1 := capacity1[0].(map[string]interface{})
	assert.Equal(t, float64(128), capacityMap1["cpu"].(float64))
	assert.Equal(t, float64(1024), capacityMap1["memory"].(float64))
	assert.Equal(t, 200, capacityMap1["pods"].(int))
}

func TestDbdcDbCustomClusterNodeResourcesDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodeResources()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "cluster_id")
	assert.Contains(t, res.Schema, "node_ids")
	assert.Contains(t, res.Schema, "result_output_file")
	assert.Contains(t, res.Schema, "node_set")

	clusterIdSchema := res.Schema["cluster_id"]
	assert.Equal(t, schema.TypeString, clusterIdSchema.Type)
	assert.True(t, clusterIdSchema.Required)

	nodeIdsSchema := res.Schema["node_ids"]
	assert.Equal(t, schema.TypeList, nodeIdsSchema.Type)
	assert.True(t, nodeIdsSchema.Optional)

	resultOutputFileSchema := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFileSchema.Type)
	assert.True(t, resultOutputFileSchema.Optional)

	nodeSetSchema := res.Schema["node_set"]
	assert.Equal(t, schema.TypeList, nodeSetSchema.Type)
	assert.True(t, nodeSetSchema.Computed)

	elemRes := nodeSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "node_id")
	assert.Contains(t, elemRes.Schema, "capacity")
	assert.Contains(t, elemRes.Schema, "allocatable")
	assert.Contains(t, elemRes.Schema, "requests")
	assert.Contains(t, elemRes.Schema, "limits")
	assert.Contains(t, elemRes.Schema, "available")

	nodeIdSchema := elemRes.Schema["node_id"]
	assert.Equal(t, schema.TypeString, nodeIdSchema.Type)
	assert.True(t, nodeIdSchema.Computed)

	for _, blockName := range []string{"capacity", "allocatable", "requests", "limits", "available"} {
		blockSchema := elemRes.Schema[blockName]
		assert.Equal(t, schema.TypeList, blockSchema.Type, "block %s should be TypeList", blockName)
		assert.Equal(t, 1, blockSchema.MaxItems, "block %s should have MaxItems 1", blockName)
		assert.True(t, blockSchema.Computed, "block %s should be Computed", blockName)

		blockElemRes := blockSchema.Elem.(*schema.Resource)
		assert.Contains(t, blockElemRes.Schema, "cpu")
		assert.Contains(t, blockElemRes.Schema, "memory")
		assert.Contains(t, blockElemRes.Schema, "pods")

		assert.Equal(t, schema.TypeFloat, blockElemRes.Schema["cpu"].Type)
		assert.True(t, blockElemRes.Schema["cpu"].Computed)
		assert.Equal(t, schema.TypeFloat, blockElemRes.Schema["memory"].Type)
		assert.True(t, blockElemRes.Schema["memory"].Computed)
		assert.Equal(t, schema.TypeInt, blockElemRes.Schema["pods"].Type)
		assert.True(t, blockElemRes.Schema["pods"].Computed)
	}
}

func TestDbdcDbCustomClusterNodeResourcesDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcClusterNodeResourcesDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterNodeResources", func(request *dbdcv20201029.DescribeDBCustomClusterNodeResourcesRequest) (*dbdcv20201029.DescribeDBCustomClusterNodeResourcesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterNodeResourcesResponse()
		// nil NodeSet triggers the service-layer NonRetryableError
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterNodeResourcesResponseParams{
			NodeSet: nil,
		}
		return resp, nil
	})

	meta := newMockMetaDbdcClusterNodeResourcesDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodeResources()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "dbcc-nmtmsew8",
	})

	err := res.Read(d, meta)
	// When response NodeSet is empty slice, the service layer NonRetryableError should trigger
	assert.Error(t, err)
}
