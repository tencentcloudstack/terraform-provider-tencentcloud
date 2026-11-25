package cos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions. Only assert the extension interfaces;
// the base resource.Resource interface is checked implicitly by the factory
// return value.
var (
	_ resource.ResourceWithConfigure   = &CosVectorBucketPolicyResource{}
	_ resource.ResourceWithImportState = &CosVectorBucketPolicyResource{}
)

// NewCosVectorBucketPolicyResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewCosVectorBucketPolicyResource() resource.Resource {
	return &CosVectorBucketPolicyResource{}
}

// CosVectorBucketPolicyResource implements resource.Resource for
// tencentcloud_cos_vector_bucket_policy.
type CosVectorBucketPolicyResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// CosVectorBucketPolicyResourceModel maps the schema attributes; each
// attribute appears exactly once. Keep the struct fields in the same order as
// the schema.
type CosVectorBucketPolicyResourceModel struct {
	Id               types.String `tfsdk:"id"`
	VectorBucketName types.String `tfsdk:"vector_bucket_name"`
	Policy           types.String `tfsdk:"policy"`
}

func (r *CosVectorBucketPolicyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cos_vector_bucket_policy"
}

func (r *CosVectorBucketPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a COS vector bucket policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource. It is the same as `vector_bucket_name`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vector_bucket_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the vector bucket the policy attaches to, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy": schema.StringAttribute{
				Required:    true,
				Description: "Policy of the vector bucket. It is a JSON-formatted access control policy serialized as a string. The new policy fully overwrites the old one.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *CosVectorBucketPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	putVectorBucketPolicyRequest := &cos.PutVectorBucketPolicyOptions{
		VectorBucketName: plan.VectorBucketName.ValueString(),
		Policy:           plan.Policy.ValueString(),
	}

	service := NewCosService(r.Client())

	if err := service.PutVectorBucketPolicy(ctx, putVectorBucketPolicyRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error creating COS vector bucket policy",
			fmt.Sprintf("PutVectorBucketPolicy failed: %v", err),
		)

		return
	}

	// The vector bucket name is the addressable identity of the resource
	// (Get/DeleteVectorBucketPolicy address it by name), so it becomes the
	// state id directly.
	plan.Id = types.StringValue(plan.VectorBucketName.ValueString())

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read-back — same request form as Read; GetVectorBucketPolicy is an
	// exact lookup, so the returned policy belongs to this resource.
	getVectorBucketPolicyRequest := &cos.GetVectorBucketPolicyOptions{
		VectorBucketName: plan.Id.ValueString(),
	}
	policy, err := service.DescribeVectorBucketPolicy(ctx, getVectorBucketPolicyRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading COS vector bucket policy after create",
			fmt.Sprintf("GetVectorBucketPolicy failed for %s: %v", plan.Id.ValueString(), err),
		)

		return
	}

	plan.Policy = cosVectorPolicyValueOrSemanticallyEqual(&policy, plan.Policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *CosVectorBucketPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vectorBucketName := state.Id.ValueString()

	getVectorBucketPolicyRequest := &cos.GetVectorBucketPolicyOptions{
		VectorBucketName: vectorBucketName,
	}

	service := NewCosService(r.Client())

	policy, err := service.DescribeVectorBucketPolicy(ctx, getVectorBucketPolicyRequest)
	if err != nil {
		// A confirmed "not found" on the cloud side (no policy configured or
		// the bucket itself is gone) drops the resource from state so
		// Terraform plans a recreation. The judgement stays INLINE here
		// (RemoveResource followed by an immediate return): moving it into a
		// shared refresh helper would let a later unconditional State.Set
		// overwrite the removal.
		if errors.Is(err, fwHelper.ErrNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("COS vector bucket policy of %s not found, removing from state.", vectorBucketName),
			)

			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading COS vector bucket policy",
			fmt.Sprintf("GetVectorBucketPolicy failed for %s: %v", vectorBucketName, err),
		)

		return
	}

	state.VectorBucketName = types.StringValue(vectorBucketName)
	state.Policy = cosVectorPolicyValueOrSemanticallyEqual(&policy, state.Policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. PutVectorBucketPolicy fully
// overwrites the policy, so it covers the `policy` attribute; the bucket name
// is the identity and is RequiresReplace.
func (r *CosVectorBucketPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vectorBucketName := state.Id.ValueString()

	service := NewCosService(r.Client())

	// Call the write API ONLY when the attribute it covers changed (by
	// semantic comparison: the cloud re-serializes the policy JSON, so two
	// strings with different key order are NOT a change). A plan reaching
	// Update with no policy change is a pointless no-change request that
	// many cloud APIs reject. The read-back below runs either way, so the
	// state always reflects the cloud.
	if !equalCosVectorPolicy(plan.Policy.ValueString(), state.Policy.ValueString()) {
		putVectorBucketPolicyRequest := &cos.PutVectorBucketPolicyOptions{
			VectorBucketName: vectorBucketName,
			Policy:           plan.Policy.ValueString(),
		}

		if err := service.PutVectorBucketPolicy(ctx, putVectorBucketPolicyRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error updating COS vector bucket policy",
				fmt.Sprintf("PutVectorBucketPolicy failed for %s: %v", vectorBucketName, err),
			)

			return
		}
	}

	getVectorBucketPolicyRequest := &cos.GetVectorBucketPolicyOptions{
		VectorBucketName: vectorBucketName,
	}
	policy, err := service.DescribeVectorBucketPolicy(ctx, getVectorBucketPolicyRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading COS vector bucket policy after update",
			fmt.Sprintf("GetVectorBucketPolicy failed for %s: %v", vectorBucketName, err),
		)

		return
	}

	plan.Id = state.Id
	plan.VectorBucketName = types.StringValue(vectorBucketName)
	plan.Policy = cosVectorPolicyValueOrSemanticallyEqual(&policy, plan.Policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion. It only deletes the access control
// policy; the vector bucket itself is untouched.
func (r *CosVectorBucketPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteVectorBucketPolicyRequest := &cos.DeleteVectorBucketPolicyOptions{
		VectorBucketName: state.Id.ValueString(),
	}

	service := NewCosService(r.Client())

	if err := service.DeleteVectorBucketPolicy(ctx, deleteVectorBucketPolicyRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting COS vector bucket policy",
			fmt.Sprintf("DeleteVectorBucketPolicy failed: %v", err),
		)

		return
	}
}

