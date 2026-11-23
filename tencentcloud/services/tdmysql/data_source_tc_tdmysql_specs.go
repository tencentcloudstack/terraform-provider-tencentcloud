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
)

var _ datasource.DataSourceWithConfigure = &TdmysqlSpecsDataSource{}

// NewTdmysqlSpecsDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTdmysqlSpecsDataSource() datasource.DataSource {
	return &TdmysqlSpecsDataSource{}
}

// TdmysqlSpecsDataSource implements datasource.DataSource for
// the tencentcloud_tdmysql_specs data source.
type TdmysqlSpecsDataSource struct {
	fw.DataSourceWithConfigure
}

// TdmysqlSpecsDataSourceModel maps the schema attributes.
type TdmysqlSpecsDataSourceModel struct {
	FullReplications      types.Int64  `tfsdk:"full_replications"`
	IsExclusiveInstance   types.Int64  `tfsdk:"is_exclusive_instance"`
	Zones                 types.Set    `tfsdk:"zones"`
	InstanceId            types.String `tfsdk:"instance_id"`
	HybridNodeSpecs       types.List   `tfsdk:"hybrid_node_specs"`
	ServerlessCcuSpec     types.List   `tfsdk:"serverless_ccu_spec"`
	ServerlessNodeNumSpec types.List   `tfsdk:"serverless_node_num_spec"`
}

// TdmysqlSpecsHybridNodeSpecModel is the element type of hybrid_node_specs.
type TdmysqlSpecsHybridNodeSpecModel struct {
	SpecCode               types.String `tfsdk:"spec_code"`
	StorageNodeCpu         types.Int64  `tfsdk:"storage_node_cpu"`
	StorageNodeMem         types.Int64  `tfsdk:"storage_node_mem"`
	StorageNodeMaxNum      types.Int64  `tfsdk:"storage_node_max_num"`
	StorageNodeMaxDisk     types.Int64  `tfsdk:"storage_node_max_disk"`
	StorageNodeMinNum      types.Int64  `tfsdk:"storage_node_min_num"`
	StorageNodeMinDisk     types.Int64  `tfsdk:"storage_node_min_disk"`
	StorageType            types.String `tfsdk:"storage_type"`
	StorageNodeDefaultDisk types.Int64  `tfsdk:"storage_node_default_disk"`
	InstanceMode           types.Set    `tfsdk:"instance_mode"`
	DiskTypeCategory       types.String `tfsdk:"disk_type_category"`
}

// TdmysqlSpecsHybridNodeSpecAttrTypes is the object type of the
// hybrid_node_specs elements.
var TdmysqlSpecsHybridNodeSpecAttrTypes = map[string]attr.Type{
	"spec_code":                 types.StringType,
	"storage_node_cpu":          types.Int64Type,
	"storage_node_mem":          types.Int64Type,
	"storage_node_max_num":      types.Int64Type,
	"storage_node_max_disk":     types.Int64Type,
	"storage_node_min_num":      types.Int64Type,
	"storage_node_min_disk":     types.Int64Type,
	"storage_type":              types.StringType,
	"storage_node_default_disk": types.Int64Type,
	"instance_mode":             types.SetType{ElemType: types.StringType},
	"disk_type_category":        types.StringType,
}

// TdmysqlSpecsServerlessCcuModel is the element type of serverless_ccu_spec.
type TdmysqlSpecsServerlessCcuModel struct {
	MinCcu types.Float64 `tfsdk:"min_ccu"`
	MaxCcu types.Set     `tfsdk:"max_ccu"`
}

// TdmysqlSpecsServerlessCcuAttrTypes is the object type of the
// serverless_ccu_spec elements.
var TdmysqlSpecsServerlessCcuAttrTypes = map[string]attr.Type{
	"min_ccu": types.Float64Type,
	"max_ccu": types.SetType{ElemType: types.Float64Type},
}

// TdmysqlSpecsServerlessNodeNumModel is the element type of
// serverless_node_num_spec.
type TdmysqlSpecsServerlessNodeNumModel struct {
	MinNodeNum types.Int64 `tfsdk:"min_node_num"`
	MaxNodeNum types.Int64 `tfsdk:"max_node_num"`
}

// TdmysqlSpecsServerlessNodeNumAttrTypes is the object type of the
// serverless_node_num_spec elements.
var TdmysqlSpecsServerlessNodeNumAttrTypes = map[string]attr.Type{
	"min_node_num": types.Int64Type,
	"max_node_num": types.Int64Type,
}

func (d *TdmysqlSpecsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_tdmysql_specs"
}

