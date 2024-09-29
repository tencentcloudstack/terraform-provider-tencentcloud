package crs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisStartupInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisStartupInstanceOperationCreate,
		Read:   resourceTencentCloudRedisStartupInstanceOperationRead,
		Delete: resourceTencentCloudRedisStartupInstanceOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},
		},
	}
}

func resourceTencentCloudRedisStartupInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_startup_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = redis.NewStartupInstanceRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().StartupInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis startupInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, err := service.DescribeRedisInstanceById(ctx, d.Id())
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		if instance == nil {
			return resource.RetryableError(fmt.Errorf("redis instance is nil, retry..."))
		}
		if *instance.Status == REDIS_STATUS_ONLINE {
			return nil
		}
		log.Printf("[DEBUG]%s api[%s] redis instance status is %v[%s], need 2[online], retry...", logId, request.GetAction(), *instance.Status, REDIS_STATUS[*instance.Status])
		return resource.RetryableError(fmt.Errorf("redis instance is %v, need 2, retry...", *instance.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis startup instance fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisStartupInstanceOperationRead(d, meta)
}

func resourceTencentCloudRedisStartupInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_startup_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisStartupInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_startup_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
