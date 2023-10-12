/*
Provide a resource to create a SSM secret.

Example Usage

Create user defined secret

```hcl
resource "tencentcloud_ssm_secret" "example" {
  secret_name             = "tf-example"
  description             = "desc."
  is_enabled              = true
  recovery_window_in_days = 0

  tags = {
    createBy = "terraform"
  }
}
```

Create redis secret

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 8
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[3].zone
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
}

resource "tencentcloud_redis_instance" "example" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[3].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[3].type_id
  password           = "Qwer@234"
  mem_size           = data.tencentcloud_redis_zone_config.zone.list[3].mem_sizes[0]
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[3].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[3].redis_replicas_nums[0]
  name               = "tf_example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_ssm_secret" "example" {
  secret_name       = "tf-example"
  description       = "redis desc."
  is_enabled        = true
  secret_type       = 4
  additional_config = jsonencode(
    {
      "Region" : "ap-guangzhou"
      "Privilege" : "r",
      "InstanceId" : tencentcloud_redis_instance.example.id
      "ReadonlyPolicy" : ["master"],
      "Remark" : "for tf test"
    }
  )
  tags = {
    createdBy = "terraform"
  }
  recovery_window_in_days = 0
}
```

Import

SSM secret can be imported using the secretName, e.g.
```
$ terraform import tencentcloud_ssm_secret.foo test
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSsmSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmSecretCreate,
		Read:   resourceTencentCloudSsmSecretRead,
		Update: resourceTencentCloudSsmSecretUpdate,
		Delete: resourceTencentCloudSsmSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"secret_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of secret which cannot be repeated in the same region. The maximum length is 128 bytes. The name can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of secret. The maximum is 2048 bytes.",
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "KMS keyId used to encrypt secret. If it is empty, it means that the CMK created by SSM for you by default is used for encryption. You can also specify the KMS CMK created by yourself in the same region for encryption.",
			},
			"secret_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Type of secret. `0`: user-defined secret. `4`: redis secret.",
			},
			"additional_config": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional config for specific secret types in JSON string format.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of secret.",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specify whether to enable secret. Default value is `true`.",
			},
			"recovery_window_in_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specify the scheduled deletion date. Default value is `0` that means to delete immediately. 1-30 means the number of days reserved, completely deleted after this date.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of secret.",
			},
		},
	}
}

func resourceTencentCloudSsmSecretCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.create")()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		ssmService    = SsmService{client: meta.(*TencentCloudClient).apiV3Conn}
		request       = ssm.NewCreateSecretRequest()
		response      = ssm.NewCreateSecretResponse()
		secretInfo    *SecretInfo
		outErr, inErr error
		secretName    string
	)

	if v, ok := d.GetOk("secret_name"); ok {
		request.SecretName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("secret_type"); ok {
		request.SecretType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("additional_config"); ok {
		request.AdditionalConfig = helper.String(v.(string))
	}

	request.SecretString = helper.String("default")
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSsmClient().CreateSecret(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if outErr != nil {
		return outErr
	}

	secretName = *response.Response.SecretName
	d.SetId(secretName)

	if isEnabled := d.Get("is_enabled").(bool); !isEnabled {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = ssmService.DisableSecret(ctx, secretName)
			if inErr != nil {
				return retryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			secretInfo, inErr = ssmService.DescribeSecretByName(ctx, secretName)
			if inErr != nil {
				return retryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudSsmSecretRead(d, meta)
}

func resourceTencentCloudSsmSecretRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.read")()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		ssmService    = SsmService{client: meta.(*TencentCloudClient).apiV3Conn}
		secretName    = d.Id()
		outErr, inErr error
		secretInfo    *SecretInfo
	)

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		secretInfo, inErr = ssmService.DescribeSecretByName(ctx, secretName)
		if inErr != nil {
			return retryError(inErr)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	if secretInfo == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("secret_name", secretInfo.secretName)
	_ = d.Set("description", secretInfo.description)
	_ = d.Set("kms_key_id", secretInfo.kmsKeyId)
	_ = d.Set("secret_type", secretInfo.secretType)
	_ = d.Set("additional_config", secretInfo.additionalConfig)
	_ = d.Set("status", secretInfo.status)

	if secretInfo.status == SSM_STATUS_ENABLED {
		_ = d.Set("is_enabled", true)
	} else {
		_ = d.Set("is_enabled", false)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "ssm", "secret", tcClient.Region, secretInfo.resourceId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudSsmSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.update")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		ssmService = SsmService{client: meta.(*TencentCloudClient).apiV3Conn}
		secretName = d.Id()
	)

	d.Partial(true)

	immutableArgs := []string{
		"secret_type",
		"additional_config",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := ssmService.UpdateSecretDescription(ctx, secretName, description)
			if e != nil {
				return retryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify SSM secret description failed, reason:%+v", logId, err)
			return err
		}

	}

	if d.HasChange("is_enabled") {
		isEnabled := d.Get("is_enabled").(bool)
		err := updateSecretIsEnabled(ctx, ssmService, secretName, isEnabled)
		if err != nil {
			log.Printf("[CRITAL]%s modify SSM secret status failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		secretInfo, err := ssmService.DescribeSecretByName(ctx, secretName)
		if err != nil {
			return err
		}

		resourceName := BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	d.Partial(false)
	return resourceTencentCloudSsmSecretRead(d, meta)
}

func resourceTencentCloudSsmSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.delete")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		ssmService = SsmService{client: meta.(*TencentCloudClient).apiV3Conn}
		secretName = d.Id()
	)

	recoveryWindowInDays := d.Get("recovery_window_in_days").(int)
	isEnabled := d.Get("is_enabled").(bool)
	if isEnabled {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := ssmService.DisableSecret(ctx, secretName)
			if e != nil {
				return retryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify SSM secret status failed, reason:%+v", logId, err)
			return err
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := ssmService.DeleteSecret(ctx, secretName, uint64(recoveryWindowInDays))
		if e != nil {
			return retryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete SSM secret failed, reason:%+v", logId, err)
		return err
	}

	return resource.Retry(readRetryTimeout, func() *resource.RetryError {
		secretInfo, e := ssmService.DescribeSecretByName(ctx, secretName)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok && sdkError.Code == "ResourceNotFound" {
				return nil
			}

			return retryError(err)
		}

		if secretInfo.status == SSM_STATUS_PENDINGDELETE {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}

func updateSecretIsEnabled(ctx context.Context, ssmService SsmService, secretName string, isEnabled bool) error {
	var err error
	if isEnabled {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := ssmService.EnableSecret(ctx, secretName)
			if e != nil {
				return retryError(e)
			}
			return nil
		})

	} else {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := ssmService.DisableSecret(ctx, secretName)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	}
	return err
}
