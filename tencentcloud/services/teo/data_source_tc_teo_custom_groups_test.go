package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStrTeoCustomGroups(s string) *string { return &s }
func ptrInt64TeoCustomGroups(v int64) *int64 { return &v }

// teoCustomGroupsSchema returns the framework data source schema used to build
// the raw config/state value.
func teoCustomGroupsSchema(t *testing.T) schema.Schema {
	d := svcteo.NewTeoCustomGroupsDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// teoCustomGroupEntryItemType builds the tftypes.Object type for a
// group_items element.
func teoCustomGroupEntryItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"content":   tftypes.String,
			"item_type": tftypes.String,
		},
	}
}

// teoCustomGroupItemType builds the tftypes.Object type for a custom_groups
// element.
func teoCustomGroupItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"group_id":               tftypes.String,
			"group_name":             tftypes.String,
			"group_type":             tftypes.String,
			"group_items":            tftypes.List{ElementType: teoCustomGroupEntryItemType()},
			"group_item_total_count": tftypes.Number,
			"ref_count":              tftypes.Number,
		},
	}
}

// teoCustomGroupsRootType builds the tftypes.Object type for the data source
// root, mirroring every schema attribute.
func teoCustomGroupsRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":           tftypes.String,
			"group_types":       tftypes.List{ElementType: tftypes.String},
			"group_total_count": tftypes.Number,
			"custom_groups":     tftypes.List{ElementType: teoCustomGroupItemType()},
		},
	}
}

// buildTeoCustomGroupsRawValue builds the raw config value with the given query
// inputs (zone_id is required; groupTypes may be nil).
func buildTeoCustomGroupsRawValue(zoneId string, groupTypes []string) tftypes.Value {
	var rawGroupTypes tftypes.Value
	if groupTypes == nil {
		rawGroupTypes = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil)
	} else {
		elems := make([]tftypes.Value, 0, len(groupTypes))
		for _, v := range groupTypes {
			elems = append(elems, tftypes.NewValue(tftypes.String, v))
		}
		rawGroupTypes = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, elems)
	}
	attrs := map[string]tftypes.Value{
		"zone_id":           tftypes.NewValue(tftypes.String, zoneId),
		"group_types":       rawGroupTypes,
		"group_total_count": tftypes.NewValue(tftypes.Number, nil),
		"custom_groups":     tftypes.NewValue(tftypes.List{ElementType: teoCustomGroupItemType()}, nil),
	}
	return tftypes.NewValue(teoCustomGroupsRootType(), attrs)
}

// setTeoCustomGroupsClient injects a mock connectivity client into the data
// source instance via the Configure method.
func setTeoCustomGroupsClient(d *svcteo.TeoCustomGroupsDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// newTeoCustomGroupsReadRequest builds a ReadRequest with the given config
// value and schema.
func newTeoCustomGroupsReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	return req
}

// newTeoCustomGroupsReadResponse builds a ReadResponse whose State carries the
// schema, so resp.State.Set can reflect the model into a raw value.
func newTeoCustomGroupsReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch
	return resp
}

// go test ./tencentcloud/services/teo/ -run "TestTeoCustomGroups" -v -count=1 -gcflags="all=-l"

