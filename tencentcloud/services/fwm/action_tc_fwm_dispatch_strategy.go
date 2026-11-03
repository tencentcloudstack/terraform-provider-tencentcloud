package fwm

import (
	"context"
	"fmt"
	"log"
	"time"

	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ action.Action is intentionally omitted: the factory function
// NewFwmDispatchStrategy returns action.Action, which already performs the
// implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &FwmDispatchStrategy{}

// NewFwmDispatchStrategy is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewFwmDispatchStrategy() action.Action {
	return &FwmDispatchStrategy{}
}

// defaultFwmDispatchStrategyInvokeTimeout is the built-in maximum duration to
// wait for the asynchronous dispatch task to finish. It tracks the provider's
// global write timeout so an env override of the write timeout is honored.
// Practitioners can override it through the `timeouts` block, e.g.
// `timeouts { invoke = "30m" }`.
var defaultFwmDispatchStrategyInvokeTimeout = tccommon.WriteRetryTimeout

// fwmDispatchStrategyPollInterval is the interval between two
// DescribeStrategyDispatchStatus polls. It is a package-level var (not a const)
// so unit tests can shorten it.
var fwmDispatchStrategyPollInterval = 5 * time.Second

// allowedFwmDispatchProducts lists the valid Product values accepted by the
// DispatchStrategy API. `all` is intentionally excluded: dispatch operates on a
// single product line at a time.
var allowedFwmDispatchProducts = map[string]bool{
	"enterprise_sg":  true,
	"cfw_vpc_acl":  true,
	"cfw_nat_acl": true,
	"cfw_edge_acl":   true,
}

// DispatchStatus exports the DescribeStrategyDispatchStatus.Response.Status
// enum for tests: DispatchStatusNoChange (0), DispatchStatusDispatching (1),
// DispatchStatusSucceeded (2), DispatchStatusFailed (3).
const (
	DispatchStatusNoChange    = fwmDispatchStatusNoChange
	DispatchStatusDispatching = fwmDispatchStatusDispatching
	DispatchStatusSucceeded   = fwmDispatchStatusSucceeded
	DispatchStatusFailed   = fwmDispatchStatusFailed
)

// FwmDispatchStrategyPollIntervalForTest exposes the poll interval for tests
// to snapshot/restore. It is only meant to be used from the action's test file.
func FwmDispatchStrategyPollIntervalForTest() time.Duration {
	return fwmDispatchStrategyPollInterval
}

// SetFwmDispatchStrategyPollIntervalForTest overrides the poll interval for
// tests and always returns true. It is only meant to be used from the
// action's test file.
func SetFwmDispatchStrategyPollIntervalForTest(d time.Duration) bool {
	fwmDispatchStrategyPollInterval = d
	return true
}

// FwmDispatchStrategy implements action.Action for
// tencentcloud_fwm_dispatch_strategy.
type FwmDispatchStrategy struct {
	fw.ActionWithConfigure
}

// FwmDispatchStrategyModel maps the schema attributes.
//
// Action schemas cannot be Computed or Sensitive: the only marks available are
// Required / Optional / WriteOnly (WriteOnly needs framework v1.17.0+).
// `product` and `uin` are required; `sub_account_uin` is optional. The cloud
// API only carries `product` (status is pinned to start=1 server-side); the
// operator UIN / SubAccountUin are injected by the API gateway, so they are
// accepted on the Terraform side for clarity/auditability but not forwarded.
type FwmDispatchStrategyModel struct {
	Product     types.String `tfsdk:"product"`
	Uin         types.String `tfsdk:"uin"`
	SubAccountUin types.String `tfsdk:"sub_account_uin"`

	// Timeouts is the value type defined by the official
	// terraform-plugin-framework-timeouts module. Only the `invoke` stage
	// exists because an action has a single lifecycle stage.
	Timeouts actiontimeouts.Value `tfsdk:"timeouts"`
}

func (a *FwmDispatchStrategy) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_dispatch_strategy"
}

func (a *FwmDispatchStrategy) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to dispatch FWM (Firewall Manager) pending strategies for a product " +
			"via the DispatchStrategy API (下发规则). This is a one-time invocation; it does not create, update or " +
			"destroy any Terraform-managed resource, and it exposes no attributes referenceable from other resources. " +
			"The dispatch dispatches every pending (new / pending-dispatch) strategy of the current account x product; " +
			"it cannot target a single rule. The cloud API has no idempotency key and does not block a running task, " +
			"so the provider first checks the latest task status and, if a dispatch is already in progress, reuses it " +
		"instead of triggering another one; dispatches of the same product should still be serialized. " +
			"Create success does NOT mean dispatch success: the action polls DescribeStrategyDispatchStatus to a " +
			"terminal state (succeeded / failed) before returning. Re-running the dispatch should be done by tainting " +
			"the action or changing a trigger field; destroy does not roll back the dispatched rules and the far-side " +
			"effect persists.",
		Attributes: map[string]schema.Attribute{
			"product": schema.StringAttribute{
				Required:    true,
				Description: "Product line to dispatch. Valid values: `enterprise_sg`, `cfw_vpc_acl`, `cfw_nat_acl`, `cfw_edge_acl`.",
			},
			"uin": schema.StringAttribute{
				Required:    true,
				Description: "Operator UIN. It is injected by the API gateway on the cloud side and not carried in the request body, but recorded here for clarity.",
			},
			"sub_account_uin": schema.StringAttribute{
				Optional:    true,
				Description: "Sub-account UIN of the operator. Optional; injected by the API gateway on the cloud side.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": fwmDispatchStrategyTimeoutsBlock(ctx),
		},
	}
}

