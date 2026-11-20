package cosvectors

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var (
	_ resource.ResourceWithConfigure      = &CosVectorBucketPolicyResource{}
	_ resource.ResourceWithImportState    = &CosVectorBucketPolicyResource{}
	_ resource.ResourceWithValidateConfig = &CosVectorBucketPolicyResource{}
)

type CosVectorBucketPolicyResource struct {
	fw.ResourceWithConfigure
	service *CosVectorService
}

type CosVectorBucketPolicyResourceModel struct {
	Id               types.String         `tfsdk:"id"`
	VectorBucketQcs  types.String         `tfsdk:"vector_bucket_qcs"`
	VectorBucketName types.String         `tfsdk:"vector_bucket_name"`
	Region           types.String         `tfsdk:"region"`
	Policy           cosVectorPolicyValue `tfsdk:"policy"`
}

func NewCosVectorBucketPolicyResource() resource.Resource { return &CosVectorBucketPolicyResource{} }
func (r *CosVectorBucketPolicyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cos_vector_bucket_policy"
}
func (r *CosVectorBucketPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the complete COS Vectors bucket policy. Create overwrites any existing policy, updates are in-place, and delete removes the policy without restoring previous contents. Do not manage the same bucket policy through multiple resources.",
		Attributes: map[string]schema.Attribute{
			"id":                 schema.StringAttribute{Computed: true, Description: "Bucket QCS identity.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"vector_bucket_qcs":  schema.StringAttribute{Required: true, Description: "Target bucket QCS. Reference the vector bucket resource to establish dependencies. Changing the bucket replaces the policy resource.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"vector_bucket_name": schema.StringAttribute{Computed: true, Description: "Full bucket name resolved from QCS.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"region":             schema.StringAttribute{Computed: true, Description: "Creation/import region. A different provider region is rejected.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"policy":             schema.StringAttribute{Required: true, CustomType: cosVectorPolicyType{}, PlanModifiers: []planmodifier.String{cosVectorPolicyPlanModifier{}}, Description: "COS policy JSON string, normally produced by jsonencode. Uses version 2.0, lowercase statement/effect/principal/action/resource keys and name/cos: actions. JSON formatting and COS policy set order are ignored. Policy content is stored in Terraform state; never include credentials."},
		},
	}
}

func (r *CosVectorBucketPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ResourceWithConfigure.Configure(ctx, req, resp)
	service, err := NewCosVectorService(r.Client())
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS Vectors configuration", err.Error())
		return
	}
	r.service = service
}

func (r *CosVectorBucketPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.Policy.IsNull() && !data.Policy.IsUnknown() {
		if err := validateCosVectorPolicy(data.Policy.ValueString()); err != nil {
			resp.Diagnostics.AddError("Invalid COS vector bucket policy", err.Error())
		}
	}
	if !data.VectorBucketQcs.IsNull() && !data.VectorBucketQcs.IsUnknown() {
		if !isCosVectorBucketQCS(data.VectorBucketQcs.ValueString()) {
			resp.Diagnostics.AddError("Invalid COS vector bucket QCS", "vector_bucket_qcs must be a COS Vectors bucket QCS matching the bucket APPID.")
		}
	}
}

func isCosVectorBucketQCS(qcs string) bool {
	parts := strings.Split(qcs, ":")
	if len(parts) != 6 || parts[0] != "qcs" || parts[1] != "" || parts[2] != "cosvector" {
		return false
	}
	_, err := parseCosVectorBucketImportID(qcs, parts[3])
	return err == nil
}

func (r *CosVectorBucketPolicyResource) target(qcs string) (string, error) {
	if r.service == nil {
		return "", errors.New("COS Vectors client is not configured")
	}
	if !isCosVectorBucketQCS(qcs) {
		return "", errors.New("policy requires a valid COS Vectors bucket QCS")
	}
	return parseCosVectorBucketImportID(qcs, r.service.region)
}

