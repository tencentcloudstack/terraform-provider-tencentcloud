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

func ResourceTencentCloudRedisSwitchMaster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisSwitchMasterCreate,
		Read:   resourceTencentCloudRedisSwitchMasterRead,
		Update: resourceTencentCloudRedisSwitchMasterUpdate,
		Delete: resourceTencentCloudRedisSwitchMasterDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Replication group ID, required for multi-AZ instances.",
			},
		},
	}
}

func resourceTencentCloudRedisSwitchMasterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_switch_master.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisSwitchMasterUpdate(d, meta)
}

func resourceTencentCloudRedisSwitchMasterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_switch_master.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()
	paramMap := make(map[string]interface{})
	paramMap["InstanceId"] = &instanceId

	switchMaster, err := service.DescribeRedisInstanceZoneInfoByFilter(ctx, paramMap)
	if err != nil {
		return err
	}

	if switchMaster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisSwitchMaster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if len(switchMaster) > 1 {
		for _, v := range switchMaster {
			if *v.Role == "master" {
				_ = d.Set("group_id", v.GroupId)
				break
			}
		}
	}

	return nil
}

func resourceTencentCloudRedisSwitchMasterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_switch_master.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := redis.NewChangeReplicaToMasterRequest()
	response := redis.NewChangeReplicaToMasterResponse()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ChangeReplicaToMaster(request)
		if e != nil {
			if _, ok := e.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(e)
			} else {
				return resource.NonRetryableError(e)
			}
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis switchMaster failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	taskId := *response.Response.TaskId
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
			return resource.RetryableError(fmt.Errorf("update redis switchMaster is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s update redis switchMaster fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisSwitchMasterRead(d, meta)
}

func resourceTencentCloudRedisSwitchMasterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_switch_master.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
