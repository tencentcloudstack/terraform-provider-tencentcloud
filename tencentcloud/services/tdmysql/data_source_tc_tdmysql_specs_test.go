package tdmysql_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svctdmysql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql"
)

func ptrStrTdmysql(s string) *string   { return &s }
func ptrInt64Tdmysql(v int64) *int64   { return &v }
func ptrFloat64Tdmysql(v float64) *float64 { return &v }

// tdmysqlSpecsSchema returns the framework data source schema used to build
// the raw config/state value.
func tdmysqlSpecsSchema(t *testing.T) schema.Schema {
	d := svctdmysql.NewTdmysqlSpecsDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// tdmysqlSpecsHybridItemType builds the tftypes.Object type for a
// hybrid_node_specs element.
func tdmysqlSpecsHybridItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"spec_code":                 tftypes.String,
			"storage_node_cpu":          tftypes.Number,
			"storage_node_mem":          tftypes.Number,
			"storage_node_max_num":      tftypes.Number,
			"storage_node_max_disk":     tftypes.Number,
			"storage_node_min_num":      tftypes.Number,
			"storage_node_min_disk":     tftypes.Number,
			"storage_type":              tftypes.String,
			"storage_node_default_disk": tftypes.Number,
			"instance_mode":             tftypes.Set{ElementType: tftypes.String},
			"disk_type_category":        tftypes.String,
		},
	}
}

// tdmysqlSpecsServerlessCcuItemType builds the tftypes.Object type for a
// serverless_ccu_spec element.
func tdmysqlSpecsServerlessCcuItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"min_ccu": tftypes.Number,
			"max_ccu": tftypes.Set{ElementType: tftypes.Number},
		},
	}
}

// tdmysqlSpecsServerlessNodeNumItemType builds the tftypes.Object type for a
// serverless_node_num_spec element.
func tdmysqlSpecsServerlessNodeNumItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"min_node_num": tftypes.Number,
			"max_node_num": tftypes.Number,
		},
	}
}

// tdmysqlSpecsRootType builds the tftypes.Object type for the data source
// root, mirroring every schema attribute.
func tdmysqlSpecsRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"full_replications":        tftypes.Number,
			"is_exclusive_instance":    tftypes.Number,
			"zones":                    tftypes.Set{ElementType: tftypes.String},
			"instance_id":              tftypes.String,
			"hybrid_node_specs":        tftypes.List{ElementType: tdmysqlSpecsHybridItemType()},
			"serverless_ccu_spec":      tftypes.List{ElementType: tdmysqlSpecsServerlessCcuItemType()},
			"serverless_node_num_spec": tftypes.List{ElementType: tdmysqlSpecsServerlessNodeNumItemType()},
		},
	}
}

// buildTdmysqlSpecsRawValue builds the raw config value with the given query
// inputs (all optional).
func buildTdmysqlSpecsRawValue(fullReplications int64, isExclusive int64, zones []string, instanceId string) tftypes.Value {
	zoneValues := make([]tftypes.Value, 0, len(zones))
	for _, z := range zones {
		zoneValues = append(zoneValues, tftypes.NewValue(tftypes.String, z))
	}
	attrs := map[string]tftypes.Value{
		"full_replications":     tftypes.NewValue(tftypes.Number, fullReplications),
		"is_exclusive_instance": tftypes.NewValue(tftypes.Number, isExclusive),
		"zones":                 tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, zoneValues),
		"instance_id":           tftypes.NewValue(tftypes.String, instanceId),
		"hybrid_node_specs":     tftypes.NewValue(tftypes.List{ElementType: tdmysqlSpecsHybridItemType()}, nil),
		"serverless_ccu_spec":   tftypes.NewValue(tftypes.List{ElementType: tdmysqlSpecsServerlessCcuItemType()}, nil),
		"serverless_node_num_spec": tftypes.NewValue(tftypes.List{ElementType: tdmysqlSpecsServerlessNodeNumItemType()}, nil),
	}
	return tftypes.NewValue(tdmysqlSpecsRootType(), attrs)
}

