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

type mockMetaDbdcDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDbdcDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDbdcDS{}

func newMockMetaDbdcDS() *mockMetaDbdcDS {
	return &mockMetaDbdcDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStr(s string) *string {
	return &s
}

func ptrInt64(n int64) *int64 {
	return &n
}

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomClustersDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomClustersDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusters", func(request *dbdcv20201029.DescribeDBCustomClustersRequest) (*dbdcv20201029.DescribeDBCustomClustersResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClustersResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClustersResponseParams{
			TotalCount: ptrInt64(2),
			ClusterSet: []*dbdcv20201029.DBCustomCluster{
				{
					ClusterId:          ptrStr("cluster-abc123"),
					ClusterName:        ptrStr("test-cluster-1"),
					Region:             ptrStr("ap-guangzhou"),
					ClusterLevel:       ptrStr("L500"),
					ClusterStatus:      ptrStr("Running"),
					ClusterVersion:     ptrStr("1.0"),
					ClusterNodeNum:     ptrInt64(3),
					ClusterDescription: ptrStr("test cluster description"),
					CreatedTime:        ptrStr("2024-01-01 00:00:00"),
					Tags: []*dbdcv20201029.Tag{
						{Key: ptrStr("env"), Value: ptrStr("production")},
					},
				},
				{
					ClusterId:     ptrStr("cluster-def456"),
					ClusterName:   ptrStr("test-cluster-2"),
					Region:        ptrStr("ap-shanghai"),
					ClusterStatus: ptrStr("Creating"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusters()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	clusterSet := d.Get("cluster_set").([]interface{})
	assert.Len(t, clusterSet, 2)

	cluster0 := clusterSet[0].(map[string]interface{})
	assert.Equal(t, "cluster-abc123", cluster0["cluster_id"].(string))
	assert.Equal(t, "test-cluster-1", cluster0["cluster_name"].(string))
	assert.Equal(t, "ap-guangzhou", cluster0["region"].(string))
	assert.Equal(t, "L500", cluster0["cluster_level"].(string))
	assert.Equal(t, "Running", cluster0["cluster_status"].(string))
	assert.Equal(t, "1.0", cluster0["cluster_version"].(string))
	assert.Equal(t, 3, cluster0["cluster_node_num"].(int))
	assert.Equal(t, "test cluster description", cluster0["cluster_description"].(string))
	assert.Equal(t, "2024-01-01 00:00:00", cluster0["created_time"].(string))

	tags0 := cluster0["tags"].([]interface{})
	assert.Len(t, tags0, 1)
	tagMap0 := tags0[0].(map[string]interface{})
	assert.Equal(t, "env", tagMap0["key"].(string))
	assert.Equal(t, "production", tagMap0["value"].(string))

	totalCount := d.Get("total_count").(int)
	assert.Equal(t, 2, totalCount)
}

func TestDbdcDbCustomClustersDS_ReadWithNilTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusters", func(request *dbdcv20201029.DescribeDBCustomClustersRequest) (*dbdcv20201029.DescribeDBCustomClustersResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClustersResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClustersResponseParams{
			TotalCount: ptrInt64(1),
			ClusterSet: []*dbdcv20201029.DBCustomCluster{
				{
					ClusterId:     ptrStr("cluster-nil-tags"),
					ClusterName:   ptrStr("cluster-no-tags"),
					ClusterStatus: ptrStr("Running"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusters()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clusterSet := d.Get("cluster_set").([]interface{})
	assert.Len(t, clusterSet, 1)

	cluster0 := clusterSet[0].(map[string]interface{})
	assert.Equal(t, "cluster-nil-tags", cluster0["cluster_id"].(string))
	// Tags field is nil in the API response, should not have tags in output
	tagsField := cluster0["tags"]
	if tagsField != nil {
		tagsList := tagsField.([]interface{})
		assert.Len(t, tagsList, 0)
	}
}

func TestDbdcDbCustomClustersDS_ReadWithClusterIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusters", func(request *dbdcv20201029.DescribeDBCustomClustersRequest) (*dbdcv20201029.DescribeDBCustomClustersResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClustersResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClustersResponseParams{
			TotalCount: ptrInt64(1),
			ClusterSet: []*dbdcv20201029.DBCustomCluster{
				{
					ClusterId:     ptrStr("cluster-abc123"),
					ClusterName:   ptrStr("test-cluster-1"),
					ClusterStatus: ptrStr("Running"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusters()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_ids": []interface{}{"cluster-abc123"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	clusterSet := d.Get("cluster_set").([]interface{})
	assert.Len(t, clusterSet, 1)
}

func TestDbdcDbCustomClustersDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusters", func(request *dbdcv20201029.DescribeDBCustomClustersRequest) (*dbdcv20201029.DescribeDBCustomClustersResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClustersResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClustersResponseParams{
			TotalCount: ptrInt64(0),
			ClusterSet: []*dbdcv20201029.DBCustomCluster{},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusters()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	// When response is empty, the datasource should return an error (NonRetryableError)
	assert.Error(t, err)
}

func TestDbdcDbCustomClustersDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusters()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "cluster_ids")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "tags")
	assert.Contains(t, res.Schema, "cluster_set")
	assert.Contains(t, res.Schema, "total_count")
	assert.Contains(t, res.Schema, "result_output_file")

	clusterIdsSchema := res.Schema["cluster_ids"]
	assert.Equal(t, schema.TypeList, clusterIdsSchema.Type)
	assert.True(t, clusterIdsSchema.Optional)

	clusterSetSchema := res.Schema["cluster_set"]
	assert.Equal(t, schema.TypeList, clusterSetSchema.Type)
	assert.True(t, clusterSetSchema.Computed)

	elemRes := clusterSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "cluster_id")
	assert.Contains(t, elemRes.Schema, "cluster_name")
	assert.Contains(t, elemRes.Schema, "region")
	assert.Contains(t, elemRes.Schema, "cluster_level")
	assert.Contains(t, elemRes.Schema, "cluster_status")
	assert.Contains(t, elemRes.Schema, "cluster_version")
	assert.Contains(t, elemRes.Schema, "cluster_node_num")
	assert.Contains(t, elemRes.Schema, "cluster_description")
	assert.Contains(t, elemRes.Schema, "created_time")
	assert.Contains(t, elemRes.Schema, "tags")

	totalCountSchema := res.Schema["total_count"]
	assert.Equal(t, schema.TypeInt, totalCountSchema.Type)
	assert.True(t, totalCountSchema.Computed)
}
