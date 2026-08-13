package ssm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var _ ephemeral.EphemeralResource = &SsmSecretVersionEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &SsmSecretVersionEphemeralResource{}

// NewSsmSecretVersionEphemeralResource is the factory referenced by
// tencentcloud/framework/registry.go to register this ephemeral resource.
func NewSsmSecretVersionEphemeralResource() ephemeral.EphemeralResource {
	return &SsmSecretVersionEphemeralResource{}
}

// SsmSecretVersionEphemeralResource implements ephemeral.EphemeralResource for
// tencentcloud_ssm_secret_version.
type SsmSecretVersionEphemeralResource struct {
	client SsmService
}

// SsmSecretVersionEphemeralResourceModel maps the schema attributes.
type SsmSecretVersionEphemeralResourceModel struct {
	SecretName   types.String `tfsdk:"secret_name"`
	VersionId    types.String `tfsdk:"version_id"`
	SecretBinary types.String `tfsdk:"secret_binary"`
	SecretString types.String `tfsdk:"secret_string"`
}

func (e *SsmSecretVersionEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "tencentcloud_ssm_secret_version"
}

func (e *SsmSecretVersionEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a secret version value from TencentCloud SSM (Secrets Manager). " +
			"This ephemeral resource reads the secret value at plan/apply time; the value " +
			"is not persisted in state and expires after use.",
		Attributes: map[string]schema.Attribute{
			"secret_name": schema.StringAttribute{
				Required:    true,
				Description: "Specifies the name of the credential to which the version to be read belongs.",
			},
			"version_id": schema.StringAttribute{
				Required:    true,
				Description: "Specifies the version ID of the secret version to read.",
			},
			"secret_binary": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The binary credential information of the secret version, encoded using Base64.",
			},
			"secret_string": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The text-based credential information of the secret version.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure (which itself reuses the SDKv2-built
// *connectivity.TencentCloudClient). We then construct the SDKv2-style
// SsmService directly from the shared client, mirroring how SDKv2 resources
// in this package construct it.
func (e *SsmSecretVersionEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData type",
			fmt.Sprintf("Expected *sharedmeta.ProviderMeta, got: %T. "+
				"This is a bug in the provider; please report it.", req.ProviderData),
		)
		return
	}

	e.client = NewSsmService(meta.Client)
}

// Open is invoked once per plan/apply. It reads the secret value via the SSM
// API and writes it back to resp.Result.
func (e *SsmSecretVersionEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data SsmSecretVersionEphemeralResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretName := data.SecretName.ValueString()
	versionId := data.VersionId.ValueString()

	if secretName == "" {
		resp.Diagnostics.AddError("Missing secret_name", "secret_name is required and cannot be empty.")
		return
	}
	if versionId == "" {
		resp.Diagnostics.AddError("Missing version_id", "version_id is required and cannot be empty.")
		return
	}

	secretVersion, err := e.client.DescribeSecretVersion(ctx, secretName, versionId)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFound") {
			resp.Diagnostics.AddError(
				"Secret version not found",
				fmt.Sprintf("Secret %q version %q not found: %v", secretName, versionId, err),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading SSM secret version",
			fmt.Sprintf("Could not read secret %q version %q: %v", secretName, versionId, err),
		)
		return
	}
	if secretVersion == nil {
		resp.Diagnostics.AddError(
			"Secret version not found",
			fmt.Sprintf("Secret %q version %q returned nil", secretName, versionId),
		)
		return
	}

	data.SecretBinary = types.StringValue(secretVersion.secretBinary)
	data.SecretString = types.StringValue(secretVersion.secretString)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
