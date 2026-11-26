package cos_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cos"
)

const (
	cosVectorIndexesTestBucketName = "tf-example-vector-1250000000"
	cosVectorIndexesTestPrefix     = "tf-example-index"
)

// cosVectorIndexesTestResponse builds a valid *cos.Response whose embedded
// *http.Response is populated (a bare &cos.Response{} dereferences nil in the
// service DEBUG log).
func cosVectorIndexesTestResponse() *cos.Response {
	return &cos.Response{Response: &http.Response{StatusCode: http.StatusOK}}
}

// cosVectorIndexesSchema returns the framework data source schema used to
// build the raw config/state value.
func cosVectorIndexesSchema(t *testing.T) schema.Schema {
	d := svcos.NewCosVectorIndexesDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// cosVectorIndexesItemType builds the tftypes.Object type for an indexes
// element.
func cosVectorIndexesItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"creation_time":      tftypes.Number,
			"index_qcs":          tftypes.String,
			"index_name":         tftypes.String,
			"vector_bucket_name": tftypes.String,
		},
	}
}

// cosVectorIndexesRootType builds the tftypes.Object type for the data source
// root, mirroring every schema attribute.
func cosVectorIndexesRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                 tftypes.String,
			"vector_bucket_name": tftypes.String,
			"prefix":             tftypes.String,
			"indexes":            tftypes.List{ElementType: cosVectorIndexesItemType()},
		},
	}
}

// buildCosVectorIndexesRawValue builds the raw config value with the given
// query inputs; every other attribute is unset (nil).
func buildCosVectorIndexesRawValue(bucketName, prefix interface{}) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":                 tftypes.NewValue(tftypes.String, nil),
		"vector_bucket_name": tftypes.NewValue(tftypes.String, bucketName),
		"prefix":             tftypes.NewValue(tftypes.String, prefix),
		"indexes":            tftypes.NewValue(tftypes.List{ElementType: cosVectorIndexesItemType()}, nil),
	}

	return tftypes.NewValue(cosVectorIndexesRootType(), attrs)
}

// patchCosVectorIndexesSdkClient stubs the vector client accessor and returns
// the mock SDK client whose methods the test patches next (the two-level
// mode).
func patchCosVectorIndexesSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *cos.Client {
	vectorClient := cos.NewClient(&cos.BaseURL{}, nil)
	patches.ApplyMethodReturn(client, "UseTencentCosVectorClient", vectorClient)

	return vectorClient
}

// newCosVectorIndexesReadRequest wires the raw config into a ReadRequest.
func newCosVectorIndexesReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch

	return req
}

// newCosVectorIndexesReadResponse prepares a ReadResponse whose State carries
// the schema so State.Set reflects into a raw value.
func newCosVectorIndexesReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	return resp
}

// TestCosVectorIndexes_MetadataAndSchema validates metadata and schema
// definition.
func TestCosVectorIndexes_MetadataAndSchema(t *testing.T) {
	d := svcos.NewCosVectorIndexesDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_cos_vector_indexes", metaResp.TypeName)

	sch := cosVectorIndexesSchema(t)
	attrs := sch.Attributes

	assert.True(t, attrs["vector_bucket_name"].(schema.StringAttribute).Required)
	assert.True(t, attrs["prefix"].(schema.StringAttribute).Optional)
	assert.True(t, attrs["indexes"].(schema.ListNestedAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "max_results")
	assert.NotContains(t, attrs, "next_token")
	assert.NotContains(t, attrs, "total_count")
}

// TestCosVectorIndexes_ReadBasic tests a successful read with one record and
// asserts every leaf field is read back.
func TestCosVectorIndexes_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vectorClient := patchCosVectorIndexesSdkClient(patches, client)
	patches.ApplyMethodFunc(vectorClient.Vector, "ListIndexes", func(ctx context.Context, opt *cos.ListIndexesOptions) (*cos.ListIndexesResult, *cos.Response, error) {
		if opt.VectorBucketName != cosVectorIndexesTestBucketName {
			t.Errorf("expected vector_bucket_name %q to be passed to the API request, got %q", cosVectorIndexesTestBucketName, opt.VectorBucketName)
		}
		if opt.Prefix != cosVectorIndexesTestPrefix {
			t.Errorf("expected prefix %q to be passed to the API request, got %q", cosVectorIndexesTestPrefix, opt.Prefix)
		}

		return &cos.ListIndexesResult{
			Indexes: []cos.IndexBrief{
				{
					CreationTime:     1735449900,
					IndexName:        "tf-example-index",
					IndexQcs:         "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000/index/tf-example-index",
					VectorBucketName: cosVectorIndexesTestBucketName,
				},
			},
		}, cosVectorIndexesTestResponse(), nil
	})

	ds := svcos.NewCosVectorIndexesDataSource().(*svcos.CosVectorIndexesDataSource)
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)

	sch := cosVectorIndexesSchema(t)
	raw := buildCosVectorIndexesRawValue(cosVectorIndexesTestBucketName, cosVectorIndexesTestPrefix)
	req := newCosVectorIndexesReadRequest(sch, raw)
	resp := newCosVectorIndexesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcos.CosVectorIndexesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())

	assert.False(t, result.Indexes.IsNull())

	var items []svcos.CosVectorIndexItemModel
	resp.Diagnostics.Append(result.Indexes.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)

	assert.Equal(t, int64(1735449900), items[0].CreationTime.ValueInt64())
	assert.Equal(t, "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000/index/tf-example-index", items[0].IndexQcs.ValueString())
	assert.Equal(t, "tf-example-index", items[0].IndexName.ValueString())
	assert.Equal(t, cosVectorIndexesTestBucketName, items[0].VectorBucketName.ValueString())
}

