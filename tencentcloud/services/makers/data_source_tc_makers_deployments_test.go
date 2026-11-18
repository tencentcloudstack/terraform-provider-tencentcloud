package makers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	makersv20261113 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcmakers "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/makers"
)

func ptrStrMakers(s string) *string { return &s }
func ptrInt64Makers(v int64) *int64 { return &v }
func ptrBoolMakers(v bool) *bool    { return &v }

// makersDeploymentsSchema returns the framework data source schema used to
// build the raw config/state value.
func makersDeploymentsSchema(t *testing.T) schema.Schema {
	d := svcmakers.NewMakersDeploymentsDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// makersDeploymentsRootType builds the tftypes.Object type for the data source
// root, mirroring every schema attribute.
func makersDeploymentsRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":        tftypes.String,
			"project_id":     tftypes.String,
			"deployment_ids": tftypes.List{ElementType: tftypes.String},
			"filters":        tftypes.List{ElementType: makersFilterItemType()},
			"order_by":       tftypes.String,
			"order_type":     tftypes.String,
			"start_time":     tftypes.String,
			"end_time":       tftypes.String,
			"limit":          tftypes.Number,
			"offset":         tftypes.Number,
			"total_count":    tftypes.Number,
			"deployments":    tftypes.List{ElementType: makersDeploymentItemType()},
		},
	}
}

// makersFilterItemType builds the tftypes.Object type for a filters element.
func makersFilterItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":   tftypes.String,
			"values": tftypes.List{ElementType: tftypes.String},
		},
	}
}

// makersDeploymentItemType builds the tftypes.Object type for a deployments
// element.
func makersDeploymentItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"deployment_id":    tftypes.String,
			"project_id":       tftypes.String,
			"status":           tftypes.String,
			"original_status":  tftypes.String,
			"code":             tftypes.Number,
			"build_cost":       tftypes.String,
			"via_meta":         tftypes.String,
			"env":              tftypes.String,
			"repo_branch":      tftypes.String,
			"repo_commit_hash": tftypes.String,
			"repo_commit_msg":  tftypes.String,
			"preview_url":      tftypes.String,
			"project_url":      tftypes.String,
			"used_in_prod":     tftypes.Bool,
			"cover_url":        tftypes.String,
			"meta_data":        tftypes.String,
			"dpl_repo_url":     tftypes.String,
			"created_on":       tftypes.String,
			"modified_on":      tftypes.String,
		},
	}
}

// buildMakersDeploymentsRawValue builds the raw config value with the given
// query inputs (zone_id and project_id are required).
func buildMakersDeploymentsRawValue(zoneId, projectId string) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"zone_id":        tftypes.NewValue(tftypes.String, zoneId),
		"project_id":     tftypes.NewValue(tftypes.String, projectId),
		"deployment_ids": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
		"filters":        tftypes.NewValue(tftypes.List{ElementType: makersFilterItemType()}, nil),
		"order_by":       tftypes.NewValue(tftypes.String, nil),
		"order_type":     tftypes.NewValue(tftypes.String, nil),
		"start_time":     tftypes.NewValue(tftypes.String, nil),
		"end_time":       tftypes.NewValue(tftypes.String, nil),
		"limit":          tftypes.NewValue(tftypes.Number, nil),
		"offset":         tftypes.NewValue(tftypes.Number, nil),
		"total_count":    tftypes.NewValue(tftypes.Number, nil),
		"deployments":    tftypes.NewValue(tftypes.List{ElementType: makersDeploymentItemType()}, nil),
	}
	return tftypes.NewValue(makersDeploymentsRootType(), attrs)
}

// setMakersDeploymentsClient injects a mock connectivity client into the data
// source instance via the Configure method.
func setMakersDeploymentsClient(d *svcmakers.MakersDeploymentsDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// newMakersDeploymentsReadRequest builds a ReadRequest with the given config
// value and schema.
func newMakersDeploymentsReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	return req
}

// newMakersDeploymentsReadResponse builds a ReadResponse whose State carries
// the schema, so resp.State.Set can reflect the model into a raw value.
func newMakersDeploymentsReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch
	return resp
}

// go test ./tencentcloud/services/makers/ -run "TestMakersDeployments" -v -count=1 -gcflags="all=-l"

