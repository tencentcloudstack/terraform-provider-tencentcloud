// Data source tencentcloud_tdmysql_specs: queries the sellable component
// specs of TDSQL Boundless (tdmysql) via DescribeSpecs. All three response
// categories (hybrid node specs, serverless CCU specs, serverless node number
// spec) are exposed as computed outputs named after the response keys.
package tdmysql

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &TdmysqlSpecsDataSource{}

// NewTdmysqlSpecsDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTdmysqlSpecsDataSource() datasource.DataSource {
	return &TdmysqlSpecsDataSource{}
}

// TdmysqlSpecsDataSource implements datasource.DataSource for
// tencentcloud_tdmysql_specs.
type TdmysqlSpecsDataSource struct {
	fw.DataSourceWithConfigure
}

// TdmysqlSpecsDataSourceModel maps the schema attributes.
type TdmysqlSpecsDataSourceModel struct {
	Id                    types.String `tfsdk:"id"`
	FullReplications      types.Int64  `tfsdk:"full_replications"`
	IsExclusiveInstance   types.Int64  `tfsdk:"is_exclusive_instance"`
	Zones                 types.List   `tfsdk:"zones"`
	InstanceId            types.String `tfsdk:"instance_id"`
	HybridNodeSpecs       types.List   `tfsdk:"hybrid_node_specs"`
	ServerlessCcuSpec     types.List   `tfsdk:"serverless_ccu_spec"`
	ServerlessNodeNumSpec types.List   `tfsdk:"serverless_node_num_spec"`
}

// TdmysqlSpecsHybridNodeSpecItemModel is the element type of hybrid_node_specs.
type TdmysqlSpecsHybridNodeSpecItemModel struct {
	SpecCode               types.String `tfsdk:"spec_code"`
	StorageNodeCpu         types.Int64  `tfsdk:"storage_node_cpu"`
	StorageNodeMem         types.Int64  `tfsdk:"storage_node_mem"`
	StorageNodeMaxNum      types.Int64  `tfsdk:"storage_node_max_num"`
	StorageNodeMaxDisk     types.Int64  `tfsdk:"storage_node_max_disk"`
	StorageNodeMinNum      types.Int64  `tfsdk:"storage_node_min_num"`
	StorageNodeMinDisk     types.Int64  `tfsdk:"storage_node_min_disk"`
	StorageType            types.String `tfsdk:"storage_type"`
	StorageNodeDefaultDisk types.Int64  `tfsdk:"storage_node_default_disk"`
	InstanceMode           types.List   `tfsdk:"instance_mode"`
	DiskTypeCategory       types.String `tfsdk:"disk_type_category"`
}

// TdmysqlSpecsHybridNodeSpecItemAttrTypes is the object type of the
// hybrid_node_specs elements.
var TdmysqlSpecsHybridNodeSpecItemAttrTypes = map[string]attr.Type{
	"spec_code":                 types.StringType,
	"storage_node_cpu":          types.Int64Type,
	"storage_node_mem":          types.Int64Type,
	"storage_node_max_num":      types.Int64Type,
	"storage_node_max_disk":     types.Int64Type,
	"storage_node_min_num":      types.Int64Type,
	"storage_node_min_disk":     types.Int64Type,
	"storage_type":              types.StringType,
	"storage_node_default_disk": types.Int64Type,
	"instance_mode":             types.ListType{ElemType: types.StringType},
	"disk_type_category":        types.StringType,
}

// TdmysqlSpecsServerlessCcuItemModel is the element type of serverless_ccu_spec.
type TdmysqlSpecsServerlessCcuItemModel struct {
	MinCcu types.Float64 `tfsdk:"min_ccu"`
	MaxCcu types.List    `tfsdk:"max_ccu"`
}

// TdmysqlSpecsServerlessCcuItemAttrTypes is the object type of the
// serverless_ccu_spec elements.
var TdmysqlSpecsServerlessCcuItemAttrTypes = map[string]attr.Type{
	"min_ccu": types.Float64Type,
	"max_ccu": types.ListType{ElemType: types.Float64Type},
}

