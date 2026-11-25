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
	cosVectorBucketsTestPrefix = "tf-example-vector"
)

// cosVectorBucketsTestResponse builds a valid *cos.Response whose embedded
// *http.Response is populated (a bare &cos.Response{} dereferences nil in the
// service DEBUG log).
func cosVectorBucketsTestResponse() *cos.Response {
	return &cos.Response{Response: &http.Response{StatusCode: http.StatusOK}}
}

// cosVectorBucketsSchema returns the framework data source schema used to
// build the raw config/state value.
func cosVectorBucketsSchema(t *testing.T) schema.Schema {
	d := svcos.NewCosVectorBucketsDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// cosVectorBucketsItemType builds the tftypes.Object type for a
// vector_buckets element.
func cosVectorBucketsItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"creation_time":      tftypes.Number,
			"vector_bucket_qcs":  tftypes.String,
			"vector_bucket_name": tftypes.String,
		},
	}
}

// cosVectorBucketsRootType builds the tftypes.Object type for the data source
// root, mirroring every schema attribute.
func cosVectorBucketsRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":             tftypes.String,
			"prefix":         tftypes.String,
			"vector_buckets": tftypes.List{ElementType: cosVectorBucketsItemType()},
		},
	}
}

// buildCosVectorBucketsRawValue builds the raw config value with the given
// query inputs; every other attribute is unset (nil).
func buildCosVectorBucketsRawValue(prefix interface{}) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":             tftypes.NewValue(tftypes.String, nil),
		"prefix":         tftypes.NewValue(tftypes.String, prefix),
		"vector_buckets": tftypes.NewValue(tftypes.List{ElementType: cosVectorBucketsItemType()}, nil),
	}

	return tftypes.NewValue(cosVectorBucketsRootType(), attrs)
}

// patchCosVectorBucketsSdkClient stubs the vector client accessor and returns
// the mock SDK client whose methods the test patches next (the two-level
// mode).
func patchCosVectorBucketsSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *cos.Client {
	vectorClient := cos.NewClient(&cos.BaseURL{}, nil)
	patches.ApplyMethodReturn(client, "UseTencentCosVectorClient", vectorClient)

	return vectorClient
}

// newCosVectorBucketsReadRequest wires the raw config into a ReadRequest.
func newCosVectorBucketsReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch

	return req
}

// newCosVectorBucketsReadResponse prepares a ReadResponse whose State carries
// the schema so State.Set reflects into a raw value.
func newCosVectorBucketsReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	return resp
}

// TestCosVectorBuckets_MetadataAndSchema validates metadata and schema
// definition.
func TestCosVectorBuckets_MetadataAndSchema(t *testing.T) {
	d := svcos.NewCosVectorBucketsDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_cos_vector_buckets", metaResp.TypeName)

	sch := cosVectorBucketsSchema(t)
	attrs := sch.Attributes

	assert.True(t, attrs["prefix"].(schema.StringAttribute).Optional)
	assert.True(t, attrs["vector_buckets"].(schema.ListNestedAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "max_results")
	assert.NotContains(t, attrs, "next_token")
	assert.NotContains(t, attrs, "total_count")
}

// TestCosVectorBuckets_ReadBasic tests a successful read with one record and
// asserts every leaf field is read back.
func TestCosVectorBuckets_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vectorClient := patchCosVectorBucketsSdkClient(patches, client)
	patches.ApplyMethodFunc(vectorClient.Vector, "ListVectorBuckets", func(ctx context.Context, opt *cos.ListVectorBucketsOptions) (*cos.ListVectorBucketsResult, *cos.Response, error) {
		if opt.Prefix != cosVectorBucketsTestPrefix {
			t.Errorf("expected prefix %q to be passed to the API request, got %q", cosVectorBucketsTestPrefix, opt.Prefix)
		}

		return &cos.ListVectorBucketsResult{
			VectorBuckets: []cos.VectorBucketBrief{
				{
					CreationTime:     1735449000,
					VectorBucketName: "tf-example-vector-1250000000",
					VectorBucketQcs:  "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000",
				},
			},
		}, cosVectorBucketsTestResponse(), nil
	})

	ds := svcos.NewCosVectorBucketsDataSource().(*svcos.CosVectorBucketsDataSource)
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)

	sch := cosVectorBucketsSchema(t)
	raw := buildCosVectorBucketsRawValue(cosVectorBucketsTestPrefix)
	req := newCosVectorBucketsReadRequest(sch, raw)
	resp := newCosVectorBucketsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcos.CosVectorBucketsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())

	assert.False(t, result.VectorBuckets.IsNull())

	var items []svcos.CosVectorBucketItemModel
	resp.Diagnostics.Append(result.VectorBuckets.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)

	assert.Equal(t, int64(1735449000), items[0].CreationTime.ValueInt64())
	assert.Equal(t, "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000", items[0].VectorBucketQcs.ValueString())
	assert.Equal(t, "tf-example-vector-1250000000", items[0].VectorBucketName.ValueString())
}

