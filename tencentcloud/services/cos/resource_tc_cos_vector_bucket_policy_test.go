package cos_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
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
	cosVectorBucketPolicyTestName = "tf-example-vector-1250000000"
	// cosVectorBucketPolicyTestPolicy is the practitioner-side serialization
	// (Statement first).
	cosVectorBucketPolicyTestPolicy = `{"Statement":[{"Action":["name/cos:GetVectorBucket"],"Effect":"Allow","Principal":{"qcs":["qcs::cam::uin/700000000000:uin/700001234567"]},"Resource":["qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000/*"],"Sid":"total"}],"Version":"2.0"}`
	// cosVectorBucketPolicyTestReorderedPolicy is the same policy as the
	// cloud would re-serialize it (Version first): semantically identical,
	// byte-different.
	cosVectorBucketPolicyTestReorderedPolicy = `{"Version":"2.0","Statement":[{"Sid":"total","Effect":"Allow","Principal":{"qcs":["qcs::cam::uin/700000000000:uin/700001234567"]},"Resource":["qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000/*"],"Action":["name/cos:GetVectorBucket"]}]}`
	cosVectorBucketPolicyTestRenamedPolicy   = `{"Statement":[{"Action":["name/cos:GetVectorBucket","name/cos:PutVectors"],"Effect":"Allow","Principal":{"qcs":["qcs::cam::uin/700000000000:uin/700001234567"]},"Resource":["qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000/*"],"Sid":"total"}],"Version":"2.0"}`
	// cosVectorBucketPolicyTestNormalizedPolicy is the cloud-side
	// re-serialization of the renamed policy.
	cosVectorBucketPolicyTestNormalizedPolicy = `{"Version":"2.0","Statement":[{"Sid":"total","Effect":"Allow","Principal":{"qcs":["qcs::cam::uin/700000000000:uin/700001234567"]},"Resource":["qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000/*"],"Action":["name/cos:GetVectorBucket","name/cos:PutVectors"]}]}`
	// cosVectorBucketPolicyTestDriftedPolicy differs in CONTENT (Deny
	// effect): a genuine drift the read-back must surface.
	cosVectorBucketPolicyTestDriftedPolicy = `{"Version":"2.0","Statement":[{"Sid":"total","Effect":"Deny","Principal":{"qcs":["qcs::cam::uin/700000000000:uin/700001234567"]},"Resource":["qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-example-vector-1250000000/*"],"Action":["name/cos:GetVectorBucket"]}]}`
)

// cosVectorBucketPolicyTestResponse builds a valid *cos.Response whose
// embedded *http.Response is populated (a bare &cos.Response{} dereferences
// nil in the service DEBUG log).
func cosVectorBucketPolicyTestResponse() *cos.Response {
	return &cos.Response{Response: &http.Response{StatusCode: http.StatusOK}}
}

// cosVectorBucketPolicyTestModel builds the terraform model used for
// plan/state. Every unset attribute carries an EXPLICIT TYPED null value —
// the raw zero value breaks tfsdk Plan/State reflection with "Value
// Conversion Error".
func cosVectorBucketPolicyTestModel() svcos.CosVectorBucketPolicyResourceModel {
	return svcos.CosVectorBucketPolicyResourceModel{
		Id:               types.StringValue(cosVectorBucketPolicyTestName),
		VectorBucketName: types.StringValue(cosVectorBucketPolicyTestName),
		Policy:           types.StringValue(cosVectorBucketPolicyTestPolicy),
	}
}

// assertCosVectorBucketPolicyModel asserts every model field.
func assertCosVectorBucketPolicyModel(t *testing.T, m svcos.CosVectorBucketPolicyResourceModel) {
	t.Helper()

	if m.Id.ValueString() != cosVectorBucketPolicyTestName {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), cosVectorBucketPolicyTestName)
	}
	if m.VectorBucketName.ValueString() != cosVectorBucketPolicyTestName {
		t.Errorf("vector_bucket_name = %q, want %q", m.VectorBucketName.ValueString(), cosVectorBucketPolicyTestName)
	}
}

// patchCosVectorBucketPolicySdkClient stubs the vector client accessor and
// returns the mock SDK client whose methods the test patches next (the
// two-level mode).
func patchCosVectorBucketPolicySdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *cos.Client {
	vectorClient := cos.NewClient(&cos.BaseURL{}, nil)
	patches.ApplyMethodReturn(client, "UseTencentCosVectorClient", vectorClient)

	return vectorClient
}

