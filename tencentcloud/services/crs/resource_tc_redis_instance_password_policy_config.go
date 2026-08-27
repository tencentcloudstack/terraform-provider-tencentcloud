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

func ResourceTencentCloudRedisInstancePasswordPolicyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisInstancePasswordPolicyConfigCreate,
		Read:   resourceTencentCloudRedisInstancePasswordPolicyConfigRead,
		Update: resourceTencentCloudRedisInstancePasswordPolicyConfigUpdate,
		Delete: resourceTencentCloudRedisInstancePasswordPolicyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of the Redis instance.",
			},

			"password_policy": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Password complexity policy configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Required:    true,
							Type:        schema.TypeBool,
							Description: "Whether to enable password complexity policy. true: enabled, false: disabled.",
						},

						"min_letter_count": {
							Optional:    true,
							Type:        schema.TypeInt,
							Description: "Minimum number of letters. Range: [1, 16]. Default: 1.",
						},

						"min_digit_count": {
							Optional:    true,
							Type:        schema.TypeInt,
							Description: "Minimum number of digits. Range: [1, 16]. Default: 1.",
						},

						"min_special_count": {
							Optional:    true,
							Type:        schema.TypeInt,
							Description: "Minimum number of special characters. Range: [1, 16]. Default: 1.",
						},

						"min_length": {
							Optional:    true,
							Type:        schema.TypeInt,
							Description: "Minimum total password length. Range: [8, 64]. Default: 8.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudRedisInstancePasswordPolicyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisInstancePasswordPolicyConfigUpdate(d, meta)
}

func resourceTencentCloudRedisInstancePasswordPolicyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	passwordPolicy, err := service.DescribeRedisInstancePasswordPolicyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if passwordPolicy == nil {
		log.Printf("[WARN]%s resource `redis_instance_password_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	policyMap := map[string]interface{}{}

	if passwordPolicy.Enabled != nil {
		policyMap["enabled"] = *passwordPolicy.Enabled
	}

	if passwordPolicy.MinLetterCount != nil {
		policyMap["min_letter_count"] = *passwordPolicy.MinLetterCount
	}

	if passwordPolicy.MinDigitCount != nil {
		policyMap["min_digit_count"] = *passwordPolicy.MinDigitCount
	}

	if passwordPolicy.MinSpecialCount != nil {
		policyMap["min_special_count"] = *passwordPolicy.MinSpecialCount
	}

	if passwordPolicy.MinLength != nil {
		policyMap["min_length"] = *passwordPolicy.MinLength
	}

	_ = d.Set("password_policy", []interface{}{policyMap})

	return nil
}

func resourceTencentCloudRedisInstancePasswordPolicyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	instanceId := d.Id()

	request := redis.NewModifyInstancePasswordPolicyRequest()
	request.InstanceId = &instanceId

	policy := redis.PasswordPolicy{}

	if v, ok := d.GetOk("password_policy"); ok {
		policyList := v.([]interface{})
		if len(policyList) > 0 {
			policyMap := policyList[0].(map[string]interface{})

			if v, ok := policyMap["enabled"]; ok {
				policy.Enabled = helper.Bool(v.(bool))
			}

			if v, ok := policyMap["min_letter_count"]; ok {
				policy.MinLetterCount = helper.IntInt64(v.(int))
			}

			if v, ok := policyMap["min_digit_count"]; ok {
				policy.MinDigitCount = helper.IntInt64(v.(int))
			}

			if v, ok := policyMap["min_special_count"]; ok {
				policy.MinSpecialCount = helper.IntInt64(v.(int))
			}

			if v, ok := policyMap["min_length"]; ok {
				policy.MinLength = helper.IntInt64(v.(int))
			}
		}
	}

	request.PasswordPolicy = &policy

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyInstancePasswordPolicy(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == RedisInstanceNotFound {
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
		log.Printf("[CRITAL]%s update redis_instance_password_policy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRedisInstancePasswordPolicyConfigRead(d, meta)
}

func resourceTencentCloudRedisInstancePasswordPolicyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
