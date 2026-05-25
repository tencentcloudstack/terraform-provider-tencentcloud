// Package framework provides the terraform-plugin-framework flavour of the
// TencentCloud provider.
//
// This package coexists with the SDKv2-based `tencentcloud.Provider()` inside
// the same provider binary via tf5muxserver. Directory layout: everything
// framework-related (entry point, registry and business implementations)
// lives under `tencentcloud/framework/`:
//
//	tencentcloud/
//	├── framework/                 # this package: framework entry + business code
//	│   ├── provider.go            # Provider Schema/Configure (this file)
//	│   ├── registry.go            # 6-type aggregator
//	│   ├── cvm/                   # product-level subdirectory
//	│   │   └── actions/           # type-level subdirectory (package cvmactions)
//	│   └── meta/                  # cross-product / not bound to any specific cloud product
//	│       ├── resources/         # package metaresources
//	│       ├── datasources/       # package metadatasources
//	│       ├── functions/         # package metafunctions
//	│       ├── ephemerals/        # package metaephemerals
//	│       └── lists/             # package metalists
//	└── provider/
//	    └── sdkv2/                 # placeholder: future home of the SDKv2 provider entry
//
// Design notes:
//   - Credentials, SDK client, UA and retry are constructed exclusively by
//     the SDKv2 provider; this provider reuses the same
//     *connectivity.TencentCloudClient via sharedmeta.GetSharedMeta().
//   - Schema fields must mirror SDKv2 (same names, same semantics, same
//     nesting). Otherwise mux will reject user-written fields when merging
//     the two schemas.
//   - Resources/DataSources/Functions/EphemeralResources/ListResources/
//     Actions are gathered by aggregator functions in registry.go, which
//     imports product-level factory subpackages directly. There is no longer
//     an intermediate `services/tcprovider/framework.go` layer.
package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/metaschema"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdk_schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// NewProvider constructs a framework provider instance to be registered as a
// secondary server inside mux. The primary parameter is the SDKv2 provider
// instance and is currently used only for metadata (such as version
// reflection); the framework provider does not invoke any runtime logic of
// the SDKv2 provider.
func NewProvider(primary *sdk_schema.Provider) provider.ProviderWithMetaSchema {
	return &Provider{
		Version: connectivity.GetReqClientVersion(),
		Primary: primary,
	}
}

// Provider is the terraform-plugin-framework Provider implementation.
type Provider struct {
	Version string
	Primary *sdk_schema.Provider
}

// Metadata exposes the provider's type name and version, used by the
// framework runtime for identification.
func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tencentcloud"
	resp.Version = p.Version
}

