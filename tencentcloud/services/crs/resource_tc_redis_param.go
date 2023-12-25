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

func ResourceTencentCloudRedisParam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisParamCreate,
		Read:   resourceTencentCloudRedisParamRead,
		Update: resourceTencentCloudRedisParamUpdate,
		Delete: resourceTencentCloudRedisParamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"instance_params": {
				Required:    true,
				Type:        schema.TypeMap,
				Description: "A list of parameters modified by the instance.",
			},
		},
	}
}

func resourceTencentCloudRedisParamCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_param.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisParamUpdate(d, meta)
}

func resourceTencentCloudRedisParamRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_param.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	param, err := service.DescribeRedisParamById(ctx, instanceId)
	if err != nil {
		return err
	}

	if len(param) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisParam` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	instanceParamsMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_params"); ok {
		for k := range v.(map[string]interface{}) {
			instanceParamsMap[k] = param[k]
		}
	} else {
		instanceParamsMap = param
	}
	_ = d.Set("instance_params", instanceParamsMap)

	return nil
}

func resourceTencentCloudRedisParamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_param.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := redis.NewModifyInstanceParamsRequest()
	response := redis.NewModifyInstanceParamsResponse()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("instance_params"); ok {
		service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		param, err := service.DescribeRedisParamById(ctx, instanceId)
		if err != nil && len(param) == 0 {
			return fmt.Errorf("[ERROR] resource `RedisParam` [%s] not found, please check if it has been deleted.\n", d.Id())
		}
		for k, v := range v.(map[string]interface{}) {
			if value, ok := param[k]; ok {
				if value != v {
					instanceParam := redis.InstanceParam{}
					instanceParam.Key = helper.String(k)
					instanceParam.Value = helper.String(v.(string))
					request.InstanceParams = append(request.InstanceParams, &instanceParam)
				}
			} else {
				return fmt.Errorf("[ERROR] The parameter name [%v] does not exist, please check the parameter name.\n", k)
			}
		}
	}

	if len(request.InstanceParams) == 0 {
		return resourceTencentCloudRedisParamRead(d, meta)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyInstanceParams(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis param failed, reason:%+v", logId, err)
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
			return resource.RetryableError(fmt.Errorf("change param is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change param fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisParamRead(d, meta)
}

func resourceTencentCloudRedisParamDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_param.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
