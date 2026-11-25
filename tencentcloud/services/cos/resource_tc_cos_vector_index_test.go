package cos_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cos"
)

const (
	cosVectorIndexTestBucketName     = "tf-example-vector-1250000000"
	cosVectorIndexTestIndexName      = "tf-example-index"
	cosVectorIndexTestQcs            = "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000/index/tf-example-index"
	cosVectorIndexTestCompositeId    = "tf-example-vector-1250000000#tf-example-index"
	cosVectorIndexTestDataType       = "float32"
	cosVectorIndexTestDimension      = int64(768)
	cosVectorIndexTestDistanceMetric = "euclidean"
	cosVectorIndexTestMetaKeyA       = "examplekey1"
	cosVectorIndexTestMetaKeyB       = "examplekey2"
	cosVectorIndexTestDriftedMetaKey = "examplekey-drifted"
)

// cosVectorIndexMetadataAttrTypes mirrors the attribute types of the
// metadata_configuration object declared by the resource schema.
func cosVectorIndexMetadataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"non_filterable_metadata_keys": types.ListType{ElemType: types.StringType},
	}
}

// cosVectorIndexTestResponse builds a valid *cos.Response whose embedded
// *http.Response is populated (a bare &cos.Response{} dereferences nil in the
// service DEBUG log).
func cosVectorIndexTestResponse() *cos.Response {
	return &cos.Response{Response: &http.Response{StatusCode: http.StatusOK}}
}

// cosVectorIndexTestItem builds the GetIndex response item.
func cosVectorIndexTestItem() *cos.IndexInfo {
	return &cos.IndexInfo{
		VectorBucketName: cosVectorIndexTestBucketName,
		IndexName:        cosVectorIndexTestIndexName,
		IndexQcs:         cosVectorIndexTestQcs,
		DataType:         cosVectorIndexTestDataType,
		Dimension:        int(cosVectorIndexTestDimension),
		DistanceMetric:   cosVectorIndexTestDistanceMetric,
	}
}

// cosVectorIndexTestModel builds the terraform model used for plan/state.
// Every unset attribute carries an EXPLICIT TYPED null value — the raw zero
// value breaks tfsdk Plan/State reflection with "Value Conversion Error".
func cosVectorIndexTestModel() svcos.CosVectorIndexResourceModel {
	return svcos.CosVectorIndexResourceModel{
		Id:                    types.StringValue(cosVectorIndexTestCompositeId),
		VectorBucketName:      types.StringValue(cosVectorIndexTestBucketName),
		IndexName:             types.StringValue(cosVectorIndexTestIndexName),
		DataType:              types.StringValue(cosVectorIndexTestDataType),
		Dimension:             types.Int64Value(cosVectorIndexTestDimension),
		DistanceMetric:        types.StringValue(cosVectorIndexTestDistanceMetric),
		MetadataConfiguration: types.ObjectNull(cosVectorIndexMetadataAttrTypes()),
	}
}

// cosVectorIndexTestModelWithMetadata builds the model with the
// metadata_configuration object set.
func cosVectorIndexTestModelWithMetadata() svcos.CosVectorIndexResourceModel {
	model := cosVectorIndexTestModel()
	model.MetadataConfiguration = types.ObjectValueMust(
		cosVectorIndexMetadataAttrTypes(),
		map[string]attr.Value{
			"non_filterable_metadata_keys": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue(cosVectorIndexTestMetaKeyA),
				types.StringValue(cosVectorIndexTestMetaKeyB),
			}),
		},
	)

	return model
}

// assertCosVectorIndexModel asserts every model field.
func assertCosVectorIndexModel(t *testing.T, m svcos.CosVectorIndexResourceModel) {
	t.Helper()

	if m.Id.ValueString() != cosVectorIndexTestCompositeId {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), cosVectorIndexTestCompositeId)
	}
	if m.VectorBucketName.ValueString() != cosVectorIndexTestBucketName {
		t.Errorf("vector_bucket_name = %q, want %q", m.VectorBucketName.ValueString(), cosVectorIndexTestBucketName)
	}
	if m.IndexName.ValueString() != cosVectorIndexTestIndexName {
		t.Errorf("index_name = %q, want %q", m.IndexName.ValueString(), cosVectorIndexTestIndexName)
	}
	if m.DataType.ValueString() != cosVectorIndexTestDataType {
		t.Errorf("data_type = %q, want %q", m.DataType.ValueString(), cosVectorIndexTestDataType)
	}
	if m.Dimension.ValueInt64() != cosVectorIndexTestDimension {
		t.Errorf("dimension = %d, want %d", m.Dimension.ValueInt64(), cosVectorIndexTestDimension)
	}
	if m.DistanceMetric.ValueString() != cosVectorIndexTestDistanceMetric {
		t.Errorf("distance_metric = %q, want %q", m.DistanceMetric.ValueString(), cosVectorIndexTestDistanceMetric)
	}
}