// MetaSchema exposes the same module_name meta attribute as SDKv2.
func (p *Provider) MetaSchema(_ context.Context, _ provider.MetaSchemaRequest, resp *provider.MetaSchemaResponse) {
	resp.Schema = metaschema.Schema{
		Attributes: map[string]metaschema.Attribute{
			"module_name": metaschema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Schema mirrors the SDKv2 provider's field set.
//
// Key principles:
//   - Field names and type semantics match SDKv2 exactly. SDKv2's
//     TypeSet/TypeList with MaxItems=1 is equivalent to a nested block at
//     the protocol v5 layer, so we use SetNestedBlock / ListNestedBlock
//     here, not nested attributes.
//   - Every field is Optional. The framework side does not actually consume
//     these fields (they are declared but not parsed); they must exist to
//     satisfy mux's schema-consistency invariant.
//   - Sensitive fields (secret_key / security_token / token_code, etc.) are
//     explicitly marked Sensitive.
func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"secret_id": schema.StringAttribute{
				Optional:    true,
				Description: "This is the TencentCloud access key. It can also be sourced from the `TENCENTCLOUD_SECRET_ID` environment variable.",
			},
			"secret_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "This is the TencentCloud secret key. It can also be sourced from the `TENCENTCLOUD_SECRET_KEY` environment variable.",
			},
			"security_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "TencentCloud Security Token of temporary access credentials. It can be sourced from the `TENCENTCLOUD_SECURITY_TOKEN` environment variable. Notice: for supported products, please refer to: [temporary key supported products](https://intl.cloud.tencent.com/document/product/598/10588).",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "This is the TencentCloud region. It can also be sourced from the `TENCENTCLOUD_REGION` environment variable. The default input value is ap-guangzhou.",
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Description: "The protocol of the API request. Valid values: `HTTP` and `HTTPS`. Default is `HTTPS`.",
			},
			"domain": schema.StringAttribute{
				Optional:    true,
				Description: "The root domain of the API request, Default is `tencentcloudapi.com`.",
			},
			"cos_domain": schema.StringAttribute{
				Optional:    true,
				Description: "The root domain of the API request, Default is `tencentcloudapi.com`.",
			},
			"enable_pod_oidc": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to enable pod oidc.",
			},
			"shared_credentials_dir": schema.StringAttribute{
				Optional:    true,
				Description: "The directory of the shared credentials. It can also be sourced from the `TENCENTCLOUD_SHARED_CREDENTIALS_DIR` environment variable. If not set this defaults to ~/.tccli.",
			},
			"profile": schema.StringAttribute{
				Optional:    true,
				Description: "The profile name as set in the shared credentials. It can also be sourced from the `TENCENTCLOUD_PROFILE` environment variable. If not set, the default profile created with `tccli configure` will be used.",
			},
			"cam_role_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the CVM instance CAM role. It can be sourced from the `TENCENTCLOUD_CAM_ROLE_NAME` environment variable.",
			},
			"allowed_account_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of allowed TencentCloud account IDs to validate against the configured TencentCloud account.",
			},
			"forbidden_account_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of forbidden TencentCloud account IDs to check against the configured TencentCloud account.",
			},
		},
		Blocks: map[string]schema.Block{
			"assume_role": schema.SetNestedBlock{
				Description: "The `assume_role` block. If provided, terraform will attempt to assume this role using the supplied credentials.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"role_arn": schema.StringAttribute{
							Optional:    true,
							Description: "The ARN of the role to assume.",
						},
						"session_name": schema.StringAttribute{
							Optional:    true,
							Description: "The session name to use when making the AssumeRole call.",
						},
						"session_duration": schema.Int64Attribute{
							Optional:    true,
							Description: "The duration of the session when making the AssumeRole call. Its value ranges from 0 to 43200(seconds).",
						},
						"policy": schema.StringAttribute{
							Optional:    true,
							Description: "A more restrictive policy when making the AssumeRole call.",
						},
						"external_id": schema.StringAttribute{
							Optional:    true,
							Description: "External role ID.",
						},
						"source_identity": schema.StringAttribute{
							Optional:    true,
							Description: "Caller identity uin.",
						},
						"serial_number": schema.StringAttribute{
							Optional:    true,
							Description: "MFA serial number.",
						},
						"token_code": schema.StringAttribute{
							Optional:    true,
							Sensitive:   true,
							Description: "MFA authentication code.",
						},
					},
				},
			},
			"assume_role_with_saml": schema.ListNestedBlock{
				Description: "The `assume_role_with_saml` block. If provided, terraform will attempt to assume this role using the supplied credentials.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"saml_assertion": schema.StringAttribute{
							Optional:    true,
							Sensitive:   true,
							Description: "SAML assertion information encoded in base64.",
						},
						"principal_arn": schema.StringAttribute{
							Optional:    true,
							Description: "Player Access Description Name.",
						},
						"role_arn": schema.StringAttribute{
							Optional:    true,
							Description: "The ARN of the role to assume.",
						},
						"session_name": schema.StringAttribute{
							Optional:    true,
							Description: "The session name to use when making the AssumeRole call.",
						},
						"session_duration": schema.Int64Attribute{
							Optional:    true,
							Description: "The duration of the session when making the AssumeRoleWithSAML call.",
						},
					},
				},
			},
			"assume_role_with_web_identity": schema.ListNestedBlock{
				Description: "The `assume_role_with_web_identity` block. If provided, terraform will attempt to assume this role using the supplied credentials.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"provider_id": schema.StringAttribute{
							Optional:    true,
							Description: "Identity provider name.",
						},
						"web_identity_token": schema.StringAttribute{
							Optional:    true,
							Sensitive:   true,
							Description: "OIDC token issued by IdP.",
						},
						"web_identity_token_file": schema.StringAttribute{
							Optional:    true,
							Description: "File containing a web identity token from an OpenID Connect (OIDC) or OAuth provider.",
						},
						"role_arn": schema.StringAttribute{
							Optional:    true,
							Description: "The ARN of the role to assume.",
						},
						"role_arn_file": schema.StringAttribute{
							Optional:    true,
							Description: "File containin the ARN of the role to assume.",
						},
						"session_name": schema.StringAttribute{
							Optional:    true,
							Description: "The session name to use when making the AssumeRole call.",
						},
						"session_duration": schema.Int64Attribute{
							Optional:    true,
							Description: "The duration of the session when making the AssumeRoleWithWebIdentity call.",
						},
					},
				},
			},
			"mfa_certification": schema.SetNestedBlock{
				Description: "The `mfa_certification` block. If provided, terraform will attempt to use the provided credentials for MFA authentication.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"serial_number": schema.StringAttribute{
							Optional:    true,
							Description: "MFA serial number.",
						},
						"token_code": schema.StringAttribute{
							Optional:    true,
							Sensitive:   true,
							Description: "MFA authentication code.",
						},
						"duration_seconds": schema.Int64Attribute{
							Optional:    true,
							Description: "MFA token duration in seconds.",
						},
					},
				},
			},
		},
	}
}

