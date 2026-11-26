package cos

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.ResourceWithConfigure   = &CosVectorIndexResource{}
	_ resource.ResourceWithImportState = &CosVectorIndexResource{}
)

// NewCosVectorIndexResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewCosVectorIndexResource() resource.Resource {
	return &CosVectorIndexResource{}
}

// CosVectorIndexResource implements resource.Resource for
// tencentcloud_cos_vector_index.
type CosVectorIndexResource struct {
	fw.ResourceWithConfigure
}

// CosVectorIndexResourceModel maps the schema attributes; each attribute
// appears exactly once. Keep the struct fields in the same order as the
// schema.
type CosVectorIndexResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	VectorBucketName      types.String `tfsdk:"vector_bucket_name"`
	IndexName             types.String `tfsdk:"index_name"`
	DataType              types.String `tfsdk:"data_type"`
	Dimension             types.Int64  `tfsdk:"dimension"`
	DistanceMetric        types.String `tfsdk:"distance_metric"`
	MetadataConfiguration types.Object `tfsdk:"metadata_configuration"`
}

// cosVectorMetadataConfigurationModel maps the nested
// metadata_configuration object.
type cosVectorMetadataConfigurationModel struct {
	NonFilterableMetadataKeys types.List `tfsdk:"non_filterable_metadata_keys"`
}

// cosVectorMetadataConfigurationAttrTypes mirrors the attribute types of the
// metadata_configuration object declared by the resource schema.
var cosVectorMetadataConfigurationAttrTypes = map[string]attr.Type{
	"non_filterable_metadata_keys": types.ListType{ElemType: types.StringType},
}

// cosVectorMetadataConfigurationAttributes declares the nested attributes of
// metadata_configuration.
var cosVectorMetadataConfigurationAttributes = map[string]schema.Attribute{
	"non_filterable_metadata_keys": schema.ListAttribute{
		Required:    true,
		ElementType: types.StringType,
		Description: "Metadata keys that cannot be used for filtering. At least 1 and at most 10 keys, each 1-63 characters in length.",
	},
}

func (r *CosVectorIndexResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cos_vector_index"
}

func (r *CosVectorIndexResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a COS vector index.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, in the format of `vector_bucket_name#index_name`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vector_bucket_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the vector bucket the index belongs to, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"index_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the vector index, composed of lowercase letters, numbers, `-` and `.`, beginning and ending with a lowercase letter or number, 3-63 characters in length. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_type": schema.StringAttribute{
				Required:    true,
				Description: "Data type of the vectors. Default value: `float32`. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dimension": schema.Int64Attribute{
				Required:    true,
				Description: "Dimension of the vectors. Value range: 1-4096. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"distance_metric": schema.StringAttribute{
				Required:    true,
				Description: "Distance metric for similarity search. Valid values: `euclidean`, `cosine`. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Attributes:  cosVectorMetadataConfigurationAttributes,
				Description: "Metadata configuration of the index. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *CosVectorIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CosVectorIndexResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createIndexRequest := &cos.CreateIndexOptions{
		VectorBucketName: plan.VectorBucketName.ValueString(),
		IndexName:        plan.IndexName.ValueString(),
		DataType:         plan.DataType.ValueString(),
		Dimension:        int(plan.Dimension.ValueInt64()),
		DistanceMetric:   plan.DistanceMetric.ValueString(),
	}

	metadataConfiguration, metadataDiags := expandCosVectorMetadataConfiguration(ctx, plan.MetadataConfiguration)
	resp.Diagnostics.Append(metadataDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createIndexRequest.MetadataConfiguration = metadataConfiguration

	service := NewCosService(r.Client())

	if err := service.CreateIndex(ctx, createIndexRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error creating COS vector index",
			fmt.Sprintf("CreateIndex failed: %v", err),
		)

		return
	}

	// The vector bucket name and the index name address the index on the
	// cloud side, so the composite id is parent-first:
	// vector_bucket_name#index_name.
	plan.Id = types.StringValue(fwHelper.IdFormat(plan.VectorBucketName.ValueString(), plan.IndexName.ValueString()))

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read-back — same request form as Read; GetIndex is an exact lookup,
	// so the returned item belongs to this resource.
	getIndexRequest := &cos.GetIndexOptions{
		VectorBucketName: plan.VectorBucketName.ValueString(),
		IndexName:        plan.IndexName.ValueString(),
	}
	item, err := service.DescribeIndex(ctx, getIndexRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading COS vector index after create",
			fmt.Sprintf("GetIndex failed for %s: %v", plan.Id.ValueString(), err),
		)

		return
	}

	resp.Diagnostics.Append(refreshCosVectorIndexModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *CosVectorIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CosVectorIndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Split the composite id: both segments address the index on the cloud
	// side (parent-first: bucket name # index name).
	idParts, err := fwHelper.IdParseErr(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading COS vector index",
			fmt.Sprintf("invalid resource id %q, expected the format `vector_bucket_name#index_name`: %v", state.Id.ValueString(), err),
		)

		return
	}

	getIndexRequest := &cos.GetIndexOptions{
		VectorBucketName: idParts[0],
		IndexName:        idParts[1],
	}

	service := NewCosService(r.Client())

	item, err := service.DescribeIndex(ctx, getIndexRequest)
	if err != nil {
		// A confirmed "not found" on the cloud side drops the resource from
		// state so Terraform plans a recreation. The judgement stays INLINE
		// here (RemoveResource followed by an immediate return): moving it
		// into the shared refresh helper would let a later unconditional
		// State.Set overwrite the removal.
		if errors.Is(err, fwHelper.ErrNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("COS vector index %s not found, removing from state.", state.Id.ValueString()),
			)

			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading COS vector index",
			fmt.Sprintf("GetIndex failed for %s: %v", state.Id.ValueString(), err),
		)

		return
	}

	resp.Diagnostics.Append(refreshCosVectorIndexModel(ctx, &state, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. The cloud provides no update API
// for vector indexes: every attribute is RequiresReplace, so Terraform never
// reaches this method with a real change — the plan is persisted as-is.
func (r *CosVectorIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CosVectorIndexResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion. It only deletes the index; the vector
// bucket itself is untouched.
func (r *CosVectorIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CosVectorIndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Split the composite id; both segments address the index.
	idParts, err := fwHelper.IdParseErr(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting COS vector index",
			fmt.Sprintf("invalid resource id %q, expected the format `vector_bucket_name#index_name`: %v", state.Id.ValueString(), err),
		)

		return
	}

	deleteIndexRequest := &cos.DeleteIndexOptions{
		VectorBucketName: idParts[0],
		IndexName:        idParts[1],
	}

	service := NewCosService(r.Client())

	if err := service.DeleteIndex(ctx, deleteIndexRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting COS vector index",
			fmt.Sprintf("DeleteIndex failed: %v", err),
		)

		return
	}
}

// ImportState supports `terraform import` with the composite id
// vector_bucket_name#index_name.
func (r *CosVectorIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts, err := fwHelper.IdParseErr(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import id",
			fmt.Sprintf("the import id must be in the format `vector_bucket_name#index_name`: %v", err),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vector_bucket_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("index_name"), idParts[1])...)
}