// cosVectorBucketPolicyTestContext wires the resource to a configured client;
// the SDK client methods are stubbed per test case.
type cosVectorBucketPolicyTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newCosVectorBucketPolicyTestContext(t *testing.T) *cosVectorBucketPolicyTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svcos.NewCosVectorBucketPolicyResource()
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

	return &cosVectorBucketPolicyTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *cosVectorBucketPolicyTestContext) plan(t *testing.T, m svcos.CosVectorBucketPolicyResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *cosVectorBucketPolicyTestContext) state(t *testing.T, m svcos.CosVectorBucketPolicyResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// TestCosVectorBucketPolicyCreateSuccess verifies that Create sends all
// request fields, tracks the vector bucket name as the id and refreshes the
// state.
func TestCosVectorBucketPolicyCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)

	var captured *cos.PutVectorBucketPolicyOptions
	patches.ApplyMethodFunc(vectorClient.Vector, "PutVectorBucketPolicy", func(ctx context.Context, opt *cos.PutVectorBucketPolicyOptions) (*cos.Response, error) {
		captured = opt

		return cosVectorBucketPolicyTestResponse(), nil
	})
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucketPolicy", func(ctx context.Context, opt *cos.GetVectorBucketPolicyOptions) (*cos.GetVectorBucketPolicyResult, *cos.Response, error) {
		// The cloud answers with its own re-serialization (Version first):
		// semantically identical, byte-different from the plan.
		return &cos.GetVectorBucketPolicyResult{
			Policy: cosVectorBucketPolicyTestReorderedPolicy,
		}, cosVectorBucketPolicyTestResponse(), nil
	})

	model := cosVectorBucketPolicyTestModel()
	// The id is Computed: it is unknown in a real create plan.
	model.Id = types.StringNull()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("PutVectorBucketPolicy request is nil")
	}
	if captured.VectorBucketName != cosVectorBucketPolicyTestName {
		t.Errorf("create request vector_bucket_name = %q, want %q", captured.VectorBucketName, cosVectorBucketPolicyTestName)
	}
	if captured.Policy != cosVectorBucketPolicyTestPolicy {
		t.Errorf("create request policy = %q, want %q", captured.Policy, cosVectorBucketPolicyTestPolicy)
	}

	var got svcos.CosVectorBucketPolicyResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorBucketPolicyModel(t, got)
	// The planned serialization is kept: overwriting it with the cloud's
	// re-serialization fails Terraform's "inconsistent result after apply".
	if got.Policy.ValueString() != cosVectorBucketPolicyTestPolicy {
		t.Errorf("policy = %q, want the planned serialization %q", got.Policy.ValueString(), cosVectorBucketPolicyTestPolicy)
	}
}

// TestCosVectorBucketPolicyCreateApiError verifies API failures surface as
// errors.
func TestCosVectorBucketPolicyCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "PutVectorBucketPolicy", func(ctx context.Context, opt *cos.PutVectorBucketPolicyOptions) (*cos.Response, error) {
		return nil, fmt.Errorf("put failed")
	})

	model := cosVectorBucketPolicyTestModel()
	model.Id = types.StringNull()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestCosVectorBucketPolicyReadSuccess verifies Read refreshes the leaf
// fields and addresses the bucket by the tracked id.
func TestCosVectorBucketPolicyReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)

	var gotName string
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucketPolicy", func(ctx context.Context, opt *cos.GetVectorBucketPolicyOptions) (*cos.GetVectorBucketPolicyResult, *cos.Response, error) {
		gotName = opt.VectorBucketName

		// The cloud re-serializes the policy; the refresh must keep the
		// state's own serialization instead of producing a permanent diff.
		return &cos.GetVectorBucketPolicyResult{
			Policy: cosVectorBucketPolicyTestReorderedPolicy,
		}, cosVectorBucketPolicyTestResponse(), nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorBucketPolicyTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorBucketPolicyTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotName != cosVectorBucketPolicyTestName {
		t.Errorf("GetVectorBucketPolicy called with %q, want %q", gotName, cosVectorBucketPolicyTestName)
	}

	var got svcos.CosVectorBucketPolicyResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorBucketPolicyModel(t, got)
	if got.Policy.ValueString() != cosVectorBucketPolicyTestPolicy {
		t.Errorf("policy = %q, want the state serialization %q kept (no reordering diff)", got.Policy.ValueString(), cosVectorBucketPolicyTestPolicy)
	}
}

