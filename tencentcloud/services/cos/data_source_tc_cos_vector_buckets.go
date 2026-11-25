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

var _ datasource.DataSourceWithConfigure = &CosVectorBucketsDataSource{}

// NewCosVectorBucketsDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewCosVectorBucketsDataSource() datasource.DataSource {
	return &CosVectorBucketsDataSource{}
}

// CosVectorBucketsDataSource implements datasource.DataSource for
// tencentcloud_cos_vector_buckets.
type CosVectorBucketsDataSource struct {
	fw.DataSourceWithConfigure
}

// CosVectorBucketsDataSourceModel maps the schema attributes.
type CosVectorBucketsDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Prefix        types.String `tfsdk:"prefix"`
	VectorBuckets types.List   `tfsdk:"vector_buckets"`
}

// CosVectorBucketItemModel is the element type of vector_buckets.
type CosVectorBucketItemModel struct {
	CreationTime     types.Int64  `tfsdk:"creation_time"`
	VectorBucketQcs  types.String `tfsdk:"vector_bucket_qcs"`
	VectorBucketName types.String `tfsdk:"vector_bucket_name"`
}

// CosVectorBucketItemAttrTypes is the object type of the vector_buckets
// elements.
var CosVectorBucketItemAttrTypes = map[string]attr.Type{
	"creation_time":      types.Int64Type,
	"vector_bucket_qcs":  types.StringType,
	"vector_bucket_name": types.StringType,
}

func (d *CosVectorBucketsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cos_vector_buckets"
}

func (d *CosVectorBucketsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of COS vector buckets.",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random token
			// generated at read time.
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query inputs ----
			"prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Limits the response to vector buckets whose names begin with the specified prefix, 1-63 characters in length.",
			},
			// ---- Query outputs: Computed ----
			// Outputs carry the resource data only — never a total_count
			// aggregate, never paging knobs.
			"vector_buckets": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of vector buckets.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"creation_time": schema.Int64Attribute{
							Computed:    true,
							Description: "Creation date and time of the vector bucket, a Unix timestamp.",
						},
						"vector_bucket_qcs": schema.StringAttribute{
							Computed:    true,
							Description: "QCS (resource name) of the vector bucket, e.g. `qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000`.",
						},
						"vector_bucket_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the vector bucket, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *CosVectorBucketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CosVectorBucketsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listVectorBucketsRequest := &cos.ListVectorBucketsOptions{
		Prefix: config.Prefix.ValueString(),
	}

	service := NewCosService(d.Client())

	// When the API pages, the service method internally loops (NextToken)
	// and returns the full list — this file never loops by itself and never
	// surfaces the total count.
	buckets, err := service.ListVectorBuckets(ctx, listVectorBucketsRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying COS vector buckets",
			fmt.Sprintf("ListVectorBuckets failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	items := make([]CosVectorBucketItemModel, 0, len(buckets))
	for i := range buckets {
		items = append(items, CosVectorBucketItemModel{
			CreationTime:     types.Int64Value(buckets[i].CreationTime),
			VectorBucketQcs:  types.StringValue(buckets[i].VectorBucketQcs),
			VectorBucketName: types.StringValue(buckets[i].VectorBucketName),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CosVectorBucketItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.VectorBuckets = list

	// The state id is a random token.
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