// cosVectorPolicyValueOrSemanticallyEqual returns the cloud policy unless it
// is semantically identical to the previous (planned/state) value.
//
// The cloud re-serializes the policy JSON when reading it back (key order,
// spacing), so a naive read-back would write a byte-different string over the
// planned value and fail Terraform's post-apply consistency check with
// "inconsistent result after apply", and produce a permanent diff in every
// plan afterwards. A semantic comparison keeps the practitioner's own
// serialization when the policy content is unchanged, while a genuinely
// different cloud policy (drift) is still recorded so it shows up in the next
// plan.
func cosVectorPolicyValueOrSemanticallyEqual(cloud *string, prev types.String) types.String {
	if cloud == nil {
		// Cloud omitted the field: keep the previous value (anti-diff),
		// null when it was never set.
		if !prev.IsNull() && !prev.IsUnknown() {
			return prev
		}

		return types.StringNull()
	}

	if !prev.IsNull() && !prev.IsUnknown() && equalCosVectorPolicy(*cloud, prev.ValueString()) {
		return prev
	}

	return types.StringValue(*cloud)
}

// equalCosVectorPolicy reports whether two policy strings carry the same JSON
// content. It compares the canonicalized forms (compact, map keys sorted, as
// produced by encoding/json), falling back to a plain string comparison when
// either side is not valid JSON.
func equalCosVectorPolicy(a, b string) bool {
	normalizedA, aOk := normalizeCosVectorPolicyJson(a)
	normalizedB, bOk := normalizeCosVectorPolicyJson(b)
	if aOk && bOk {
		return normalizedA == normalizedB
	}

	return a == b
}

// normalizeCosVectorPolicyJson parses the policy JSON and re-serializes it in
// a canonical form. It reports false (and echoes the input) when the string
// is not valid JSON.
func normalizeCosVectorPolicyJson(policy string) (string, bool) {
	var value interface{}
	if err := json.Unmarshal([]byte(policy), &value); err != nil {
		return policy, false
	}

	canonical, err := json.Marshal(value)
	if err != nil {
		return policy, false
	}

	return string(canonical), true
}
