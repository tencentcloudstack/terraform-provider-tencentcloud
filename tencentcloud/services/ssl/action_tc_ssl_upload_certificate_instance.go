package ssl

import (
	"context"
	"fmt"
	"time"

	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

// _ action.Action is intentionally omitted: the factory function
// NewSslUploadCertificateInstance returns action.Action, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &SslUploadCertificateInstance{}

// NewSslUploadCertificateInstance is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewSslUploadCertificateInstance() action.Action {
	return &SslUploadCertificateInstance{}
}

// defaultSslUploadCertificateInstanceInvokeTimeout is the built-in maximum
// duration to wait for the asynchronous upload-update deploy task to finish.
// Practitioners can override it through the `timeouts` block, e.g.
// `timeouts { invoke = "30m" }`.
const defaultSslUploadCertificateInstanceInvokeTimeout = 15 * time.Minute

// SslUploadCertificateInstance implements action.Action for
// tencentcloud_ssl_upload_certificate_instance.
type SslUploadCertificateInstance struct {
	fw.ActionWithConfigure
}

// SslUploadCertificateInstanceModel maps the schema attributes.
type SslUploadCertificateInstanceModel struct {
	OldCertificateId      types.String `tfsdk:"old_certificate_id"`
	ResourceTypes         types.Set    `tfsdk:"resource_types"`
	CertificatePublicKey  types.String `tfsdk:"certificate_public_key"`
	CertificatePrivateKey types.String `tfsdk:"certificate_private_key"`
	CertificateId         types.String `tfsdk:"certificate_id"`

	// ResourceTypesRegions maps the resource_types_regions nested block.
	ResourceTypesRegions types.List `tfsdk:"resource_types_regions"`

	// Timeouts is the value type defined by the official
	// terraform-plugin-framework-timeouts module. Only the `invoke` stage
	// exists because an action has a single lifecycle stage.
	Timeouts actiontimeouts.Value `tfsdk:"timeouts"`
}

// SslUploadCertificateInstanceRegionModel maps the resource_types_regions
// nested block.
type SslUploadCertificateInstanceRegionModel struct {
	ResourceType types.String `tfsdk:"resource_type"`
	Regions      types.Set    `tfsdk:"regions"`
}

func (a *SslUploadCertificateInstance) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_ssl_upload_certificate_instance"
}

func (a *SslUploadCertificateInstance) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to upload new certificate content and update the associated cloud " +
			"resources without changing the certificate ID via the UploadUpdateCertificateInstance API. " +
			"The new certificate is given either by an existing `certificate_id` or by a pair of " +
			"`certificate_public_key`/`certificate_private_key`. " +
			"Currently, only the `clb` resource type is supported. " +
			"The action waits for the asynchronous deploy task to complete. This is a one-time operation; no cloud-side state is persisted.",
		Attributes: map[string]schema.Attribute{
			"old_certificate_id": schema.StringAttribute{
				Required:    true,
				Description: "Update the original certificate ID. The cloud resources bound to this certificate ID are queried, and then these cloud resources are updated with the new certificate.",
			},
			"resource_types": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "The resource type that needs to be deployed. Currently, only `clb` is supported.",
			},
			"certificate_public_key": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate public key content. Either specify `certificate_id`, or specify both `certificate_public_key` and `certificate_private_key`. Ignored if `certificate_id` is set. This is a write-only attribute: the value is only used when the action runs and is never persisted or shown in the plan output.",
			},
			"certificate_private_key": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate private key content. Either specify `certificate_id`, or specify both `certificate_public_key` and `certificate_private_key`. Ignored if `certificate_id` is set. This is a write-only attribute: the value is only used when the action runs and is never persisted or shown in the plan output.",
			},
			"certificate_id": schema.StringAttribute{
				Optional:    true,
				Description: "New certificate ID. Either specify `certificate_id`, or specify both `certificate_public_key` and `certificate_private_key`. If set, it takes precedence over the key content.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": sslUploadCertificateInstanceTimeoutsBlock(ctx),
			"resource_types_regions": schema.ListNestedBlock{
				Description: "List of regions where cloud resources need to be deployed.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"resource_type": schema.StringAttribute{
							Optional:    true,
							Description: "Cloud resource type.",
						},
						"regions": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Region list.",
						},
					},
				},
			},
		},
	}
}

// sslUploadCertificateInstanceTimeoutsBlock builds the `timeouts` block from
// the official terraform-plugin-framework-timeouts module, whose single
// optional `invoke` attribute is validated as a time.Duration at plan time.
// That module leaves the block itself undescribed (only the attribute carries
// a description), so the block description is restored here to keep the schema
// - and the documentation generated from it by gendoc - self-explanatory.
func sslUploadCertificateInstanceTimeoutsBlock(ctx context.Context) schema.Block {
	block := actiontimeouts.BlockWithOpts(ctx, actiontimeouts.Opts{
		InvokeDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `15m`. " +
			"It limits how long the action waits for the asynchronous deploy task to complete. Default is `15m`.",
	})

	if nested, ok := block.(schema.SingleNestedBlock); ok {
		nested.Description = "The timeouts block allows you to specify the timeout for the invoke operation."
		return nested
	}

	return block
}

