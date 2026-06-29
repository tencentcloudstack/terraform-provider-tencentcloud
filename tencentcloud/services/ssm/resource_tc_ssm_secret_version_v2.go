package ssm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var (
	_ resource.Resource                = &SsmSecretVersionV2Resource{}
	_ resource.ResourceWithConfigure   = &SsmSecretVersionV2Resource{}
	_ resource.ResourceWithImportState = &SsmSecretVersionV2Resource{}
)

// NewSsmSecretVersionV2Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewSsmSecretVersionV2Resource() resource.Resource {
	return &SsmSecretVersionV2Resource{}
}

// SsmSecretVersionV2Resource implements resource.Resource for
// tencentcloud_ssm_secret_version_v2.
type SsmSecretVersionV2Resource struct {
	client SsmService
}

// SsmSecretVersionV2ResourceModel maps the schema attributes.
//
// SecretStringWo is the write-only counterpart of SecretString. Terraform
// never persists a write-only attribute in plan or state, so this field
// only carries data during a single Create/Update operation read from
// req.Config and is always types.StringNull() in plan/state.
type SsmSecretVersionV2ResourceModel struct {
	Id             types.String `tfsdk:"id"`
	SecretName     types.String `tfsdk:"secret_name"`
	VersionId      types.String `tfsdk:"version_id"`
	SecretBinary   types.String `tfsdk:"secret_binary"`
	SecretString   types.String `tfsdk:"secret_string"`
	SecretStringWo types.String `tfsdk:"secret_string_wo"`
}

func (r *SsmSecretVersionV2Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ssm_secret_version_v2"
}

