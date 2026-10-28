package kms

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
)

var _ action.ActionWithConfigure = &KmsInstanceOperation{}

// NewKmsInstanceOperation is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewKmsInstanceOperation() action.Action {
	return &KmsInstanceOperation{}
}

// KmsInstanceOperation implements action.Action for
// tencentcloud_kms_instance_operation.
type KmsInstanceOperation struct {
	fw.ActionWithConfigure
}

// KmsInstanceOperationModel maps the schema attributes.
type KmsInstanceOperationModel struct {
	InstanceChargeType    types.String `tfsdk:"instance_charge_type"`
	ProductVersion        types.String `tfsdk:"product_version"`
	ClientToken           types.String `tfsdk:"client_token"`
	InstanceChargePrepaid types.Object `tfsdk:"instance_charge_prepaid"`
	Tags                  types.List   `tfsdk:"tags"`
}

// KmsInstanceChargePrepaidModel maps the instance_charge_prepaid nested block.
type KmsInstanceChargePrepaidModel struct {
	Period    types.Int64 `tfsdk:"period"`
	RenewFlag types.Int64 `tfsdk:"renew_flag"`
}

// KmsTagModel maps the tags nested block.
type KmsTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

func (a *KmsInstanceOperation) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_kms_instance_operation"
}

func (a *KmsInstanceOperation) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to purchase and activate a prepaid KMS (Key Management Service) instance " +
			"via the CreateKmsInstance API. This is a one-time operation; no cloud-side state is persisted after the action completes.",
		Attributes: map[string]schema.Attribute{
			"instance_charge_type": schema.StringAttribute{
				Required:    true,
				Description: "Instance charge type. Valid value: `PREPAID`.",
			},
			"product_version": schema.StringAttribute{
				Optional:    true,
				Description: "Product version. Valid value: `professional` (professional edition).",
			},
			"client_token": schema.StringAttribute{
				Optional:    true,
				Description: "Reentrant client id for idempotency.",
			},
		},
		Blocks: map[string]schema.Block{
			"instance_charge_prepaid": schema.SingleNestedBlock{
				Description: "Prepaid instance charge info.",
				Attributes: map[string]schema.Attribute{
					"period": schema.Int64Attribute{
						Required:    true,
						Description: "Purchase duration in months. Valid range: `1` to `100`.",
					},
					"renew_flag": schema.Int64Attribute{
						Required:    true,
						Description: "Auto-renewal flag. `0`: not set, `1`: auto-renew, `2`: no auto-renew.",
					},
				},
			},
			"tags": schema.ListNestedBlock{
				Description: "Instance tags.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"tag_key": schema.StringAttribute{
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": schema.StringAttribute{
							Required:    true,
							Description: "tag value.",
						},
					},
				},
			},
		},
	}
}

// Invoke is called to run the logic of the action.
func (a *KmsInstanceOperation) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data KmsInstanceOperationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.InstanceChargeType.IsNull() || data.InstanceChargeType.IsUnknown() || data.InstanceChargeType.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing instance_charge_type",
			"instance_charge_type is required and cannot be empty.",
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

	logId := tccommon.GetLogId(ctx)
	request := kms.NewCreateKmsInstanceRequest()
	request.InstanceChargeType = helper.String(data.InstanceChargeType.ValueString())

	if !data.ProductVersion.IsNull() && !data.ProductVersion.IsUnknown() {
		request.ProductVersion = helper.String(data.ProductVersion.ValueString())
	}

	if !data.ClientToken.IsNull() && !data.ClientToken.IsUnknown() {
		request.ClientToken = helper.String(data.ClientToken.ValueString())
	}

	if !data.InstanceChargePrepaid.IsNull() && !data.InstanceChargePrepaid.IsUnknown() {
		var prepaid KmsInstanceChargePrepaidModel
		resp.Diagnostics.Append(data.InstanceChargePrepaid.As(ctx, &prepaid, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		sdkPrepaid := &kms.InstanceChargePrepaid{}
		if !prepaid.Period.IsNull() && !prepaid.Period.IsUnknown() {
			sdkPrepaid.Period = helper.Int64Uint64(prepaid.Period.ValueInt64())
		}
		if !prepaid.RenewFlag.IsNull() && !prepaid.RenewFlag.IsUnknown() {
			sdkPrepaid.RenewFlag = helper.Int64Uint64(prepaid.RenewFlag.ValueInt64())
		}
		request.InstanceChargePrepaid = sdkPrepaid
	}

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tagModels []KmsTagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tagModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, tag := range tagModels {
			t := &kms.Tag{}
			if !tag.TagKey.IsNull() && !tag.TagKey.IsUnknown() {
				t.TagKey = helper.String(tag.TagKey.ValueString())
			}
			if !tag.TagValue.IsNull() && !tag.TagValue.IsUnknown() {
				t.TagValue = helper.String(tag.TagValue.ValueString())
			}
			request.Tags = append(request.Tags, t)
		}
	}

	service := NewKmsService(a.Client())
	response, err := service.CreateKmsInstance(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("creating kms instance_operation"),
			err.Error(),
		)
		return
	}

	// Framework actions cannot expose output attributes, so the API response
	// fields are logged for observability only.
	if response != nil && response.Response != nil {
		instanceId := ""
		if response.Response.InstanceId != nil {
			instanceId = *response.Response.InstanceId
		}
		dealName := ""
		if response.Response.DealName != nil {
			dealName = *response.Response.DealName
		}
		billId := ""
		if response.Response.BillId != nil {
			billId = *response.Response.BillId
		}
		log.Printf("[DEBUG]%s kms instance_operation invoked, instance_id=%s, deal_name=%s, bill_id=%s", logId, instanceId, dealName, billId)
	}
}