func (r *CosVectorBucketPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bucket, err := r.target(data.VectorBucketQcs.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS vector bucket policy target", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	err = r.service.PutBucketPolicy(ctx, bucket, data.Policy.ValueString())
	if err != nil && !isCosVectorBucketCreateAccepted(err) {
		resp.Diagnostics.AddError("Error creating COS vector bucket policy", err.Error())
		return
	}
	data.Id = data.VectorBucketQcs
	data.VectorBucketName = types.StringValue(bucket)
	data.Region = types.StringValue(r.service.region)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.service.WaitBucketPolicy(ctx, bucket, data.Policy.ValueString(), true); err != nil {
		resp.Diagnostics.AddError("COS vector bucket policy written but refresh failed", "Policy identity remains tracked. "+err.Error())
	}
}

func (r *CosVectorBucketPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bucket, err := r.checkIdentity(data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS vector bucket policy identity", err.Error())
		return
	}
	policy, err := r.service.GetBucketPolicy(ctx, bucket)
	if isCosVectorBucketNotFound(err) {
		resp.Diagnostics.AddWarning("COS vector bucket policy no longer exists", "The policy is being removed from state; a subsequent plan may recreate it.")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading COS vector bucket policy", err.Error())
		return
	}
	data.Policy = cosVectorPolicyString(policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CosVectorBucketPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var old, data CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bucket, err := r.checkIdentity(old)
	if err != nil || !old.VectorBucketQcs.Equal(data.VectorBucketQcs) {
		resp.Diagnostics.AddError("Invalid COS vector bucket policy update", "Policy target must remain unchanged and match the provider region.")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := r.service.PutBucketPolicy(ctx, bucket, data.Policy.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating COS vector bucket policy", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.service.WaitBucketPolicy(ctx, bucket, data.Policy.ValueString(), true); err != nil {
		resp.Diagnostics.AddError("COS vector bucket policy update not confirmed", err.Error())
	}
}

func (r *CosVectorBucketPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CosVectorBucketPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bucket, err := r.checkIdentity(data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS vector bucket policy identity", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	err = r.service.DeleteBucketPolicy(ctx, bucket)
	if isCosVectorBucketNotFound(err) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error deleting COS vector bucket policy", err.Error())
		return
	}
	if err := r.service.WaitBucketPolicy(ctx, bucket, "", false); err != nil {
		resp.Diagnostics.AddError("COS vector bucket policy deletion not confirmed", err.Error())
	}
}

func (r *CosVectorBucketPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.service == nil {
		resp.Diagnostics.AddError("COS Vectors client unavailable", "Configure the provider before importing a policy.")
		return
	}
	bucket, err := parseCosVectorBucketImportID(req.ID, r.service.region)
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS vector bucket policy import ID", err.Error())
		return
	}
	qcs := req.ID
	if !isCosVectorBucketQCS(qcs) {
		info, err := r.service.GetBucket(ctx, bucket)
		if err != nil {
			resp.Diagnostics.AddError("Error resolving COS vector bucket policy import", err.Error())
			return
		}
		qcs = info.VectorBucketQcs
	}
	data := CosVectorBucketPolicyResourceModel{Id: types.StringValue(qcs), VectorBucketQcs: types.StringValue(qcs), VectorBucketName: types.StringValue(bucket), Region: types.StringValue(r.service.region), Policy: cosVectorPolicyNull()}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CosVectorBucketPolicyResource) checkIdentity(data CosVectorBucketPolicyResourceModel) (string, error) {
	bucket, err := r.target(data.VectorBucketQcs.ValueString())
	if err != nil {
		return "", err
	}
	if !data.Id.Equal(data.VectorBucketQcs) || data.VectorBucketName.ValueString() != bucket || data.Region.ValueString() != r.service.region {
		return "", errors.New("policy state identity or region mismatch")
	}
	return bucket, nil
}