// TestTeoCustomGroups_MetadataAndSchema validates metadata and schema
// definition.
func TestTeoCustomGroups_MetadataAndSchema(t *testing.T) {
	d := svcteo.NewTeoCustomGroupsDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_custom_groups", metaResp.TypeName)

	sch := teoCustomGroupsSchema(t)
	attrs := sch.Attributes

	// Required input
	assert.True(t, attrs["zone_id"].(schema.StringAttribute).Required)

	// Optional input
	assert.True(t, attrs["group_types"].(schema.ListAttribute).Optional)

	// limit and offset are intentionally not exposed: the service layer
	// loops through all pages internally.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")

	// Computed outputs
	assert.True(t, attrs["group_total_count"].(schema.Int64Attribute).Computed)
	assert.True(t, attrs["custom_groups"].(schema.ListNestedAttribute).Computed)

	// Nested group attributes
	nested := attrs["custom_groups"].(schema.ListNestedAttribute).NestedObject.Attributes
	for _, leaf := range []string{"group_id", "group_name", "group_type"} {
		nestedAttr, ok := nested[leaf].(schema.StringAttribute)
		assert.True(t, ok, "nested attribute %s should be StringAttribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
	assert.True(t, nested["group_item_total_count"].(schema.Int64Attribute).Computed)
	assert.True(t, nested["ref_count"].(schema.Int64Attribute).Computed)
	assert.True(t, nested["group_items"].(schema.ListNestedAttribute).Computed)

	// Nested entry attributes
	entryNested := nested["group_items"].(schema.ListNestedAttribute).NestedObject.Attributes
	for _, leaf := range []string{"content", "item_type"} {
		entryAttr, ok := entryNested[leaf].(schema.StringAttribute)
		assert.True(t, ok, "nested entry attribute %s should be StringAttribute", leaf)
		assert.True(t, entryAttr.Computed, "nested entry attribute %s should be Computed", leaf)
	}
}

// TestTeoCustomGroups_ReadBasic tests a successful read with two custom groups
// and asserts every leaf field is read back.
func TestTeoCustomGroups_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &svcteo.TeoService{}
	patches.ApplyMethodFunc(service, "DescribeTeoCustomGroups", func(_ context.Context, request *teov20220901.DescribeCustomGroupsRequest) ([]*teov20220901.CustomGroup, *int64, error) {
		if request.ZoneId == nil || *request.ZoneId != "zone-27h0vbm5w1e" {
			t.Errorf("expected zone_id to be passed to the API request")
		}
		if len(request.GroupTypes) != 2 {
			t.Errorf("expected 2 group types to be passed to the API request, got %d", len(request.GroupTypes))
		}
		customGroups := []*teov20220901.CustomGroup{
			{
				GroupId:   ptrStrTeoCustomGroups("cg-abcdefg"),
				GroupName: ptrStrTeoCustomGroups("tls-fingerprints"),
				GroupType: ptrStrTeoCustomGroups("TLSFingerprint"),
				GroupItems: []*teov20220901.CustomGroupItem{
					{
						Content:  ptrStrTeoCustomGroups("cd568b9c5d9dfa6f2b1f4d1c0f0e9a8b"),
						ItemType: ptrStrTeoCustomGroups("JA3"),
					},
					{
						Content:  ptrStrTeoCustomGroups("t13d1516h2_8daaf6cb27e1a1d3"),
						ItemType: ptrStrTeoCustomGroups("JA4"),
					},
				},
				GroupItemTotalCount: ptrInt64TeoCustomGroups(2),
				RefCount:            ptrInt64TeoCustomGroups(1),
			},
			{
				GroupId:   ptrStrTeoCustomGroups("cg-hijklmn"),
				GroupName: ptrStrTeoCustomGroups("blocked-uas"),
				GroupType: ptrStrTeoCustomGroups("UA"),
				GroupItems: []*teov20220901.CustomGroupItem{
					{
						Content:  ptrStrTeoCustomGroups("Mozilla/5.0 (compatible; badbot/1.0)"),
						ItemType: ptrStrTeoCustomGroups("UA"),
					},
				},
				GroupItemTotalCount: ptrInt64TeoCustomGroups(1),
				RefCount:            ptrInt64TeoCustomGroups(0),
			},
		}
		return customGroups, ptrInt64TeoCustomGroups(2), nil
	})

	ds := svcteo.NewTeoCustomGroupsDataSource().(*svcteo.TeoCustomGroupsDataSource)
	setTeoCustomGroupsClient(ds, client)

	sch := teoCustomGroupsSchema(t)
	raw := buildTeoCustomGroupsRawValue("zone-27h0vbm5w1e", []string{"TLSFingerprint", "UA"})
	req := newTeoCustomGroupsReadRequest(sch, raw)
	resp := newTeoCustomGroupsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcteo.TeoCustomGroupsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, "zone-27h0vbm5w1e", result.ZoneId.ValueString())
	assert.Equal(t, int64(2), result.GroupTotalCount.ValueInt64())
	assert.False(t, result.CustomGroups.IsNull())

	var items []svcteo.TeoCustomGroupItemModel
	resp.Diagnostics.Append(result.CustomGroups.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 2)

	// assert every leaf field of the first group
	assert.Equal(t, "cg-abcdefg", items[0].GroupId.ValueString())
	assert.Equal(t, "tls-fingerprints", items[0].GroupName.ValueString())
	assert.Equal(t, "TLSFingerprint", items[0].GroupType.ValueString())
	assert.Equal(t, int64(2), items[0].GroupItemTotalCount.ValueInt64())
	assert.Equal(t, int64(1), items[0].RefCount.ValueInt64())

	assert.False(t, items[0].GroupItems.IsNull())
	var entries []svcteo.TeoCustomGroupEntryModel
	resp.Diagnostics.Append(items[0].GroupItems.ElementsAs(context.Background(), &entries, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, entries, 2)
	assert.Equal(t, "cd568b9c5d9dfa6f2b1f4d1c0f0e9a8b", entries[0].Content.ValueString())
	assert.Equal(t, "JA3", entries[0].ItemType.ValueString())
	assert.Equal(t, "t13d1516h2_8daaf6cb27e1a1d3", entries[1].Content.ValueString())
	assert.Equal(t, "JA4", entries[1].ItemType.ValueString())

	// assert every leaf field of the second group
	assert.Equal(t, "cg-hijklmn", items[1].GroupId.ValueString())
	assert.Equal(t, "blocked-uas", items[1].GroupName.ValueString())
	assert.Equal(t, "UA", items[1].GroupType.ValueString())
	assert.Equal(t, int64(1), items[1].GroupItemTotalCount.ValueInt64())
	assert.Equal(t, int64(0), items[1].RefCount.ValueInt64())

	assert.False(t, items[1].GroupItems.IsNull())
	var entries1 []svcteo.TeoCustomGroupEntryModel
	resp.Diagnostics.Append(items[1].GroupItems.ElementsAs(context.Background(), &entries1, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, entries1, 1)
	assert.Equal(t, "Mozilla/5.0 (compatible; badbot/1.0)", entries1[0].Content.ValueString())
	assert.Equal(t, "UA", entries1[0].ItemType.ValueString())
}

// TestTeoCustomGroups_ReadEmpty covers a site that has no custom groups. An
// empty result is a legitimate answer, not a failure: the data source must
// succeed with group_total_count = 0 and an empty custom_groups list.
func TestTeoCustomGroups_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &svcteo.TeoService{}
	patches.ApplyMethodFunc(service, "DescribeTeoCustomGroups", func(_ context.Context, request *teov20220901.DescribeCustomGroupsRequest) ([]*teov20220901.CustomGroup, *int64, error) {
		return []*teov20220901.CustomGroup{}, ptrInt64TeoCustomGroups(0), nil
	})

	ds := svcteo.NewTeoCustomGroupsDataSource().(*svcteo.TeoCustomGroupsDataSource)
	setTeoCustomGroupsClient(ds, client)

	sch := teoCustomGroupsSchema(t)
	raw := buildTeoCustomGroupsRawValue("zone-27h0vbm5w1e", nil)
	req := newTeoCustomGroupsReadRequest(sch, raw)
	resp := newTeoCustomGroupsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcteo.TeoCustomGroupsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, int64(0), result.GroupTotalCount.ValueInt64())
	assert.False(t, result.CustomGroups.IsNull())

	var items []svcteo.TeoCustomGroupItemModel
	resp.Diagnostics.Append(result.CustomGroups.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

// TestTeoCustomGroups_ReadApiError covers an API-level failure: the data source
// must surface the error through diagnostics instead of silently persisting
// empty state.
func TestTeoCustomGroups_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &svcteo.TeoService{}
	patches.ApplyMethodFunc(service, "DescribeTeoCustomGroups", func(_ context.Context, request *teov20220901.DescribeCustomGroupsRequest) ([]*teov20220901.CustomGroup, *int64, error) {
		return nil, nil, fmt.Errorf("mocked API error")
	})

	ds := svcteo.NewTeoCustomGroupsDataSource().(*svcteo.TeoCustomGroupsDataSource)
	setTeoCustomGroupsClient(ds, client)

	sch := teoCustomGroupsSchema(t)
	raw := buildTeoCustomGroupsRawValue("zone-27h0vbm5w1e", nil)
	req := newTeoCustomGroupsReadRequest(sch, raw)
	resp := newTeoCustomGroupsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

// TestTeoCustomGroupsService_ReadAllPages verifies that the service layer loops
// through every DescribeCustomGroups page with the page size 100 and returns
// the complete group set plus the total count from the cloud side.
func TestTeoCustomGroupsService_ReadAllPages(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	callCount := 0
	patches.ApplyMethodFunc(teoClient, "DescribeCustomGroupsWithContext", func(_ context.Context, request *teov20220901.DescribeCustomGroupsRequest) (*teov20220901.DescribeCustomGroupsResponse, error) {
		callCount++
		// The service must drive pagination itself with the page size 100.
		if request.Limit == nil || *request.Limit != 100 {
			t.Errorf("expected Limit 100 on call %d, got %v", callCount, request.Limit)
		}
		var offset int64
		if request.Offset != nil {
			offset = *request.Offset
		}
		switch callCount {
		case 1:
			if offset != 0 {
				t.Errorf("expected offset 0 on first call, got %d", offset)
			}
			customGroups := make([]*teov20220901.CustomGroup, 0, 100)
			for i := 0; i < 100; i++ {
				customGroups = append(customGroups, &teov20220901.CustomGroup{
					GroupId:   ptrStrTeoCustomGroups(fmt.Sprintf("cg-page1-%d", i)),
					GroupType: ptrStrTeoCustomGroups("UA"),
				})
			}
			return &teov20220901.DescribeCustomGroupsResponse{Response: &teov20220901.DescribeCustomGroupsResponseParams{
				GroupTotalCount: ptrInt64TeoCustomGroups(101),
				CustomGroups:    customGroups,
				RequestId:       ptrStrTeoCustomGroups("fake-request-id"),
			}}, nil
		case 2:
			if offset != 100 {
				t.Errorf("expected offset 100 on second call, got %d", offset)
			}
			return &teov20220901.DescribeCustomGroupsResponse{Response: &teov20220901.DescribeCustomGroupsResponseParams{
				GroupTotalCount: ptrInt64TeoCustomGroups(101),
				CustomGroups: []*teov20220901.CustomGroup{
					{GroupId: ptrStrTeoCustomGroups("cg-page2-0"), GroupType: ptrStrTeoCustomGroups("Referer")},
				},
				RequestId: ptrStrTeoCustomGroups("fake-request-id"),
			}}, nil
		default:
			t.Errorf("unexpected extra API call %d", callCount)
			return &teov20220901.DescribeCustomGroupsResponse{}, nil
		}
	})

	service := svcteo.NewTeoService(client)
	request := teov20220901.NewDescribeCustomGroupsRequest()
	request.ZoneId = ptrStrTeoCustomGroups("zone-27h0vbm5w1e")

	customGroups, groupTotalCount, err := service.DescribeTeoCustomGroups(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
	if groupTotalCount == nil || *groupTotalCount != 101 {
		t.Errorf("expected group total count 101, got %v", groupTotalCount)
	}
	if len(customGroups) != 101 {
		t.Fatalf("expected 101 custom groups across pages, got %d", len(customGroups))
	}
	if customGroups[0].GroupId == nil || *customGroups[0].GroupId != "cg-page1-0" {
		t.Errorf("unexpected first custom group id: %v", customGroups[0].GroupId)
	}
	if customGroups[100].GroupId == nil || *customGroups[100].GroupId != "cg-page2-0" {
		t.Errorf("unexpected last custom group id: %v", customGroups[100].GroupId)
	}
}
