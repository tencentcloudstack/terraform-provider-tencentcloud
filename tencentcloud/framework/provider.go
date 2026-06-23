// Package framework provides the terraform-plugin-framework flavour of the
// TencentCloud provider.
//
// This package coexists with the SDKv2-based `tencentcloud.Provider()` inside
// the same provider binary via tf5muxserver. Directory layout: framework
// entry-points (provider, registry, tests) live under
// `tencentcloud/framework/`, while every business reference (resource,
// data source, function, ephemeral, list, action) is co-located with the
// SDKv2 implementations under `tencentcloud/services/<product>/`:
//
//	tencentcloud/
//	├── framework/                 # this package: framework entry only
//	│   ├── provider.go            # Provider Schema/Configure (this file)
//	│   ├── registry.go            # 6-type aggregator
//	│   ├── acctest/               # ProtoV5 test factories
//	│   └── internal/              # framework-only helpers
//	└── services/
//	    ├── common/                # cross-product / provider-meta references
//	    │   ├── data_source_tc_provider_runtime.go
//	    │   ├── resource_tc_local_note.go
//	    │   ├── function_tc_parse_resource_id.go
//	    │   ├── ephemeral_tc_temp_credential.go
//	    │   └── list_tc_region.go
//	    ├── cvm/                   # CVM product (SDKv2 + framework mixed)
//	    │   ├── resource_tc_instance.go        # SDKv2
//	    │   └── action_tc_cvm_reboot_instance.go  # framework
//	    └── <product>/             # other products follow the same pattern
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
//     imports product-level packages directly from
//     `tencentcloud/services/`.
package framework

