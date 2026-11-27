package tdmysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svtdmysql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql"
)

func tdmysqlSpecsStr(s string) *string     { return &s }
func tdmysqlSpecsInt64(v int64) *int64     { return &v }
func tdmysqlSpecsFloat(v float64) *float64 { return &v }

// tdmysqlSpecsSchema returns the framework data source schema used to build
// the raw config/state value.
func tdmysqlSpecsSchema(t *testing.T) schema.Schema {
	d := svtdmysql.NewTdmysqlSpecsDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// tdmysqlSpecsRootType builds the tftypes.Object type for the data source
// root, mirroring every schema attribute.
func tdmysqlSpecsRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                       tftypes.String,
			"full_replications":        tftypes.Number,
			"is_exclusive_instance":    tftypes.Number,
			"zones":                    tftypes.List{ElementType: tftypes.String},
			"instance_id":              tftypes.String,
			"hybrid_node_specs":        tftypes.List{ElementType: tdmysqlSpecsHybridItemType()},
			"serverless_ccu_spec":      tftypes.List{ElementType: tdmysqlSpecsServerlessCcuItemType()},
			"serverless_node_num_spec": tftypes.List{ElementType: tdmysqlSpecsServerlessNodeNumItemType()},
		},
	}
}

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
			"instance_mode":             tftypes.List{ElementType: tftypes.String},
			"disk_type_category":        tftypes.String,
		},
	}
}

func tdmysqlSpecsServerlessCcuItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"min_ccu": tftypes.Number,
			"max_ccu": tftypes.List{ElementType: tftypes.Number},
		},
	}
}

func tdmysqlSpecsServerlessNodeNumItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"min_node_num": tftypes.Number,
			"max_node_num": tftypes.Number,
		},
	}
}

// buildTdmysqlSpecsRawValue builds the raw config value with the given query
// inputs; every other attribute is unset (nil).
func buildTdmysqlSpecsRawValue(fullReplications int64) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":                       tftypes.NewValue(tftypes.String, nil),
		"full_replications":        tftypes.NewValue(tftypes.Number, fullReplications),
		"is_exclusive_instance":    tftypes.NewValue(tftypes.Number, nil),
		"zones":                    tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
		"instance_id":              tftypes.NewValue(tftypes.String, nil),
		"hybrid_node_specs":        tftypes.NewValue(tftypes.List{ElementType: tdmysqlSpecsHybridItemType()}, nil),
		"serverless_ccu_spec":      tftypes.NewValue(tftypes.List{ElementType: tdmysqlSpecsServerlessCcuItemType()}, nil),
		"serverless_node_num_spec": tftypes.NewValue(tftypes.List{ElementType: tdmysqlSpecsServerlessNodeNumItemType()}, nil),
	}

	return tftypes.NewValue(tdmysqlSpecsRootType(), attrs)
}

// setTdmysqlSpecsClient injects a mock connectivity client into the data
// source instance via the Configure method.
func setTdmysqlSpecsClient(d *svtdmysql.TdmysqlSpecsDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

func newTdmysqlSpecsReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch

	return req
}

func newTdmysqlSpecsReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	return resp
}

