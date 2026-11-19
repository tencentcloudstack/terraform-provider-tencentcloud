package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ action.Action is intentionally omitted: the factory function
// NewTeoApplyFreeCertificate returns action.Action, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &TeoApplyFreeCertificate{}

// NewTeoApplyFreeCertificate is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewTeoApplyFreeCertificate() action.Action {
	return &TeoApplyFreeCertificate{}
}

// TeoApplyFreeCertificate implements action.Action for
// tencentcloud_teo_apply_free_certificate.
type TeoApplyFreeCertificate struct {
	fw.ActionWithConfigure
}

// TeoApplyFreeCertificateModel maps the schema attributes.
type TeoApplyFreeCertificateModel struct {
	ZoneId             types.String `tfsdk:"zone_id"`
	Domain             types.String `tfsdk:"domain"`
	VerificationMethod types.String `tfsdk:"verification_method"`
}

func (a *TeoApplyFreeCertificate) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_apply_free_certificate"
}

func (a *TeoApplyFreeCertificate) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to apply for a TEO (EdgeOne) free certificate for a domain via the ApplyFreeCertificate API. " +
			"After the apply call succeeds, the action polls CheckFreeCertificateVerification until the certificate application succeeds, fails, or times out. " +
			"This is a one-time operation; no cloud-side state is persisted after the action completes.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Zone ID.",
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "The target domain for which to apply for a free certificate.",
			},
			"verification_method": schema.StringAttribute{
				Required:    true,
				Description: "The verification method for the free certificate application. Valid values: `http_challenge` (HTTP file verification), `dns_challenge` (DNS delegation verification).",
			},
		},
	}
}

// Invoke is called to run the logic of the action.
func (a *TeoApplyFreeCertificate) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TeoApplyFreeCertificateModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ZoneId.IsNull() || data.ZoneId.IsUnknown() || data.ZoneId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing zone_id",
			"zone_id is required and cannot be empty.",
		)
		return
	}

	if data.Domain.IsNull() || data.Domain.IsUnknown() || data.Domain.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing domain",
			"domain is required and cannot be empty.",
		)
		return
	}

	if data.VerificationMethod.IsNull() || data.VerificationMethod.IsUnknown() || data.VerificationMethod.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing verification_method",
			"verification_method is required and cannot be empty.",
		)
		return
	}

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	zoneId := data.ZoneId.ValueString()
	domain := data.Domain.ValueString()
	verificationMethod := data.VerificationMethod.ValueString()
	service := NewTeoService(a.Client())

	dnsVerification, fileVerification, err := service.ApplyFreeCertificate(ctx, zoneId, domain, verificationMethod)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("applying TEO free certificate for domain (%s)", domain),
			err.Error(),
		)
		return
	}

	if dnsVerification != nil {
		subdomain := ""
		if dnsVerification.Subdomain != nil {
			subdomain = *dnsVerification.Subdomain
		}
		recordType := ""
		if dnsVerification.RecordType != nil {
			recordType = *dnsVerification.RecordType
		}
		recordValue := ""
		if dnsVerification.RecordValue != nil {
			recordValue = *dnsVerification.RecordValue
		}
		log.Printf("[DEBUG] teo apply_free_certificate dns verification: subdomain=%s, recordType=%s, recordValue=%s\n",
			subdomain, recordType, recordValue)
	}
	if fileVerification != nil {
		path := ""
		if fileVerification.Path != nil {
			path = *fileVerification.Path
		}
		content := ""
		if fileVerification.Content != nil {
			content = *fileVerification.Content
		}
		log.Printf("[DEBUG] teo apply_free_certificate file verification: path=%s, content=%s\n",
			path, content)
	}

	if err := service.CheckFreeCertificateVerification(ctx, zoneId, domain, 3*tccommon.ReadRetryTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("checking TEO free certificate verification for domain (%s)", domain),
			err.Error(),
		)
		return
	}
}
