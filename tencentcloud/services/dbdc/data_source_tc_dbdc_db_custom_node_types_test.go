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

type mockMetaDbdcNodeTypesDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDbdcNodeTypesDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDbdcNodeTypesDS{}

func newMockMetaDbdcNodeTypesDS() *mockMetaDbdcNodeTypesDS {
	return &mockMetaDbdcNodeTypesDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrNT(s string) *string {
	return &s
}

func ptrUint64NT(n uint64) *uint64 {
	return &n
}

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomNodeTypesDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomNodeTypesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcNodeTypesDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeTypes", func(request *dbdcv20201029.DescribeDBCustomNodeTypesRequest) (*dbdcv20201029.DescribeDBCustomNodeTypesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeTypesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeTypesResponseParams{
			NodeTypeSet: []*dbdcv20201029.DBCustomNodeTypeInfo{
				{
					Zone:       ptrStrNT("ap-guangzhou-6"),
					NodeType:   ptrStrNT("DB.SA5.2XLARGE32"),
					NodeFamily: ptrStrNT("DB.SA5"),
					CPU:        ptrUint64NT(32),
					Memory:     ptrUint64NT(128),
					Status:     ptrStrNT("SELL"),
					SystemDiskTypes: []*string{
						ptrStrNT("CLOUD_BSSD"),
						ptrStrNT("CLOUD_HSSD"),
					},
					DataDiskTypes: []*string{
						ptrStrNT("CLOUD_BSSD"),
						ptrStrNT("LOCAL_NVME"),
					},
				},
				{
					Zone:       ptrStrNT("ap-guangzhou-3"),
					NodeType:   ptrStrNT("DB.AT5.8XLARGE128"),
					NodeFamily: ptrStrNT("DB.AT5"),
					CPU:        ptrUint64NT(32),
					Memory:     ptrUint64NT(128),
					Status:     ptrStrNT("SOLD_OUT"),
					SystemDiskTypes: []*string{
						ptrStrNT("CLOUD_HSSD"),
					},
					DataDiskTypes: []*string{
						ptrStrNT("CLOUD_HSSD"),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcNodeTypesDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	nodeTypeSet := d.Get("node_type_set").([]interface{})
	assert.Len(t, nodeTypeSet, 2)

	nodeType0 := nodeTypeSet[0].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou-6", nodeType0["zone"].(string))
	assert.Equal(t, "DB.SA5.2XLARGE32", nodeType0["node_type"].(string))
	assert.Equal(t, "DB.SA5", nodeType0["node_family"].(string))
	assert.Equal(t, 32, nodeType0["cpu"].(int))
	assert.Equal(t, 128, nodeType0["memory"].(int))
	assert.Equal(t, "SELL", nodeType0["status"].(string))

	systemDiskTypes0 := nodeType0["system_disk_types"].([]interface{})
	assert.Len(t, systemDiskTypes0, 2)
	assert.Equal(t, "CLOUD_BSSD", systemDiskTypes0[0].(string))
	assert.Equal(t, "CLOUD_HSSD", systemDiskTypes0[1].(string))

	dataDiskTypes0 := nodeType0["data_disk_types"].([]interface{})
	assert.Len(t, dataDiskTypes0, 2)
	assert.Equal(t, "CLOUD_BSSD", dataDiskTypes0[0].(string))
	assert.Equal(t, "LOCAL_NVME", dataDiskTypes0[1].(string))

	nodeType1 := nodeTypeSet[1].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou-3", nodeType1["zone"].(string))
	assert.Equal(t, "DB.AT5.8XLARGE128", nodeType1["node_type"].(string))
	assert.Equal(t, "DB.AT5", nodeType1["node_family"].(string))
	assert.Equal(t, 32, nodeType1["cpu"].(int))
	assert.Equal(t, 128, nodeType1["memory"].(int))
	assert.Equal(t, "SOLD_OUT", nodeType1["status"].(string))

	systemDiskTypes1 := nodeType1["system_disk_types"].([]interface{})
	assert.Len(t, systemDiskTypes1, 1)
	assert.Equal(t, "CLOUD_HSSD", systemDiskTypes1[0].(string))

	dataDiskTypes1 := nodeType1["data_disk_types"].([]interface{})
	assert.Len(t, dataDiskTypes1, 1)
	assert.Equal(t, "CLOUD_HSSD", dataDiskTypes1[0].(string))
}

func TestDbdcDbCustomNodeTypesDS_ReadWithNilDiskTypes(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcNodeTypesDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeTypes", func(request *dbdcv20201029.DescribeDBCustomNodeTypesRequest) (*dbdcv20201029.DescribeDBCustomNodeTypesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeTypesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeTypesResponseParams{
			NodeTypeSet: []*dbdcv20201029.DBCustomNodeTypeInfo{
				{
					Zone:       ptrStrNT("ap-guangzhou-6"),
					NodeType:   ptrStrNT("DB.SA5.2XLARGE32"),
					NodeFamily: ptrStrNT("DB.SA5"),
					CPU:        ptrUint64NT(32),
					Memory:     ptrUint64NT(128),
					Status:     ptrStrNT("SELL"),
					// SystemDiskTypes and DataDiskTypes are nil
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcNodeTypesDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	nodeTypeSet := d.Get("node_type_set").([]interface{})
	assert.Len(t, nodeTypeSet, 1)

	nodeType0 := nodeTypeSet[0].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou-6", nodeType0["zone"].(string))
	assert.Equal(t, "DB.SA5.2XLARGE32", nodeType0["node_type"].(string))
	assert.Equal(t, 32, nodeType0["cpu"].(int))
	assert.Equal(t, 128, nodeType0["memory"].(int))
	// system_disk_types and data_disk_types are nil in the API response;
	// Terraform SDK defaults them to empty lists
	systemDiskTypes0 := nodeType0["system_disk_types"].([]interface{})
	assert.Len(t, systemDiskTypes0, 0)
	dataDiskTypes0 := nodeType0["data_disk_types"].([]interface{})
	assert.Len(t, dataDiskTypes0, 0)
}

func TestDbdcDbCustomNodeTypesDS_ReadWithNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcNodeTypesDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeTypes", func(request *dbdcv20201029.DescribeDBCustomNodeTypesRequest) (*dbdcv20201029.DescribeDBCustomNodeTypesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeTypesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeTypesResponseParams{
			NodeTypeSet: []*dbdcv20201029.DBCustomNodeTypeInfo{
				{
					NodeType: ptrStrNT("DB.SA5.2XLARGE32"),
					Status:   ptrStrNT("SELL"),
					// Zone, NodeFamily, CPU, Memory, SystemDiskTypes, DataDiskTypes are nil
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcNodeTypesDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	nodeTypeSet := d.Get("node_type_set").([]interface{})
	assert.Len(t, nodeTypeSet, 1)

	nodeType0 := nodeTypeSet[0].(map[string]interface{})
	assert.Equal(t, "DB.SA5.2XLARGE32", nodeType0["node_type"].(string))
	assert.Equal(t, "SELL", nodeType0["status"].(string))
	// nil pointer fields default to zero values
	assert.Equal(t, "", nodeType0["zone"].(string))
	assert.Equal(t, "", nodeType0["node_family"].(string))
	assert.Equal(t, 0, nodeType0["cpu"].(int))
	assert.Equal(t, 0, nodeType0["memory"].(int))
}

func TestDbdcDbCustomNodeTypesDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcNodeTypesDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeTypes", func(request *dbdcv20201029.DescribeDBCustomNodeTypesRequest) (*dbdcv20201029.DescribeDBCustomNodeTypesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeTypesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeTypesResponseParams{
			NodeTypeSet: []*dbdcv20201029.DBCustomNodeTypeInfo{},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcNodeTypesDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	// When NodeTypeSet is empty (non-nil but len 0), the service helper guards
	// `result.Response.NodeTypeSet == nil` only; an empty slice passes the guard.
	// The Read still succeeds and sets an empty list. This case documents that behavior.
	assert.NoError(t, err)
	nodeTypeSet := d.Get("node_type_set").([]interface{})
	assert.Len(t, nodeTypeSet, 0)
}

func TestDbdcDbCustomNodeTypesDS_ReadWithNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcNodeTypesDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeTypes", func(request *dbdcv20201029.DescribeDBCustomNodeTypesRequest) (*dbdcv20201029.DescribeDBCustomNodeTypesResponse, error) {
		// Return a response with nil NodeTypeSet — simulates API returning null
		resp := dbdcv20201029.NewDescribeDBCustomNodeTypesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeTypesResponseParams{
			NodeTypeSet: nil,
		}
		return resp, nil
	})

	meta := newMockMetaDbdcNodeTypesDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	// When NodeTypeSet is nil, the service helper returns NonRetryableError and
	// the Read returns an error WITHOUT clearing the id (NonRetryableError behavior)
	assert.Error(t, err)
	// id is NOT set (still empty) because the helper returned an error before SetId
	assert.Empty(t, d.Id())
}

func TestDbdcDbCustomNodeTypesDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeTypes()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "node_type_set")
	assert.Contains(t, res.Schema, "result_output_file")
	assert.Contains(t, res.Schema, "filters")

	filtersSchema := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filtersSchema.Type)
	assert.True(t, filtersSchema.Optional)

	filtersElemRes := filtersSchema.Elem.(*schema.Resource)
	assert.Contains(t, filtersElemRes.Schema, "name")
	assert.Contains(t, filtersElemRes.Schema, "values")

	nameSchema := filtersElemRes.Schema["name"]
	assert.Equal(t, schema.TypeString, nameSchema.Type)
	assert.True(t, nameSchema.Required)

	valuesSchema := filtersElemRes.Schema["values"]
	assert.Equal(t, schema.TypeList, valuesSchema.Type)
	assert.True(t, valuesSchema.Required)

	nodeTypeSetSchema := res.Schema["node_type_set"]
	assert.Equal(t, schema.TypeList, nodeTypeSetSchema.Type)
	assert.True(t, nodeTypeSetSchema.Computed)

	elemRes := nodeTypeSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "zone")
	assert.Contains(t, elemRes.Schema, "node_type")
	assert.Contains(t, elemRes.Schema, "node_family")
	assert.Contains(t, elemRes.Schema, "cpu")
	assert.Contains(t, elemRes.Schema, "memory")
	assert.Contains(t, elemRes.Schema, "status")
	assert.Contains(t, elemRes.Schema, "system_disk_types")
	assert.Contains(t, elemRes.Schema, "data_disk_types")

	cpuSchema := elemRes.Schema["cpu"]
	assert.Equal(t, schema.TypeInt, cpuSchema.Type)
	assert.True(t, cpuSchema.Computed)

	memorySchema := elemRes.Schema["memory"]
	assert.Equal(t, schema.TypeInt, memorySchema.Type)
	assert.True(t, memorySchema.Computed)

	systemDiskTypesSchema := elemRes.Schema["system_disk_types"]
	assert.Equal(t, schema.TypeList, systemDiskTypesSchema.Type)
	assert.True(t, systemDiskTypesSchema.Computed)

	dataDiskTypesSchema := elemRes.Schema["data_disk_types"]
	assert.Equal(t, schema.TypeList, dataDiskTypesSchema.Type)
	assert.True(t, dataDiskTypesSchema.Computed)
}
