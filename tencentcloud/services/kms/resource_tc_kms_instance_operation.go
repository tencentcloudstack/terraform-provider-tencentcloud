package kms

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKmsInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsInstanceOperationCreate,
		Read:   resourceTencentCloudKmsInstanceOperationRead,
		Delete: resourceTencentCloudKmsInstanceOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_charge_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance charge type. Valid value: `PREPAID`.",
			},
			"product_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Product version. Valid value: `professional` (professional edition).",
			},
			"instance_charge_prepaid": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Prepaid instance charge info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Purchase duration in months. Valid range: `1` to `100`.",
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Auto-renewal flag. `0`: not set, `1`: auto-renew, `2`: no auto-renew.",
						},
					},
				},
			},
			"client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reentrant client id for idempotency.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance id.",
			},
			"deal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Billing order no.",
			},
			"bill_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Billing frozen serial no.",
			},
		},
	}
}

func resourceTencentCloudKmsInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = kms.NewCreateKmsInstanceRequest()
	)

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product_version"); ok {
		request.ProductVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_prepaid"); ok {
		prepaidList := v.([]interface{})
		if len(prepaidList) > 0 {
			prepaidMap := prepaidList[0].(map[string]interface{})
			prepaid := kms.InstanceChargePrepaid{}
			if period, ok := prepaidMap["period"]; ok {
				prepaid.Period = helper.IntUint64(period.(int))
			}
			if renewFlag, ok := prepaidMap["renew_flag"]; ok {
				prepaid.RenewFlag = helper.IntUint64(renewFlag.(int))
			}
			request.InstanceChargePrepaid = &prepaid
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagMap := item.(map[string]interface{})
			tag := kms.Tag{}
			if tagKey, ok := tagMap["tag_key"]; ok {
				tag.TagKey = helper.String(tagKey.(string))
			}
			if tagValue, ok := tagMap["tag_value"]; ok {
				tag.TagValue = helper.String(tagValue.(string))
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	var response *kms.CreateKmsInstanceResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseKmsClient().CreateKmsInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			response = result
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create kms kms_instance_operation failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil || response.Response.InstanceId == nil || *response.Response.InstanceId == "" {
		log.Printf("[CRITAL]%s create kms kms_instance_operation return empty, logId=%s, d.Id()=%s", logId, logId, d.Id())
		return fmt.Errorf("create kms kms_instance_operation return empty")
	}

	instanceId := *response.Response.InstanceId
	d.SetId(instanceId)

	_ = d.Set("instance_id", instanceId)
	if response.Response.DealName != nil {
		_ = d.Set("deal_name", *response.Response.DealName)
	}
	if response.Response.BillId != nil {
		_ = d.Set("bill_id", *response.Response.BillId)
	}

	return resourceTencentCloudKmsInstanceOperationRead(d, meta)
}

func resourceTencentCloudKmsInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudKmsInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