// TestMakersDeployments_MetadataAndSchema validates metadata and schema
// definition.
func TestMakersDeployments_MetadataAndSchema(t *testing.T) {
	d := svcmakers.NewMakersDeploymentsDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_makers_deployments", metaResp.TypeName)

	sch := makersDeploymentsSchema(t)
	attrs := sch.Attributes

	// Required inputs
	assert.True(t, attrs["zone_id"].(schema.StringAttribute).Required)
	assert.True(t, attrs["project_id"].(schema.StringAttribute).Required)

	// Optional inputs
	for _, opt := range []string{"order_by", "order_type", "start_time", "end_time"} {
		attr, ok := attrs[opt].(schema.StringAttribute)
		assert.True(t, ok, "attribute %s should be StringAttribute", opt)
		assert.True(t, attr.Optional, "attribute %s should be Optional", opt)
	}
	assert.True(t, attrs["deployment_ids"].(schema.ListAttribute).Optional)
	assert.True(t, attrs["limit"].(schema.Int64Attribute).Optional)
	assert.True(t, attrs["offset"].(schema.Int64Attribute).Optional)
	assert.True(t, attrs["filters"].(schema.ListNestedAttribute).Optional)

	// Computed outputs
	assert.True(t, attrs["total_count"].(schema.Int64Attribute).Computed)
	assert.True(t, attrs["deployments"].(schema.ListNestedAttribute).Computed)

	// Nested item attributes
	nested := attrs["deployments"].(schema.ListNestedAttribute).NestedObject.Attributes
	for _, leaf := range []string{"deployment_id", "project_id", "status", "original_status", "build_cost", "via_meta", "env", "repo_branch", "repo_commit_hash", "repo_commit_msg", "preview_url", "project_url", "cover_url", "meta_data", "dpl_repo_url", "created_on", "modified_on"} {
		nestedAttr, ok := nested[leaf].(schema.StringAttribute)
		assert.True(t, ok, "nested attribute %s should be StringAttribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
	assert.True(t, nested["code"].(schema.Int64Attribute).Computed)
	assert.True(t, nested["used_in_prod"].(schema.BoolAttribute).Computed)

	// Nested filter attributes
	filterNested := attrs["filters"].(schema.ListNestedAttribute).NestedObject.Attributes
	assert.True(t, filterNested["name"].(schema.StringAttribute).Required)
	assert.True(t, filterNested["values"].(schema.ListAttribute).Required)
}

// TestMakersDeployments_ReadBasic tests a successful read with one deployment
// record and asserts every leaf field is read back.
func TestMakersDeployments_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	makersClient := &makersv20261113.Client{}
	patches.ApplyMethodReturn(client, "UseMakersClient", makersClient)

	patches.ApplyMethodFunc(makersClient, "DescribeDeploymentsWithContext", func(_ context.Context, request *makersv20261113.DescribeDeploymentsRequest) (*makersv20261113.DescribeDeploymentsResponse, error) {
		if request.ZoneId == nil || *request.ZoneId != "zone-abcdefg" {
			t.Errorf("expected zone_id to be passed to the API request")
		}
		if request.ProjectId == nil || *request.ProjectId != "prj-123456" {
			t.Errorf("expected project_id to be passed to the API request")
		}
		resp := &makersv20261113.DescribeDeploymentsResponse{Response: &makersv20261113.DescribeDeploymentsResponseParams{
			TotalCount: ptrInt64Makers(1),
			Deployments: []*makersv20261113.Deployment{
				{
					DeploymentId:   ptrStrMakers("dpl-abcdefg"),
					ProjectId:      ptrStrMakers("prj-123456"),
					Status:         ptrStrMakers("Success"),
					OriginalStatus: ptrStrMakers("Process"),
					Code:           ptrInt64Makers(0),
					BuildCost:      ptrStrMakers("1200"),
					ViaMeta:        ptrStrMakers("Github"),
					Env:            ptrStrMakers("production"),
					RepoBranch:     ptrStrMakers("main"),
					RepoCommitHash: ptrStrMakers("abc123"),
					RepoCommitMsg:  ptrStrMakers("fix bug"),
					PreviewUrl:     ptrStrMakers("https://preview.example.com"),
					ProjectUrl:     ptrStrMakers("https://example.com"),
					UsedInProd:     ptrBoolMakers(true),
					CoverUrl:       ptrStrMakers("https://cover.example.com"),
					MetaData:       ptrStrMakers("{\"routes\":[]}"),
					DplRepoUrl:     ptrStrMakers("https://github.com/example/repo"),
					CreatedOn:      ptrStrMakers("2024-01-01 00:00:00"),
					ModifiedOn:     ptrStrMakers("2024-01-02 00:00:00"),
				},
			},
			RequestId: ptrStrMakers("fake-request-id"),
		}}
		return resp, nil
	})

	ds := svcmakers.NewMakersDeploymentsDataSource().(*svcmakers.MakersDeploymentsDataSource)
	setMakersDeploymentsClient(ds, client)

	sch := makersDeploymentsSchema(t)
	raw := buildMakersDeploymentsRawValue("zone-abcdefg", "prj-123456")
	req := newMakersDeploymentsReadRequest(sch, raw)
	resp := newMakersDeploymentsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcmakers.MakersDeploymentsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, "zone-abcdefg", result.ZoneId.ValueString())
	assert.Equal(t, "prj-123456", result.ProjectId.ValueString())
	assert.Equal(t, int64(1), result.TotalCount.ValueInt64())
	assert.False(t, result.Deployments.IsNull())

	var items []svcmakers.MakersDeploymentItemModel
	resp.Diagnostics.Append(result.Deployments.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)

	// assert every leaf field
	assert.Equal(t, "dpl-abcdefg", items[0].DeploymentId.ValueString())
	assert.Equal(t, "prj-123456", items[0].ProjectId.ValueString())
	assert.Equal(t, "Success", items[0].Status.ValueString())
	assert.Equal(t, "Process", items[0].OriginalStatus.ValueString())
	assert.Equal(t, int64(0), items[0].Code.ValueInt64())
	assert.Equal(t, "1200", items[0].BuildCost.ValueString())
	assert.Equal(t, "Github", items[0].ViaMeta.ValueString())
	assert.Equal(t, "production", items[0].Env.ValueString())
	assert.Equal(t, "main", items[0].RepoBranch.ValueString())
	assert.Equal(t, "abc123", items[0].RepoCommitHash.ValueString())
	assert.Equal(t, "fix bug", items[0].RepoCommitMsg.ValueString())
	assert.Equal(t, "https://preview.example.com", items[0].PreviewUrl.ValueString())
	assert.Equal(t, "https://example.com", items[0].ProjectUrl.ValueString())
	assert.Equal(t, true, items[0].UsedInProd.ValueBool())
	assert.Equal(t, "https://cover.example.com", items[0].CoverUrl.ValueString())
	assert.Equal(t, "{\"routes\":[]}", items[0].MetaData.ValueString())
	assert.Equal(t, "https://github.com/example/repo", items[0].DplRepoUrl.ValueString())
	assert.Equal(t, "2024-01-01 00:00:00", items[0].CreatedOn.ValueString())
	assert.Equal(t, "2024-01-02 00:00:00", items[0].ModifiedOn.ValueString())
}

