package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &TeoCustomGroupReferencesDataSource{}

// NewTeoCustomGroupReferencesDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTeoCustomGroupReferencesDataSource() datasource.DataSource {
	return &TeoCustomGroupReferencesDataSource{}
}

// TeoCustomGroupReferencesDataSource implements datasource.DataSource for
// the tencentcloud_teo_custom_group_references data source.
type TeoCustomGroupReferencesDataSource struct {
	fw.DataSourceWithConfigure
}

// TeoCustomGroupReferencesDataSourceModel maps the schema attributes.
type TeoCustomGroupReferencesDataSourceModel struct {
	ZoneId  types.String `tfsdk:"zone_id"`
	GroupId types.String `tfsdk:"group_id"`

	References types.List `tfsdk:"references"`
}

// TeoCustomGroupReferenceItemModel is the element type of references.
type TeoCustomGroupReferenceItemModel struct {
	ZoneId        types.String `tfsdk:"zone_id"`
	EntityType    types.String `tfsdk:"entity_type"`
	EntityId      types.String `tfsdk:"entity_id"`
	EntityName    types.String `tfsdk:"entity_name"`
	SubEntityType types.String `tfsdk:"sub_entity_type"`
	SubEntityId   types.String `tfsdk:"sub_entity_id"`
	SubEntityName types.String `tfsdk:"sub_entity_name"`
}

// TeoCustomGroupReferenceItemAttrTypes is the object type of the references elements.
var TeoCustomGroupReferenceItemAttrTypes = map[string]attr.Type{
	"zone_id":         types.StringType,
	"entity_type":     types.StringType,
	"entity_id":       types.StringType,
	"entity_name":     types.StringType,
	"sub_entity_type": types.StringType,
	"sub_entity_id":   types.StringType,
	"sub_entity_name": types.StringType,
}

func (d *TeoCustomGroupReferencesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_custom_group_references"
}

func (d *TeoCustomGroupReferencesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the reference information of a specified TEO (EdgeOne) custom group.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs ----
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Site ID.",
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "Custom group ID.",
			},
			// ---- Query outputs: Computed ----
			"references": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of configurations that reference the custom group.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"zone_id": schema.StringAttribute{
							Computed:    true,
							Description: "Site ID.",
						},
						"entity_type": schema.StringAttribute{
							Computed: true,
							Description: "Entity type. Valid values: `WebSec.ZonePolicy` (zone-level security policy), " +
								"`WebSec.HostPolicy` (host-level security policy), `WebSec.Template` (policy template), " +
								"`WebSec.APIResource` (API resource), `WebSec.JSInjectionRule` (JavaScript injection rule), " +
								"`DDoS.L4Proxy` (layer-4 proxy DDoS protection), `DDoS.L3Transit` (layer-3 transit DDoS protection).",
						},
						"entity_id": schema.StringAttribute{
							Computed: true,
							Description: "Entity ID. Its meaning depends on EntityType: zone ID for `WebSec.ZonePolicy`, " +
								"domain for `WebSec.HostPolicy`, template ID for `WebSec.Template`, API resource ID for " +
								"`WebSec.APIResource`, JavaScript injection rule ID for `WebSec.JSInjectionRule`, " +
								"instance ID for `DDoS.L4Proxy` and `DDoS.L3Transit`.",
						},
						"entity_name": schema.StringAttribute{
							Computed: true,
							Description: "Entity name. Its meaning depends on EntityType: template name for " +
								"`WebSec.Template`, API resource name for `WebSec.APIResource`, JavaScript injection " +
								"rule name for `WebSec.JSInjectionRule`; empty for `WebSec.ZonePolicy`, " +
								"`WebSec.HostPolicy`, `DDoS.L4Proxy` and `DDoS.L3Transit`.",
						},
						"sub_entity_type": schema.StringAttribute{
							Computed: true,
							Description: "Sub entity type. Valid values: `WebSec.ExceptionRule` (protection exception rule), " +
								"`WebSec.BasicAccessRule` (basic access control), `WebSec.PreciseMatchRule` " +
								"(precise match rule), `WebSec.RateLimitRule` (rate limit rule), `WebSec.BotCustomRule` " +
								"(custom bot rule), `DDoS.L4Proxy.IpAccessControl` (layer-4 proxy IP allowlist/blocklist), " +
								"`DDoS.L3Transit.IpAccessControl` (layer-3 transit IP allowlist/blocklist).",
						},
						"sub_entity_id": schema.StringAttribute{
							Computed: true,
							Description: "Sub entity ID, i.e. the rule ID of the corresponding SubEntityType; empty for " +
								"`DDoS.L4Proxy.IpAccessControl` and `DDoS.L3Transit.IpAccessControl`.",
						},
						"sub_entity_name": schema.StringAttribute{
							Computed: true,
							Description: "Sub entity name, i.e. the rule name of the corresponding SubEntityType; empty for " +
								"`DDoS.L4Proxy.IpAccessControl` and `DDoS.L3Transit.IpAccessControl`.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *TeoCustomGroupReferencesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoCustomGroupReferencesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := teov20220901.NewDescribeCustomGroupReferencesRequest()
	if !config.ZoneId.IsNull() && !config.ZoneId.IsUnknown() {
		request.ZoneId = config.ZoneId.ValueStringPointer()
	}
	if !config.GroupId.IsNull() && !config.GroupId.IsUnknown() {
		request.GroupId = config.GroupId.ValueStringPointer()
	}

	service := NewTeoService(d.Client())
	references, err := service.DescribeTeoCustomGroupReferences(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TEO custom group references",
			fmt.Sprintf("DescribeCustomGroupReferences failed: %v", err),
		)
		return
	}

	// A nil response is reported by the service layer above; an empty
	// references list is a genuinely empty result, not a failure.
	items := make([]TeoCustomGroupReferenceItemModel, 0, len(references))
	for _, reference := range references {
		items = append(items, TeoCustomGroupReferenceItemModel{
			ZoneId:        types.StringPointerValue(reference.ZoneId),
			EntityType:    types.StringPointerValue(reference.EntityType),
			EntityId:      types.StringPointerValue(reference.EntityId),
			EntityName:    types.StringPointerValue(reference.EntityName),
			SubEntityType: types.StringPointerValue(reference.SubEntityType),
			SubEntityId:   types.StringPointerValue(reference.SubEntityId),
			SubEntityName: types.StringPointerValue(reference.SubEntityName),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: TeoCustomGroupReferenceItemAttrTypes,
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.References = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