import (
	"context"
	"os"

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

// MetaSchema MUST mirror SDKv2 exactly. The SDKv2 provider in
// tencentcloud/provider.go does NOT define a ProviderMetaSchemaFunc, which
// means SDKv2 reports an empty (zero-attribute) provider meta schema to
// the protocol. tf5muxserver compares the meta schema across underlying
// providers, so framework MUST also report an empty schema. Any attribute
// declared here that is not declared in SDKv2 will fail provider startup
// with "Invalid Provider Server Combination".
func (p *Provider) MetaSchema(_ context.Context, _ provider.MetaSchemaRequest, resp *provider.MetaSchemaResponse) {
	resp.Schema = metaschema.Schema{}
}

// Environment-variable keys that the SDKv2 provider binds to the following
// nested-block attributes via schema.EnvDefaultFunc(key, nil). These MUST stay
// byte-identical to the constants declared in tencentcloud/provider.go.
const (
	envAssumeRoleArn                = "TENCENTCLOUD_ASSUME_ROLE_ARN"
	envAssumeRoleSessionName        = "TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME"
	envAssumeRoleSamlAssertion      = "TENCENTCLOUD_ASSUME_ROLE_SAML_ASSERTION"
	envAssumeRolePrincipalArn       = "TENCENTCLOUD_ASSUME_ROLE_PRINCIPAL_ARN"
	envMfaCertificationSerialNumber = "TENCENTCLOUD_MFA_CERTIFICATION_SERIAL_NUMBER"
	envMfaCertificationTokenCode    = "TENCENTCLOUD_MFA_CERTIFICATION_TOKEN_CODE"
)

// requiredUnlessEnv builds a StringAttribute whose Required/Optional flags
// mirror the SDKv2 provider's *runtime* protocol schema.
//
// In SDKv2 these attributes are declared Required:true together with
// schema.EnvDefaultFunc(key, nil). helper/schema.coreConfigSchemaAttribute
// downgrades such an attribute to Optional whenever the DefaultFunc returns a
// non-nil value, i.e. whenever the environment variable is set (see
// terraform-plugin-sdk/v2/helper/schema/core_schema.go). When the variable is
// unset the attribute stays Required.
//
// Because tf5muxserver requires the framework protocol schema to be identical
// to SDKv2's, and both schemas are produced within the same GetProviderSchema
// RPC (observing the same environment), we must replicate this env-dependent
// decision here. A static Required/Optional value would fail schema-consistency
// validation in exactly the half of cases where the variable's set-ness does
// not match the hard-coded flag.
func requiredUnlessEnv(envKey, description string) schema.StringAttribute {
	attr := schema.StringAttribute{Description: description}
	if os.Getenv(envKey) != "" {
		attr.Optional = true
	} else {
		attr.Required = true
	}
	return attr
}

// Schema mirrors the SDKv2 provider's field set byte-for-byte.
//
// Key principles:
//   - Field names, types, Required/Optional, Sensitive flags and
//     Description strings MUST match SDKv2 (tencentcloud/provider.go)
//     exactly. tf5muxserver compares the protocol-level schema across
//     underlying providers; any divergence (including a single character
//     in Description) will fail provider startup with
//     "Invalid Provider Server Combination".
//   - SDKv2's TypeSet/TypeList with MaxItems=1 is equivalent to a nested
//     block at the protocol v5 layer, so we use SetNestedBlock /
//     ListNestedBlock here, not nested attributes.
//   - The framework side does not consume any of these fields at runtime;
//     credential parsing lives exclusively in the SDKv2 providerConfigure.
//     This Schema only exists to satisfy mux's schema-consistency invariant.
//   - SDKv2 ConflictsWith / ValidateFunc are NOT serialized into the
//     protocol schema, so they are intentionally not mirrored here. SDKv2
//     DefaultFunc is likewise not serialized as a value, but a Required
//     attribute combined with EnvDefaultFunc is downgraded to Optional at the
//     protocol layer when the env var is set; that env-dependent behaviour is
//     reproduced via requiredUnlessEnv (see its doc comment).
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
				Description: "This is the TencentCloud region. It can also be sourced from the `TENCENTCLOUD_REGION` environment variables. The default input value is ap-guangzhou.",
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
				Description: "The cos domain of the API request, Default is `https://cos.{region}.myqcloud.com`, Other Examples: `https://cluster-123456.cos-cdc.ap-guangzhou.myqcloud.com`.",
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
				Description: "List of allowed TencentCloud account IDs to prevent you from mistakenly using the wrong one (and potentially end up destroying a live environment). Conflicts with `forbidden_account_ids`, If use `assume_role_with_saml` or `assume_role_with_web_identity`, it is not supported.",
			},
			"forbidden_account_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of forbidden TencentCloud account IDs to prevent you from mistakenly using the wrong one (and potentially end up destroying a live environment). Conflicts with `allowed_account_ids`, If use `assume_role_with_saml` or `assume_role_with_web_identity`, it is not supported.",
			},
		},
		Blocks: map[string]schema.Block{
			"assume_role": schema.SetNestedBlock{
				Description: "The `assume_role` block. If provided, terraform will attempt to assume this role using the supplied credentials.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"role_arn": requiredUnlessEnv(envAssumeRoleArn,
							"The ARN of the role to assume. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN`."),
						"session_name": requiredUnlessEnv(envAssumeRoleSessionName,
							"The session name to use when making the AssumeRole call. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`."),
						"session_duration": schema.Int64Attribute{
							Optional:    true,
							Description: "The duration of the session when making the AssumeRole call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`.",
						},
						"policy": schema.StringAttribute{
							Optional:    true,
							Description: "A more restrictive policy when making the AssumeRole call. Its content must not contains `principal` elements. Notice: more syntax references, please refer to: [policies syntax logic](https://intl.cloud.tencent.com/document/product/598/10603).",
						},
						"external_id": schema.StringAttribute{
							Optional:    true,
							Description: "External role ID, which can be obtained by clicking the role name in the CAM console. It can contain 2-128 letters, digits, and symbols (=,.@:/-). Regex: [\\w+=,.@:/-]*. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_EXTERNAL_ID`.",
						},
						"source_identity": schema.StringAttribute{
							Optional:    true,
							Description: "Caller identity uin. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SOURCE_IDENTITY`.",
						},
						"serial_number": schema.StringAttribute{
							Optional:    true,
							Description: "MFA serial number, the identification number of the MFA device associated with the calling CAM user. Format qcs: cam:uin/${ownerUin}::mfa/${mfaType}. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SERIAL_NUMBER`.",
						},
						"token_code": schema.StringAttribute{
							Optional:    true,
							Description: "MFA authentication code. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_TOKEN_CODE`.",
						},
					},
				},
			},
			"assume_role_with_saml": schema.ListNestedBlock{
				Description: "The `assume_role_with_saml` block. If provided, terraform will attempt to assume this role using the supplied credentials.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"saml_assertion": requiredUnlessEnv(envAssumeRoleSamlAssertion,
							"SAML assertion information encoded in base64. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SAML_ASSERTION`."),
						"principal_arn": requiredUnlessEnv(envAssumeRolePrincipalArn,
							"Player Access Description Name. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_PRINCIPAL_ARN`."),
						"role_arn": requiredUnlessEnv(envAssumeRoleArn,
							"The ARN of the role to assume. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN`."),
						"session_name": requiredUnlessEnv(envAssumeRoleSessionName,
							"The session name to use when making the AssumeRole call. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`."),
						"session_duration": schema.Int64Attribute{
							Optional:    true,
							Description: "The duration of the session when making the AssumeRoleWithSAML call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`.",
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
							Description: "Identity provider name. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_PROVIDER_ID`, Default is OIDC.",
						},
						"web_identity_token": schema.StringAttribute{
							Optional:    true,
							Description: "OIDC token issued by IdP. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN`. One of `web_identity_token` or `web_identity_token_file` is required.",
						},
						"web_identity_token_file": schema.StringAttribute{
							Optional:    true,
							Description: "File containing a web identity token from an OpenID Connect (OIDC) or OAuth provider. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN_FILE`. One of `web_identity_token` or `web_identity_token_file` is required.",
						},
						"role_arn": schema.StringAttribute{
							Optional:    true,
							Description: "The ARN of the role to assume. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN`. One of `role_arn` or `role_arn_file` is required.",
						},
						"role_arn_file": schema.StringAttribute{
							Optional:    true,
							Description: "File containin the ARN of the role to assume. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARNN_FILE`. One of `role_arn` or `role_arn_file` is required.",
						},
						"session_name": requiredUnlessEnv(envAssumeRoleSessionName,
							"The session name to use when making the AssumeRole call. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`."),
						"session_duration": schema.Int64Attribute{
							Optional:    true,
							Description: "The duration of the session when making the AssumeRoleWithWebIdentity call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`.",
						},
					},
				},
			},
			"mfa_certification": schema.SetNestedBlock{
				Description: "The `mfa_certification` block. If provided, terraform will attempt to use the provided credentials for MFA authentication.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"serial_number": requiredUnlessEnv(envMfaCertificationSerialNumber,
							"MFA serial number, the identification number of the MFA device associated with the calling CAM user. Format qcs: cam:uin/${ownerUin}::mfa/${mfaType}. It can be sourced from the `TENCENTCLOUD_MFA_CERTIFICATION_SERIAL_NUMBER`."),
						"token_code": requiredUnlessEnv(envMfaCertificationTokenCode,
							"MFA authentication code. It can be sourced from the `TENCENTCLOUD_MFA_CERTIFICATION_TOKEN_CODE`."),
						"duration_seconds": schema.Int64Attribute{
							Optional:    true,
							Description: "Specify the validity period of the temporary certificate. The main account can be set to a maximum validity period of 7200 seconds, and the sub account can be set to a maximum validity period of 129600 seconds, and default is 1800 seconds. It can be sourced from the `TENCENTCLOUD_MFA_CERTIFICATION_DURATION_SECONDS`.",
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
// are gathered by frameworkResources() in tencentcloud/framework/registry.go,
// which imports factories directly from `tencentcloud/services/<product>/`
// (e.g. `services/common/` for cross-product references).
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
// tencentcloud/services/<product>/ packages (for example,
// `services/cvm/action_tc_cvm_reboot_instance.go`) and are gathered by
// frameworkActions() in registry.go.
func (p *Provider) Actions(_ context.Context) []func() action.Action {
	return frameworkActions()
}

// GenerateResourceConfig is required by the framework interface but is not
// currently part of any configuration-generation flow.
func (p *Provider) GenerateResourceConfig(context.Context, any) (any, error) {
	return nil, nil
}
