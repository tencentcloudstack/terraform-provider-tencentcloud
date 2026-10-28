package ssm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSsmServiceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmServiceOperationCreate,
		Read:   resourceTencentCloudSsmServiceOperationRead,
		Delete: resourceTencentCloudSsmServiceOperationDelete,
		Schema: map[string]*schema.Schema{
			"pay_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Billing mode for opening SSM. Only `PostPay` (pay-as-you-go) is supported; the value is case-sensitive.",
			},
		},
	}
}

func resourceTencentCloudSsmServiceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_service_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		payMode = d.Get("pay_mode").(string)
	)
	// SSM CreateService requires an idempotent client token (1-64 visible ASCII
	// chars). Generate it internally so it is not exposed as a schema field.
	clientToken := helper.BuildToken()
	ssmService := SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	createResp, err := ssmService.CreateService(ctx, payMode, clientToken)
	if err != nil {
		log.Printf("[CRITAL]%s ssm_service_operation create failed, reason:%+v", logId, err)
		return err
	}

	// Log the CreateService response fields for observability.
	if createResp != nil && createResp.Response != nil {
		resourceId := ""
		if createResp.Response.ResourceId != nil {
			resourceId = *createResp.Response.ResourceId
		}
		orderId := ""
		if createResp.Response.OrderId != nil {
			orderId = *createResp.Response.OrderId
		}
		log.Printf("[DEBUG] ssm_service_operation CreateService invoked, resource_id=%s, order_id=%s", resourceId, orderId)
	}

	// CreateService is asynchronous; set the id before polling so the state
	// retains it even if the polling ultimately fails, allowing terraform
	// destroy to clean up.
	d.SetId(helper.BuildToken())

	// Poll GetServiceStatus until the service is enabled, bounded by
	// ReadRetryTimeout.
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		status, e := ssmService.GetServiceStatus(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if status == nil || status.ServiceEnabled == nil {
			return resource.NonRetryableError(fmt.Errorf("ssm_service_operation status ServiceEnabled is nil"))
		}
		if *status.ServiceEnabled {
			invalidType := int64(0)
			if status.InvalidType != nil {
				invalidType = *status.InvalidType
			}
			log.Printf("[DEBUG] ssm_service_operation enabled confirmed, service_enabled=%v, invalid_type=%d",
				*status.ServiceEnabled, invalidType)
			return nil
		}
		log.Printf("[DEBUG] ssm_service_operation still not enabled, keep polling")
		return resource.RetryableError(fmt.Errorf("ssm_service_operation not enabled yet"))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s polling ssm_service_operation status failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudSsmServiceOperationRead(d, meta)
}

func resourceTencentCloudSsmServiceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_service_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSsmServiceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_service_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
