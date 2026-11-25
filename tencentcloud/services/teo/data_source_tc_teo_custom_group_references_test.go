package teo_test

import (
	"context"
	"errors"
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

func ptrStringTeoCustomGroupRefs(s string) *string { return &s }

// teoCustomGroupReferencesSchema returns the framework data source schema
// used to build the raw config/state value.
func teoCustomGroupReferencesSchema(t *testing.T) schema.Schema {
	d := svcteo.NewTeoCustomGroupReferencesDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// teoCustomGroupReferenceItemType builds the tftypes.Object type for a
// references element.
func teoCustomGroupReferenceItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":         tftypes.String,
			"entity_type":     tftypes.String,
			"entity_id":       tftypes.String,
			"entity_name":     tftypes.String,
			"sub_entity_type": tftypes.String,
			"sub_entity_id":   tftypes.String,
			"sub_entity_name": tftypes.String,
		},
	}
}

// teoCustomGroupReferencesRootType builds the tftypes.Object type for the
// data source root, mirroring every schema attribute.
func teoCustomGroupReferencesRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":    tftypes.String,
			"group_id":   tftypes.String,
			"references": tftypes.List{ElementType: teoCustomGroupReferenceItemType()},
		},
	}
}

// buildTeoCustomGroupReferencesRawValue builds the raw config value with the
// given query inputs.
func buildTeoCustomGroupReferencesRawValue(zoneId, groupId string) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"zone_id":    tftypes.NewValue(tftypes.String, zoneId),
		"group_id":   tftypes.NewValue(tftypes.String, groupId),
		"references": tftypes.NewValue(tftypes.List{ElementType: teoCustomGroupReferenceItemType()}, nil),
	}
	return tftypes.NewValue(teoCustomGroupReferencesRootType(), attrs)
}

// setTeoCustomGroupReferencesClient injects a mock connectivity client into
// the data source instance via the Configure method.
func setTeoCustomGroupReferencesClient(d *svcteo.TeoCustomGroupReferencesDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// newTeoCustomGroupReferencesReadRequest builds a ReadRequest with the given
// config value and schema.
func newTeoCustomGroupReferencesReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	return req
}

// newTeoCustomGroupReferencesReadResponse builds a ReadResponse whose State
// carries the schema, so resp.State.Set can reflect the model into a raw value.
func newTeoCustomGroupReferencesReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch
	return resp
}

