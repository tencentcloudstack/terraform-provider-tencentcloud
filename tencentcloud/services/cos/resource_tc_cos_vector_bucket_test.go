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

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cos"
)

const (
	cosVectorBucketTestName     = "tf-example-vector-1250000000"
	cosVectorBucketTestQcs      = "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000"
	cosVectorBucketTestSseType  = "AES256"
	cosVectorBucketTestCreation = int64(1735445700)
)

// cosVectorBucketEncryptionAttrTypes mirrors the attribute types of the
// encryption_configuration object declared by the resource schema.
func cosVectorBucketEncryptionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"sse_type": types.StringType,
	}
}

// cosVectorBucketTestItem builds the GetVectorBucket response item.
func cosVectorBucketTestItem() *cos.VectorBucketInfo {
	return &cos.VectorBucketInfo{
		VectorBucketName: cosVectorBucketTestName,
		VectorBucketQcs:  cosVectorBucketTestQcs,
		CreationTime:     cosVectorBucketTestCreation,
	}
}

// cosVectorBucketTestModel builds the terraform model used for plan/state.
// Every unset attribute carries an EXPLICIT TYPED null value — the raw zero
// value breaks tfsdk Plan/State reflection with "Value Conversion Error".
func cosVectorBucketTestModel() svcos.CosVectorBucketResourceModel {
	return svcos.CosVectorBucketResourceModel{
		Id:                      types.StringValue(cosVectorBucketTestName),
		VectorBucketName:        types.StringValue(cosVectorBucketTestName),
		EncryptionConfiguration: types.ObjectNull(cosVectorBucketEncryptionAttrTypes()),
	}
}

// cosVectorBucketTestModelEncrypted builds the model with the
// encryption_configuration object set.
func cosVectorBucketTestModelEncrypted() svcos.CosVectorBucketResourceModel {
	model := cosVectorBucketTestModel()
	model.EncryptionConfiguration = types.ObjectValueMust(
		cosVectorBucketEncryptionAttrTypes(),
		map[string]attr.Value{
			"sse_type": types.StringValue(cosVectorBucketTestSseType),
		},
	)

	return model
}

// assertCosVectorBucketModel asserts every model field.
func assertCosVectorBucketModel(t *testing.T, m svcos.CosVectorBucketResourceModel) {
	t.Helper()

	if m.Id.ValueString() != cosVectorBucketTestName {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), cosVectorBucketTestName)
	}
	if m.VectorBucketName.ValueString() != cosVectorBucketTestName {
		t.Errorf("vector_bucket_name = %q, want %q", m.VectorBucketName.ValueString(), cosVectorBucketTestName)
	}
}

// assertCosVectorBucketEncryption asserts the encryption_configuration object.
func assertCosVectorBucketEncryption(t *testing.T, m svcos.CosVectorBucketResourceModel, wantSseType string) {
	t.Helper()

	if m.EncryptionConfiguration.IsNull() {
		t.Errorf("encryption_configuration is null, want sse_type %q", wantSseType)
		return
	}

	sseType, ok := m.EncryptionConfiguration.Attributes()["sse_type"].(types.String)
	if !ok {
		t.Fatalf("encryption_configuration.sse_type is not a string value")
	}
	if sseType.ValueString() != wantSseType {
		t.Errorf("encryption_configuration.sse_type = %q, want %q", sseType.ValueString(), wantSseType)
	}
}

// cosVectorBucketTestResponse builds a valid *cos.Response whose embedded
// *http.Response is populated (a bare &cos.Response{} dereferences nil in the
// service DEBUG log).
func cosVectorBucketTestResponse() *cos.Response {
	return &cos.Response{Response: &http.Response{StatusCode: http.StatusOK}}
}

// patchCosVectorBucketSdkClient stubs the vector client accessor and returns
// the mock SDK client whose methods the test patches next (the two-level
// mode).
func patchCosVectorBucketSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *cos.Client {
	vectorClient := cos.NewClient(&cos.BaseURL{}, nil)
	patches.ApplyMethodReturn(client, "UseTencentCosVectorClient", vectorClient)

	return vectorClient
}

