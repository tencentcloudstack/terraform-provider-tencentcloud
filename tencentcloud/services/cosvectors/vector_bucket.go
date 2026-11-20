package cosvectors

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

var cosVectorEncryptionTypes = map[string]attr.Type{"sse_type": types.StringType}

var (
	_ resource.ResourceWithConfigure      = &CosVectorBucketResource{}
	_ resource.ResourceWithImportState    = &CosVectorBucketResource{}
	_ resource.ResourceWithValidateConfig = &CosVectorBucketResource{}
)

type CosVectorBucketResource struct {
	fw.ResourceWithConfigure
	service *CosVectorService
}

type CosVectorBucketResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	VectorBucketName        types.String `tfsdk:"vector_bucket_name"`
	VectorBucketQcs         types.String `tfsdk:"vector_bucket_qcs"`
	Region                  types.String `tfsdk:"region"`
	CreationTime            types.Int64  `tfsdk:"creation_time"`
	EncryptionConfiguration types.List   `tfsdk:"encryption_configuration"`
	ForceDestroy            types.Bool   `tfsdk:"force_destroy"`
}

func NewCosVectorBucketResource() resource.Resource { return &CosVectorBucketResource{} }

func (r *CosVectorBucketResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cos_vector_bucket"
}

func (r *CosVectorBucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a COS Vectors bucket, separate from a normal COS bucket or Terraform state backend. " +
			"Name and encryption changes replace the bucket. Indexes and policies use separate resources. " +
			"Uses the provider credentials and regional HTTPS Vectors endpoint. Tags and KMS are not supported.",
		Attributes: map[string]schema.Attribute{
			"force_destroy":      schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(false), Description: "Defaults to false. If applied as true, bucket deletion first removes all indexes, including those outside Terraform and all their vectors. This is irreversible. Apply changes to this flag before destroy."},
			"id":                 schema.StringAttribute{Computed: true, Description: "Full vector bucket name including APPID.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"vector_bucket_name": schema.StringAttribute{Required: true, Description: "Full bucket name, 3-63 lowercase letters, digits or hyphens, ending in -APPID.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"vector_bucket_qcs":  schema.StringAttribute{Computed: true, Description: "Bucket QCS returned by COS.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"region":             schema.StringAttribute{Computed: true, Description: "Creation/import region. Changing the provider region is rejected to protect existing resources.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"creation_time":      schema.Int64Attribute{Computed: true, Description: "Raw creationTime integer returned by the API. Null when the SDK reports no nonzero timestamp.", PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"encryption_configuration": schema.ListAttribute{
				Optional: true, Computed: true,
				Description:   "Single SSE-COS encryption block with sse_type = AES256. New buckets default to AES256; importing preserves actual encryption.",
				ElementType:   types.ObjectType{AttrTypes: cosVectorEncryptionTypes},
				PlanModifiers: []planmodifier.List{cosVectorEncryptionDefault{}, listplanmodifier.RequiresReplaceIfConfigured(), listplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *CosVectorBucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CosVectorBucketResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data CosVectorBucketResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.VectorBucketName.IsNull() && !data.VectorBucketName.IsUnknown() {
		if err := validateCosVectorBucketName(data.VectorBucketName.ValueString()); err != nil {
			resp.Diagnostics.AddError("Invalid COS vector bucket name", err.Error())
		}
	}
	if data.EncryptionConfiguration.IsNull() || data.EncryptionConfiguration.IsUnknown() {
		return
	}
	configuration, err := cosVectorConfigurationObject(data.EncryptionConfiguration)
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS vector encryption", err.Error())
		return
	}
	if configuration.IsUnknown() {
		return
	}
	sse, ok := configuration.Attributes()["sse_type"].(types.String)
	if !ok || sse.IsNull() || (!sse.IsUnknown() && sse.ValueString() != "AES256") {
		resp.Diagnostics.AddError("Unsupported COS vector encryption", "encryption_configuration.sse_type must be AES256. KMS is not supported.")
	}
}

func (r *CosVectorBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CosVectorBucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.service == nil {
		resp.Diagnostics.AddError("COS Vectors client unavailable", "Configure the TencentCloud provider before creating a bucket.")
		return
	}
	if data.EncryptionConfiguration.IsNull() || data.EncryptionConfiguration.IsUnknown() {
		data.EncryptionConfiguration = cosVectorEncryptionValue("AES256")
	}
	configuration, err := cosVectorConfigurationObject(data.EncryptionConfiguration)
	if err != nil || configuration.IsUnknown() || configuration.IsNull() {
		resp.Diagnostics.AddError("Invalid COS vector encryption", "A single known encryption block is required.")
		return
	}
	sse, ok := configuration.Attributes()["sse_type"].(types.String)
	if !ok || sse.IsNull() || sse.IsUnknown() || sse.ValueString() != "AES256" {
		resp.Diagnostics.AddError("Unsupported COS vector encryption", "encryption_configuration.sse_type must be the known value AES256.")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	qcs, err := r.service.CreateBucket(ctx, data.VectorBucketName.ValueString(), &cosSDK.VectorEncryptionConfig{SseType: sse.ValueString()})
	if err != nil && !isCosVectorBucketCreateAccepted(err) {
		resp.Diagnostics.AddError("Error creating COS vector bucket", err.Error())
		return
	}
	data.Id = data.VectorBucketName
	data.Region = types.StringValue(r.service.region)
	data.VectorBucketQcs = types.StringNull()
	if qcs != "" {
		data.VectorBucketQcs = types.StringValue(qcs)
	}
	data.CreationTime = types.Int64Null()
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bucket, err := r.service.WaitBucket(ctx, data.Id.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("COS vector bucket created but refresh failed", "Bucket identity remains tracked in state. Inspect the state and taint/replacement plan before retrying. "+err.Error())
		return
	}
	if bucket.EncryptionConfiguration == nil || bucket.EncryptionConfiguration.SseType != sse.ValueString() {
		resp.Diagnostics.AddError("COS vector bucket encryption not confirmed", "GetVectorBucket did not return the requested AES256 encryption. Bucket identity remains tracked; verify encryption before retrying.")
		return
	}
	data.setBucket(bucket)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CosVectorBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CosVectorBucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.checkIdentity(data); err != nil {
		resp.Diagnostics.AddError("Invalid COS vector bucket identity", err.Error())
		return
	}
	bucket, err := r.service.GetBucket(ctx, data.Id.ValueString())
	if isCosVectorBucketNotFound(err) {
		resp.Diagnostics.AddWarning("COS vector bucket no longer exists", "The bucket was removed outside Terraform; it is being removed from state and a subsequent plan may recreate it.")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading COS vector bucket", err.Error())
		return
	}
	data.setBucket(bucket)
	if data.ForceDestroy.IsNull() || data.ForceDestroy.IsUnknown() {
		data.ForceDestroy = types.BoolValue(false)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CosVectorBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CosVectorBucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CosVectorBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CosVectorBucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.checkIdentity(data); err != nil {
		resp.Diagnostics.AddError("Invalid COS vector bucket identity", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	if data.ForceDestroy.ValueBool() {
		if err := r.service.EmptyBucket(ctx, data.Id.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error force deleting COS vector indexes", "Cleanup stopped; some indexes may already have been deleted. Bucket remains tracked. "+err.Error())
			return
		}
	}
	err := r.service.DeleteBucketAfterCleanup(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting COS vector bucket", "Without force_destroy the bucket is not emptied; with force_destroy some indexes may already have been deleted. "+err.Error())
		return
	}
	if _, err := r.service.WaitBucket(ctx, data.Id.ValueString(), false); err != nil {
		resp.Diagnostics.AddError("COS vector bucket deletion not confirmed", err.Error())
	}
}

func (r *CosVectorBucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.service == nil {
		resp.Diagnostics.AddError("COS Vectors client unavailable", "Configure the TencentCloud provider before importing a bucket.")
		return
	}
	name, err := parseCosVectorBucketImportID(req.ID, r.service.region)
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS vector bucket import ID", err.Error())
		return
	}
	data := CosVectorBucketResourceModel{Id: types.StringValue(name), VectorBucketName: types.StringValue(name), Region: types.StringValue(r.service.region),
		VectorBucketQcs: types.StringNull(), CreationTime: types.Int64Null(), EncryptionConfiguration: types.ListNull(types.ObjectType{AttrTypes: cosVectorEncryptionTypes}), ForceDestroy: types.BoolValue(false)}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CosVectorBucketResource) checkIdentity(data CosVectorBucketResourceModel) error {
	if r.service == nil {
		return errors.New("COS Vectors client is not configured")
	}
	if data.Region.IsNull() || data.Region.IsUnknown() || data.Region.ValueString() != r.service.region {
		return errors.New("state region differs from provider region; restore the original region or use a separate provider alias and import")
	}
	if data.Id.IsNull() || data.Id.IsUnknown() || data.Id.ValueString() != data.VectorBucketName.ValueString() {
		return errors.New("state ID must match the full vector bucket name")
	}
	return validateCosVectorBucketName(data.Id.ValueString())
}

func (data *CosVectorBucketResourceModel) setBucket(bucket *cosSDK.VectorBucketInfo) {
	data.Id = types.StringValue(bucket.VectorBucketName)
	data.VectorBucketName = data.Id
	data.VectorBucketQcs = types.StringValue(bucket.VectorBucketQcs)
	data.CreationTime = types.Int64Null()
	if bucket.CreationTime != 0 {
		data.CreationTime = types.Int64Value(bucket.CreationTime)
	}
	data.EncryptionConfiguration = types.ListNull(types.ObjectType{AttrTypes: cosVectorEncryptionTypes})
	if bucket.EncryptionConfiguration != nil {
		data.EncryptionConfiguration = cosVectorEncryptionValue(bucket.EncryptionConfiguration.SseType)
	}
}

// cosVectorEncryptionDefault 只默认新建且省略的配置，不覆盖导入或未知表达式。
type cosVectorEncryptionDefault struct{}

func (cosVectorEncryptionDefault) Description(context.Context) string {
	return "New buckets default to AES256; existing bucket encryption is preserved."
}
func (m cosVectorEncryptionDefault) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}
func (cosVectorEncryptionDefault) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.Plan.Raw.IsNull() || !req.State.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		resp.PlanValue = cosVectorEncryptionValue("AES256")
	}
}

func cosVectorEncryptionValue(sse string) types.List {
	object := types.ObjectValueMust(cosVectorEncryptionTypes, map[string]attr.Value{"sse_type": types.StringValue(sse)})
	return types.ListValueMust(types.ObjectType{AttrTypes: cosVectorEncryptionTypes}, []attr.Value{object})
}
