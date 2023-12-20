package crs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisRenewInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisRenewInstanceOperationCreate,
		Read:   resourceTencentCloudRedisRenewInstanceOperationRead,
		Delete: resourceTencentCloudRedisRenewInstanceOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"period": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Purchase duration, in months.",
			},

			"modify_pay_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Identifies whether the billing model is modified:The current instance billing mode is pay-as-you-go, which is prepaid and renewed.The billing mode of the current instance is subscription and you can not set this parameter.",
			},
		},
	}
}

func resourceTencentCloudRedisRenewInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_renew_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = redis.NewRenewInstanceRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("modify_pay_mode"); ok {
		request.ModifyPayMode = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().RenewInstance(request)
		if e != nil {
			if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "FailedOperation.PayFailed" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis renewInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	_, _, _, err = service.CheckRedisOnlineOk(ctx, instanceId, 20*tccommon.ReadRetryTimeout)
	if err != nil {
		log.Printf("[CRITAL]%s redis upgradeVersionOperation fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisRenewInstanceOperationRead(d, meta)
}

func resourceTencentCloudRedisRenewInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_renew_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisRenewInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_renew_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
