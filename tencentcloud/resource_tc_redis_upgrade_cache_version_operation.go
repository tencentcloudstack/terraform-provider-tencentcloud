package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisUpgradeCacheVersionOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisUpgradeCacheVersionOperationCreate,
		Read:   resourceTencentCloudRedisUpgradeCacheVersionOperationRead,
		Delete: resourceTencentCloudRedisUpgradeCacheVersionOperationDelete,
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

			"current_redis_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Current redis version.",
			},

			"upgrade_redis_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Upgradeable redis version.",
			},

			"instance_type_upgrade_now": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Switch mode:1 - Upgrade now0 - Maintenance window upgrade.",
			},
		},
	}
}

func resourceTencentCloudRedisUpgradeCacheVersionOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade_cache_version_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = redis.NewUpgradeSmallVersionRequest()
		response   = redis.NewUpgradeSmallVersionResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("current_redis_version"); ok {
		request.CurrentRedisVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upgrade_redis_version"); ok {
		request.UpgradeRedisVersion = helper.String(v.(string))
	}

	if v, _ := d.GetOk("instance_type_upgrade_now"); v != nil {
		request.InstanceTypeUpgradeNow = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().UpgradeSmallVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis upgradeCacheVersion failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	taskId := *response.Response.FlowId
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
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
			return resource.RetryableError(fmt.Errorf("upgrade cache version is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis upgrade cache version fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisUpgradeCacheVersionOperationRead(d, meta)
}

func resourceTencentCloudRedisUpgradeCacheVersionOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade_cache_version_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisUpgradeCacheVersionOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade_cache_version_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
