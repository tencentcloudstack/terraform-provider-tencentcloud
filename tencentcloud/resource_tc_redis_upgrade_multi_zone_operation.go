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

func resourceTencentCloudRedisUpgradeMultiZoneOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisUpgradeMultiZoneOperationCreate,
		Read:   resourceTencentCloudRedisUpgradeMultiZoneOperationRead,
		Delete: resourceTencentCloudRedisUpgradeMultiZoneOperationDelete,
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

			"upgrade_proxy_and_redis_server": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "After you upgrade Multi-AZ, whether the nearby access feature is supported.true: Supports nearby access.The upgrade process, which requires upgrading both the proxy version and the Redis kernel minor version, involves data migration and can take several hours.false: No need to support nearby access.Upgrading Multi-AZ only involves managing metadata migration, with no service impact, and the upgrade process typically completes within 3 minutes.",
			},
		},
	}
}

func resourceTencentCloudRedisUpgradeMultiZoneOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade_multi_zone_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = redis.NewUpgradeVersionToMultiAvailabilityZonesRequest()
		response   = redis.NewUpgradeVersionToMultiAvailabilityZonesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("upgrade_proxy_and_redis_server"); v != nil {
		request.UpgradeProxyAndRedisServer = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().UpgradeVersionToMultiAvailabilityZones(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis upgradeMultiZoneOperation failed, reason:%+v", logId, err)
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
			return resource.RetryableError(fmt.Errorf("upgrade multi zone is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis upgrade multi zone fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisUpgradeMultiZoneOperationRead(d, meta)
}

func resourceTencentCloudRedisUpgradeMultiZoneOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade_multi_zone_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisUpgradeMultiZoneOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade_multi_zone_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