// TdmysqlSpecsServerlessNodeNumItemModel is the element type of
// serverless_node_num_spec.
type TdmysqlSpecsServerlessNodeNumItemModel struct {
	MinNodeNum types.Int64 `tfsdk:"min_node_num"`
	MaxNodeNum types.Int64 `tfsdk:"max_node_num"`
}

// TdmysqlSpecsServerlessNodeNumItemAttrTypes is the object type of the
// serverless_node_num_spec elements.
var TdmysqlSpecsServerlessNodeNumItemAttrTypes = map[string]attr.Type{
	"min_node_num": types.Int64Type,
	"max_node_num": types.Int64Type,
}

func (d *TdmysqlSpecsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_tdmysql_specs"
}

func (d *TdmysqlSpecsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the sellable component specs of TDSQL Boundless (tdmysql).",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random token
			// generated at read time (see the end of Read).
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query inputs ----
			"full_replications": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of replicas of the all-round (hybrid) node, currently supported range: 1-3, default is 3.",
			},
			"is_exclusive_instance": schema.Int64Attribute{
				Optional:    true,
				Description: "Exclusive instance.",
			},
			"zones": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of availability zone combinations.",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "Instance ID.",
			},
			// ---- Query outputs: Computed ----
			"hybrid_node_specs": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of hybrid node sellable specs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"spec_code": schema.StringAttribute{
							Computed:    true,
							Description: "Spec code.",
						},
						"storage_node_cpu": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of CPU cores of the storage node.",
						},
						"storage_node_mem": schema.Int64Attribute{
							Computed:    true,
							Description: "Memory size of the storage node.",
						},
						"storage_node_max_num": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum number of storage nodes.",
						},
						"storage_node_max_disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Upper limit of the storage node disk size.",
						},
						"storage_node_min_num": schema.Int64Attribute{
							Computed:    true,
							Description: "Minimum number of storage nodes.",
						},
						"storage_node_min_disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Lower limit of the storage node disk size.",
						},
						"storage_type": schema.StringAttribute{
							Computed:    true,
							Description: "Disk type, CLOUD_HSSD enhanced SSD, CLOUD_TCS local SSD disk.",
						},
						"storage_node_default_disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Default disk size of the storage node, used for frontend display.",
						},
						"instance_mode": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "List of billing modes supported by the spec.",
						},
						"disk_type_category": schema.StringAttribute{
							Computed:    true,
							Description: "Disk type category: CLOUD_DISK for cloud disk, LOCAL_DISK for local disk.",
						},
					},
				},
			},
			"serverless_ccu_spec": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of serverless node sellable CCU specs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"min_ccu": schema.Float64Attribute{
							Computed:    true,
							Description: "Minimum CCU value.",
						},
						"max_ccu": schema.ListAttribute{
							Computed:    true,
							ElementType: types.Float64Type,
							Description: "Range of maximum CCU values.",
						},
					},
				},
			},
			"serverless_node_num_spec": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Serverless node number configuration.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"min_node_num": schema.Int64Attribute{
							Computed:    true,
							Description: "Minimum number of nodes.",
						},
						"max_node_num": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum number of nodes.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *TdmysqlSpecsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TdmysqlSpecsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeSpecsRequest := tdmysqlv20211122.NewDescribeSpecsRequest()
	if !config.FullReplications.IsNull() && !config.FullReplications.IsUnknown() {
		describeSpecsRequest.FullReplications = fwHelper.Int64PointerFromValue(config.FullReplications)
	}
	if !config.IsExclusiveInstance.IsNull() && !config.IsExclusiveInstance.IsUnknown() {
		describeSpecsRequest.IsExclusiveInstance = fwHelper.Int64PointerFromValue(config.IsExclusiveInstance)
	}
	// An optional list filter that the SDK takes as []*string — null/unknown
	// maps to nil, i.e. no filter.
	describeSpecsRequest.Zones = fwHelper.PointersFromList[string](ctx, config.Zones, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !config.InstanceId.IsNull() && !config.InstanceId.IsUnknown() {
		describeSpecsRequest.InstanceId = config.InstanceId.ValueStringPointer()
	}

	service := NewTdmysqlService(d.Client())

	// The API has no pagination fields, so the service method makes a single
	// retried call and returns the three response categories as-is.
	hybridNodeSpecs, serverlessCcuSpec, serverlessNodeNumSpec, err := service.DescribeSpecs(ctx, describeSpecsRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TDSQL Boundless specs",
			fmt.Sprintf("DescribeSpecs failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	hybridItems := make([]TdmysqlSpecsHybridNodeSpecItemModel, 0, len(hybridNodeSpecs))
	for _, item := range hybridNodeSpecs {
		if item == nil {
			continue
		}

		hybridItems = append(hybridItems, TdmysqlSpecsHybridNodeSpecItemModel{
			SpecCode:               fwHelper.StringPointerValueOrNull(item.SpecCode),
			StorageNodeCpu:         fwHelper.Int64ValueOrNull(item.StorageNodeCpu),
			StorageNodeMem:         fwHelper.Int64ValueOrNull(item.StorageNodeMem),
			StorageNodeMaxNum:      fwHelper.Int64ValueOrNull(item.StorageNodeMaxNum),
			StorageNodeMaxDisk:     fwHelper.Int64ValueOrNull(item.StorageNodeMaxDisk),
			StorageNodeMinNum:      fwHelper.Int64ValueOrNull(item.StorageNodeMinNum),
			StorageNodeMinDisk:     fwHelper.Int64ValueOrNull(item.StorageNodeMinDisk),
			StorageType:            fwHelper.StringPointerValueOrNull(item.StorageType),
			StorageNodeDefaultDisk: fwHelper.Int64ValueOrNull(item.StorageNodeDefaultDisk),
			InstanceMode:           fwHelper.ListFromPointers(ctx, types.StringType, item.InstanceMode, &resp.Diagnostics),
			DiskTypeCategory:       fwHelper.StringPointerValueOrNull(item.DiskTypeCategory),
		})
	}
	if resp.Diagnostics.HasError() {
		return
	}

	hybridList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlSpecsHybridNodeSpecItemAttrTypes}, hybridItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverlessCcuItems := make([]TdmysqlSpecsServerlessCcuItemModel, 0, len(serverlessCcuSpec))
	for _, item := range serverlessCcuSpec {
		if item == nil {
			continue
		}

		serverlessCcuItems = append(serverlessCcuItems, TdmysqlSpecsServerlessCcuItemModel{
			MinCcu: types.Float64PointerValue(item.MinCcu),
			MaxCcu: fwHelper.ListFromPointers(ctx, types.Float64Type, item.MaxCcu, &resp.Diagnostics),
		})
	}
	if resp.Diagnostics.HasError() {
		return
	}

	serverlessCcuList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlSpecsServerlessCcuItemAttrTypes}, serverlessCcuItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ServerlessNodeNumSpec is a single object in the API response — it is
	// exposed as a zero-or-one element list.
	serverlessNodeNumItems := make([]TdmysqlSpecsServerlessNodeNumItemModel, 0, 1)
	if serverlessNodeNumSpec != nil {
		serverlessNodeNumItems = append(serverlessNodeNumItems, TdmysqlSpecsServerlessNodeNumItemModel{
			MinNodeNum: fwHelper.Int64ValueOrNull(serverlessNodeNumSpec.MinNodeNum),
			MaxNodeNum: fwHelper.Int64ValueOrNull(serverlessNodeNumSpec.MaxNodeNum),
		})
	}

	serverlessNodeNumList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlSpecsServerlessNodeNumItemAttrTypes}, serverlessNodeNumItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.HybridNodeSpecs = hybridList
	config.ServerlessCcuSpec = serverlessCcuList
	config.ServerlessNodeNumSpec = serverlessNodeNumList

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