func (d *TdmysqlSpecsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query TDSQL Boundless (TEO Tdmysql) salable component specs.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs ----
			"full_replications": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of replicas of hybrid instances. Valid range: 1-3, default is 3.",
			},
			"is_exclusive_instance": schema.Int64Attribute{
				Optional:    true,
				Description: "Whether the instance is exclusive (dedicated).",
			},
			"zones": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of availability zone combinations.",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the instance.",
			},
			// ---- Query outputs: Computed ----
			"hybrid_node_specs": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of salable specs of hybrid (peer) nodes.",
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
							Description: "Maximum disk size of the storage node.",
						},
						"storage_node_min_num": schema.Int64Attribute{
							Computed:    true,
							Description: "Minimum number of storage nodes.",
						},
						"storage_node_min_disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Minimum disk size of the storage node.",
						},
						"storage_type": schema.StringAttribute{
							Computed:    true,
							Description: "Disk type. Valid values: CLOUD_HSSD (enhanced SSD), CLOUD_TCS (local SSD).",
						},
						"storage_node_default_disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Default disk size of the storage node, for display purposes.",
						},
						"instance_mode": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "List of supported billing modes of the spec.",
						},
						"disk_type_category": schema.StringAttribute{
							Computed:    true,
							Description: "Disk type category. Valid values: CLOUD_DISK, LOCAL_DISK.",
						},
					},
				},
			},
			"serverless_ccu_spec": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of salable specs of serverless nodes.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"min_ccu": schema.Float64Attribute{
							Computed:    true,
							Description: "Minimum CCU value.",
						},
						"max_ccu": schema.SetAttribute{
							Computed:    true,
							ElementType: types.Float64Type,
							Description: "Range of maximum CCU values.",
						},
					},
				},
			},
			"serverless_node_num_spec": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Serverless node number configuration. Contains at most one element because the API returns a single object.",
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

	service := NewTdmysqlService(d.Client())
	request := tdmysqlv20211122.NewDescribeSpecsRequest()

	if !config.FullReplications.IsNull() && !config.FullReplications.IsUnknown() {
		request.FullReplications = config.FullReplications.ValueInt64Pointer()
	}
	if !config.IsExclusiveInstance.IsNull() && !config.IsExclusiveInstance.IsUnknown() {
		request.IsExclusiveInstance = config.IsExclusiveInstance.ValueInt64Pointer()
	}
	if !config.Zones.IsNull() && !config.Zones.IsUnknown() {
		var zones []types.String
		resp.Diagnostics.Append(config.Zones.ElementsAs(ctx, &zones, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		zoneList := make([]*string, 0, len(zones))
		for _, v := range zones {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			val := v.ValueString()
			zoneList = append(zoneList, &val)
		}
		request.Zones = zoneList
	}
	if !config.InstanceId.IsNull() && !config.InstanceId.IsUnknown() {
		request.InstanceId = config.InstanceId.ValueStringPointer()
	}

	response, err := service.DescribeTdmysqlSpecs(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying tdmysql specs",
			fmt.Sprintf("DescribeSpecs failed: %v", err),
		)
		return
	}
	if response == nil {
		resp.Diagnostics.AddError(
			"Error querying tdmysql specs",
			"DescribeSpecs returned an empty response.",
		)
		return
	}

	// hybrid_node_specs
	hybridItems := make([]TdmysqlSpecsHybridNodeSpecModel, 0, len(response.HybridNodeSpecs))
	for _, item := range response.HybridNodeSpecs {
		if item == nil {
			continue
		}
		instanceMode := types.SetNull(types.StringType)
		if item.InstanceMode != nil {
			modes := make([]string, 0, len(item.InstanceMode))
			for _, m := range item.InstanceMode {
				if m == nil {
					continue
				}
				modes = append(modes, *m)
			}
			set, diags := types.SetValueFrom(ctx, types.StringType, modes)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			instanceMode = set
		}
		hybridItems = append(hybridItems, TdmysqlSpecsHybridNodeSpecModel{
			SpecCode:               types.StringPointerValue(item.SpecCode),
			StorageNodeCpu:         types.Int64PointerValue(item.StorageNodeCpu),
			StorageNodeMem:         types.Int64PointerValue(item.StorageNodeMem),
			StorageNodeMaxNum:      types.Int64PointerValue(item.StorageNodeMaxNum),
			StorageNodeMaxDisk:     types.Int64PointerValue(item.StorageNodeMaxDisk),
			StorageNodeMinNum:      types.Int64PointerValue(item.StorageNodeMinNum),
			StorageNodeMinDisk:     types.Int64PointerValue(item.StorageNodeMinDisk),
			StorageType:            types.StringPointerValue(item.StorageType),
			StorageNodeDefaultDisk: types.Int64PointerValue(item.StorageNodeDefaultDisk),
			InstanceMode:           instanceMode,
			DiskTypeCategory:       types.StringPointerValue(item.DiskTypeCategory),
		})
	}
	hybridList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlSpecsHybridNodeSpecAttrTypes}, hybridItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// serverless_ccu_spec
	serverlessCcuItems := make([]TdmysqlSpecsServerlessCcuModel, 0, len(response.ServerlessCcuSpec))
	for _, item := range response.ServerlessCcuSpec {
		if item == nil {
			continue
		}
		maxCcu := types.SetNull(types.Float64Type)
		if item.MaxCcu != nil {
			values := make([]float64, 0, len(item.MaxCcu))
			for _, v := range item.MaxCcu {
				if v == nil {
					continue
				}
				values = append(values, *v)
			}
			set, diags := types.SetValueFrom(ctx, types.Float64Type, values)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			maxCcu = set
		}
		serverlessCcuItems = append(serverlessCcuItems, TdmysqlSpecsServerlessCcuModel{
			MinCcu: types.Float64PointerValue(item.MinCcu),
			MaxCcu: maxCcu,
		})
	}
	ccuList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlSpecsServerlessCcuAttrTypes}, serverlessCcuItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// serverless_node_num_spec: the API returns a single object, exposed as a
	// 0/1-element list to keep the mapped schema shape.
	nodeNumItems := make([]TdmysqlSpecsServerlessNodeNumModel, 0, 1)
	if response.ServerlessNodeNumSpec != nil {
		nodeNumItems = append(nodeNumItems, TdmysqlSpecsServerlessNodeNumModel{
			MinNodeNum: types.Int64PointerValue(response.ServerlessNodeNumSpec.MinNodeNum),
			MaxNodeNum: types.Int64PointerValue(response.ServerlessNodeNumSpec.MaxNodeNum),
		})
	}
	nodeNumList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlSpecsServerlessNodeNumAttrTypes}, nodeNumItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.HybridNodeSpecs = hybridList
	config.ServerlessCcuSpec = ccuList
	config.ServerlessNodeNumSpec = nodeNumList
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
