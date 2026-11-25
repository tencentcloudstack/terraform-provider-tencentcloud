package cos

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &CosVectorIndexesDataSource{}

// NewCosVectorIndexesDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewCosVectorIndexesDataSource() datasource.DataSource {
	return &CosVectorIndexesDataSource{}
}

// CosVectorIndexesDataSource implements datasource.DataSource for
// tencentcloud_cos_vector_indexes.
type CosVectorIndexesDataSource struct {
	fw.DataSourceWithConfigure
}

// CosVectorIndexesDataSourceModel maps the schema attributes.
type CosVectorIndexesDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	VectorBucketName types.String `tfsdk:"vector_bucket_name"`
	Prefix           types.String `tfsdk:"prefix"`
	Indexes          types.List   `tfsdk:"indexes"`
}

// CosVectorIndexItemModel is the element type of indexes.
type CosVectorIndexItemModel struct {
	CreationTime     types.Int64  `tfsdk:"creation_time"`
	IndexQcs         types.String `tfsdk:"index_qcs"`
	IndexName        types.String `tfsdk:"index_name"`
	VectorBucketName types.String `tfsdk:"vector_bucket_name"`
}

// CosVectorIndexItemAttrTypes is the object type of the indexes elements.
var CosVectorIndexItemAttrTypes = map[string]attr.Type{
	"creation_time":      types.Int64Type,
	"index_qcs":          types.StringType,
	"index_name":         types.StringType,
	"vector_bucket_name": types.StringType,
}

func (d *CosVectorIndexesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cos_vector_indexes"
}

func (d *CosVectorIndexesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of indexes of a COS vector bucket.",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random token
			// generated at read time.
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query inputs ----
			"vector_bucket_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the vector bucket whose indexes are queried, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`.",
			},
			"prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Limits the response to indexes whose names begin with the specified prefix, 1-63 characters in length.",
			},
			// ---- Query outputs: Computed ----
			// Outputs carry the resource data only — never a total_count
			// aggregate, never paging knobs.
			"indexes": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of indexes.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"creation_time": schema.Int64Attribute{
							Computed:    true,
							Description: "Creation date and time of the index, a Unix timestamp.",
						},
						"index_qcs": schema.StringAttribute{
							Computed:    true,
							Description: "QCS (resource name) of the index, e.g. `qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000/index/exampleindex`.",
						},
						"index_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the index, 3-63 characters in length.",
						},
						"vector_bucket_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the vector bucket the index belongs to, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *CosVectorIndexesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CosVectorIndexesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listIndexesRequest := &cos.ListIndexesOptions{
		VectorBucketName: config.VectorBucketName.ValueString(),
		Prefix:           config.Prefix.ValueString(),
	}

	service := NewCosService(d.Client())

	// When the API pages, the service method internally loops (NextToken)
	// and returns the full list — this file never loops by itself and never
	// surfaces the total count.
	indexes, err := service.ListIndexes(ctx, listIndexesRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying COS vector indexes",
			fmt.Sprintf("ListIndexes failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	items := make([]CosVectorIndexItemModel, 0, len(indexes))
	for i := range indexes {
		items = append(items, CosVectorIndexItemModel{
			CreationTime:     types.Int64Value(indexes[i].CreationTime),
			IndexQcs:         types.StringValue(indexes[i].IndexQcs),
			IndexName:        types.StringValue(indexes[i].IndexName),
			VectorBucketName: types.StringValue(indexes[i].VectorBucketName),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CosVectorIndexItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Indexes = list

	// The state id is a random token.
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