// expandCosVectorMetadataConfiguration converts the plan's
// metadata_configuration object into the SDK request form. A null/unknown
// object is transmitted as "absent".
func expandCosVectorMetadataConfiguration(ctx context.Context, object types.Object) (*cos.MetadataConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model cosVectorMetadataConfigurationModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	var keys []string
	if !model.NonFilterableMetadataKeys.IsNull() && !model.NonFilterableMetadataKeys.IsUnknown() {
		keyDiags := model.NonFilterableMetadataKeys.ElementsAs(ctx, &keys, false)
		diags.Append(keyDiags...)
		if diags.HasError() {
			return nil, diags
		}
	}

	return &cos.MetadataConfiguration{
		NonFilterableMetadataKeys: keys,
	}, diags
}

// flattenCosVectorMetadataConfiguration converts the SDK metadata
// configuration into the framework object form. A cloud-omitted configuration
// keeps the previous value (a configured value must not silently become
// null), a cloud-reported one is recorded as-is so drift shows up in the next
// plan.
func flattenCosVectorMetadataConfiguration(ctx context.Context, item *cos.MetadataConfiguration, prev types.Object) (types.Object, diag.Diagnostics) {
	if item == nil {
		return prev, nil
	}

	keys, keyDiags := types.ListValueFrom(ctx, types.StringType, item.NonFilterableMetadataKeys)
	if keyDiags.HasError() {
		return prev, keyDiags
	}

	return types.ObjectValue(cosVectorMetadataConfigurationAttrTypes, map[string]attr.Value{
		"non_filterable_metadata_keys": keys,
	})
}

// refreshCosVectorIndexModel copies cloud-owned fields into the model.
func refreshCosVectorIndexModel(ctx context.Context, m *CosVectorIndexResourceModel, item *cos.IndexInfo) diag.Diagnostics {
	if item == nil {
		return nil
	}

	var diags diag.Diagnostics

	m.VectorBucketName = fwHelper.StringValueOrPrev(&item.VectorBucketName, m.VectorBucketName)
	m.IndexName = fwHelper.StringValueOrPrev(&item.IndexName, m.IndexName)
	m.DataType = fwHelper.StringValueOrPrev(&item.DataType, m.DataType)
	dimension := int64(item.Dimension)
	m.Dimension = fwHelper.Int64ValueOrPrev(&dimension, m.Dimension)
	m.DistanceMetric = fwHelper.StringValueOrPrev(&item.DistanceMetric, m.DistanceMetric)

	metadataConfiguration, metadataDiags := flattenCosVectorMetadataConfiguration(ctx, item.MetadataConfiguration, m.MetadataConfiguration)
	diags.Append(metadataDiags...)
	if diags.HasError() {
		return diags
	}
	m.MetadataConfiguration = metadataConfiguration

	return diags
}