// Invoke is called to run the logic of the action.
func (a *SslUploadCertificateInstance) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SslUploadCertificateInstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.OldCertificateId.IsNull() || data.OldCertificateId.IsUnknown() || data.OldCertificateId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing old_certificate_id",
			"old_certificate_id is required and cannot be empty.",
		)
		return
	}

	// The API requires the new certificate to be given either by an existing
	// certificate ID or by a pair of certificate public/private keys.
	certificateId := ""
	if !data.CertificateId.IsNull() && !data.CertificateId.IsUnknown() {
		certificateId = data.CertificateId.ValueString()
	}
	certificatePublicKey := ""
	if !data.CertificatePublicKey.IsNull() && !data.CertificatePublicKey.IsUnknown() {
		certificatePublicKey = data.CertificatePublicKey.ValueString()
	}
	certificatePrivateKey := ""
	if !data.CertificatePrivateKey.IsNull() && !data.CertificatePrivateKey.IsUnknown() {
		certificatePrivateKey = data.CertificatePrivateKey.ValueString()
	}
	if certificateId == "" && (certificatePublicKey == "" || certificatePrivateKey == "") {
		resp.Diagnostics.AddError(
			"Missing certificate",
			"Either `certificate_id` or both `certificate_public_key` and `certificate_private_key` must be specified.",
		)
		return
	}

	var resourceTypes []*string
	if !data.ResourceTypes.IsNull() && !data.ResourceTypes.IsUnknown() {
		elems := make([]types.String, 0, len(data.ResourceTypes.Elements()))
		resp.Diagnostics.Append(data.ResourceTypes.ElementsAs(ctx, &elems, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range elems {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			resourceType := v.ValueString()
			resourceTypes = append(resourceTypes, &resourceType)
		}
	}
	if len(resourceTypes) == 0 {
		resp.Diagnostics.AddError(
			"Missing resource_types",
			"resource_types is required and cannot be empty.",
		)
		return
	}

	var resourceTypesRegions []*ssl.ResourceTypeRegions
	if !data.ResourceTypesRegions.IsNull() && !data.ResourceTypesRegions.IsUnknown() {
		var regions []SslUploadCertificateInstanceRegionModel
		resp.Diagnostics.Append(data.ResourceTypesRegions.ElementsAs(ctx, &regions, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, item := range regions {
			resourceTypeRegions := ssl.ResourceTypeRegions{}
			if !item.ResourceType.IsNull() && !item.ResourceType.IsUnknown() {
				resourceType := item.ResourceType.ValueString()
				resourceTypeRegions.ResourceType = &resourceType
			}
			if !item.Regions.IsNull() && !item.Regions.IsUnknown() {
				elems := make([]types.String, 0, len(item.Regions.Elements()))
				resp.Diagnostics.Append(item.Regions.ElementsAs(ctx, &elems, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				for _, v := range elems {
					if v.IsNull() || v.IsUnknown() {
						continue
					}
					region := v.ValueString()
					resourceTypeRegions.Regions = append(resourceTypeRegions.Regions, &region)
				}
			}
			resourceTypesRegions = append(resourceTypesRegions, &resourceTypeRegions)
		}
	}

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	// An unset/null timeouts block falls back to the built-in default, and a
	// malformed duration is reported as an error diagnostic. The schema
	// validator normally rejects malformed values at plan time; this is the
	// runtime safety net. The resolved timeout bounds the polling loop that
	// waits for the asynchronous deploy task to complete.
	timeout, timeoutDiags := data.Timeouts.Invoke(ctx, defaultSslUploadCertificateInstanceInvokeTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	oldCertificateId := data.OldCertificateId.ValueString()

	request := ssl.NewUploadUpdateCertificateInstanceRequest()
	request.OldCertificateId = &oldCertificateId
	request.ResourceTypes = resourceTypes
	if certificateId != "" {
		request.CertificateId = &certificateId
		// If CertificateId is provided, the public/private key content is not required. Clear them to avoid conflicts.
		// The SSL R&D team has advised passing an empty string.
		request.CertificatePublicKey = helper.String("")
		request.CertificatePrivateKey = helper.String("")
	} else {
		request.CertificatePublicKey = helper.String(certificatePublicKey)
		request.CertificatePrivateKey = helper.String(certificatePrivateKey)
	}
	request.ResourceTypesRegions = resourceTypesRegions

	service := NewSSLService(a.Client())

	deployRecordId, err := service.UploadUpdateCertificateInstance(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("uploading ssl certificate instance for old certificate (%s)", oldCertificateId),
			err.Error(),
		)
		return
	}

	if err := service.DescribeHostUploadUpdateRecordStatus(ctx, oldCertificateId, deployRecordId, timeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("waiting for ssl upload certificate instance deploy record (%d) to complete", deployRecordId),
			err.Error(),
		)
		return
	}
}