// TestMakersDeployments_ReadEmpty covers a project that has no deployment
// records. An empty result is a legitimate answer, not a failure: the data
// source must succeed with total_count = 0 and an empty deployments list.
func TestMakersDeployments_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	makersClient := &makersv20261113.Client{}
	patches.ApplyMethodReturn(client, "UseMakersClient", makersClient)

	patches.ApplyMethodFunc(makersClient, "DescribeDeploymentsWithContext", func(_ context.Context, request *makersv20261113.DescribeDeploymentsRequest) (*makersv20261113.DescribeDeploymentsResponse, error) {
		resp := &makersv20261113.DescribeDeploymentsResponse{Response: &makersv20261113.DescribeDeploymentsResponseParams{
			TotalCount:  ptrInt64Makers(0),
			Deployments: []*makersv20261113.Deployment{},
			RequestId:   ptrStrMakers("fake-request-id"),
		}}
		return resp, nil
	})

	ds := svcmakers.NewMakersDeploymentsDataSource().(*svcmakers.MakersDeploymentsDataSource)
	setMakersDeploymentsClient(ds, client)

	sch := makersDeploymentsSchema(t)
	raw := buildMakersDeploymentsRawValue("zone-abcdefg", "prj-empty")
	req := newMakersDeploymentsReadRequest(sch, raw)
	resp := newMakersDeploymentsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcmakers.MakersDeploymentsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, int64(0), result.TotalCount.ValueInt64())
	assert.False(t, result.Deployments.IsNull())

	var items []svcmakers.MakersDeploymentItemModel
	resp.Diagnostics.Append(result.Deployments.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

// TestMakersDeployments_ReadApiError covers an API-level failure: the data
// source must surface the error through diagnostics instead of silently
// persisting empty state.
func TestMakersDeployments_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	makersClient := &makersv20261113.Client{}
	patches.ApplyMethodReturn(client, "UseMakersClient", makersClient)

	patches.ApplyMethodFunc(makersClient, "DescribeDeploymentsWithContext", func(_ context.Context, request *makersv20261113.DescribeDeploymentsRequest) (*makersv20261113.DescribeDeploymentsResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svcmakers.NewMakersDeploymentsDataSource().(*svcmakers.MakersDeploymentsDataSource)
	setMakersDeploymentsClient(ds, client)

	sch := makersDeploymentsSchema(t)
	raw := buildMakersDeploymentsRawValue("zone-abcdefg", "prj-123456")
	req := newMakersDeploymentsReadRequest(sch, raw)
	resp := newMakersDeploymentsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

// TestMakersDeployments_ReadNilResponse covers a nil response returned without
// an error: this must be reported as an error, not treated as an empty result.
func TestMakersDeployments_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	makersClient := &makersv20261113.Client{}
	patches.ApplyMethodReturn(client, "UseMakersClient", makersClient)

	patches.ApplyMethodFunc(makersClient, "DescribeDeploymentsWithContext", func(_ context.Context, request *makersv20261113.DescribeDeploymentsRequest) (*makersv20261113.DescribeDeploymentsResponse, error) {
		return nil, nil
	})

	ds := svcmakers.NewMakersDeploymentsDataSource().(*svcmakers.MakersDeploymentsDataSource)
	setMakersDeploymentsClient(ds, client)

	sch := makersDeploymentsSchema(t)
	raw := buildMakersDeploymentsRawValue("zone-abcdefg", "prj-123456")
	req := newMakersDeploymentsReadRequest(sch, raw)
	resp := newMakersDeploymentsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for nil response")
}