// TestCosVectorBuckets_ReadEmpty covers an empty result: a legitimate answer,
// not a failure — the data source must succeed with an empty list.
func TestCosVectorBuckets_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vectorClient := patchCosVectorBucketsSdkClient(patches, client)
	patches.ApplyMethodFunc(vectorClient.Vector, "ListVectorBuckets", func(ctx context.Context, opt *cos.ListVectorBucketsOptions) (*cos.ListVectorBucketsResult, *cos.Response, error) {
		return &cos.ListVectorBucketsResult{}, cosVectorBucketsTestResponse(), nil
	})

	ds := svcos.NewCosVectorBucketsDataSource().(*svcos.CosVectorBucketsDataSource)
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)

	sch := cosVectorBucketsSchema(t)
	raw := buildCosVectorBucketsRawValue(nil)
	req := newCosVectorBucketsReadRequest(sch, raw)
	resp := newCosVectorBucketsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcos.CosVectorBucketsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.VectorBuckets.IsNull())

	var items []svcos.CosVectorBucketItemModel
	resp.Diagnostics.Append(result.VectorBuckets.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

// TestCosVectorBuckets_ReadApiError covers an API-level failure: the error
// must surface through diagnostics instead of silently persisting empty
// state.
func TestCosVectorBuckets_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vectorClient := patchCosVectorBucketsSdkClient(patches, client)
	patches.ApplyMethodFunc(vectorClient.Vector, "ListVectorBuckets", func(ctx context.Context, opt *cos.ListVectorBucketsOptions) (*cos.ListVectorBucketsResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("mocked API error")
	})

	ds := svcos.NewCosVectorBucketsDataSource().(*svcos.CosVectorBucketsDataSource)
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)

	sch := cosVectorBucketsSchema(t)
	raw := buildCosVectorBucketsRawValue(nil)
	req := newCosVectorBucketsReadRequest(sch, raw)
	resp := newCosVectorBucketsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

// TestCosVectorBucketsService_ReadAllPages verifies that the service layer
// loops through every page with the fixed page size and returns the complete
// set (the total count and the paging token stay internal).
func TestCosVectorBucketsService_ReadAllPages(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vectorClient := patchCosVectorBucketsSdkClient(patches, client)

	callCount := 0
	patches.ApplyMethodFunc(vectorClient.Vector, "ListVectorBuckets", func(ctx context.Context, opt *cos.ListVectorBucketsOptions) (*cos.ListVectorBucketsResult, *cos.Response, error) {
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

			buckets := make([]cos.VectorBucketBrief, 0, 100)
			for i := 0; i < 100; i++ {
				buckets = append(buckets, cos.VectorBucketBrief{
					VectorBucketName: fmt.Sprintf("tf-example-vector-page1-%d", i),
				})
			}

			return &cos.ListVectorBucketsResult{
				VectorBuckets: buckets,
				NextToken:     "token-page-2",
			}, cosVectorBucketsTestResponse(), nil
		case 2:
			if opt.NextToken != "token-page-2" {
				t.Errorf("expected the previous NextToken on second call, got %q", opt.NextToken)
			}

			return &cos.ListVectorBucketsResult{
				VectorBuckets: []cos.VectorBucketBrief{
					{VectorBucketName: "tf-example-vector-page2-0"},
				},
			}, cosVectorBucketsTestResponse(), nil
		default:
			t.Errorf("unexpected extra API call %d", callCount)

			return &cos.ListVectorBucketsResult{}, cosVectorBucketsTestResponse(), nil
		}
	})

	service := svcos.NewCosService(client)
	listVectorBucketsRequest := &cos.ListVectorBucketsOptions{
		Prefix: cosVectorBucketsTestPrefix,
	}

	buckets, err := service.ListVectorBuckets(context.Background(), listVectorBucketsRequest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
	if len(buckets) != 101 {
		t.Fatalf("expected 101 vector buckets across pages, got %d", len(buckets))
	}
	if buckets[0].VectorBucketName != "tf-example-vector-page1-0" {
		t.Errorf("unexpected first bucket name: %q", buckets[0].VectorBucketName)
	}
	if buckets[100].VectorBucketName != "tf-example-vector-page2-0" {
		t.Errorf("unexpected last bucket name: %q", buckets[100].VectorBucketName)
	}
}