// assertCosVectorIndexMetadataKeys asserts the metadata_configuration keys.
func assertCosVectorIndexMetadataKeys(t *testing.T, m svcos.CosVectorIndexResourceModel, wantKeys ...string) {
	t.Helper()

	if m.MetadataConfiguration.IsNull() {
		t.Errorf("metadata_configuration is null, want keys %v", wantKeys)
		return
	}

	keysValue, ok := m.MetadataConfiguration.Attributes()["non_filterable_metadata_keys"].(types.List)
	if !ok {
		t.Fatalf("metadata_configuration.non_filterable_metadata_keys is not a list value")
	}

	var keys []string
	if diags := keysValue.ElementsAs(context.Background(), &keys, false); diags.HasError() {
		t.Fatalf("non_filterable_metadata_keys: %v", diags.Errors())
	}
	if len(keys) != len(wantKeys) {
		t.Fatalf("non_filterable_metadata_keys = %v, want %v", keys, wantKeys)
	}
	for i, key := range wantKeys {
		if keys[i] != key {
			t.Errorf("non_filterable_metadata_keys[%d] = %q, want %q", i, keys[i], key)
		}
	}
}

// patchCosVectorIndexSdkClient stubs the vector client accessor and returns
// the mock SDK client whose methods the test patches next (the two-level
// mode).
func patchCosVectorIndexSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *cos.Client {
	vectorClient := cos.NewClient(&cos.BaseURL{}, nil)
	patches.ApplyMethodReturn(client, "UseTencentCosVectorClient", vectorClient)

	return vectorClient
}

// cosVectorIndexTestContext wires the resource to a configured client; the
// SDK client methods are stubbed per test case.
type cosVectorIndexTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newCosVectorIndexTestContext(t *testing.T) *cosVectorIndexTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svcos.NewCosVectorIndexResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}

	var configureResp resource.ConfigureResponse
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}},
		&configureResp,
	)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("failed to configure resource: %v", configureResp.Diagnostics.Errors())
	}

	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	return &cosVectorIndexTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *cosVectorIndexTestContext) plan(t *testing.T, m svcos.CosVectorIndexResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *cosVectorIndexTestContext) state(t *testing.T, m svcos.CosVectorIndexResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// TestCosVectorIndexCreateSuccess verifies that Create sends all request
// fields (including the nested metadata configuration), tracks the composite
// id and refreshes the state.
func TestCosVectorIndexCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)

	var captured *cos.CreateIndexOptions
	patches.ApplyMethodFunc(vectorClient.Vector, "CreateIndex", func(ctx context.Context, opt *cos.CreateIndexOptions) (*cos.CreateIndexResult, *cos.Response, error) {
		captured = opt

		return &cos.CreateIndexResult{
			IndexQcs: cosVectorIndexTestQcs,
		}, cosVectorIndexTestResponse(), nil
	})
	patches.ApplyMethodFunc(vectorClient.Vector, "GetIndex", func(ctx context.Context, opt *cos.GetIndexOptions) (*cos.GetIndexResult, *cos.Response, error) {
		item := cosVectorIndexTestItem()
		item.MetadataConfiguration = &cos.MetadataConfiguration{
			NonFilterableMetadataKeys: []string{cosVectorIndexTestMetaKeyA, cosVectorIndexTestMetaKeyB},
		}

		return &cos.GetIndexResult{
			Index: item,
		}, cosVectorIndexTestResponse(), nil
	})

	model := cosVectorIndexTestModelWithMetadata()
	// The id is Computed: it is unknown in a real create plan.
	model.Id = types.StringNull()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateIndex request is nil")
	}
	if captured.VectorBucketName != cosVectorIndexTestBucketName {
		t.Errorf("create request vector_bucket_name = %q, want %q", captured.VectorBucketName, cosVectorIndexTestBucketName)
	}
	if captured.IndexName != cosVectorIndexTestIndexName {
		t.Errorf("create request index_name = %q, want %q", captured.IndexName, cosVectorIndexTestIndexName)
	}
	if captured.DataType != cosVectorIndexTestDataType {
		t.Errorf("create request data_type = %q, want %q", captured.DataType, cosVectorIndexTestDataType)
	}
	if captured.Dimension != int(cosVectorIndexTestDimension) {
		t.Errorf("create request dimension = %d, want %d", captured.Dimension, cosVectorIndexTestDimension)
	}
	if captured.DistanceMetric != cosVectorIndexTestDistanceMetric {
		t.Errorf("create request distance_metric = %q, want %q", captured.DistanceMetric, cosVectorIndexTestDistanceMetric)
	}
	if captured.MetadataConfiguration == nil || len(captured.MetadataConfiguration.NonFilterableMetadataKeys) != 2 ||
		captured.MetadataConfiguration.NonFilterableMetadataKeys[0] != cosVectorIndexTestMetaKeyA ||
		captured.MetadataConfiguration.NonFilterableMetadataKeys[1] != cosVectorIndexTestMetaKeyB {
		t.Errorf("create request metadata_configuration = %v, want keys [%q %q]", captured.MetadataConfiguration, cosVectorIndexTestMetaKeyA, cosVectorIndexTestMetaKeyB)
	}

	var got svcos.CosVectorIndexResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorIndexModel(t, got)
	assertCosVectorIndexMetadataKeys(t, got, cosVectorIndexTestMetaKeyA, cosVectorIndexTestMetaKeyB)
}

