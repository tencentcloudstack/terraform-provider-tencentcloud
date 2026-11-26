// Data source for TDMSQL (TDSQL Boundless) DB engine versions, backed by
// the DescribeDBEngines API (tdmysql v20211122).

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

var _ datasource.DataSourceWithConfigure = &TdmysqlDbEnginesDataSource{}

// NewTdmysqlDbEnginesDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTdmysqlDbEnginesDataSource() datasource.DataSource {
	return &TdmysqlDbEnginesDataSource{}
}

// TdmysqlDbEnginesDataSource implements datasource.DataSource for
// tencentcloud_tdmysql_db_engines.
type TdmysqlDbEnginesDataSource struct {
	fw.DataSourceWithConfigure
}

// TdmysqlDbEnginesDataSourceModel maps the schema attributes.
type TdmysqlDbEnginesDataSourceModel struct {
	Id    types.String `tfsdk:"id"`
	Items types.List   `tfsdk:"items"`
}

// TdmysqlDbEnginesItemModel is the element type of items.
type TdmysqlDbEnginesItemModel struct {
	Type                        types.String `tfsdk:"type"`
	Version                     types.String `tfsdk:"version"`
	Name                        types.String `tfsdk:"name"`
	Description                 types.String `tfsdk:"description"`
	New                         types.Bool   `tfsdk:"new"`
	SQLMode                     types.List   `tfsdk:"sql_mode"`
	IsSupportParamTemplate      types.Bool   `tfsdk:"is_support_param_template"`
	IsSupportServerless         types.Bool   `tfsdk:"is_support_serverless"`
	IsSupportTDE                types.Bool   `tfsdk:"is_support_tde"`
	IsSupportPasswordComplexity types.Bool   `tfsdk:"is_support_password_complexity"`
}

// TdmysqlDbEnginesItemAttrTypes is the object type of the items elements.
var TdmysqlDbEnginesItemAttrTypes = map[string]attr.Type{
	"type":                           types.StringType,
	"version":                        types.StringType,
	"name":                           types.StringType,
	"description":                    types.StringType,
	"new":                            types.BoolType,
	"sql_mode":                       types.ListType{ElemType: types.StringType},
	"is_support_param_template":      types.BoolType,
	"is_support_serverless":          types.BoolType,
	"is_support_tde":                 types.BoolType,
	"is_support_password_complexity": types.BoolType,
}

func (d *TdmysqlDbEnginesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_tdmysql_db_engines"
}

func (d *TdmysqlDbEnginesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of TDSQL Boundless (tdmysql) DB engine versions.",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random token
			// generated at read time (see the end of Read).
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query outputs: Computed ----
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of DB engine information.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Engine type. Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "Engine version. Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Engine name. Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Engine description.",
						},
						"new": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether it is the latest version. Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"sql_mode": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Supported compatibility modes, separated by ,.",
						},
						"is_support_param_template": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the parameter template is supported.",
						},
						"is_support_serverless": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the Serverless mode is supported.",
						},
						"is_support_tde": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether transparent encryption (TDE) is supported.",
						},
						"is_support_password_complexity": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether password complexity is supported.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *TdmysqlDbEnginesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TdmysqlDbEnginesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeDbEnginesRequest := tdmysqlv20211122.NewDescribeDBEnginesRequest()

	service := NewTdmysqlService(d.Client())

	items, err := service.DescribeDBEngines(ctx, describeDbEnginesRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TDSQL Boundless DB engines",
			fmt.Sprintf("DescribeDBEngines failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	engines := make([]TdmysqlDbEnginesItemModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		engine := TdmysqlDbEnginesItemModel{
			Type:                        fwHelper.StringPointerValueOrNull(item.Type),
			Version:                     fwHelper.StringPointerValueOrNull(item.Version),
			Name:                        fwHelper.StringPointerValueOrNull(item.Name),
			Description:                 fwHelper.StringPointerValueOrNull(item.Description),
			New:                         fwHelper.BoolValueOrNull(item.New),
			SQLMode:                     fwHelper.ListFromPointers(ctx, types.StringType, item.SQLMode, &resp.Diagnostics),
			IsSupportParamTemplate:      fwHelper.BoolValueOrNull(item.IsSupportParamTemplate),
			IsSupportServerless:         fwHelper.BoolValueOrNull(item.IsSupportServerless),
			IsSupportTDE:                fwHelper.BoolValueOrNull(item.IsSupportTDE),
			IsSupportPasswordComplexity: fwHelper.BoolValueOrNull(item.IsSupportPasswordComplexity),
		}

		engines = append(engines, engine)
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlDbEnginesItemAttrTypes}, engines)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Items = list

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