// cosVectorBucketTestContext wires the resource to a configured client; the
// SDK client methods are stubbed per test case.
type cosVectorBucketTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newCosVectorBucketTestContext(t *testing.T) *cosVectorBucketTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svcos.NewCosVectorBucketResource()
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

	return &cosVectorBucketTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *cosVectorBucketTestContext) plan(t *testing.T, m svcos.CosVectorBucketResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *cosVectorBucketTestContext) state(t *testing.T, m svcos.CosVectorBucketResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// TestCosVectorBucketCreateSuccess verifies that Create sends all request
// fields (including the nested encryption configuration), tracks the vector
// bucket name as the id and refreshes the state.
func TestCosVectorBucketCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)

	var captured *cos.CreateVectorBucketOptions
	patches.ApplyMethodFunc(vectorClient.Vector, "CreateVectorBucket", func(ctx context.Context, opt *cos.CreateVectorBucketOptions) (*cos.CreateVectorBucketResult, *cos.Response, error) {
		captured = opt

		return &cos.CreateVectorBucketResult{
			VectorBucketQcs: cosVectorBucketTestQcs,
		}, cosVectorBucketTestResponse(), nil
	})
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucket", func(ctx context.Context, opt *cos.GetVectorBucketOptions) (*cos.GetVectorBucketResult, *cos.Response, error) {
		item := cosVectorBucketTestItem()
		item.EncryptionConfiguration = &cos.VectorEncryptionConfig{
			SseType: cosVectorBucketTestSseType,
		}

		return &cos.GetVectorBucketResult{
			VectorBucket: item,
		}, cosVectorBucketTestResponse(), nil
	})

	model := cosVectorBucketTestModelEncrypted()
	// The id is Computed: it is unknown in a real create plan.
	model.Id = types.StringNull()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateVectorBucket request is nil")
	}
	if captured.VectorBucketName != cosVectorBucketTestName {
		t.Errorf("create request vector_bucket_name = %q, want %q", captured.VectorBucketName, cosVectorBucketTestName)
	}
	if captured.EncryptionConfiguration == nil || captured.EncryptionConfiguration.SseType != cosVectorBucketTestSseType {
		t.Errorf("create request encryption_configuration = %v, want sse_type %q", captured.EncryptionConfiguration, cosVectorBucketTestSseType)
	}

	var got svcos.CosVectorBucketResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorBucketModel(t, got)
	assertCosVectorBucketEncryption(t, got, cosVectorBucketTestSseType)
}

// TestCosVectorBucketCreateWithoutEncryption verifies that a null
// encryption_configuration is transmitted as "absent" (nil SDK pointer).
func TestCosVectorBucketCreateWithoutEncryption(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)

	var captured *cos.CreateVectorBucketOptions
	patches.ApplyMethodFunc(vectorClient.Vector, "CreateVectorBucket", func(ctx context.Context, opt *cos.CreateVectorBucketOptions) (*cos.CreateVectorBucketResult, *cos.Response, error) {
		captured = opt

		return &cos.CreateVectorBucketResult{
			VectorBucketQcs: cosVectorBucketTestQcs,
		}, cosVectorBucketTestResponse(), nil
	})
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucket", func(ctx context.Context, opt *cos.GetVectorBucketOptions) (*cos.GetVectorBucketResult, *cos.Response, error) {
		return &cos.GetVectorBucketResult{
			VectorBucket: cosVectorBucketTestItem(),
		}, cosVectorBucketTestResponse(), nil
	})

	model := cosVectorBucketTestModel()
	model.Id = types.StringNull()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateVectorBucket request is nil")
	}
	if captured.EncryptionConfiguration != nil {
		t.Errorf("create request encryption_configuration = %v, want nil", captured.EncryptionConfiguration)
	}

	var got svcos.CosVectorBucketResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorBucketModel(t, got)
	if !got.EncryptionConfiguration.IsNull() {
		t.Error("encryption_configuration should stay null when the cloud reports none")
	}
}

// TestCosVectorBucketCreateApiError verifies API failures surface as errors.
func TestCosVectorBucketCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "CreateVectorBucket", func(ctx context.Context, opt *cos.CreateVectorBucketOptions) (*cos.CreateVectorBucketResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("create failed")
	})

	model := cosVectorBucketTestModel()
	model.Id = types.StringNull()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestCosVectorBucketReadSuccess verifies Read refreshes the leaf fields.
func TestCosVectorBucketReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)

	var gotName string
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucket", func(ctx context.Context, opt *cos.GetVectorBucketOptions) (*cos.GetVectorBucketResult, *cos.Response, error) {
		gotName = opt.VectorBucketName

		return &cos.GetVectorBucketResult{
			VectorBucket: cosVectorBucketTestItem(),
		}, cosVectorBucketTestResponse(), nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorBucketTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorBucketTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotName != cosVectorBucketTestName {
		t.Errorf("GetVectorBucket called with %q, want %q", gotName, cosVectorBucketTestName)
	}

	var got svcos.CosVectorBucketResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorBucketModel(t, got)
}

// TestCosVectorBucketReadDetectsEncryptionDrift verifies a user-configurable
// nested object is refreshed from the cloud: the state holds a null
// encryption_configuration, the cloud now answers AES256, and Read must
// record the cloud value (not keep the previous one), so a change made
// outside Terraform shows up in the next plan.
func TestCosVectorBucketReadDetectsEncryptionDrift(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucket", func(ctx context.Context, opt *cos.GetVectorBucketOptions) (*cos.GetVectorBucketResult, *cos.Response, error) {
		drifted := cosVectorBucketTestItem()
		drifted.EncryptionConfiguration = &cos.VectorEncryptionConfig{
			SseType: cosVectorBucketTestSseType,
		}

		return &cos.GetVectorBucketResult{
			VectorBucket: drifted,
		}, cosVectorBucketTestResponse(), nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorBucketTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorBucketTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got svcos.CosVectorBucketResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorBucketEncryption(t, got, cosVectorBucketTestSseType)
}

// TestCosVectorBucketReadNotFound verifies that a documented HTTP 404 answer
// drops the resource from state (confirmed not-found).
func TestCosVectorBucketReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucket", func(ctx context.Context, opt *cos.GetVectorBucketOptions) (*cos.GetVectorBucketResult, *cos.Response, error) {
		return nil, nil, &cos.VectorErrorResponse{
			Response: &http.Response{StatusCode: http.StatusNotFound},
		}
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorBucketTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorBucketTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the vector bucket is confirmed gone")
	}
}

// TestCosVectorBucketReadApiError verifies non-404 API failures surface as
// errors instead of dropping state.
func TestCosVectorBucketReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucket", func(ctx context.Context, opt *cos.GetVectorBucketOptions) (*cos.GetVectorBucketResult, *cos.Response, error) {
		return nil, nil, &cos.VectorErrorResponse{
			Response: &http.Response{StatusCode: http.StatusForbidden},
		}
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorBucketTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorBucketTestModel())}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on a non-404 API error")
	}
}

// TestCosVectorBucketUpdateNoop verifies that Update persists the plan without
// calling any write API: the cloud has no update API for vector buckets and
// every attribute is RequiresReplace.
func TestCosVectorBucketUpdateNoop(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "CreateVectorBucket", func(ctx context.Context, opt *cos.CreateVectorBucketOptions) (*cos.CreateVectorBucketResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("the create API must not be called")
	})
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucket", func(ctx context.Context, opt *cos.GetVectorBucketOptions) (*cos.GetVectorBucketResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("the read API must not be called")
	})

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, cosVectorBucketTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, cosVectorBucketTestModel()),
			State: ctx.state(t, cosVectorBucketTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	var got svcos.CosVectorBucketResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorBucketModel(t, got)
}

// TestCosVectorBucketDeleteSuccess verifies Delete sends the delete request
// with the tracked vector bucket name.
func TestCosVectorBucketDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)

	var captured *cos.DeleteVectorBucketOptions
	patches.ApplyMethodFunc(vectorClient.Vector, "DeleteVectorBucket", func(ctx context.Context, opt *cos.DeleteVectorBucketOptions) (*cos.Response, error) {
		captured = opt

		return cosVectorBucketTestResponse(), nil
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, cosVectorBucketTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteVectorBucket request is nil")
	}
	if captured.VectorBucketName != cosVectorBucketTestName {
		t.Errorf("delete request vector_bucket_name = %q, want %q", captured.VectorBucketName, cosVectorBucketTestName)
	}
}

// TestCosVectorBucketDeleteApiError verifies delete failures surface as
// errors.
func TestCosVectorBucketDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketTestContext(t)
	vectorClient := patchCosVectorBucketSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "DeleteVectorBucket", func(ctx context.Context, opt *cos.DeleteVectorBucketOptions) (*cos.Response, error) {
		return nil, fmt.Errorf("delete failed")
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, cosVectorBucketTestModel())}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}