// TestCosVectorIndexCreateWithoutMetadata verifies that a null
// metadata_configuration is transmitted as "absent" (nil SDK pointer).
func TestCosVectorIndexCreateWithoutMetadata(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)

	var captured *cos.CreateIndexOptions
	patches.ApplyMethodFunc(vectorClient.Vector, "CreateIndex", func(ctx context.Context, opt *cos.CreateIndexOptions) (*cos.CreateIndexResult, *cos.Response, error) {
		captured = opt

		return &cos.CreateIndexResult{
			IndexQcs: cosVectorIndexTestQcs,
		}, cosVectorIndexTestResponse(), nil
	})
	patches.ApplyMethodFunc(vectorClient.Vector, "GetIndex", func(ctx context.Context, opt *cos.GetIndexOptions) (*cos.GetIndexResult, *cos.Response, error) {
		return &cos.GetIndexResult{
			Index: cosVectorIndexTestItem(),
		}, cosVectorIndexTestResponse(), nil
	})

	model := cosVectorIndexTestModel()
	model.Id = types.StringNull()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateIndex request is nil")
	}
	if captured.MetadataConfiguration != nil {
		t.Errorf("create request metadata_configuration = %v, want nil", captured.MetadataConfiguration)
	}

	var got svcos.CosVectorIndexResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorIndexModel(t, got)
	if !got.MetadataConfiguration.IsNull() {
		t.Error("metadata_configuration should stay null when the cloud reports none")
	}
}

// TestCosVectorIndexCreateApiError verifies API failures surface as errors.
func TestCosVectorIndexCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "CreateIndex", func(ctx context.Context, opt *cos.CreateIndexOptions) (*cos.CreateIndexResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("create failed")
	})

	model := cosVectorIndexTestModel()
	model.Id = types.StringNull()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestCosVectorIndexReadSuccess verifies Read addresses the index by both id
// segments and refreshes the leaf fields.
func TestCosVectorIndexReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)

	var gotBucketName, gotIndexName string
	patches.ApplyMethodFunc(vectorClient.Vector, "GetIndex", func(ctx context.Context, opt *cos.GetIndexOptions) (*cos.GetIndexResult, *cos.Response, error) {
		gotBucketName = opt.VectorBucketName
		gotIndexName = opt.IndexName

		return &cos.GetIndexResult{
			Index: cosVectorIndexTestItem(),
		}, cosVectorIndexTestResponse(), nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorIndexTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorIndexTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotBucketName != cosVectorIndexTestBucketName || gotIndexName != cosVectorIndexTestIndexName {
		t.Errorf("GetIndex called with (%q, %q), want (%q, %q)", gotBucketName, gotIndexName, cosVectorIndexTestBucketName, cosVectorIndexTestIndexName)
	}

	var got svcos.CosVectorIndexResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorIndexModel(t, got)
}

// TestCosVectorIndexReadDetectsMetadataDrift verifies a user-configurable
// nested object is refreshed from the cloud: the state holds a null
// metadata_configuration, the cloud now answers two keys, and Read must
// record the cloud value (not keep the previous one), so a change made
// outside Terraform shows up in the next plan.
func TestCosVectorIndexReadDetectsMetadataDrift(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetIndex", func(ctx context.Context, opt *cos.GetIndexOptions) (*cos.GetIndexResult, *cos.Response, error) {
		drifted := cosVectorIndexTestItem()
		drifted.MetadataConfiguration = &cos.MetadataConfiguration{
			NonFilterableMetadataKeys: []string{cosVectorIndexTestDriftedMetaKey},
		}

		return &cos.GetIndexResult{
			Index: drifted,
		}, cosVectorIndexTestResponse(), nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorIndexTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorIndexTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got svcos.CosVectorIndexResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorIndexMetadataKeys(t, got, cosVectorIndexTestDriftedMetaKey)
}

// TestCosVectorIndexReadNotFound verifies that a documented HTTP 404 answer
// drops the resource from state (confirmed not-found).
func TestCosVectorIndexReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetIndex", func(ctx context.Context, opt *cos.GetIndexOptions) (*cos.GetIndexResult, *cos.Response, error) {
		return nil, nil, &cos.VectorErrorResponse{
			Response: &http.Response{StatusCode: http.StatusNotFound},
		}
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorIndexTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorIndexTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the index is confirmed gone")
	}
}

