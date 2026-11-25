package cos

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions. Only assert the extension interfaces;
// the base resource.Resource interface is checked implicitly by the factory
// return value.
var (
	_ resource.ResourceWithConfigure   = &CosVectorBucketResource{}
	_ resource.ResourceWithImportState = &CosVectorBucketResource{}
)

// NewCosVectorBucketResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewCosVectorBucketResource() resource.Resource {
	return &CosVectorBucketResource{}
}

// CosVectorBucketResource implements resource.Resource for
// tencentcloud_cos_vector_bucket.
type CosVectorBucketResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// CosVectorBucketResourceModel maps the schema attributes; each attribute
// appears exactly once. Keep the struct fields in the same order as the
// schema.
type CosVectorBucketResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	VectorBucketName        types.String `tfsdk:"vector_bucket_name"`
	EncryptionConfiguration types.Object `tfsdk:"encryption_configuration"`
}

// cosVectorEncryptionModel maps the nested encryption_configuration object.
type cosVectorEncryptionModel struct {
	SseType types.String `tfsdk:"sse_type"`
}

// cosVectorEncryptionAttrTypes mirrors the attribute types of the
// encryption_configuration object declared by the resource schema.
var cosVectorEncryptionAttrTypes = map[string]attr.Type{
	"sse_type": types.StringType,
}

// cosVectorEncryptionAttributes declares the nested attributes of
// encryption_configuration.
var cosVectorEncryptionAttributes = map[string]schema.Attribute{
	"sse_type": schema.StringAttribute{
		Required:    true,
		Description: "Server-side encryption type of the vector bucket. Only `AES256` (SSE-COS) is supported.",
	},
}

func (r *CosVectorBucketResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cos_vector_bucket"
}

func (r *CosVectorBucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a COS vector bucket.",
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
				Description: "Name of the vector bucket, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`. It supports lowercase letters, numbers and `-`, with a length of 3-63 characters. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"encryption_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Attributes:  cosVectorEncryptionAttributes,
				Description: "Server-side encryption configuration of the vector bucket. Only SSE-COS is supported. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *CosVectorBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CosVectorBucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createVectorBucketRequest := &cos.CreateVectorBucketOptions{
		VectorBucketName: plan.VectorBucketName.ValueString(),
	}

	encryptionConfiguration, encryptionDiags := expandCosVectorEncryptionConfiguration(ctx, plan.EncryptionConfiguration)
	resp.Diagnostics.Append(encryptionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createVectorBucketRequest.EncryptionConfiguration = encryptionConfiguration

	service := NewCosService(r.Client())

	if _, err := service.CreateVectorBucket(ctx, createVectorBucketRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error creating COS vector bucket",
			fmt.Sprintf("CreateVectorBucket failed: %v", err),
		)

		return
	}

	// The vector bucket name is the addressable identity of the resource
	// (GetVectorBucket/DeleteVectorBucket address it by name), so it becomes
	// the state id directly.
	plan.Id = types.StringValue(plan.VectorBucketName.ValueString())

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read-back — same request form as Read; GetVectorBucket is an exact
	// lookup, so the returned item belongs to this resource.
	getVectorBucketRequest := &cos.GetVectorBucketOptions{
		VectorBucketName: plan.Id.ValueString(),
	}
	item, err := service.DescribeVectorBucket(ctx, getVectorBucketRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading COS vector bucket after create",
			fmt.Sprintf("GetVectorBucket failed for %s: %v", plan.Id.ValueString(), err),
		)

		return
	}

	resp.Diagnostics.Append(refreshCosVectorBucketModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *CosVectorBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CosVectorBucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vectorBucketName := state.Id.ValueString()

	getVectorBucketRequest := &cos.GetVectorBucketOptions{
		VectorBucketName: vectorBucketName,
	}

	service := NewCosService(r.Client())

	item, err := service.DescribeVectorBucket(ctx, getVectorBucketRequest)
	if err != nil {
		// A confirmed "not found" on the cloud side drops the resource from
		// state so Terraform plans a recreation. The judgement stays INLINE
		// here (RemoveResource followed by an immediate return): moving it
		// into the shared refresh helper would let a later unconditional
		// State.Set overwrite the removal.
		if errors.Is(err, fwHelper.ErrNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("COS vector bucket %s not found, removing from state.", vectorBucketName),
			)

			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading COS vector bucket",
			fmt.Sprintf("GetVectorBucket failed for %s: %v", vectorBucketName, err),
		)

		return
	}

	resp.Diagnostics.Append(refreshCosVectorBucketModel(ctx, &state, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. The cloud provides no update API for
// vector buckets: every attribute is RequiresReplace, so Terraform never
// reaches this method with a real change — the plan is persisted as-is.
func (r *CosVectorBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CosVectorBucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *CosVectorBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CosVectorBucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteVectorBucketRequest := &cos.DeleteVectorBucketOptions{
		VectorBucketName: state.Id.ValueString(),
	}

	service := NewCosService(r.Client())

	if err := service.DeleteVectorBucket(ctx, deleteVectorBucketRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting COS vector bucket",
			fmt.Sprintf("DeleteVectorBucket failed: %v", err),
		)

		return
	}
}

// expandCosVectorEncryptionConfiguration converts the plan's
// encryption_configuration object into the SDK request form. A null/unknown
// object is transmitted as "absent".
func expandCosVectorEncryptionConfiguration(ctx context.Context, object types.Object) (*cos.VectorEncryptionConfig, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model cosVectorEncryptionModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	return &cos.VectorEncryptionConfig{
		SseType: model.SseType.ValueString(),
	}, diags
}

// flattenCosVectorEncryptionConfiguration converts the SDK encryption
// configuration into the framework object form. A cloud-omitted configuration
// keeps the previous value (a configured value must not silently become
// null), a cloud-reported one is recorded as-is so drift shows up in the next
// plan.
func flattenCosVectorEncryptionConfiguration(ctx context.Context, item *cos.VectorEncryptionConfig, prev types.Object) (types.Object, diag.Diagnostics) {
	if item == nil {
		return prev, nil
	}

	return types.ObjectValue(cosVectorEncryptionAttrTypes, map[string]attr.Value{
		"sse_type": types.StringValue(item.SseType),
	})
}

// refreshCosVectorBucketModel copies cloud-owned fields into the model.
func refreshCosVectorBucketModel(ctx context.Context, m *CosVectorBucketResourceModel, item *cos.VectorBucketInfo) diag.Diagnostics {
	if item == nil {
		return nil
	}

	var diags diag.Diagnostics

	m.VectorBucketName = types.StringValue(item.VectorBucketName)

	encryptionConfiguration, encryptionDiags := flattenCosVectorEncryptionConfiguration(ctx, item.EncryptionConfiguration, m.EncryptionConfiguration)
	diags.Append(encryptionDiags...)
	if diags.HasError() {
		return diags
	}
	m.EncryptionConfiguration = encryptionConfiguration

	return diags
}
