package ssm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ action.ActionWithConfigure = &SsmServiceAction{}

// NewSsmServiceAction is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewSsmServiceAction() action.Action {
	return &SsmServiceAction{}
}

// SsmServiceAction implements action.Action for
// tencentcloud_ssm_service.
type SsmServiceAction struct {
	fw.ActionWithConfigure
}

// SsmServiceActionModel maps the schema attributes.
type SsmServiceActionModel struct {
	PayMode     types.String `tfsdk:"pay_mode"`
	ClientToken types.String `tfsdk:"client_token"`
}

func (a *SsmServiceAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_ssm_service"
}

func (a *SsmServiceAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to activate (open) the SSM (Secrets Manager) service (PostPay / " +
			"pay-as-you-go) via the SSM CreateService API. Because CreateService is asynchronous, the action " +
			"polls GetServiceStatus until the service is enabled. This is a one-time operation; no cloud-side " +
			"state is persisted and no output attributes are exposed.",
		Attributes: map[string]schema.Attribute{
			"pay_mode": schema.StringAttribute{
				Required:    true,
				Description: "Billing mode for opening SSM. Only `PostPay` (pay-as-you-go) is supported; the value is case-sensitive.",
			},
			"client_token": schema.StringAttribute{
				Required:    true,
				Description: "Client-generated idempotent create token. 1-64 visible ASCII characters (no spaces, Chinese, or control characters); recommended `ssm-<UUID>`.",
			},
		},
	}
}

// Invoke is called to run the logic of the action.
func (a *SsmServiceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SsmServiceActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payMode := ""
	if !data.PayMode.IsNull() && !data.PayMode.IsUnknown() {
		payMode = data.PayMode.ValueString()
	}
	if payMode == "" {
		resp.Diagnostics.AddError(
			"Missing pay_mode",
			"pay_mode is required and cannot be empty. Only 'PostPay' is supported.",
		)
		return
	}

	clientToken := ""
	if !data.ClientToken.IsNull() && !data.ClientToken.IsUnknown() {
		clientToken = data.ClientToken.ValueString()
	}
	if clientToken == "" {
		resp.Diagnostics.AddError(
			"Missing client_token",
			"client_token is required and cannot be empty. Use a unique token such as 'ssm-<UUID>'.",
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

	service := NewSsmService(a.Client())
	createResp, err := service.CreateService(ctx, payMode, clientToken)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("opening SSM service"),
			err.Error(),
		)
		return
	}

	// Framework actions cannot expose output attributes, so the CreateService
	// response fields are logged for observability only.
	if createResp != nil && createResp.Response != nil {
		resourceId := ""
		if createResp.Response.ResourceId != nil {
			resourceId = *createResp.Response.ResourceId
		}
		orderId := ""
		if createResp.Response.OrderId != nil {
			orderId = *createResp.Response.OrderId
		}
		log.Printf("[DEBUG] ssm service CreateService invoked, resource_id=%s, order_id=%s", resourceId, orderId)
	}

	// CreateService is asynchronous — poll GetServiceStatus until the service
	// is enabled, bounded by ReadRetryTimeout.
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		status, e := service.GetServiceStatus(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if status == nil || status.ServiceEnabled == nil {
			return resource.NonRetryableError(fmt.Errorf("SSM service status ServiceEnabled is nil"))
		}
		if *status.ServiceEnabled {
			invalidType := int64(0)
			if status.InvalidType != nil {
				invalidType = *status.InvalidType
			}
			log.Printf("[DEBUG] ssm service enabled confirmed, service_enabled=%v, invalid_type=%d",
				*status.ServiceEnabled, invalidType)
			return nil
		}
		log.Printf("[DEBUG] ssm service still not enabled, keep polling")
		return resource.RetryableError(fmt.Errorf("SSM service not enabled yet"))
	})

	if reqErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("polling SSM service status"),
			reqErr.Error(),
		)
		return
	}
}