// TestTdmysqlSpecs_MetadataAndSchema validates metadata and schema definition.
func TestTdmysqlSpecs_MetadataAndSchema(t *testing.T) {
	d := svtdmysql.NewTdmysqlSpecsDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_tdmysql_specs", metaResp.TypeName)

	sch := tdmysqlSpecsSchema(t)
	attrs := sch.Attributes

	assert.True(t, attrs["full_replications"].(schema.Int64Attribute).Optional)
	assert.True(t, attrs["is_exclusive_instance"].(schema.Int64Attribute).Optional)
	assert.True(t, attrs["zones"].(schema.ListAttribute).Optional)
	assert.True(t, attrs["instance_id"].(schema.StringAttribute).Optional)
	assert.True(t, attrs["hybrid_node_specs"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["serverless_ccu_spec"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["serverless_node_num_spec"].(schema.ListNestedAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "total_count")
}

// TestTdmysqlSpecs_ReadBasic tests a successful read with fake spec records
// and asserts every leaf field is read back.
func TestTdmysqlSpecs_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeSpecs", func(request *tdmysqlv20211122.DescribeSpecsRequest) (*tdmysqlv20211122.DescribeSpecsResponse, error) {
		if request.FullReplications == nil || *request.FullReplications != int64(3) {
			t.Errorf("expected full_replications to be passed to the API request")
		}

		return &tdmysqlv20211122.DescribeSpecsResponse{Response: &tdmysqlv20211122.DescribeSpecsResponseParams{
			HybridNodeSpecs: []*tdmysqlv20211122.StorageNodeSpec{
				{
					SpecCode:               tdmysqlSpecsStr("tcl.c4.xl.mock"),
					StorageNodeCpu:         tdmysqlSpecsInt64(4),
					StorageNodeMem:         tdmysqlSpecsInt64(16),
					StorageNodeMaxNum:      tdmysqlSpecsInt64(6),
					StorageNodeMaxDisk:     tdmysqlSpecsInt64(4000),
					StorageNodeMinNum:      tdmysqlSpecsInt64(2),
					StorageNodeMinDisk:     tdmysqlSpecsInt64(100),
					StorageType:            tdmysqlSpecsStr("CLOUD_HSSD"),
					StorageNodeDefaultDisk: tdmysqlSpecsInt64(200),
					InstanceMode:           []*string{tdmysqlSpecsStr("PREPAID"), tdmysqlSpecsStr("POSTPAID")},
					DiskTypeCategory:       tdmysqlSpecsStr("CLOUD_DISK"),
				},
			},
			ServerlessCcuSpec: []*tdmysqlv20211122.ServerlessCcu{
				{
					MinCcu: tdmysqlSpecsFloat(1.5),
					MaxCcu: []*float64{tdmysqlSpecsFloat(2.5), tdmysqlSpecsFloat(4.0)},
				},
			},
			ServerlessNodeNumSpec: &tdmysqlv20211122.ServerlessNodeNumSpec{
				MinNodeNum: tdmysqlSpecsInt64(1),
				MaxNodeNum: tdmysqlSpecsInt64(4),
			},
			RequestId: tdmysqlSpecsStr("fake-request-id"),
		}}, nil
	})

	ds := svtdmysql.NewTdmysqlSpecsDataSource().(*svtdmysql.TdmysqlSpecsDataSource)
	setTdmysqlSpecsClient(ds, client)

	sch := tdmysqlSpecsSchema(t)
	raw := buildTdmysqlSpecsRawValue(3)
	req := newTdmysqlSpecsReadRequest(sch, raw)
	resp := newTdmysqlSpecsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svtdmysql.TdmysqlSpecsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())

	// Hybrid node specs.
	assert.False(t, result.HybridNodeSpecs.IsNull())
	var hybridItems []svtdmysql.TdmysqlSpecsHybridNodeSpecItemModel
	resp.Diagnostics.Append(result.HybridNodeSpecs.ElementsAs(context.Background(), &hybridItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, hybridItems, 1)
	assert.Equal(t, "tcl.c4.xl.mock", hybridItems[0].SpecCode.ValueString())
	assert.Equal(t, int64(4), hybridItems[0].StorageNodeCpu.ValueInt64())
	assert.Equal(t, int64(16), hybridItems[0].StorageNodeMem.ValueInt64())
	assert.Equal(t, int64(6), hybridItems[0].StorageNodeMaxNum.ValueInt64())
	assert.Equal(t, int64(4000), hybridItems[0].StorageNodeMaxDisk.ValueInt64())
	assert.Equal(t, int64(2), hybridItems[0].StorageNodeMinNum.ValueInt64())
	assert.Equal(t, int64(100), hybridItems[0].StorageNodeMinDisk.ValueInt64())
	assert.Equal(t, "CLOUD_HSSD", hybridItems[0].StorageType.ValueString())
	assert.Equal(t, int64(200), hybridItems[0].StorageNodeDefaultDisk.ValueInt64())
	assert.Equal(t, "CLOUD_DISK", hybridItems[0].DiskTypeCategory.ValueString())

	var instanceModes []string
	resp.Diagnostics.Append(hybridItems[0].InstanceMode.ElementsAs(context.Background(), &instanceModes, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"PREPAID", "POSTPAID"}, instanceModes)

	// Serverless CCU specs.
	assert.False(t, result.ServerlessCcuSpec.IsNull())
	var ccuItems []svtdmysql.TdmysqlSpecsServerlessCcuItemModel
	resp.Diagnostics.Append(result.ServerlessCcuSpec.ElementsAs(context.Background(), &ccuItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, ccuItems, 1)
	assert.Equal(t, 1.5, ccuItems[0].MinCcu.ValueFloat64())

	var maxCcu []float64
	resp.Diagnostics.Append(ccuItems[0].MaxCcu.ElementsAs(context.Background(), &maxCcu, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []float64{2.5, 4.0}, maxCcu)

	// Serverless node number spec (zero-or-one element list).
	assert.False(t, result.ServerlessNodeNumSpec.IsNull())
	var nodeNumItems []svtdmysql.TdmysqlSpecsServerlessNodeNumItemModel
	resp.Diagnostics.Append(result.ServerlessNodeNumSpec.ElementsAs(context.Background(), &nodeNumItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, nodeNumItems, 1)
	assert.Equal(t, int64(1), nodeNumItems[0].MinNodeNum.ValueInt64())
	assert.Equal(t, int64(4), nodeNumItems[0].MaxNodeNum.ValueInt64())
}

// TestTdmysqlSpecs_ReadEmpty covers an empty result: a legitimate answer, not
// a failure — the data source must succeed with empty lists.
func TestTdmysqlSpecs_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeSpecs", func(request *tdmysqlv20211122.DescribeSpecsRequest) (*tdmysqlv20211122.DescribeSpecsResponse, error) {
		return &tdmysqlv20211122.DescribeSpecsResponse{Response: &tdmysqlv20211122.DescribeSpecsResponseParams{}}, nil
	})

	ds := svtdmysql.NewTdmysqlSpecsDataSource().(*svtdmysql.TdmysqlSpecsDataSource)
	setTdmysqlSpecsClient(ds, client)

	sch := tdmysqlSpecsSchema(t)
	raw := buildTdmysqlSpecsRawValue(0)
	req := newTdmysqlSpecsReadRequest(sch, raw)
	resp := newTdmysqlSpecsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svtdmysql.TdmysqlSpecsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.HybridNodeSpecs.IsNull())
	var hybridItems []svtdmysql.TdmysqlSpecsHybridNodeSpecItemModel
	resp.Diagnostics.Append(result.HybridNodeSpecs.ElementsAs(context.Background(), &hybridItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, hybridItems, 0)

	assert.False(t, result.ServerlessCcuSpec.IsNull())
	var ccuItems []svtdmysql.TdmysqlSpecsServerlessCcuItemModel
	resp.Diagnostics.Append(result.ServerlessCcuSpec.ElementsAs(context.Background(), &ccuItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, ccuItems, 0)

	assert.False(t, result.ServerlessNodeNumSpec.IsNull())
	var nodeNumItems []svtdmysql.TdmysqlSpecsServerlessNodeNumItemModel
	resp.Diagnostics.Append(result.ServerlessNodeNumSpec.ElementsAs(context.Background(), &nodeNumItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, nodeNumItems, 0)
}

// TestTdmysqlSpecs_ReadApiError covers an API-level failure: the error must
// surface through diagnostics instead of silently persisting empty state.
func TestTdmysqlSpecs_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeSpecs", func(request *tdmysqlv20211122.DescribeSpecsRequest) (*tdmysqlv20211122.DescribeSpecsResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svtdmysql.NewTdmysqlSpecsDataSource().(*svtdmysql.TdmysqlSpecsDataSource)
	setTdmysqlSpecsClient(ds, client)

	sch := tdmysqlSpecsSchema(t)
	raw := buildTdmysqlSpecsRawValue(3)
	req := newTdmysqlSpecsReadRequest(sch, raw)
	resp := newTdmysqlSpecsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}