// TestCosVectorBucketPolicyReadDetectsPolicyDrift verifies that a
// CONTENT-different cloud policy (not a mere re-serialization) is recorded:
// the state holds the Allow policy, the cloud now answers Deny, and Read must
// refresh the cloud value so the change shows up in the next plan.
func TestCosVectorBucketPolicyReadDetectsPolicyDrift(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucketPolicy", func(ctx context.Context, opt *cos.GetVectorBucketPolicyOptions) (*cos.GetVectorBucketPolicyResult, *cos.Response, error) {
		return &cos.GetVectorBucketPolicyResult{
			Policy: cosVectorBucketPolicyTestDriftedPolicy,
		}, cosVectorBucketPolicyTestResponse(), nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorBucketPolicyTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorBucketPolicyTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got svcos.CosVectorBucketPolicyResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Policy.ValueString() != cosVectorBucketPolicyTestDriftedPolicy {
		t.Errorf("policy = %q, want the cloud value %q", got.Policy.ValueString(), cosVectorBucketPolicyTestDriftedPolicy)
	}
}

// TestCosVectorBucketPolicyReadNotFound verifies that a documented HTTP 404
// answer drops the resource from state (confirmed not-found).
func TestCosVectorBucketPolicyReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucketPolicy", func(ctx context.Context, opt *cos.GetVectorBucketPolicyOptions) (*cos.GetVectorBucketPolicyResult, *cos.Response, error) {
		return nil, nil, &cos.VectorErrorResponse{
			Response: &http.Response{StatusCode: http.StatusNotFound},
		}
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorBucketPolicyTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorBucketPolicyTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the policy is confirmed gone")
	}
}

// TestCosVectorBucketPolicyReadEmptyPolicy verifies that an empty policy
// answer (no policy configured) drops the resource from state.
func TestCosVectorBucketPolicyReadEmptyPolicy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucketPolicy", func(ctx context.Context, opt *cos.GetVectorBucketPolicyOptions) (*cos.GetVectorBucketPolicyResult, *cos.Response, error) {
		return &cos.GetVectorBucketPolicyResult{}, cosVectorBucketPolicyTestResponse(), nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorBucketPolicyTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorBucketPolicyTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on an empty policy answer: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed when no policy is configured")
	}
}

// TestCosVectorBucketPolicyReadApiError verifies non-404 API failures surface
// as errors instead of dropping state.
func TestCosVectorBucketPolicyReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucketPolicy", func(ctx context.Context, opt *cos.GetVectorBucketPolicyOptions) (*cos.GetVectorBucketPolicyResult, *cos.Response, error) {
		return nil, nil, &cos.VectorErrorResponse{
			Response: &http.Response{StatusCode: http.StatusForbidden},
		}
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cosVectorBucketPolicyTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cosVectorBucketPolicyTestModel())}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on a non-404 API error")
	}
}

// TestCosVectorBucketPolicyUpdateSuccess verifies Update sends the PLANNED
// policy to the put request, and that the read-back keeps the planned
// serialization when the cloud answers with a semantically identical but
// re-serialized policy (an overwrite would fail Terraform's "inconsistent
// result after apply").
func TestCosVectorBucketPolicyUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)

	var captured *cos.PutVectorBucketPolicyOptions
	patches.ApplyMethodFunc(vectorClient.Vector, "PutVectorBucketPolicy", func(ctx context.Context, opt *cos.PutVectorBucketPolicyOptions) (*cos.Response, error) {
		captured = opt

		return cosVectorBucketPolicyTestResponse(), nil
	})
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucketPolicy", func(ctx context.Context, opt *cos.GetVectorBucketPolicyOptions) (*cos.GetVectorBucketPolicyResult, *cos.Response, error) {
		// The cloud answers with its re-serialization of the renamed policy.
		return &cos.GetVectorBucketPolicyResult{
			Policy: cosVectorBucketPolicyTestNormalizedPolicy,
		}, cosVectorBucketPolicyTestResponse(), nil
	})

	plannedModel := cosVectorBucketPolicyTestModel()
	plannedModel.Policy = types.StringValue(cosVectorBucketPolicyTestRenamedPolicy)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, cosVectorBucketPolicyTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, cosVectorBucketPolicyTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("PutVectorBucketPolicy request is nil")
	}
	if captured.VectorBucketName != cosVectorBucketPolicyTestName {
		t.Errorf("update request vector_bucket_name = %q, want %q", captured.VectorBucketName, cosVectorBucketPolicyTestName)
	}
	if captured.Policy != cosVectorBucketPolicyTestRenamedPolicy {
		t.Errorf("update request policy = %q, want the planned %q", captured.Policy, cosVectorBucketPolicyTestRenamedPolicy)
	}

	var got svcos.CosVectorBucketPolicyResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != cosVectorBucketPolicyTestName {
		t.Errorf("id = %q, want the unchanged %q", got.Id.ValueString(), cosVectorBucketPolicyTestName)
	}
	if got.Policy.ValueString() != cosVectorBucketPolicyTestRenamedPolicy {
		t.Errorf("policy = %q, want the planned serialization %q (the re-serialized cloud answer is semantically identical)", got.Policy.ValueString(), cosVectorBucketPolicyTestRenamedPolicy)
	}
}