// TestCosVectorIndexes_ReadEmpty covers an empty result: a legitimate answer,
// not a failure — the data source must succeed with an empty list.
func TestCosVectorIndexes_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vectorClient := patchCosVectorIndexesSdkClient(patches, client)
	patches.ApplyMethodFunc(vectorClient.Vector, "ListIndexes", func(ctx context.Context, opt *cos.ListIndexesOptions) (*cos.ListIndexesResult, *cos.Response, error) {
		return &cos.ListIndexesResult{}, cosVectorIndexesTestResponse(), nil
	})

	ds := svcos.NewCosVectorIndexesDataSource().(*svcos.CosVectorIndexesDataSource)
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)

	sch := cosVectorIndexesSchema(t)
	raw := buildCosVectorIndexesRawValue(cosVectorIndexesTestBucketName, nil)
	req := newCosVectorIndexesReadRequest(sch, raw)
	resp := newCosVectorIndexesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcos.CosVectorIndexesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.Indexes.IsNull())

	var items []svcos.CosVectorIndexItemModel
	resp.Diagnostics.Append(result.Indexes.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

// TestCosVectorIndexes_ReadApiError covers an API-level failure: the error
// must surface through diagnostics instead of silently persisting empty
// state.
func TestCosVectorIndexes_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vectorClient := patchCosVectorIndexesSdkClient(patches, client)
	patches.ApplyMethodFunc(vectorClient.Vector, "ListIndexes", func(ctx context.Context, opt *cos.ListIndexesOptions) (*cos.ListIndexesResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("mocked API error")
	})

	ds := svcos.NewCosVectorIndexesDataSource().(*svcos.CosVectorIndexesDataSource)
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)

	sch := cosVectorIndexesSchema(t)
	raw := buildCosVectorIndexesRawValue(cosVectorIndexesTestBucketName, nil)
	req := newCosVectorIndexesReadRequest(sch, raw)
	resp := newCosVectorIndexesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

// TestCosVectorIndexesService_ReadAllPages verifies that the service layer
// loops through every page with the fixed page size and returns the complete
// set (the total count and the paging token stay internal).
func TestCosVectorIndexesService_ReadAllPages(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vectorClient := patchCosVectorIndexesSdkClient(patches, client)

	callCount := 0
	patches.ApplyMethodFunc(vectorClient.Vector, "ListIndexes", func(ctx context.Context, opt *cos.ListIndexesOptions) (*cos.ListIndexesResult, *cos.Response, error) {
		callCount++
		// The service must drive pagination itself with the page size 100.
		if opt.MaxResults != 100 {
			t.Errorf("expected MaxResults 100 on call %d, got %d", callCount, opt.MaxResults)
		}
		switch callCount {
		case 1:
			if opt.NextToken != "" {
				t.Errorf("expected empty NextToken on first call, got %q", opt.NextToken)
			}

			indexes := make([]cos.IndexBrief, 0, 100)
			for i := 0; i < 100; i++ {
				indexes = append(indexes, cos.IndexBrief{
					IndexName:        fmt.Sprintf("tf-example-index-page1-%d", i),
					VectorBucketName: cosVectorIndexesTestBucketName,
				})
			}

			return &cos.ListIndexesResult{
				Indexes:   indexes,
				NextToken: "token-page-2",
			}, cosVectorIndexesTestResponse(), nil
		case 2:
			if opt.NextToken != "token-page-2" {
				t.Errorf("expected the previous NextToken on second call, got %q", opt.NextToken)
			}

			return &cos.ListIndexesResult{
				Indexes: []cos.IndexBrief{
					{
						IndexName:        "tf-example-index-page2-0",
						VectorBucketName: cosVectorIndexesTestBucketName,
					},
				},
			}, cosVectorIndexesTestResponse(), nil
		default:
			t.Errorf("unexpected extra API call %d", callCount)

			return &cos.ListIndexesResult{}, cosVectorIndexesTestResponse(), nil
		}
	})

	service := svcos.NewCosService(client)
	listIndexesRequest := &cos.ListIndexesOptions{
		VectorBucketName: cosVectorIndexesTestBucketName,
	}

	indexes, err := service.ListIndexes(context.Background(), listIndexesRequest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
	if len(indexes) != 101 {
		t.Fatalf("expected 101 indexes across pages, got %d", len(indexes))
	}
	if indexes[0].IndexName != "tf-example-index-page1-0" {
		t.Errorf("unexpected first index name: %q", indexes[0].IndexName)
	}
	if indexes[100].IndexName != "tf-example-index-page2-0" {
		t.Errorf("unexpected last index name: %q", indexes[100].IndexName)
	}
}