// Configure reuses the *connectivity.TencentCloudClient that SDKv2 has
// already constructed.
//
// Key constraints:
//   - Credentials are never re-parsed here; credential logic must live in
//     the single SDKv2 providerConfigure function.
//   - Inside mux, SDKv2 is Configure-d before framework (see registration
//     order in main.go); sharedmeta.GetSharedMeta() should therefore return
//     non-nil at this point.
//   - As a defensive measure, nil only adds an Error diagnostic instead of
//     panicking, to make troubleshooting easier.
func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	client := sharedmeta.GetSharedMeta()
	if client == nil {
		resp.Diagnostics.AddError(
			"TencentCloud provider not configured",
			"The framework provider was Configure-d before the SDKv2 provider populated the shared client. "+
				"This is unexpected because the muxer registers SDKv2 first. "+
				"Please file an issue if you see this in production.",
		)
		return
	}

	meta := &sharedmeta.ProviderMeta{Client: client}
	resp.ResourceData = meta
	resp.DataSourceData = meta
	resp.EphemeralResourceData = meta
	resp.ActionData = meta
}

// Resources aggregates every framework-side resource. The concrete entries
// are gathered by frameworkResources() in tencentcloud/framework/registry.go
// using a two-level "product (service) -> type" directory layout (e.g.
// framework/meta/resources/).
func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	return frameworkResources()
}

// DataSources aggregates every framework-side data source.
func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return frameworkDataSources()
}

// Functions aggregates every framework-side provider-defined function.
func (p *Provider) Functions(_ context.Context) []func() function.Function {
	return frameworkFunctions()
}

// EphemeralResources aggregates every framework-side ephemeral resource.
func (p *Provider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return frameworkEphemeralResources()
}

// ListResources aggregates every framework-side list resource.
func (p *Provider) ListResources(_ context.Context) []func() list.ListResource {
	return frameworkListResources()
}

// Actions aggregates every framework-side action. Implementations live in
// tencentcloud/framework/<product>/actions/ product subpackages (for
// example, framework/cvm/actions/) and are gathered by frameworkActions() in
// registry.go.
func (p *Provider) Actions(_ context.Context) []func() action.Action {
	return frameworkActions()
}

// GenerateResourceConfig is required by the framework interface but is not
// currently part of any configuration-generation flow.
func (p *Provider) GenerateResourceConfig(context.Context, any) (any, error) {
	return nil, nil
}