func (r *SsmSecretVersionV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a SSM secret version (terraform-plugin-framework implementation). " +
			"Exactly one of secret_binary, secret_string or secret_string_wo must be set. " +
			"secret_string_wo is a write-only attribute whose value is never written to plan or state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<secret_name>#<version_id>`.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
			"secret_name": schema.StringAttribute{
				Required:    true,
				Description: "Specifies the name of the credential to which the new version is to be added.",
				PlanModifiers: []planmodifier.String{
					stringRequiresReplace{},
				},
			},
			"version_id": schema.StringAttribute{
				Required: true,
				Description: "Specifies the version ID for the newly added version. It can be up to 64 bytes in length and " +
					"must consist of a combination of letters, numbers, and the characters `-`, `_`, or `.`, " +
					"starting with a letter or a number.",
				PlanModifiers: []planmodifier.String{
					stringRequiresReplace{},
				},
			},
			"secret_binary": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Binary credential information, encoded using Base64. Conflicts with secret_string and secret_string_wo.",
			},
			"secret_string": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Text-based credential information in plaintext (Base64 encoding is not required). Conflicts with secret_binary and secret_string_wo.",
			},
			"secret_string_wo": schema.StringAttribute{
				Optional:    true,
				WriteOnly:   true,
				Sensitive:   true,
				Description: "Write-only text-based credential information. Its value is sent to the API but never persisted in Terraform plan or state. Conflicts with secret_binary and secret_string.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure and constructs the SDKv2-style SsmService
// from the shared *connectivity.TencentCloudClient.
func (r *SsmSecretVersionV2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = NewSsmService(meta.Client)
}

// Create handles `terraform apply` for new resources.
func (r *SsmSecretVersionV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SsmSecretVersionV2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the write-only attribute from Config — it is never present in plan/state.
	var configWo types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("secret_string_wo"), &configWo)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretName := plan.SecretName.ValueString()
	versionId := plan.VersionId.ValueString()

	param, err := r.buildSecretValueParam(secretName, versionId, plan.SecretBinary, plan.SecretString, configWo)
	if err != nil {
		resp.Diagnostics.AddError("Invalid secret value combination", err.Error())
		return
	}

	gotName, gotVersion, err := r.client.PutSecretValue(ctx, param)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SSM secret version",
			fmt.Sprintf("PutSecretValue failed for secret %q version %q: %v", secretName, versionId, err),
		)
		return
	}

	plan.Id = types.StringValue(strings.Join([]string{gotName, gotVersion}, tccommon.FILED_SP))

	// Persist plan; a write-only attribute is automatically null-ified in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *SsmSecretVersionV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SsmSecretVersionV2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretName, versionId, ok := splitSsmVersionV2Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("SSM secret version id is broken, id is %q (expected `<secret_name>#<version_id>`)", state.Id.ValueString()),
		)
		return
	}

	secretInfo, err := r.client.DescribeSecretByName(ctx, secretName)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFound") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading SSM secret",
			fmt.Sprintf("DescribeSecretByName failed for %q: %v", secretName, err),
		)
		return
	}

	versionIds, err := r.client.DescribeSecretVersionIdsByName(ctx, secretName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing SSM secret versions",
			fmt.Sprintf("DescribeSecretVersionIdsByName failed for %q: %v", secretName, err),
		)
		return
	}

	hasVersion := false
	for _, id := range versionIds {
		if id == versionId {
			hasVersion = true
			break
		}
	}
	if !hasVersion {
		resp.State.RemoveResource(ctx)
		return
	}

	state.SecretName = types.StringValue(secretName)
	state.VersionId = types.StringValue(versionId)

	if secretInfo.Status() == SSM_STATUS_ENABLED {
		secretVersion, err := r.client.DescribeSecretVersion(ctx, secretName, versionId)
		if err != nil {
			if strings.Contains(err.Error(), "ResourceNotFound") {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(
				"Error reading SSM secret version",
				fmt.Sprintf("DescribeSecretVersion failed for %q version %q: %v", secretName, versionId, err),
			)
			return
		}

		state.SecretBinary = stringValueOrNull(secretVersion.secretBinary)
		state.SecretString = stringValueOrNull(secretVersion.secretString)
	}

	// secret_string_wo is write-only — never refresh it from API into state.
	state.SecretStringWo = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
func (r *SsmSecretVersionV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SsmSecretVersionV2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SsmSecretVersionV2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretName, versionId, ok := splitSsmVersionV2Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("SSM secret version id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	// Read write-only value from Config.
	var configWo types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("secret_string_wo"), &configWo)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretInfo, err := r.client.DescribeSecretByName(ctx, secretName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SSM secret",
			fmt.Sprintf("DescribeSecretByName failed for %q: %v", secretName, err),
		)
		return
	}

	if secretInfo.Status() == SSM_STATUS_ENABLED {
		// Decide whether the secret content changed. For secret_string_wo we
		// detect "field is set in config" because the value is not in plan/state.
		binaryChanged := !plan.SecretBinary.Equal(state.SecretBinary)
		stringChanged := !plan.SecretString.Equal(state.SecretString)
		woProvided := !configWo.IsNull() && !configWo.IsUnknown()

		if binaryChanged || stringChanged || woProvided {
			param, perr := r.buildSecretValueParam(secretName, versionId, plan.SecretBinary, plan.SecretString, configWo)
			if perr != nil {
				resp.Diagnostics.AddError("Invalid secret value combination", perr.Error())
				return
			}

			if err := r.client.UpdateSecret(ctx, param); err != nil {
				resp.Diagnostics.AddError(
					"Error updating SSM secret version",
					fmt.Sprintf("UpdateSecret failed for %q version %q: %v", secretName, versionId, err),
				)
				return
			}
		}
	}

	plan.Id = state.Id
	plan.SecretStringWo = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`.
func (r *SsmSecretVersionV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SsmSecretVersionV2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretName, versionId, ok := splitSsmVersionV2Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("SSM secret version id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	if err := r.client.DeleteSecretVersion(ctx, secretName, versionId); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting SSM secret version",
			fmt.Sprintf("DeleteSecretVersion failed for %q version %q: %v", secretName, versionId, err),
		)
		return
	}
}

// ImportState supports `terraform import` using the `<secret_name>#<version_id>` form.
func (r *SsmSecretVersionV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := splitSsmVersionV2Id(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<secret_name>%s<version_id>`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// buildSecretValueParam validates the "exactly one of" rule across
// secret_binary / secret_string / secret_string_wo and produces the
// param map consumed by SsmService.PutSecretValue / UpdateSecret.
func (r *SsmSecretVersionV2Resource) buildSecretValueParam(
	secretName, versionId string,
	secretBinary, secretString, secretStringWo types.String,
) (map[string]interface{}, error) {
	param := map[string]interface{}{
		"secret_name": secretName,
		"version_id":  versionId,
	}

	count := 0
	if !secretBinary.IsNull() && !secretBinary.IsUnknown() {
		param["secret_binary"] = secretBinary.ValueString()
		count++
	}
	if !secretString.IsNull() && !secretString.IsUnknown() {
		param["secret_string"] = secretString.ValueString()
		count++
	}
	if !secretStringWo.IsNull() && !secretStringWo.IsUnknown() {
		// On the wire, secret_string_wo is sent as the regular SecretString field.
		param["secret_string"] = secretStringWo.ValueString()
		count++
	}

	if count != 1 {
		return nil, fmt.Errorf(
			"exactly one of secret_binary, secret_string or secret_string_wo must be set, got %d",
			count,
		)
	}
	return param, nil
}

// splitSsmVersionV2Id parses `<secret_name>#<version_id>`.
func splitSsmVersionV2Id(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func stringValueOrNull(v string) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}

// stringRequiresReplace forces resource replacement when a string attribute changes.
type stringRequiresReplace struct{}

func (stringRequiresReplace) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (stringRequiresReplace) MarkdownDescription(ctx context.Context) string {
	return stringRequiresReplace{}.Description(ctx)
}

func (stringRequiresReplace) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		// Resource is being created; nothing to replace.
		return
	}
	if req.PlanValue.IsNull() {
		// Resource is being destroyed; nothing to replace.
		return
	}
	if !req.PlanValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}

// stringUseStateForUnknown copies the prior state value into the plan when
// the new plan value is unknown — used for the computed `id` attribute so it
// stays stable across refreshes.
type stringUseStateForUnknown struct{}

func (stringUseStateForUnknown) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

func (stringUseStateForUnknown) MarkdownDescription(ctx context.Context) string {
	return stringUseStateForUnknown{}.Description(ctx)
}

func (stringUseStateForUnknown) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = req.StateValue
}
