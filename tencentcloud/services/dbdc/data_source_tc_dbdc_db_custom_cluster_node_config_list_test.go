package dbdc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

type mockMetaDbdcClusterNodeConfigListDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDbdcClusterNodeConfigListDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDbdcClusterNodeConfigListDS{}

func newMockMetaDbdcClusterNodeConfigListDS() *mockMetaDbdcClusterNodeConfigListDS {
	return &mockMetaDbdcClusterNodeConfigListDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrCNCL(s string) *string {
	return &s
}

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomClusterNodeConfigListDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomClusterNodeConfigListDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcClusterNodeConfigListDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterNodeConfigWithContext", func(ctx context.Context, request *dbdcv20201029.DescribeDBCustomClusterNodeConfigRequest) (*dbdcv20201029.DescribeDBCustomClusterNodeConfigResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterNodeConfigResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterNodeConfigResponseParams{
			NodeSet: []*dbdcv20201029.DBCustomClusterNodeConfig{
				{
					NodeId: ptrStrCNCL("node-abc123"),
					Labels: []*dbdcv20201029.Label{
						{
							Key:   ptrStrCNCL("app"),
							Value: ptrStrCNCL("mysql"),
						},
						{
							Key:   ptrStrCNCL("tier"),
							Value: ptrStrCNCL("db"),
						},
					},
					Taints: []*dbdcv20201029.Taint{
						{
							Key:    ptrStrCNCL("dedicated"),
							Value:  ptrStrCNCL("dbdc"),
							Effect: ptrStrCNCL("NoSchedule"),
						},
					},
				},
				{
					NodeId: ptrStrCNCL("node-def456"),
					Labels: []*dbdcv20201029.Label{},
					Taints: []*dbdcv20201029.Taint{
						{
							Key:    ptrStrCNCL("maintenance"),
							Value:  ptrStrCNCL(""),
							Effect: ptrStrCNCL("NoExecute"),
						},
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcClusterNodeConfigListDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodeConfigList()
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

	labels0 := node0["labels"].([]interface{})
	assert.Len(t, labels0, 2)
	label0 := labels0[0].(map[string]interface{})
	assert.Equal(t, "app", label0["key"].(string))
	assert.Equal(t, "mysql", label0["value"].(string))
	label1 := labels0[1].(map[string]interface{})
	assert.Equal(t, "tier", label1["key"].(string))
	assert.Equal(t, "db", label1["value"].(string))

	taints0 := node0["taints"].([]interface{})
	assert.Len(t, taints0, 1)
	taint0 := taints0[0].(map[string]interface{})
	assert.Equal(t, "dedicated", taint0["key"].(string))
	assert.Equal(t, "dbdc", taint0["value"].(string))
	assert.Equal(t, "NoSchedule", taint0["effect"].(string))

	node1 := nodeSet[1].(map[string]interface{})
	assert.Equal(t, "node-def456", node1["node_id"].(string))
	assert.Len(t, node1["labels"].([]interface{}), 0)

	taints1 := node1["taints"].([]interface{})
	assert.Len(t, taints1, 1)
	taint1 := taints1[0].(map[string]interface{})
	assert.Equal(t, "maintenance", taint1["key"].(string))
	assert.Equal(t, "", taint1["value"].(string))
	assert.Equal(t, "NoExecute", taint1["effect"].(string))
}

func TestDbdcDbCustomClusterNodeConfigListDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodeConfigList()

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
	assert.Contains(t, elemRes.Schema, "labels")
	assert.Contains(t, elemRes.Schema, "taints")

	nodeIdElem := elemRes.Schema["node_id"]
	assert.Equal(t, schema.TypeString, nodeIdElem.Type)
	assert.True(t, nodeIdElem.Computed)

	labelsElem := elemRes.Schema["labels"]
	assert.Equal(t, schema.TypeList, labelsElem.Type)
	assert.True(t, labelsElem.Computed)
	labelRes := labelsElem.Elem.(*schema.Resource)
	assert.Contains(t, labelRes.Schema, "key")
	assert.Contains(t, labelRes.Schema, "value")

	taintsElem := elemRes.Schema["taints"]
	assert.Equal(t, schema.TypeList, taintsElem.Type)
	assert.True(t, taintsElem.Computed)
	taintRes := taintsElem.Elem.(*schema.Resource)
	assert.Contains(t, taintRes.Schema, "key")
	assert.Contains(t, taintRes.Schema, "value")
	assert.Contains(t, taintRes.Schema, "effect")
}

func TestDbdcDbCustomClusterNodeConfigListDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcClusterNodeConfigListDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterNodeConfigWithContext", func(ctx context.Context, request *dbdcv20201029.DescribeDBCustomClusterNodeConfigRequest) (*dbdcv20201029.DescribeDBCustomClusterNodeConfigResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterNodeConfigResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterNodeConfigResponseParams{
			NodeSet: []*dbdcv20201029.DBCustomClusterNodeConfig{},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcClusterNodeConfigListDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodeConfigList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "dbcc-nmtmsew8",
	})

	err := res.Read(d, meta)
	// When response NodeSet is empty slice, the service layer NonRetryableError should trigger
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}
