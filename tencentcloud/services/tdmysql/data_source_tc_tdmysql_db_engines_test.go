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

func tdmysqlDbEnginesStr(s string) *string { return &s }
func tdmysqlDbEnginesBool(v bool) *bool    { return &v }

// tdmysqlDbEnginesSchema returns the framework data source schema used to
// build the raw config/state value.
func tdmysqlDbEnginesSchema(t *testing.T) schema.Schema {
	d := svtdmysql.NewTdmysqlDbEnginesDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// tdmysqlDbEnginesRootType builds the tftypes.Object type for the data source
// root, mirroring every schema attribute.
func tdmysqlDbEnginesRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":    tftypes.String,
			"items": tftypes.List{ElementType: tdmysqlDbEnginesItemType()},
		},
	}
}

// tdmysqlDbEnginesItemType builds the tftypes.Object type for an items element.
func tdmysqlDbEnginesItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"type":                           tftypes.String,
			"version":                        tftypes.String,
			"name":                           tftypes.String,
			"description":                    tftypes.String,
			"new":                            tftypes.Bool,
			"sql_mode":                       tftypes.List{ElementType: tftypes.String},
			"is_support_param_template":      tftypes.Bool,
			"is_support_serverless":          tftypes.Bool,
			"is_support_tde":                 tftypes.Bool,
			"is_support_password_complexity": tftypes.Bool,
		},
	}
}

// buildTdmysqlDbEnginesRawValue builds the raw config value; every attribute
// is unset (nil) — the data source takes no query inputs.
func buildTdmysqlDbEnginesRawValue() tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":    tftypes.NewValue(tftypes.String, nil),
		"items": tftypes.NewValue(tftypes.List{ElementType: tdmysqlDbEnginesItemType()}, nil),
	}

	return tftypes.NewValue(tdmysqlDbEnginesRootType(), attrs)
}

// setTdmysqlDbEnginesClient injects a mock connectivity client into the data
// source instance via the Configure method.
func setTdmysqlDbEnginesClient(d *svtdmysql.TdmysqlDbEnginesDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

func newTdmysqlDbEnginesReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch

	return req
}

func newTdmysqlDbEnginesReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	return resp
}

func tdmysqlDbEnginesStubItems() []*tdmysqlv20211122.DBEngineInfo {
	return []*tdmysqlv20211122.DBEngineInfo{
		{
			Type:                        tdmysqlDbEnginesStr("TDSQL-C"),
			Version:                     tdmysqlDbEnginesStr("8.0"),
			Name:                        tdmysqlDbEnginesStr("tf-test-engine"),
			Description:                 tdmysqlDbEnginesStr("engine for terraform test"),
			New:                         tdmysqlDbEnginesBool(true),
			SQLMode:                     []*string{tdmysqlDbEnginesStr("STRICT_TRANS_TABLES"), tdmysqlDbEnginesStr("NO_ENGINE_SUBSTITUTION")},
			IsSupportParamTemplate:      tdmysqlDbEnginesBool(true),
			IsSupportServerless:         tdmysqlDbEnginesBool(false),
			IsSupportTDE:                tdmysqlDbEnginesBool(false),
			IsSupportPasswordComplexity: tdmysqlDbEnginesBool(true),
		},
	}
}

