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
				Description: "The ID of redis instance.",
			},

			"enabled": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable the instance-level password complexity policy. true: enable; false: disable.",
			},

			"min_letter_count": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The minimum number of letters (uppercase and lowercase). Value range: [1,16].",
			},

			"min_digit_count": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The minimum number of digit characters. Value range: [1,16].",
			},

			"min_special_count": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The minimum number of special characters. Value range: [1,16].",
			},

			"min_length": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The minimum total length of the password. Value range: [8,64].",
			},
		},
	}
}

func resourceTencentCloudRedisInstancePasswordPolicyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy_config.create")()
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
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	request := redis.NewDescribeInstancePasswordPolicyRequest()
	request.InstanceId = helper.String(instanceId)

	var passwordPolicy *redis.PasswordPolicy
	var notFound bool
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().DescribeInstancePasswordPolicyWithContext(ctx, request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound.InstanceNotExists" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil || result.Response.PasswordPolicy == nil {
			notFound = true
			return nil
		}
		passwordPolicy = result.Response.PasswordPolicy
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read redis_instance_password_policy_config failed, reason:%+v", logId, err)
		return err
	}

	if notFound {
		log.Printf("[CRUD] redis_instance_password_policy_config id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if passwordPolicy != nil {
		if passwordPolicy.Enabled != nil {
			_ = d.Set("enabled", passwordPolicy.Enabled)
		}
		if passwordPolicy.MinLetterCount != nil {
			_ = d.Set("min_letter_count", passwordPolicy.MinLetterCount)
		}
		if passwordPolicy.MinDigitCount != nil {
			_ = d.Set("min_digit_count", passwordPolicy.MinDigitCount)
		}
		if passwordPolicy.MinSpecialCount != nil {
			_ = d.Set("min_special_count", passwordPolicy.MinSpecialCount)
		}
		if passwordPolicy.MinLength != nil {
			_ = d.Set("min_length", passwordPolicy.MinLength)
		}
	}

	return nil
}

func resourceTencentCloudRedisInstancePasswordPolicyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	request := redis.NewModifyInstancePasswordPolicyRequest()
	request.InstanceId = helper.String(instanceId)

	passwordPolicy := &redis.PasswordPolicy{}

	if v, ok := d.GetOk("enabled"); ok {
		passwordPolicy.Enabled = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("min_letter_count"); ok {
		passwordPolicy.MinLetterCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("min_digit_count"); ok {
		passwordPolicy.MinDigitCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("min_special_count"); ok {
		passwordPolicy.MinSpecialCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("min_length"); ok {
		passwordPolicy.MinLength = helper.IntInt64(v.(int))
	}

	request.PasswordPolicy = passwordPolicy

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyInstancePasswordPolicyWithContext(ctx, request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound.InstanceNotExists" {
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
		log.Printf("[CRITAL]%s update redis_instance_password_policy_config failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRedisInstancePasswordPolicyConfigRead(d, meta)
}

func resourceTencentCloudRedisInstancePasswordPolicyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
