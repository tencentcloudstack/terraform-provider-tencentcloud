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

func ResourceTencentCloudRedisUpgradeProxyVersionOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisUpgradeProxyVersionOperationCreate,
		Read:   resourceTencentCloudRedisUpgradeProxyVersionOperationRead,
		Delete: resourceTencentCloudRedisUpgradeProxyVersionOperationDelete,
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

			"current_proxy_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Current proxy version.",
			},

			"upgrade_proxy_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Upgradeable redis proxy version.",
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

func resourceTencentCloudRedisUpgradeProxyVersionOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_upgrade_proxy_version_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = redis.NewUpgradeProxyVersionRequest()
		response   = redis.NewUpgradeProxyVersionResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("current_proxy_version"); ok {
		request.CurrentProxyVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upgrade_proxy_version"); ok {
		request.UpgradeProxyVersion = helper.String(v.(string))
	}

	if v, _ := d.GetOk("instance_type_upgrade_now"); v != nil {
		request.InstanceTypeUpgradeNow = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().UpgradeProxyVersion(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis upgradeProxyVersionOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	taskId := *response.Response.FlowId
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
			return resource.RetryableError(fmt.Errorf("upgrade proxy version is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis upgrade proxy version fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisUpgradeProxyVersionOperationRead(d, meta)
}

func resourceTencentCloudRedisUpgradeProxyVersionOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_upgrade_proxy_version_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisUpgradeProxyVersionOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_upgrade_proxy_version_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