// TestTdmysqlDbEngines_MetadataAndSchema validates metadata and schema
// definition.
func TestTdmysqlDbEngines_MetadataAndSchema(t *testing.T) {
	d := svtdmysql.NewTdmysqlDbEnginesDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_tdmysql_db_engines", metaResp.TypeName)

	sch := tdmysqlDbEnginesSchema(t)
	attrs := sch.Attributes

	assert.True(t, attrs["id"].(schema.StringAttribute).Computed)
	assert.True(t, attrs["items"].(schema.ListNestedAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "total_count")
}

// TestTdmysqlDbEngines_ReadBasic tests a successful read with one record and
// asserts every leaf field is read back.
func TestTdmysqlDbEngines_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBEngines", func(request *tdmysqlv20211122.DescribeDBEnginesRequest) (*tdmysqlv20211122.DescribeDBEnginesResponse, error) {
		return &tdmysqlv20211122.DescribeDBEnginesResponse{Response: &tdmysqlv20211122.DescribeDBEnginesResponseParams{
			TotalCount: func() *int64 { v := int64(1); return &v }(),
			Items:      tdmysqlDbEnginesStubItems(),
			RequestId:  tdmysqlDbEnginesStr("fake-request-id"),
		}}, nil
	})

	ds := svtdmysql.NewTdmysqlDbEnginesDataSource().(*svtdmysql.TdmysqlDbEnginesDataSource)
	setTdmysqlDbEnginesClient(ds, client)

	sch := tdmysqlDbEnginesSchema(t)
	raw := buildTdmysqlDbEnginesRawValue()
	req := newTdmysqlDbEnginesReadRequest(sch, raw)
	resp := newTdmysqlDbEnginesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svtdmysql.TdmysqlDbEnginesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())

	assert.False(t, result.Items.IsNull())

	var items []svtdmysql.TdmysqlDbEnginesItemModel
	resp.Diagnostics.Append(result.Items.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)

	assert.Equal(t, "TDSQL-C", items[0].Type.ValueString())
	assert.Equal(t, "8.0", items[0].Version.ValueString())
	assert.Equal(t, "tf-test-engine", items[0].Name.ValueString())
	assert.Equal(t, "engine for terraform test", items[0].Description.ValueString())
	assert.True(t, items[0].New.ValueBool())

	var sqlModes []string
	resp.Diagnostics.Append(items[0].SQLMode.ElementsAs(context.Background(), &sqlModes, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"STRICT_TRANS_TABLES", "NO_ENGINE_SUBSTITUTION"}, sqlModes)

	assert.True(t, items[0].IsSupportParamTemplate.ValueBool())
	assert.False(t, items[0].IsSupportServerless.ValueBool())
	assert.False(t, items[0].IsSupportTDE.ValueBool())
	assert.True(t, items[0].IsSupportPasswordComplexity.ValueBool())
}

// TestTdmysqlDbEngines_ReadEmpty covers an empty result: a legitimate answer,
// not a failure — the data source must succeed with an empty list.
func TestTdmysqlDbEngines_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBEngines", func(request *tdmysqlv20211122.DescribeDBEnginesRequest) (*tdmysqlv20211122.DescribeDBEnginesResponse, error) {
		return &tdmysqlv20211122.DescribeDBEnginesResponse{Response: &tdmysqlv20211122.DescribeDBEnginesResponseParams{
			TotalCount: func() *int64 { v := int64(0); return &v }(),
			Items:      []*tdmysqlv20211122.DBEngineInfo{},
			RequestId:  tdmysqlDbEnginesStr("fake-request-id"),
		}}, nil
	})

	ds := svtdmysql.NewTdmysqlDbEnginesDataSource().(*svtdmysql.TdmysqlDbEnginesDataSource)
	setTdmysqlDbEnginesClient(ds, client)

	sch := tdmysqlDbEnginesSchema(t)
	raw := buildTdmysqlDbEnginesRawValue()
	req := newTdmysqlDbEnginesReadRequest(sch, raw)
	resp := newTdmysqlDbEnginesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svtdmysql.TdmysqlDbEnginesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.Items.IsNull())

	var items []svtdmysql.TdmysqlDbEnginesItemModel
	resp.Diagnostics.Append(result.Items.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

// TestTdmysqlDbEngines_ReadApiError covers an API-level failure: the error
// must surface through diagnostics instead of silently persisting empty
// state.
func TestTdmysqlDbEngines_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBEngines", func(request *tdmysqlv20211122.DescribeDBEnginesRequest) (*tdmysqlv20211122.DescribeDBEnginesResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svtdmysql.NewTdmysqlDbEnginesDataSource().(*svtdmysql.TdmysqlDbEnginesDataSource)
	setTdmysqlDbEnginesClient(ds, client)

	sch := tdmysqlDbEnginesSchema(t)
	raw := buildTdmysqlDbEnginesRawValue()
	req := newTdmysqlDbEnginesReadRequest(sch, raw)
	resp := newTdmysqlDbEnginesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}
