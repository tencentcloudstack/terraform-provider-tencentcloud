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

type mockMetaDbdcClusterNodesDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDbdcClusterNodesDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDbdcClusterNodesDS{}

func newMockMetaDbdcClusterNodesDS() *mockMetaDbdcClusterNodesDS {
	return &mockMetaDbdcClusterNodesDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrCN(s string) *string {
	return &s
}

func ptrInt64CN(n int64) *int64 {
	return &n
}

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomClusterNodesDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomClusterNodesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcClusterNodesDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterNodes", func(request *dbdcv20201029.DescribeDBCustomClusterNodesRequest) (*dbdcv20201029.DescribeDBCustomClusterNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterNodesResponseParams{
			TotalCount: ptrInt64CN(2),
			NodeSet: []*dbdcv20201029.DBCustomClusterNode{
				{
					NodeId:      ptrStrCN("node-abc123"),
					NodeName:    ptrStrCN("test-node-1"),
					LanIP:       ptrStrCN("10.0.0.1"),
					SSHEndpoint: ptrStrCN("10.0.0.1:22"),
					Status:      ptrStrCN("Running"),
					Zone:        ptrStrCN("ap-guangzhou-3"),
					NodeType:    ptrStrCN("DB.AT5.32XLARGE512"),
				},
				{
					NodeId:      ptrStrCN("node-def456"),
					NodeName:    ptrStrCN("test-node-2"),
					LanIP:       ptrStrCN("10.0.0.2"),
					SSHEndpoint: ptrStrCN("10.0.0.2:22"),
					Status:      ptrStrCN("Creating"),
					Zone:        ptrStrCN("ap-guangzhou-3"),
					NodeType:    ptrStrCN("DB.AT5.64XLARGE1152"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcClusterNodesDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodes()
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
	assert.Equal(t, "test-node-1", node0["node_name"].(string))
	assert.Equal(t, "10.0.0.1", node0["lan_ip"].(string))
	assert.Equal(t, "10.0.0.1:22", node0["ssh_endpoint"].(string))
	assert.Equal(t, "Running", node0["status"].(string))
	assert.Equal(t, "ap-guangzhou-3", node0["zone"].(string))
	assert.Equal(t, "DB.AT5.32XLARGE512", node0["node_type"].(string))

	node1 := nodeSet[1].(map[string]interface{})
	assert.Equal(t, "node-def456", node1["node_id"].(string))
	assert.Equal(t, "test-node-2", node1["node_name"].(string))
	assert.Equal(t, "10.0.0.2", node1["lan_ip"].(string))
	assert.Equal(t, "10.0.0.2:22", node1["ssh_endpoint"].(string))
	assert.Equal(t, "Creating", node1["status"].(string))
	assert.Equal(t, "ap-guangzhou-3", node1["zone"].(string))
	assert.Equal(t, "DB.AT5.64XLARGE1152", node1["node_type"].(string))

	totalCount := d.Get("total_count").(int)
	assert.Equal(t, 2, totalCount)
}

func TestDbdcDbCustomClusterNodesDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodes()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "cluster_id")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "result_output_file")
	assert.Contains(t, res.Schema, "node_set")
	assert.Contains(t, res.Schema, "total_count")

	clusterIdSchema := res.Schema["cluster_id"]
	assert.Equal(t, schema.TypeString, clusterIdSchema.Type)
	assert.True(t, clusterIdSchema.Required)

	filtersSchema := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filtersSchema.Type)
	assert.True(t, filtersSchema.Optional)

	nodeSetSchema := res.Schema["node_set"]
	assert.Equal(t, schema.TypeList, nodeSetSchema.Type)
	assert.True(t, nodeSetSchema.Computed)

	elemRes := nodeSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "node_id")
	assert.Contains(t, elemRes.Schema, "node_name")
	assert.Contains(t, elemRes.Schema, "lan_ip")
	assert.Contains(t, elemRes.Schema, "ssh_endpoint")
	assert.Contains(t, elemRes.Schema, "status")
	assert.Contains(t, elemRes.Schema, "zone")
	assert.Contains(t, elemRes.Schema, "node_type")

	totalCountSchema := res.Schema["total_count"]
	assert.Equal(t, schema.TypeInt, totalCountSchema.Type)
	assert.True(t, totalCountSchema.Computed)
}

func TestDbdcDbCustomClusterNodesDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcClusterNodesDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterNodes", func(request *dbdcv20201029.DescribeDBCustomClusterNodesRequest) (*dbdcv20201029.DescribeDBCustomClusterNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterNodesResponseParams{
			TotalCount: ptrInt64CN(0),
			NodeSet:    []*dbdcv20201029.DBCustomClusterNode{},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcClusterNodesDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "dbcc-nmtmsew8",
	})

	err := res.Read(d, meta)
	// When response NodeSet is empty slice, the service layer NonRetryableError should trigger
	assert.Error(t, err)
}
