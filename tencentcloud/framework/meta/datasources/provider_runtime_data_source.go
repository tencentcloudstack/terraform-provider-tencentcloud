package metadatasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// Compile-time assertions: ProviderRuntimeDataSource satisfies
// datasource.DataSource and datasource.DataSourceWithConfigure.
var (
	_ datasource.DataSource              = &ProviderRuntimeDataSource{}
	_ datasource.DataSourceWithConfigure = &ProviderRuntimeDataSource{}
)

// NewProviderRuntimeDataSource is the data source factory expected by the
// framework registration convention. It is referenced from the DataSources
// aggregator slice in tencentcloud/framework/registry.go.
func NewProviderRuntimeDataSource() datasource.DataSource {
	return &ProviderRuntimeDataSource{}
}

// ProviderRuntimeDataSource exposes the provider's current runtime
// information: region, SDK client version and whether the provider is
// running in the SDKv2 + framework dual-stack mode.
//
// The data source never calls any cloud API; it only reads the
// in-process provider runtime context, making it a suitable first
// reference implementation for the dual-stack architecture.
type ProviderRuntimeDataSource struct {
	// client is fetched from *sharedmeta.ProviderMeta during Configure;
	// it is the same instance shared with the SDKv2 resources.
	client *connectivity.TencentCloudClient
}

// providerRuntimeModel is the Go counterpart of the data source schema.
// All fields are computed inside Read and are therefore Computed.
type providerRuntimeModel struct {
	ID              types.String `tfsdk:"id"`
	Region          types.String `tfsdk:"region"`
	ClientVersion   types.String `tfsdk:"client_version"`
	StackMode       types.String `tfsdk:"stack_mode"`
	Protocol        types.String `tfsdk:"protocol"`
	Domain          types.String `tfsdk:"domain"`
	CosDomain       types.String `tfsdk:"cos_domain"`
	SecretIDPresent types.Bool   `tfsdk:"secret_id_present"`
}

// Metadata sets the Terraform type name. Users reference this data source
// from HCL as:
//
//	data "tencentcloud_provider_runtime" "this" {}
func (d *ProviderRuntimeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_runtime"
}

// Schema declares the data source attributes.
func (d *ProviderRuntimeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read-only provider runtime metadata. " +
			"Useful for debugging which region the provider currently targets " +
			"and which client version is running. Does not call any TencentCloud API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Synthetic id, equal to the current region.",
			},
			"region": schema.StringAttribute{
				Computed:    true,
				Description: "The region the provider is currently configured to use.",
			},
			"client_version": schema.StringAttribute{
				Computed:    true,
				Description: "The X-TC-RequestClient version reported by the provider.",
			},
			"stack_mode": schema.StringAttribute{
				Computed: true,
				Description: "Always returns \"sdkv2+framework\" for this provider, " +
					"indicating that resources can be implemented in either stack and are served via tf5muxserver.",
			},
			"protocol": schema.StringAttribute{
				Computed:    true,
				Description: "The API request protocol the provider currently uses (typically `HTTPS`).",
			},
			"domain": schema.StringAttribute{
				Computed:    true,
				Description: "The API root domain the provider currently uses. Empty string when the default TencentCloud public domain is used.",
			},
			"cos_domain": schema.StringAttribute{
				Computed:    true,
				Description: "The COS root domain the provider currently uses. Empty string when the default public COS domain is used.",
			},
			"secret_id_present": schema.BoolAttribute{
				Computed: true,
				Description: "Whether a non-empty SecretId is currently configured on the provider. " +
					"**This field only indicates whether SecretId is configured, NOT whether the credential is valid.** " +
					"No real credential value is ever exposed by this data source.",
			},
		},
	}
}

// Configure retrieves the shared client injected by the framework provider
// during its own Configure phase. When ProviderData is nil or has the wrong
// type, the function appends a diagnostic per the framework convention
// rather than panicking.
func (d *ProviderRuntimeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// The framework may invoke Configure once with a nil ProviderData
		// during early lifecycle; silently return and wait for the next
		// callback after configuration is fully established.
		return
	}

	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			"Expected *sharedmeta.ProviderMeta, please report this issue to the provider maintainers.",
		)
		return
	}
	d.client = meta.Client
}

// Read computes and returns provider runtime metadata. This method does
// not call any cloud API.
func (d *ProviderRuntimeDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		// Defensive: in theory this should never be nil after Configure;
		// guard against a nil deref by returning an empty region.
		resp.Diagnostics.AddWarning(
			"Provider client unavailable",
			"The provider runtime data source was Read before the framework provider completed Configure.",
		)
	}

	region := ""
	protocol := ""
	domain := ""
	cosDomain := ""
	secretIDPresent := false
	if d.client != nil {
		region = d.client.Region
		protocol = d.client.Protocol
		domain = d.client.Domain
		cosDomain = d.client.CosDomain
		secretIDPresent = d.client.Credential != nil && d.client.Credential.SecretId != ""
	}

	state := providerRuntimeModel{
		ID:              helper.StringValueOrNull(region),
		Region:          helper.StringValueOrNull(region),
		ClientVersion:   types.StringValue(connectivity.GetReqClientVersion()),
		StackMode:       types.StringValue("sdkv2+framework"),
		Protocol:        helper.StringValueOrNull(protocol),
		Domain:          helper.StringValueOrNull(domain),
		CosDomain:       helper.StringValueOrNull(cosDomain),
		SecretIDPresent: types.BoolValue(secretIDPresent),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