// patchTeoCustomGroupReferences patches the SDK client method
// DescribeCustomGroupReferencesWithContext with the given fake.
func patchTeoCustomGroupReferences(patches *gomonkey.Patches, fn func(ctx context.Context, req *teov20220901.DescribeCustomGroupReferencesRequest) (*teov20220901.DescribeCustomGroupReferencesResponse, error)) {
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeCustomGroupReferencesWithContext", fn)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoCustomGroupReferences" -v -count=1 -gcflags="all=-l"

// TestTeoCustomGroupReferences_MetadataAndSchema validates metadata and
// schema definition.
func TestTeoCustomGroupReferences_MetadataAndSchema(t *testing.T) {
	d := svcteo.NewTeoCustomGroupReferencesDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_custom_group_references", metaResp.TypeName)

	sch := teoCustomGroupReferencesSchema(t)
	attrs := sch.Attributes

	// Required inputs
	for _, req := range []string{"zone_id", "group_id"} {
		attr, ok := attrs[req].(schema.StringAttribute)
		assert.True(t, ok, "attribute %s should be StringAttribute", req)
		assert.True(t, attr.Required, "attribute %s should be Required", req)
	}

	// Computed outputs
	references, ok := attrs["references"].(schema.ListNestedAttribute)
	assert.True(t, ok, "attribute references should be ListNestedAttribute")
	assert.True(t, references.Computed, "attribute references should be Computed")

	// Nested reference attributes
	nested := references.NestedObject.Attributes
	for _, leaf := range []string{"zone_id", "entity_type", "entity_id", "entity_name", "sub_entity_type", "sub_entity_id", "sub_entity_name"} {
		nestedAttr, ok := nested[leaf].(schema.StringAttribute)
		assert.True(t, ok, "nested attribute %s should be StringAttribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
}

// TestTeoCustomGroupReferences_ReadBasic tests a successful read with two
// references and asserts every leaf field is read back.
func TestTeoCustomGroupReferences_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchTeoCustomGroupReferences(patches, func(_ context.Context, request *teov20220901.DescribeCustomGroupReferencesRequest) (*teov20220901.DescribeCustomGroupReferencesResponse, error) {
		if request.ZoneId == nil || *request.ZoneId != "zone-2qtuhspy7cr6" {
			t.Errorf("expected zone_id to be passed to the API request")
		}
		if request.GroupId == nil || *request.GroupId != "cg-2q1pkwfyagzz" {
			t.Errorf("expected group_id to be passed to the API request")
		}
		resp := teov20220901.NewDescribeCustomGroupReferencesResponse()
		resp.Response = &teov20220901.DescribeCustomGroupReferencesResponseParams{
			References: []*teov20220901.CustomGroupReference{
				{
					ZoneId:        ptrStringTeoCustomGroupRefs("zone-2qtuhspy7cr6"),
					EntityType:    ptrStringTeoCustomGroupRefs("WebSec.ZonePolicy"),
					EntityId:      ptrStringTeoCustomGroupRefs("zone-2qtuhspy7cr6"),
					EntityName:    ptrStringTeoCustomGroupRefs(""),
					SubEntityType: ptrStringTeoCustomGroupRefs("WebSec.ExceptionRule"),
					SubEntityId:   ptrStringTeoCustomGroupRefs("rule-001"),
					SubEntityName: ptrStringTeoCustomGroupRefs("exception rule"),
				},
				{
					ZoneId:        ptrStringTeoCustomGroupRefs("zone-2qtuhspy7cr6"),
					EntityType:    ptrStringTeoCustomGroupRefs("DDoS.L4Proxy"),
					EntityId:      ptrStringTeoCustomGroupRefs("proxy-xxx"),
					EntityName:    ptrStringTeoCustomGroupRefs(""),
					SubEntityType: ptrStringTeoCustomGroupRefs("DDoS.L4Proxy.IpAccessControl"),
					SubEntityId:   ptrStringTeoCustomGroupRefs(""),
					SubEntityName: ptrStringTeoCustomGroupRefs(""),
				},
			},
			TotalCount: func(v int64) *int64 { return &v }(2),
			RequestId:  ptrStringTeoCustomGroupRefs("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoCustomGroupReferencesDataSource().(*svcteo.TeoCustomGroupReferencesDataSource)
	setTeoCustomGroupReferencesClient(ds, &connectivity.TencentCloudClient{})

	sch := teoCustomGroupReferencesSchema(t)
	raw := buildTeoCustomGroupReferencesRawValue("zone-2qtuhspy7cr6", "cg-2q1pkwfyagzz")
	req := newTeoCustomGroupReferencesReadRequest(sch, raw)
	resp := newTeoCustomGroupReferencesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcteo.TeoCustomGroupReferencesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, "zone-2qtuhspy7cr6", result.ZoneId.ValueString())
	assert.Equal(t, "cg-2q1pkwfyagzz", result.GroupId.ValueString())
	assert.False(t, result.References.IsNull())

	var items []svcteo.TeoCustomGroupReferenceItemModel
	resp.Diagnostics.Append(result.References.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 2)

	// assert every leaf field of the first reference
	assert.Equal(t, "zone-2qtuhspy7cr6", items[0].ZoneId.ValueString())
	assert.Equal(t, "WebSec.ZonePolicy", items[0].EntityType.ValueString())
	assert.Equal(t, "zone-2qtuhspy7cr6", items[0].EntityId.ValueString())
	assert.Equal(t, "", items[0].EntityName.ValueString())
	assert.Equal(t, "WebSec.ExceptionRule", items[0].SubEntityType.ValueString())
	assert.Equal(t, "rule-001", items[0].SubEntityId.ValueString())
	assert.Equal(t, "exception rule", items[0].SubEntityName.ValueString())

	// assert every leaf field of the second reference
	assert.Equal(t, "zone-2qtuhspy7cr6", items[1].ZoneId.ValueString())
	assert.Equal(t, "DDoS.L4Proxy", items[1].EntityType.ValueString())
	assert.Equal(t, "proxy-xxx", items[1].EntityId.ValueString())
	assert.Equal(t, "", items[1].EntityName.ValueString())
	assert.Equal(t, "DDoS.L4Proxy.IpAccessControl", items[1].SubEntityType.ValueString())
	assert.Equal(t, "", items[1].SubEntityId.ValueString())
	assert.Equal(t, "", items[1].SubEntityName.ValueString())
}

// TestTeoCustomGroupReferences_ReadPaginated tests that the service layer
// keeps fetching pages until the cloud side returns fewer records than the
// page size.
func TestTeoCustomGroupReferences_ReadPaginated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	callCount := 0
	patchTeoCustomGroupReferences(patches, func(_ context.Context, request *teov20220901.DescribeCustomGroupReferencesRequest) (*teov20220901.DescribeCustomGroupReferencesResponse, error) {
		callCount++
		resp := teov20220901.NewDescribeCustomGroupReferencesResponse()
		resp.Response = &teov20220901.DescribeCustomGroupReferencesResponseParams{
			RequestId: ptrStringTeoCustomGroupRefs("fake-request-id"),
		}
		if callCount == 1 {
			// first page: a full page (limit=200), so the loop continues
			resp.Response.References = make([]*teov20220901.CustomGroupReference, 200)
			for i := range resp.Response.References {
				resp.Response.References[i] = &teov20220901.CustomGroupReference{
					ZoneId:     ptrStringTeoCustomGroupRefs("zone-paginated"),
					EntityType: ptrStringTeoCustomGroupRefs("WebSec.HostPolicy"),
					EntityId:   ptrStringTeoCustomGroupRefs("example.com"),
				}
			}
			resp.Response.TotalCount = func(v int64) *int64 { return &v }(201)
		} else {
			// second page: fewer records than the page size, the loop ends
			resp.Response.References = []*teov20220901.CustomGroupReference{
				{
					ZoneId:     ptrStringTeoCustomGroupRefs("zone-paginated"),
					EntityType: ptrStringTeoCustomGroupRefs("WebSec.Template"),
					EntityId:   ptrStringTeoCustomGroupRefs("template-001"),
				},
			}
			resp.Response.TotalCount = func(v int64) *int64 { return &v }(201)
		}
		return resp, nil
	})

	ds := svcteo.NewTeoCustomGroupReferencesDataSource().(*svcteo.TeoCustomGroupReferencesDataSource)
	setTeoCustomGroupReferencesClient(ds, &connectivity.TencentCloudClient{})

	sch := teoCustomGroupReferencesSchema(t)
	raw := buildTeoCustomGroupReferencesRawValue("zone-paginated", "cg-paginated")
	req := newTeoCustomGroupReferencesReadRequest(sch, raw)
	resp := newTeoCustomGroupReferencesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.Equal(t, 2, callCount)

	var result svcteo.TeoCustomGroupReferencesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.False(t, result.References.IsNull())
	assert.Equal(t, 201, len(result.References.Elements()))
}

// TestTeoCustomGroupReferences_ReadEmpty covers a custom group that is not
// referenced by any configuration. An empty result is a legitimate answer, not
// a failure: the data source must succeed with an empty references list.
func TestTeoCustomGroupReferences_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchTeoCustomGroupReferences(patches, func(_ context.Context, request *teov20220901.DescribeCustomGroupReferencesRequest) (*teov20220901.DescribeCustomGroupReferencesResponse, error) {
		resp := teov20220901.NewDescribeCustomGroupReferencesResponse()
		resp.Response = &teov20220901.DescribeCustomGroupReferencesResponseParams{
			References: []*teov20220901.CustomGroupReference{},
			TotalCount: func(v int64) *int64 { return &v }(0),
			RequestId:  ptrStringTeoCustomGroupRefs("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoCustomGroupReferencesDataSource().(*svcteo.TeoCustomGroupReferencesDataSource)
	setTeoCustomGroupReferencesClient(ds, &connectivity.TencentCloudClient{})

	sch := teoCustomGroupReferencesSchema(t)
	raw := buildTeoCustomGroupReferencesRawValue("zone-empty", "cg-empty")
	req := newTeoCustomGroupReferencesReadRequest(sch, raw)
	resp := newTeoCustomGroupReferencesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "an unreferenced custom group must not be reported as an error")

	var result svcteo.TeoCustomGroupReferencesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.References.IsNull())
	assert.Equal(t, 0, len(result.References.Elements()))
}

// TestTeoCustomGroupReferences_ReadNilReferences covers the case where the
// API returns a nil References slice. The flatten logic must not panic and
// must produce an empty list.
func TestTeoCustomGroupReferences_ReadNilReferences(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchTeoCustomGroupReferences(patches, func(_ context.Context, request *teov20220901.DescribeCustomGroupReferencesRequest) (*teov20220901.DescribeCustomGroupReferencesResponse, error) {
		resp := teov20220901.NewDescribeCustomGroupReferencesResponse()
		resp.Response = &teov20220901.DescribeCustomGroupReferencesResponseParams{
			References: nil,
			TotalCount: func(v int64) *int64 { return &v }(0),
			RequestId:  ptrStringTeoCustomGroupRefs("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoCustomGroupReferencesDataSource().(*svcteo.TeoCustomGroupReferencesDataSource)
	setTeoCustomGroupReferencesClient(ds, &connectivity.TencentCloudClient{})

	sch := teoCustomGroupReferencesSchema(t)
	raw := buildTeoCustomGroupReferencesRawValue("zone-nil", "cg-nil")
	req := newTeoCustomGroupReferencesReadRequest(sch, raw)
	resp := newTeoCustomGroupReferencesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "a nil references list must not be reported as an error")

	var result svcteo.TeoCustomGroupReferencesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.References.IsNull())
	assert.Equal(t, 0, len(result.References.Elements()))
}

// TestTeoCustomGroupReferences_ReadApiError covers an API failure: the data
// source must report an error diagnostic instead of silently returning an
// empty result.
func TestTeoCustomGroupReferences_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchTeoCustomGroupReferences(patches, func(_ context.Context, request *teov20220901.DescribeCustomGroupReferencesRequest) (*teov20220901.DescribeCustomGroupReferencesResponse, error) {
		return nil, errors.New("InternalError.RouteError")
	})

	ds := svcteo.NewTeoCustomGroupReferencesDataSource().(*svcteo.TeoCustomGroupReferencesDataSource)
	setTeoCustomGroupReferencesClient(ds, &connectivity.TencentCloudClient{})

	sch := teoCustomGroupReferencesSchema(t)
	raw := buildTeoCustomGroupReferencesRawValue("zone-error", "cg-error")
	req := newTeoCustomGroupReferencesReadRequest(sch, raw)
	resp := newTeoCustomGroupReferencesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "an API failure must be reported as an error")
}

// TestTeoCustomGroupReferences_ReadNilResponse covers the case where the
// API returns a nil response with a nil error: the service layer must report
// it as an error instead of treating it as an empty result.
func TestTeoCustomGroupReferences_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchTeoCustomGroupReferences(patches, func(_ context.Context, request *teov20220901.DescribeCustomGroupReferencesRequest) (*teov20220901.DescribeCustomGroupReferencesResponse, error) {
		// nil response with nil error; the service layer converts this into
		// an error via the resource.Retry wrapper.
		return nil, nil
	})

	ds := svcteo.NewTeoCustomGroupReferencesDataSource().(*svcteo.TeoCustomGroupReferencesDataSource)
	setTeoCustomGroupReferencesClient(ds, &connectivity.TencentCloudClient{})

	sch := teoCustomGroupReferencesSchema(t)
	raw := buildTeoCustomGroupReferencesRawValue("zone-nil-response", "cg-nil-response")
	req := newTeoCustomGroupReferencesReadRequest(sch, raw)
	resp := newTeoCustomGroupReferencesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "a nil response must be reported as an error")
}
