// Package metaephemerals provides framework-side ephemeral resource
// implementations for the "meta" product family — i.e. short-lived
// credential resources that are cross-product or not bound to any specific
// cloud product.
//
// Currently includes:
//   - tencentcloud_temp_credential: locally constructs a 5-minute
//     placeholder credential (a fake STS token); it does NOT call any
//     STS / CAM API.
//
// Wiring: append metaephemerals.NewTempCredentialEphemeralResource to
// frameworkEphemeralResources() in tencentcloud/framework/registry.go to
// expose the ephemeral resource through the framework provider.
package metaephemerals

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// Compile-time assertions: tempCredentialEphemeralResource implements both
// ephemeral.EphemeralResource and ephemeral.EphemeralResourceWithConfigure.
var (
	_ ephemeral.EphemeralResource              = &tempCredentialEphemeralResource{}
	_ ephemeral.EphemeralResourceWithConfigure = &tempCredentialEphemeralResource{}
)

// fakeCredentialTTL is the lifetime of the placeholder credential. It is
// kept on the same order of magnitude as a real STS short-term credential
// (real STS defaults to 30 minutes; this reference uses 5 minutes purely
// to illustrate the expires_at field calculation).
const fakeCredentialTTL = 5 * time.Minute

// NewTempCredentialEphemeralResource is the framework ephemeral factory.
func NewTempCredentialEphemeralResource() ephemeral.EphemeralResource {
	return &tempCredentialEphemeralResource{}
}

// tempCredentialEphemeralResource is the in-repo reference implementation
// of the framework ephemeral resource type, covering the full Metadata /
// Schema / Configure / Open lifecycle.
//
// Behavioural contract:
//   - Open returns a **locally-constructed placeholder credential**; no
//     remote API is called.
//   - All sensitive fields (secret_key, token) are marked Sensitive so the
//     Terraform UI does not print them in plaintext.
//   - Renew / Close are not implemented (those belong to the optional
//     EphemeralResourceWith{Renew,Close} interfaces); the framework
//     transparently handles TTL expiry.
type tempCredentialEphemeralResource struct {
	// region is the provider's default region read from sharedmeta during
	// Configure. When the user does not specify a region in HCL, Open falls
	// back to this value.
	region string
}

// tempCredentialModel is the Go counterpart of the ephemeral resource
// schema.
type tempCredentialModel struct {
	Region    types.String `tfsdk:"region"`
	SecretID  types.String `tfsdk:"secret_id"`
	SecretKey types.String `tfsdk:"secret_key"`
	Token     types.String `tfsdk:"token"`
	ExpiresAt types.String `tfsdk:"expires_at"`
}

// Metadata sets the Terraform type name to tencentcloud_temp_credential.
func (r *tempCredentialEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_temp_credential"
}

// Schema declares the ephemeral resource's attribute set.
//
// Fields:
//   - region:     input; falls back to the provider's default region when
//     unspecified.
//   - secret_id:  output; locally constructed, not sourced from real STS.
//   - secret_key: output; Sensitive.
//   - token:      output; Sensitive.
//   - expires_at: output; RFC3339 timestamp = now + 5min.
func (r *tempCredentialEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Locally-constructed placeholder STS credential. " +
			"This ephemeral resource does NOT call any TencentCloud STS / CAM API; " +
			"it serves as a reference implementation of a framework EphemeralResource " +
			"(Metadata / Schema / Configure / Open lifecycle). " +
			"Do NOT rely on the returned secret_id / secret_key / token for real API calls.",
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Region the placeholder credential is bound to. Falls back to the provider's configured region when omitted.",
			},
			"secret_id": schema.StringAttribute{
				Computed:    true,
				Description: "Locally-constructed secret id, prefixed with \"STS-fake-\". Not a real credential.",
			},
			"secret_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Locally-constructed random hex string. Not a real credential.",
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Locally-constructed random hex string. Not a real credential.",
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "RFC3339 timestamp 5 minutes after Open. Mirrors the shape of a real short-term credential.",
			},
		},
	}
}

// Configure retrieves the shared client injected by the framework provider
// during its own Configure phase and reads the provider's default region
// out of it.
func (r *tempCredentialEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
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
	if meta.Client != nil {
		r.region = meta.Client.Region
	}
}

// Open constructs a 5-minute placeholder credential and writes it to
// resp.Result.
//
// No remote calls are performed; secret_id / secret_key / token are all
// locally generated random placeholders.
func (r *tempCredentialEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var plan tempCredentialModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.region
	}

	now := time.Now().UTC()
	model := tempCredentialModel{
		Region:    types.StringValue(region),
		SecretID:  types.StringValue("STS-fake-" + randomHex(8)),
		SecretKey: types.StringValue(randomHex(16)),
		Token:     types.StringValue(randomHex(32)),
		ExpiresAt: types.StringValue(now.Add(fakeCredentialTTL).Format(time.RFC3339)),
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &model)...)
}

// randomHex returns a random hex string of length 2*n. When the underlying
// crypto/rand fails, the function falls back to a timestamp; this still
// guarantees a different value across calls.
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b)
}
