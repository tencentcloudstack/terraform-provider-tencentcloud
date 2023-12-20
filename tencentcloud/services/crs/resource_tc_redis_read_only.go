package crs

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisReadOnly() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisReadOnlyCreate,
		Read:   resourceTencentCloudRedisReadOnlyRead,
		Update: resourceTencentCloudRedisReadOnlyUpdate,
		Delete: resourceTencentCloudRedisReadOnlyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"input_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance input mode: `0`: read-write; `1`: read-only.",
			},
		},
	}
}

func resourceTencentCloudRedisReadOnlyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_read_only.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisReadOnlyUpdate(d, meta)
}

func resourceTencentCloudRedisReadOnlyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_read_only.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	setInputMode, err := service.DescribeRedisInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if setInputMode == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisSetInputMode` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if setInputMode.InstanceId != nil {
		_ = d.Set("instance_id", setInputMode.InstanceId)
	}

	if setInputMode.ReadOnly != nil {
		_ = d.Set("input_mode", strconv.FormatInt(*setInputMode.ReadOnly, 10))
	}

	return nil
}

func resourceTencentCloudRedisReadOnlyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_read_only.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := redis.NewModifyInstanceReadOnlyRequest()
	response := redis.NewModifyInstanceReadOnlyResponse()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("input_mode"); ok {
		request.InputMode = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyInstanceReadOnly(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis setInputMode failed, reason:%+v", logId, err)
		return err
	}

	taskId := *response.Response.TaskId
	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
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
			return resource.RetryableError(fmt.Errorf("change inputMode is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change inputMode fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisReadOnlyRead(d, meta)
}

func resourceTencentCloudRedisReadOnlyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_read_only.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