// TestCosVectorBucketPolicyUpdateNoChange verifies that a plan equal to the
// state calls NO write API: no attribute the put API covers changed, so the
// call would be a pointless no-change request. The read-back still runs and
// the state is refreshed.
func TestCosVectorBucketPolicyUpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)

	putCalls := 0
	patches.ApplyMethodFunc(vectorClient.Vector, "PutVectorBucketPolicy", func(ctx context.Context, opt *cos.PutVectorBucketPolicyOptions) (*cos.Response, error) {
		putCalls++

		return nil, fmt.Errorf("the put API must not be called")
	})
	patches.ApplyMethodFunc(vectorClient.Vector, "GetVectorBucketPolicy", func(ctx context.Context, opt *cos.GetVectorBucketPolicyOptions) (*cos.GetVectorBucketPolicyResult, *cos.Response, error) {
		return &cos.GetVectorBucketPolicyResult{
			Policy: cosVectorBucketPolicyTestPolicy,
		}, cosVectorBucketPolicyTestResponse(), nil
	})

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, cosVectorBucketPolicyTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, cosVectorBucketPolicyTestModel()),
			State: ctx.state(t, cosVectorBucketPolicyTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if putCalls != 0 {
		t.Errorf("PutVectorBucketPolicy was called %d time(s) for a no-change plan", putCalls)
	}

	var got svcos.CosVectorBucketPolicyResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCosVectorBucketPolicyModel(t, got)
}

// TestCosVectorBucketPolicyUpdateApiError verifies update failures surface
// as errors.
func TestCosVectorBucketPolicyUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "PutVectorBucketPolicy", func(ctx context.Context, opt *cos.PutVectorBucketPolicyOptions) (*cos.Response, error) {
		return nil, fmt.Errorf("put failed")
	})

	plannedModel := cosVectorBucketPolicyTestModel()
	plannedModel.Policy = types.StringValue(cosVectorBucketPolicyTestRenamedPolicy)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, cosVectorBucketPolicyTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, cosVectorBucketPolicyTestModel()),
		}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should fail on an API error")
	}
}

// TestCosVectorBucketPolicyDeleteSuccess verifies Delete sends the delete
// request with the tracked vector bucket name.
func TestCosVectorBucketPolicyDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)

	var captured *cos.DeleteVectorBucketPolicyOptions
	patches.ApplyMethodFunc(vectorClient.Vector, "DeleteVectorBucketPolicy", func(ctx context.Context, opt *cos.DeleteVectorBucketPolicyOptions) (*cos.Response, error) {
		captured = opt

		return cosVectorBucketPolicyTestResponse(), nil
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, cosVectorBucketPolicyTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteVectorBucketPolicy request is nil")
	}
	if captured.VectorBucketName != cosVectorBucketPolicyTestName {
		t.Errorf("delete request vector_bucket_name = %q, want %q", captured.VectorBucketName, cosVectorBucketPolicyTestName)
	}
}

// TestCosVectorBucketPolicyDeleteApiError verifies delete failures surface
// as errors.
func TestCosVectorBucketPolicyDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCosVectorBucketPolicyTestContext(t)
	vectorClient := patchCosVectorBucketPolicySdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(vectorClient.Vector, "DeleteVectorBucketPolicy", func(ctx context.Context, opt *cos.DeleteVectorBucketPolicyOptions) (*cos.Response, error) {
		return nil, fmt.Errorf("delete failed")
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, cosVectorBucketPolicyTestModel())}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}
