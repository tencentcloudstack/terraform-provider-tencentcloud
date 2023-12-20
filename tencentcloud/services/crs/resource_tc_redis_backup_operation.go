package crs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisBackupOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisBackupOperationCreate,
		Read:   resourceTencentCloudRedisBackupOperationRead,
		Delete: resourceTencentCloudRedisBackupOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"remark": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Notes information for the backup.",
			},

			"storage_days": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Number of days to store.0 specifies the default retention time.",
			},
		},
	}
}

func resourceTencentCloudRedisBackupOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = redis.NewManualBackupInstanceRequest()
		response   = redis.NewManualBackupInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, _ := d.GetOk("storage_days"); v != nil {
		request.StorageDays = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ManualBackupInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis backupOperation failed, reason:%+v", logId, err)
		return err
	}

	taskId := *response.Response.TaskId
	d.SetId(instanceId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if taskId > 0 {
		err := resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ok, err := service.DescribeTaskInfo(ctx, instanceId, taskId)
			if err != nil {
				if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if ok {
				return nil
			} else {
				return resource.RetryableError(fmt.Errorf("redis backupOperation is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis backupOperation fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudRedisBackupOperationRead(d, meta)
}

func resourceTencentCloudRedisBackupOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisBackupOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
