package dbdc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomClusterResourcesDS" -v -count=1 -gcflags="all=-l"

func ptrUint64(n uint64) *uint64 {
	return &n
}

func ptrFloat64(f float64) *float64 {
	return &f
}

func TestDbdcDbCustomClusterResourcesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterResourcesWithContext", func(ctx context.Context, request *dbdcv20201029.DescribeDBCustomClusterResourcesRequest) (*dbdcv20201029.DescribeDBCustomClusterResourcesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterResourcesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterResourcesResponseParams{
			NodeCount: ptrUint64(3),
			Capacity: &dbdcv20201029.MetaResource{
				Cpu:    ptrFloat64(96),
				Memory: ptrFloat64(384),
				Pods:   ptrUint64(240),
			},
			Allocatable: &dbdcv20201029.MetaResource{
				Cpu:    ptrFloat64(90),
				Memory: ptrFloat64(360),
				Pods:   ptrUint64(220),
			},
			Requests: &dbdcv20201029.MetaResource{
				Cpu:    ptrFloat64(45),
				Memory: ptrFloat64(120),
				Pods:   ptrUint64(60),
			},
			Limits: &dbdcv20201029.MetaResource{
				Cpu:    ptrFloat64(60),
				Memory: ptrFloat64(180),
				Pods:   ptrUint64(0),
			},
			Available: &dbdcv20201029.MetaResource{
				Cpu:    ptrFloat64(45),
				Memory: ptrFloat64(240),
				Pods:   ptrUint64(160),
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterResources()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "dbcc-abc123",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	nodeCount := d.Get("node_count").(int)
	assert.Equal(t, 3, nodeCount)

	capacity := d.Get("capacity").([]interface{})
	assert.Len(t, capacity, 1)
	capacityMap := capacity[0].(map[string]interface{})
	assert.Equal(t, 96.0, capacityMap["cpu"].(float64))
	assert.Equal(t, 384.0, capacityMap["memory"].(float64))
	assert.Equal(t, 240, capacityMap["pods"].(int))

	allocatable := d.Get("allocatable").([]interface{})
	assert.Len(t, allocatable, 1)
	allocatableMap := allocatable[0].(map[string]interface{})
	assert.Equal(t, 90.0, allocatableMap["cpu"].(float64))
	assert.Equal(t, 360.0, allocatableMap["memory"].(float64))
	assert.Equal(t, 220, allocatableMap["pods"].(int))

	requests := d.Get("requests").([]interface{})
	assert.Len(t, requests, 1)
	requestsMap := requests[0].(map[string]interface{})
	assert.Equal(t, 45.0, requestsMap["cpu"].(float64))
	assert.Equal(t, 120.0, requestsMap["memory"].(float64))
	assert.Equal(t, 60, requestsMap["pods"].(int))

	limits := d.Get("limits").([]interface{})
	assert.Len(t, limits, 1)
	limitsMap := limits[0].(map[string]interface{})
	assert.Equal(t, 60.0, limitsMap["cpu"].(float64))
	assert.Equal(t, 180.0, limitsMap["memory"].(float64))
	assert.Equal(t, 0, limitsMap["pods"].(int))

	available := d.Get("available").([]interface{})
	assert.Len(t, available, 1)
	availableMap := available[0].(map[string]interface{})
	assert.Equal(t, 45.0, availableMap["cpu"].(float64))
	assert.Equal(t, 240.0, availableMap["memory"].(float64))
	assert.Equal(t, 160, availableMap["pods"].(int))
}

func TestDbdcDbCustomClusterResourcesDS_ReadWithNilNestedObjects(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterResourcesWithContext", func(ctx context.Context, request *dbdcv20201029.DescribeDBCustomClusterResourcesRequest) (*dbdcv20201029.DescribeDBCustomClusterResourcesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterResourcesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomClusterResourcesResponseParams{
			NodeCount:   ptrUint64(2),
			Capacity:    nil,
			Allocatable: &dbdcv20201029.MetaResource{Cpu: ptrFloat64(10)},
			Requests:    nil,
			Limits:      nil,
			Available:   nil,
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterResources()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "dbcc-nested-nil",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	nodeCount := d.Get("node_count").(int)
	assert.Equal(t, 2, nodeCount)

	allocatable := d.Get("allocatable").([]interface{})
	assert.Len(t, allocatable, 1)
	allocatableMap := allocatable[0].(map[string]interface{})
	assert.Equal(t, 10.0, allocatableMap["cpu"].(float64))
}

func TestDbdcDbCustomClusterResourcesDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomClusterResourcesWithContext", func(ctx context.Context, request *dbdcv20201029.DescribeDBCustomClusterResourcesRequest) (*dbdcv20201029.DescribeDBCustomClusterResourcesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomClusterResourcesResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterResources()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "dbcc-empty",
	})

	err := res.Read(d, meta)
	// When response is nil, the datasource should return an error (NonRetryableError)
	assert.Error(t, err)
}

func TestDbdcDbCustomClusterResourcesDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomClusterResources()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "cluster_id")
	assert.Contains(t, res.Schema, "result_output_file")
	assert.Contains(t, res.Schema, "node_count")
	assert.Contains(t, res.Schema, "capacity")
	assert.Contains(t, res.Schema, "allocatable")
	assert.Contains(t, res.Schema, "requests")
	assert.Contains(t, res.Schema, "limits")
	assert.Contains(t, res.Schema, "available")

	clusterIdSchema := res.Schema["cluster_id"]
	assert.Equal(t, schema.TypeString, clusterIdSchema.Type)
	assert.True(t, clusterIdSchema.Required)

	nodeCountSchema := res.Schema["node_count"]
	assert.Equal(t, schema.TypeInt, nodeCountSchema.Type)
	assert.True(t, nodeCountSchema.Computed)

	capacitySchema := res.Schema["capacity"]
	assert.Equal(t, schema.TypeList, capacitySchema.Type)
	assert.True(t, capacitySchema.Computed)
	assert.Equal(t, 1, capacitySchema.MaxItems)
	capElemRes := capacitySchema.Elem.(*schema.Resource)
	assert.Contains(t, capElemRes.Schema, "cpu")
	assert.Contains(t, capElemRes.Schema, "memory")
	assert.Contains(t, capElemRes.Schema, "pods")
}
