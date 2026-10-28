package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var (
	_ action.Action              = &TeoIncreasePlanQuotaAction{}
	_ action.ActionWithConfigure = &TeoIncreasePlanQuotaAction{}
)

// NewTeoIncreasePlanQuotaAction is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewTeoIncreasePlanQuotaAction() action.Action {
	return &TeoIncreasePlanQuotaAction{}
}

// TeoIncreasePlanQuotaAction implements action.Action for
// tencentcloud_teo_increase_plan_quota.
type TeoIncreasePlanQuotaAction struct {
	client *connectivity.TencentCloudClient
}

// TeoIncreasePlanQuotaActionModel maps the schema attributes.
type TeoIncreasePlanQuotaActionModel struct {
	PlanId      types.String `tfsdk:"plan_id"`
	QuotaType   types.String `tfsdk:"quota_type"`
	QuotaNumber types.Int64  `tfsdk:"quota_number"`
	DealName    types.String `tfsdk:"deal_name"`
}

func (a *TeoIncreasePlanQuotaAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_increase_plan_quota"
}

func (a *TeoIncreasePlanQuotaAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Increase plan quota for TEO (EdgeOne) plans. " +
			"This action calls the TEO `IncreasePlanQuota` API to purchase additional plan quotas " +
			"(e.g., site count, precise access control rules, rate limiting rules). " +
			"It is a one-time operation that does not persist any state.",
		Attributes: map[string]schema.Attribute{
			"plan_id": schema.StringAttribute{
				Required:    true,
				Description: "Plan ID, in the format of edgeone-xxxxxxxx.",
			},
			"quota_type": schema.StringAttribute{
				Required:    true,
				Description: "Quota type to increase. Valid values: `site`, `precise_access_control_rule`, `rate_limiting_rule`.",
			},
			"quota_number": schema.Int64Attribute{
				Required:    true,
				Description: "Number of quotas to increase. Maximum 100 per request.",
			},
			"deal_name": schema.StringAttribute{
				Optional:    true,
				Description: "Order number returned by the API.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure (which itself reuses the SDKv2-built
// *connectivity.TencentCloudClient).
func (a *TeoIncreasePlanQuotaAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
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

	a.client = meta.Client
}

// SetClientForTest injects a *connectivity.TencentCloudClient directly into
// the action, bypassing the Configure flow. It is intended ONLY for unit
// tests that need to pre-set the client before invoking the action.
func (a *TeoIncreasePlanQuotaAction) SetClientForTest(client *connectivity.TencentCloudClient) {
	a.client = client
}

// Invoke is called to run the logic of the action. It reads the config,
// calls the TEO IncreasePlanQuota API, and reports the deal_name back to the
// user via a warning diagnostic.
func (a *TeoIncreasePlanQuotaAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TeoIncreasePlanQuotaActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planId := data.PlanId.ValueString()
	quotaType := data.QuotaType.ValueString()
	quotaNumber := data.QuotaNumber.ValueInt64()

	if planId == "" {
		resp.Diagnostics.AddError(
			"Missing plan_id",
			"plan_id is required and cannot be empty.",
		)
		return
	}
	if quotaType == "" {
		resp.Diagnostics.AddError(
			"Missing quota_type",
			"quota_type is required and cannot be empty.",
		)
		return
	}

	if a.client == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The teo_increase_plan_quota client is nil. This usually means the provider was not configured correctly.",
		)
		return
	}

	teoClient := a.client.UseTeoV20220901Client()

	request := teov20220901.NewIncreasePlanQuotaRequest()
	request.PlanId = helper.String(planId)
	request.QuotaType = helper.String(quotaType)
	request.QuotaNumber = helper.Int64(quotaNumber)

	var dealName string
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := teoClient.IncreasePlanQuotaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil || result.Response.DealName == nil {
			return resource.NonRetryableError(fmt.Errorf("IncreasePlanQuota API returned empty response"))
		}
		dealName = *result.Response.DealName
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error invoking teo_increase_plan_quota action",
			fmt.Sprintf("IncreasePlanQuota failed for plan_id %q: %v", planId, err),
		)
		return
	}

	log.Printf("[INFO] teo_increase_plan_quota action success, plan_id: %s, deal_name: %s", planId, dealName)

	resp.Diagnostics.AddWarning(
		"teo_increase_plan_quota action completed",
		fmt.Sprintf("deal_name: %s", dealName),
	)
}