// buildTdmysqlSpecsNullRawValue builds the raw config value with all query
// inputs unset.
func buildTdmysqlSpecsNullRawValue() tftypes.Value {
	attrs := map[string]tftypes.Value{
		"full_replications":        tftypes.NewValue(tftypes.Number, nil),
		"is_exclusive_instance":    tftypes.NewValue(tftypes.Number, nil),
		"zones":                    tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"instance_id":              tftypes.NewValue(tftypes.String, nil),
		"hybrid_node_specs":        tftypes.NewValue(tftypes.List{ElementType: tdmysqlSpecsHybridItemType()}, nil),
		"serverless_ccu_spec":      tftypes.NewValue(tftypes.List{ElementType: tdmysqlSpecsServerlessCcuItemType()}, nil),
		"serverless_node_num_spec": tftypes.NewValue(tftypes.List{ElementType: tdmysqlSpecsServerlessNodeNumItemType()}, nil),
	}
	return tftypes.NewValue(tdmysqlSpecsRootType(), attrs)
}

// setTdmysqlSpecsClient injects a mock connectivity client into the data
// source instance via the Configure method.
func setTdmysqlSpecsClient(d *svctdmysql.TdmysqlSpecsDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// newTdmysqlSpecsReadRequest builds a ReadRequest with the given config value
// and schema.
func newTdmysqlSpecsReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	return req
}

// newTdmysqlSpecsReadResponse builds a ReadResponse whose State carries the
// schema, so resp.State.Set can reflect the model into a raw value.
func newTdmysqlSpecsReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch
	return resp
}

// go test ./tencentcloud/services/tdmysql/ -run "TestTdmysqlSpecs" -v -count=1 -gcflags="all=-l"