// fwmDispatchStrategyTimeoutsBlock builds the `timeouts` block from the
// official terraform-plugin-framework-timeouts module, whose single optional
// `invoke` attribute is validated as a time.Duration at plan time.
func fwmDispatchStrategyTimeoutsBlock(ctx context.Context) schema.Block {
	block := actiontimeouts.BlockWithOpts(ctx, actiontimeouts.Opts{
		InvokeDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `5m`. " +
			"It limits how long the action waits for the dispatch task to reach a terminal state. Default follows the provider write timeout.",
	})

	if nested, ok := block.(schema.SingleNestedBlock); ok {
		nested.Description = "The timeouts block allows you to specify the timeout for the invoke operation."
		return nested
	}

	return block
}

// Invoke is called to run the logic of the action.
func (a *FwmDispatchStrategy) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data FwmDispatchStrategyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Product.IsNull() || data.Product.IsUnknown() || data.Product.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing product",
			"product is required and cannot be empty.",
		)
		return
	}

	if data.Uin.IsNull() || data.Uin.IsUnknown() || data.Uin.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing uin",
			"uin is required and cannot be empty.",
		)
		return
	}

	if !allowedFwmDispatchProducts[data.Product.ValueString()] {
		resp.Diagnostics.AddError(
			"Invalid product",
			fmt.Sprintf("product %q is not one of the allowed values (enterprise_sg, cfw_vpc_acl, cfw_nat_acl, cfw_edge_acl).", data.Product.ValueString()),
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

	// An unset/null timeouts block falls back to the built-in default, and a
	// malformed duration is reported as an error diagnostic. The schema
	// validator normally rejects malformed values at plan time; this is the
	// runtime safety net.
	timeout, timeoutDiags := data.Timeouts.Invoke(ctx, defaultFwmDispatchStrategyInvokeTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	progress := func(msg string) {
		if resp.SendProgress != nil {
			resp.SendProgress(action.InvokeProgressEvent{Message: msg})
		}
	}

	product := data.Product.ValueString()
	service := NewFwmService(a.Client())

	// Dedup guard: the API has no ClientToken and does not block a running
	// task, so before dispatching we look up the latest task. If one is
	// already in progress (status 0=no change / 1=dispatching) we reuse it
	// instead of inserting a new dispatch task.
	progress(fmt.Sprintf("checking latest dispatch status for product %s", product))
	latest, err := service.DescribeFwmStrategyDispatchStatus(ctx, product)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("checking FWM dispatch status for product (%s)", product),
			err.Error(),
		)
		return
	}

	if latest != nil && latest.Status != nil &&
		(*latest.Status == fwmDispatchStatusNoChange || *latest.Status == fwmDispatchStatusDispatching) {
		log.Printf("[DEBUG] reuse in-progress dispatch task for product %s, current status %d", product, *latest.Status)
		progress(fmt.Sprintf("an in-progress dispatch task exists for product %s, waiting for it to finish", product))
	} else {
		progress(fmt.Sprintf("dispatching FWM strategy for product %s", product))
		if err := service.DispatchFwmStrategy(ctx, product); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("dispatching FWM strategy for product (%s)", product),
				err.Error(),
			)
			return
		}
	}

	// Poll to a terminal state: success (2) returns nil, failure (3) returns an
	// error carrying ErrorMsg, otherwise keep polling until the budget is up.
	last, err := service.WaitFwmStrategyDispatchComplete(ctx, product, timeout, fwmDispatchStrategyPollInterval)
	if err != nil {
		snapshot := "no status response"
		if last != nil {
			statusStr := "nil"
			if last.Status != nil {
				statusStr = fmt.Sprintf("%d", *last.Status)
			}
			progressStr := "nil"
			if last.Progress != nil {
				progressStr = fmt.Sprintf("%v", *last.Progress)
			}
			snapshot = fmt.Sprintf("status=%s, progress=%s", statusStr, progressStr)
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("waiting for FWM dispatch strategy to complete for product (%s)", product),
			fmt.Sprintf("%s; last observed: %s", err.Error(), snapshot),
		)
		return
	}

	progress(fmt.Sprintf("FWM strategy dispatch completed for product %s", product))
	// No state to write: actions have no terraform state and no output attributes.
}