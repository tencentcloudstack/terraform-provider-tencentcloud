package cosvectors

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

var cosVectorMetadataTypes = map[string]attr.Type{"non_filterable_metadata_keys": types.SetType{ElemType: types.StringType}}

var (
	_ resource.ResourceWithConfigure      = &CosVectorIndexResource{}
	_ resource.ResourceWithImportState    = &CosVectorIndexResource{}
	_ resource.ResourceWithValidateConfig = &CosVectorIndexResource{}
)

type CosVectorIndexResource struct {
	fw.ResourceWithConfigure
	service *CosVectorService
}

type CosVectorIndexResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	VectorBucketName      types.String `tfsdk:"vector_bucket_name"`
	IndexName             types.String `tfsdk:"index_name"`
	IndexQcs              types.String `tfsdk:"index_qcs"`
	Region                types.String `tfsdk:"region"`
	CreationTime          types.Int64  `tfsdk:"creation_time"`
	DataType              types.String `tfsdk:"data_type"`
	Dimension             types.Int64  `tfsdk:"dimension"`
	DistanceMetric        types.String `tfsdk:"distance_metric"`
	MetadataConfiguration types.List   `tfsdk:"metadata_configuration"`
}

func NewCosVectorIndexResource() resource.Resource { return &CosVectorIndexResource{} }
func (r *CosVectorIndexResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cos_vector_index"
}
func (r *CosVectorIndexResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a COS Vectors index. Index configuration changes replace it; deleting an index irreversibly deletes all vectors stored in it. Vector data operations are not managed by Terraform.",
		Attributes: map[string]schema.Attribute{
			"id":                 schema.StringAttribute{Computed: true, Description: "Full bucket-name/index-name identity.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"vector_bucket_name": schema.StringAttribute{Required: true, Description: "Full vector bucket name including APPID. Reference the bucket resource to establish creation/deletion ordering.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"index_name":         schema.StringAttribute{Required: true, Description: "Index name. Changing it replaces the index.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"index_qcs":          schema.StringAttribute{Computed: true, Description: "Actual index QCS returned by COS.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"region":             schema.StringAttribute{Computed: true, Description: "Creation/import region. A different provider region is rejected.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"creation_time":      schema.Int64Attribute{Computed: true, Description: "Raw creationTime returned by COS, null when no nonzero value is reported.", PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"data_type":          schema.StringAttribute{Required: true, Description: "Vector data type. Only float32 is supported.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"dimension":          schema.Int64Attribute{Required: true, Description: "Vector dimension, between 1 and 4096 per the COS SDK contract.", PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}},
			"distance_metric":    schema.StringAttribute{Required: true, Description: "cosine or euclidean.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"metadata_configuration": schema.ListAttribute{
				Optional: true, Computed: true, Description: "Single immutable metadata block with 1-10 distinct non_filterable_metadata_keys. Keys match ^[a-zA-Z_][a-zA-Z0-9_]{0,62}$; order is ignored. Omit to preserve imported configuration.",
				ElementType:   types.ObjectType{AttrTypes: cosVectorMetadataTypes},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplaceIfConfigured(), listplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *CosVectorIndexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CosVectorIndexResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data CosVectorIndexResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.VectorBucketName.IsNull() && !data.VectorBucketName.IsUnknown() {
		if err := validateCosVectorBucketName(data.VectorBucketName.ValueString()); err != nil {
			resp.Diagnostics.AddError("Invalid COS vector bucket name", err.Error())
		}
	}
	if !data.IndexName.IsNull() && !data.IndexName.IsUnknown() {
		if err := validateCosVectorIndexName(data.IndexName.ValueString()); err != nil {
			resp.Diagnostics.AddError("Invalid COS vector index name", err.Error())
		}
	}
	if !data.DataType.IsNull() && !data.DataType.IsUnknown() && data.DataType.ValueString() != "float32" {
		resp.Diagnostics.AddError("Invalid COS vector data type", "data_type must be float32.")
	}
	if !data.Dimension.IsNull() && !data.Dimension.IsUnknown() && (data.Dimension.ValueInt64() < 1 || data.Dimension.ValueInt64() > 4096) {
		resp.Diagnostics.AddError("Invalid COS vector dimension", "dimension must be between 1 and 4096.")
	}
	if !data.DistanceMetric.IsNull() && !data.DistanceMetric.IsUnknown() && data.DistanceMetric.ValueString() != "cosine" && data.DistanceMetric.ValueString() != "euclidean" {
		resp.Diagnostics.AddError("Invalid COS vector distance metric", "distance_metric must be cosine or euclidean.")
	}
	if data.MetadataConfiguration.IsNull() || data.MetadataConfiguration.IsUnknown() {
		return
	}
	configuration, err := cosVectorConfigurationObject(data.MetadataConfiguration)
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS vector metadata", err.Error())
		return
	}
	if configuration.IsUnknown() {
		return
	}
	keys, ok := configuration.Attributes()["non_filterable_metadata_keys"].(types.Set)
	if ok && keys.IsUnknown() {
		return
	}
	if !ok || keys.IsNull() {
		resp.Diagnostics.AddError("Invalid COS vector metadata", "non_filterable_metadata_keys is required.")
		return
	}
	known := 0
	for _, value := range keys.Elements() {
		key, ok := value.(types.String)
		if !ok || key.IsNull() || (!key.IsUnknown() && !cosVectorMetadataKeyPattern.MatchString(key.ValueString())) {
			resp.Diagnostics.AddError("Invalid COS vector metadata", "Metadata keys must match ^[a-zA-Z_][a-zA-Z0-9_]{0,62}$.")
		}
		if !value.IsUnknown() {
			known++
		}
	}
	if len(keys.Elements()) < 1 || known > 10 {
		resp.Diagnostics.AddError("Invalid COS vector metadata", "non_filterable_metadata_keys must contain 1-10 keys.")
	}
}

func (data *CosVectorIndexResourceModel) options() (*cosSDK.CreateIndexOptions, error) {
	opt := &cosSDK.CreateIndexOptions{VectorBucketName: data.VectorBucketName.ValueString(), IndexName: data.IndexName.ValueString(), DataType: data.DataType.ValueString(), Dimension: int(data.Dimension.ValueInt64()), DistanceMetric: data.DistanceMetric.ValueString()}
	if !data.MetadataConfiguration.IsNull() && !data.MetadataConfiguration.IsUnknown() {
		configuration, err := cosVectorConfigurationObject(data.MetadataConfiguration)
		if err != nil {
			return nil, err
		}
		if configuration.IsUnknown() {
			return nil, errors.New("metadata configuration must be known before creating an index")
		}
		keys, ok := configuration.Attributes()["non_filterable_metadata_keys"].(types.Set)
		if !ok || keys.IsNull() || keys.IsUnknown() {
			return nil, errors.New("metadata keys must be known before creating an index")
		}
		metadata := &cosSDK.MetadataConfiguration{}
		for _, value := range keys.Elements() {
			key, ok := value.(types.String)
			if !ok || key.IsNull() || key.IsUnknown() {
				return nil, errors.New("metadata keys must be known non-empty strings")
			}
			metadata.NonFilterableMetadataKeys = append(metadata.NonFilterableMetadataKeys, key.ValueString())
		}
		opt.MetadataConfiguration = metadata
	}
	return opt, validateCosVectorIndexOptions(opt)
}

func (r *CosVectorIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CosVectorIndexResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.service == nil {
		resp.Diagnostics.AddError("COS Vectors client unavailable", "Configure the provider before creating an index.")
		return
	}
	opt, err := data.options()
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS vector index", err.Error())
		return
	}
	if data.MetadataConfiguration.IsUnknown() {
		data.MetadataConfiguration = types.ListNull(types.ObjectType{AttrTypes: cosVectorMetadataTypes})
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	qcs, err := r.service.CreateIndex(ctx, opt)
	if err != nil && !isCosVectorBucketCreateAccepted(err) {
		resp.Diagnostics.AddError("Error creating COS vector index", err.Error())
		return
	}
	data.Id = types.StringValue(opt.VectorBucketName + "/" + opt.IndexName)
	data.Region = types.StringValue(r.service.region)
	data.IndexQcs = types.StringNull()
	if qcs != "" {
		data.IndexQcs = types.StringValue(qcs)
	}
	data.CreationTime = types.Int64Null()
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	index, err := r.service.WaitIndex(ctx, opt.VectorBucketName, opt.IndexName, true)
	if err != nil {
		resp.Diagnostics.AddError("COS vector index created but refresh failed", "Identity remains tracked. Inspect state and taint/replacement plans before recovery. "+err.Error())
		return
	}
	actual := data
	actual.setIndex(index)
	if actual.DataType != data.DataType || actual.Dimension != data.Dimension || actual.DistanceMetric != data.DistanceMetric || !actual.MetadataConfiguration.Equal(data.MetadataConfiguration) {
		resp.Diagnostics.AddError("COS vector index configuration not confirmed", "GetIndex did not return the requested immutable configuration; identity remains tracked.")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &actual)...)
}

func (r *CosVectorIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CosVectorIndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.checkIdentity(data); err != nil {
		resp.Diagnostics.AddError("Invalid COS vector index identity", err.Error())
		return
	}
	index, err := r.service.GetIndex(ctx, data.VectorBucketName.ValueString(), data.IndexName.ValueString())
	if isCosVectorBucketNotFound(err) {
		resp.Diagnostics.AddWarning("COS vector index no longer exists", "The index is being removed from state; a subsequent plan may recreate it.")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading COS vector index", err.Error())
		return
	}
	data.setIndex(index)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CosVectorIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CosVectorIndexResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CosVectorIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CosVectorIndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.checkIdentity(data); err != nil {
		resp.Diagnostics.AddError("Invalid COS vector index identity", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	err := r.service.DeleteIndex(ctx, data.VectorBucketName.ValueString(), data.IndexName.ValueString())
	if isCosVectorBucketNotFound(err) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error deleting COS vector index", err.Error())
		return
	}
	if _, err := r.service.WaitIndex(ctx, data.VectorBucketName.ValueString(), data.IndexName.ValueString(), false); err != nil {
		resp.Diagnostics.AddError("COS vector index deletion not confirmed", err.Error())
	}
}

func (r *CosVectorIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.service == nil {
		resp.Diagnostics.AddError("COS Vectors client unavailable", "Configure the provider before importing an index.")
		return
	}
	bucket, index, err := parseCosVectorIndexImportID(req.ID, r.service.region)
	if err != nil {
		resp.Diagnostics.AddError("Invalid COS vector index import ID", err.Error())
		return
	}
	data := CosVectorIndexResourceModel{Id: types.StringValue(bucket + "/" + index), VectorBucketName: types.StringValue(bucket), IndexName: types.StringValue(index), Region: types.StringValue(r.service.region), IndexQcs: types.StringNull(), DataType: types.StringNull(), Dimension: types.Int64Null(), DistanceMetric: types.StringNull(), CreationTime: types.Int64Null(), MetadataConfiguration: types.ListNull(types.ObjectType{AttrTypes: cosVectorMetadataTypes})}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CosVectorIndexResource) checkIdentity(data CosVectorIndexResourceModel) error {
	if r.service == nil {
		return errors.New("COS Vectors client is not configured")
	}
	if data.Region.IsNull() || data.Region.ValueString() != r.service.region {
		return errors.New("index state region differs from provider region")
	}
	bucket, index, err := parseCosVectorIndexImportID(data.Id.ValueString(), r.service.region)
	if err != nil || bucket != data.VectorBucketName.ValueString() || index != data.IndexName.ValueString() {
		return errors.New("index state ID must match bucket and index names")
	}
	return nil
}

func (data *CosVectorIndexResourceModel) setIndex(index *cosSDK.IndexInfo) {
	data.Id = types.StringValue(index.VectorBucketName + "/" + index.IndexName)
	data.VectorBucketName = types.StringValue(index.VectorBucketName)
	data.IndexName = types.StringValue(index.IndexName)
	data.IndexQcs = types.StringValue(index.IndexQcs)
	data.DataType = types.StringValue(index.DataType)
	data.Dimension = types.Int64Value(int64(index.Dimension))
	data.DistanceMetric = types.StringValue(index.DistanceMetric)
	data.CreationTime = types.Int64Null()
	if index.CreationTime != 0 {
		data.CreationTime = types.Int64Value(index.CreationTime)
	}
	data.MetadataConfiguration = types.ListNull(types.ObjectType{AttrTypes: cosVectorMetadataTypes})
	if index.MetadataConfiguration != nil && len(index.MetadataConfiguration.NonFilterableMetadataKeys) > 0 {
		keys := make([]attr.Value, 0, len(index.MetadataConfiguration.NonFilterableMetadataKeys))
		for _, key := range index.MetadataConfiguration.NonFilterableMetadataKeys {
			keys = append(keys, types.StringValue(key))
		}
		object := types.ObjectValueMust(cosVectorMetadataTypes, map[string]attr.Value{"non_filterable_metadata_keys": types.SetValueMust(types.StringType, keys)})
		data.MetadataConfiguration = types.ListValueMust(types.ObjectType{AttrTypes: cosVectorMetadataTypes}, []attr.Value{object})
	}
}