// TestCosVectorIndexReadApiError verifies non-404 API failures surface as
// errors instead of dropping state.
func TestCosVectorIndexReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetIndex", func(ctx context.Context, opt *cos.GetIndexOptions) (*cos.GetIndexResult, *cos.Response, error) {
		return nil, nil, &cos.VectorErrorResponse{
			Response: &http.Response{StatusCode: http.StatusForbidden},
		}
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorIndexTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorIndexTestModel())}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on a non-404 API error")
	}
}

// TestCosVectorIndexUpdateNoop verifies that Update persists the plan without
// calling any write API: the cloud has no update API for vector indexes and
// every attribute is RequiresReplace.
func TestCosVectorIndexUpdateNoop(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "CreateIndex", func(ctx context.Context, opt *cos.CreateIndexOptions) (*cos.CreateIndexResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("the create API must not be called")
	})
	patches.ApplyMethodFunc(vectorClient.Vector, "GetIndex", func(ctx context.Context, opt *cos.GetIndexOptions) (*cos.GetIndexResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("the read API must not be called")
	})

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, cosVectorIndexTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, cosVectorIndexTestModel()),
			State: ctx.state(t, cosVectorIndexTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	var got svcos.CosVectorIndexResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorIndexModel(t, got)
}

// TestCosVectorIndexDeleteSuccess verifies Delete sends the delete request
// addressed by both id segments.
func TestCosVectorIndexDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)

	var captured *cos.DeleteIndexOptions
	patches.ApplyMethodFunc(vectorClient.Vector, "DeleteIndex", func(ctx context.Context, opt *cos.DeleteIndexOptions) (*cos.Response, error) {
		captured = opt

		return cosVectorIndexTestResponse(), nil
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, cosVectorIndexTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteIndex request is nil")
	}
	if captured.VectorBucketName != cosVectorIndexTestBucketName || captured.IndexName != cosVectorIndexTestIndexName {
		t.Errorf("delete request = (%q, %q), want (%q, %q)", captured.VectorBucketName, captured.IndexName, cosVectorIndexTestBucketName, cosVectorIndexTestIndexName)
	}
}

// TestCosVectorIndexDeleteApiError verifies delete failures surface as
// errors.
func TestCosVectorIndexDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)
	vectorClient := patchCosVectorIndexSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "DeleteIndex", func(ctx context.Context, opt *cos.DeleteIndexOptions) (*cos.Response, error) {
		return nil, fmt.Errorf("delete failed")
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, cosVectorIndexTestModel())}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}

// TestCosVectorIndexImportState verifies the composite id is split into the
// id segments on import.
func TestCosVectorIndexImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)

	importer, ok := ctx.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}

	var importResp resource.ImportStateResponse
	importResp.State = tfsdk.State{
		Schema: ctx.schema,
		Raw:    tftypes.NewValue(ctx.schema.Type().TerraformType(context.Background()), nil),
	}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: cosVectorIndexTestCompositeId},
		&importResp,
	)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", importResp.Diagnostics.Errors())
	}

	var got svcos.CosVectorIndexResourceModel
	if diags := importResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != cosVectorIndexTestCompositeId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), cosVectorIndexTestCompositeId)
	}
	if got.VectorBucketName.ValueString() != cosVectorIndexTestBucketName {
		t.Errorf("vector_bucket_name = %q, want %q", got.VectorBucketName.ValueString(), cosVectorIndexTestBucketName)
	}
	if got.IndexName.ValueString() != cosVectorIndexTestIndexName {
		t.Errorf("index_name = %q, want %q", got.IndexName.ValueString(), cosVectorIndexTestIndexName)
	}
}

// TestCosVectorIndexImportStateInvalidId verifies a malformed import id is
// rejected with an error instead of a panic.
func TestCosVectorIndexImportStateInvalidId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorIndexTestContext(t)

	importer, ok := ctx.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}

	var importResp resource.ImportStateResponse
	importResp.State = tfsdk.State{
		Schema: ctx.schema,
		Raw:    tftypes.NewValue(ctx.schema.Type().TerraformType(context.Background()), nil),
	}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: cosVectorIndexTestIndexName},
		&importResp,
	)
	if !importResp.Diagnostics.HasError() {
		t.Fatal("ImportState should fail on an id without the composite format")
	}
}