// TestTdmysqlSpecs_MetadataAndSchema validates metadata and schema
// definition.
func TestTdmysqlSpecs_MetadataAndSchema(t *testing.T) {
	d := svctdmysql.NewTdmysqlSpecsDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_tdmysql_specs", metaResp.TypeName)

	sch := tdmysqlSpecsSchema(t)
	attrs := sch.Attributes

	// Optional inputs
	assert.True(t, attrs["full_replications"].(schema.Int64Attribute).Optional)
	assert.True(t, attrs["is_exclusive_instance"].(schema.Int64Attribute).Optional)
	assert.True(t, attrs["zones"].(schema.SetAttribute).Optional)
	assert.True(t, attrs["instance_id"].(schema.StringAttribute).Optional)

	// Computed outputs
	assert.True(t, attrs["hybrid_node_specs"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["serverless_ccu_spec"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["serverless_node_num_spec"].(schema.ListNestedAttribute).Computed)

	// Nested hybrid spec attributes
	nested := attrs["hybrid_node_specs"].(schema.ListNestedAttribute).NestedObject.Attributes
	for _, leaf := range []string{"spec_code", "storage_type", "disk_type_category"} {
		nestedAttr, ok := nested[leaf].(schema.StringAttribute)
		assert.True(t, ok, "nested attribute %s should be StringAttribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
	for _, leaf := range []string{"storage_node_cpu", "storage_node_mem", "storage_node_max_num", "storage_node_max_disk", "storage_node_min_num", "storage_node_min_disk", "storage_node_default_disk"} {
		nestedAttr, ok := nested[leaf].(schema.Int64Attribute)
		assert.True(t, ok, "nested attribute %s should be Int64Attribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
	assert.True(t, nested["instance_mode"].(schema.SetAttribute).Computed)

	// Nested serverless ccu attributes
	ccuNested := attrs["serverless_ccu_spec"].(schema.ListNestedAttribute).NestedObject.Attributes
	assert.True(t, ccuNested["min_ccu"].(schema.Float64Attribute).Computed)
	assert.True(t, ccuNested["max_ccu"].(schema.SetAttribute).Computed)

	// Nested serverless node num attributes
	nodeNumNested := attrs["serverless_node_num_spec"].(schema.ListNestedAttribute).NestedObject.Attributes
	assert.True(t, nodeNumNested["min_node_num"].(schema.Int64Attribute).Computed)
	assert.True(t, nodeNumNested["max_node_num"].(schema.Int64Attribute).Computed)
}

// TestTdmysqlSpecs_ReadBasic tests a successful read with a full response and
// asserts every leaf field is read back.
func TestTdmysqlSpecs_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &svctdmysql.TdmysqlService{}
	patches.ApplyMethodFunc(service, "DescribeTdmysqlSpecs", func(_ context.Context, request *tdmysqlv20211122.DescribeSpecsRequest) (*tdmysqlv20211122.DescribeSpecsResponseParams, error) {
		if request.FullReplications == nil || *request.FullReplications != 2 {
			t.Errorf("expected full_replications to be passed to the API request")
		}
		if request.IsExclusiveInstance == nil || *request.IsExclusiveInstance != 1 {
			t.Errorf("expected is_exclusive_instance to be passed to the API request")
		}
		if len(request.Zones) != 1 || *request.Zones[0] != "ap-guangzhou-1" {
			t.Errorf("expected zones to be passed to the API request")
		}
		if request.InstanceId == nil || *request.InstanceId != "tdsql-o2mtest" {
			t.Errorf("expected instance_id to be passed to the API request")
		}
		return &tdmysqlv20211122.DescribeSpecsResponseParams{
			HybridNodeSpecs: []*tdmysqlv20211122.StorageNodeSpec{
				{
					SpecCode:               ptrStrTdmysql("tcs-tcsp.ssd.2c8g"),
					StorageNodeCpu:         ptrInt64Tdmysql(2),
					StorageNodeMem:         ptrInt64Tdmysql(8),
					StorageNodeMaxNum:      ptrInt64Tdmysql(4),
					StorageNodeMaxDisk:     ptrInt64Tdmysql(3000),
					StorageNodeMinNum:      ptrInt64Tdmysql(2),
					StorageNodeMinDisk:     ptrInt64Tdmysql(100),
					StorageType:            ptrStrTdmysql("CLOUD_HSSD"),
					StorageNodeDefaultDisk: ptrInt64Tdmysql(200),
					InstanceMode:           []*string{ptrStrTdmysql("PREPAID"), ptrStrTdmysql("POSTPAID")},
					DiskTypeCategory:       ptrStrTdmysql("CLOUD_DISK"),
				},
			},
			ServerlessCcuSpec: []*tdmysqlv20211122.ServerlessCcu{
				{
					MinCcu: ptrFloat64Tdmysql(2.5),
					MaxCcu: []*float64{ptrFloat64Tdmysql(16), ptrFloat64Tdmysql(32)},
				},
			},
			ServerlessNodeNumSpec: &tdmysqlv20211122.ServerlessNodeNumSpec{
				MinNodeNum: ptrInt64Tdmysql(1),
				MaxNodeNum: ptrInt64Tdmysql(3),
			},
		}, nil
	})

	ds := svctdmysql.NewTdmysqlSpecsDataSource().(*svctdmysql.TdmysqlSpecsDataSource)
	setTdmysqlSpecsClient(ds, client)

	sch := tdmysqlSpecsSchema(t)
	raw := buildTdmysqlSpecsRawValue(2, 1, []string{"ap-guangzhou-1"}, "tdsql-o2mtest")
	req := newTdmysqlSpecsReadRequest(sch, raw)
	resp := newTdmysqlSpecsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svctdmysql.TdmysqlSpecsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, int64(2), result.FullReplications.ValueInt64())
	assert.Equal(t, int64(1), result.IsExclusiveInstance.ValueInt64())
	assert.Equal(t, "tdsql-o2mtest", result.InstanceId.ValueString())

	var hybrid []svctdmysql.TdmysqlSpecsHybridNodeSpecModel
	resp.Diagnostics.Append(result.HybridNodeSpecs.ElementsAs(context.Background(), &hybrid, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, hybrid, 1)

	assert.Equal(t, "tcs-tcsp.ssd.2c8g", hybrid[0].SpecCode.ValueString())
	assert.Equal(t, int64(2), hybrid[0].StorageNodeCpu.ValueInt64())
	assert.Equal(t, int64(8), hybrid[0].StorageNodeMem.ValueInt64())
	assert.Equal(t, int64(4), hybrid[0].StorageNodeMaxNum.ValueInt64())
	assert.Equal(t, int64(3000), hybrid[0].StorageNodeMaxDisk.ValueInt64())
	assert.Equal(t, int64(2), hybrid[0].StorageNodeMinNum.ValueInt64())
	assert.Equal(t, int64(100), hybrid[0].StorageNodeMinDisk.ValueInt64())
	assert.Equal(t, "CLOUD_HSSD", hybrid[0].StorageType.ValueString())
	assert.Equal(t, int64(200), hybrid[0].StorageNodeDefaultDisk.ValueInt64())
	assert.Equal(t, "CLOUD_DISK", hybrid[0].DiskTypeCategory.ValueString())

	var modes []string
	resp.Diagnostics.Append(hybrid[0].InstanceMode.ElementsAs(context.Background(), &modes, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"PREPAID", "POSTPAID"}, modes)

	var ccu []svctdmysql.TdmysqlSpecsServerlessCcuModel
	resp.Diagnostics.Append(result.ServerlessCcuSpec.ElementsAs(context.Background(), &ccu, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, ccu, 1)
	assert.Equal(t, 2.5, ccu[0].MinCcu.ValueFloat64())

	var maxCcu []float64
	resp.Diagnostics.Append(ccu[0].MaxCcu.ElementsAs(context.Background(), &maxCcu, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []float64{16, 32}, maxCcu)

	var nodeNum []svctdmysql.TdmysqlSpecsServerlessNodeNumModel
	resp.Diagnostics.Append(result.ServerlessNodeNumSpec.ElementsAs(context.Background(), &nodeNum, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, nodeNum, 1)
	assert.Equal(t, int64(1), nodeNum[0].MinNodeNum.ValueInt64())
	assert.Equal(t, int64(3), nodeNum[0].MaxNodeNum.ValueInt64())
}

// TestTdmysqlSpecs_ReadEmpty tests a successful read where the API returns no
// specs: every output list must be empty (not null) and no error diagnostics.
func TestTdmysqlSpecs_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &svctdmysql.TdmysqlService{}
	patches.ApplyMethodFunc(service, "DescribeTdmysqlSpecs", func(_ context.Context, request *tdmysqlv20211122.DescribeSpecsRequest) (*tdmysqlv20211122.DescribeSpecsResponseParams, error) {
		if request.FullReplications != nil || request.IsExclusiveInstance != nil || request.Zones != nil || request.InstanceId != nil {
			t.Errorf("expected unset inputs to be omitted from the API request")
		}
		return &tdmysqlv20211122.DescribeSpecsResponseParams{}, nil
	})

	ds := svctdmysql.NewTdmysqlSpecsDataSource().(*svctdmysql.TdmysqlSpecsDataSource)
	setTdmysqlSpecsClient(ds, client)

	sch := tdmysqlSpecsSchema(t)
	raw := buildTdmysqlSpecsNullRawValue()
	req := newTdmysqlSpecsReadRequest(sch, raw)
	resp := newTdmysqlSpecsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svctdmysql.TdmysqlSpecsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.HybridNodeSpecs.IsNull(), "hybrid_node_specs should be an empty list, not null")
	assert.False(t, result.ServerlessCcuSpec.IsNull(), "serverless_ccu_spec should be an empty list, not null")
	assert.False(t, result.ServerlessNodeNumSpec.IsNull(), "serverless_node_num_spec should be an empty list, not null")

	var hybrid []svctdmysql.TdmysqlSpecsHybridNodeSpecModel
	resp.Diagnostics.Append(result.HybridNodeSpecs.ElementsAs(context.Background(), &hybrid, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, hybrid, 0)

	var ccu []svctdmysql.TdmysqlSpecsServerlessCcuModel
	resp.Diagnostics.Append(result.ServerlessCcuSpec.ElementsAs(context.Background(), &ccu, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, ccu, 0)

	var nodeNum []svctdmysql.TdmysqlSpecsServerlessNodeNumModel
	resp.Diagnostics.Append(result.ServerlessNodeNumSpec.ElementsAs(context.Background(), &nodeNum, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, nodeNum, 0)
}

// TestTdmysqlSpecs_ReadApiError tests that an API failure surfaces as an error
// diagnostic.
func TestTdmysqlSpecs_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &svctdmysql.TdmysqlService{}
	patches.ApplyMethodFunc(service, "DescribeTdmysqlSpecs", func(_ context.Context, request *tdmysqlv20211122.DescribeSpecsRequest) (*tdmysqlv20211122.DescribeSpecsResponseParams, error) {
		return nil, assert.AnError
	})

	ds := svctdmysql.NewTdmysqlSpecsDataSource().(*svctdmysql.TdmysqlSpecsDataSource)
	setTdmysqlSpecsClient(ds, client)

	sch := tdmysqlSpecsSchema(t)
	raw := buildTdmysqlSpecsNullRawValue()
	req := newTdmysqlSpecsReadRequest(sch, raw)
	resp := newTdmysqlSpecsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}
