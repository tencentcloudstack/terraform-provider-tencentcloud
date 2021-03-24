/*
Provide a resource to create a SSM secret.
Example Usage
```hcl
resource "tencentcloud_ssm_secret" "foo" {
  secret_name = "test"
  description = "test secret"
  recovery_window_in_days = 0
  is_enabled = true

  tags = {
    test-tag = "test"
  }
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specify whether to enable secret. Default value is `true`.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of secret. The maximum is 2048 bytes.",
			},
			"recovery_window_in_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specify the scheduled deletion date. Default value is `0` that means to delete immediately. 1-30 means the number of days reserved, completely deleted after this date.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of secret.",
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "KMS keyId used to encrypt secret. If it is empty, it means that the CMK created by SSM for you by default is used for encryption. You can also specify the KMS CMK created by yourself in the same region for encryption.",
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	param := make(map[string]interface{})
	param["secret_name"] = d.Get("secret_name").(string)
	if v, ok := d.GetOk("description"); ok {
		param["description"] = v.(string)
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		param["kms_key_id"] = v.(string)
	}
	//use a default version info, after create secret will delete this version
	//because sdk do not support create secret without version
	param["version_id"] = "default"
	param["secret_string"] = "default"

	var outErr, inErr error
	var secretName string
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		secretName, inErr = ssmService.CreateSecret(ctx, param)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	d.SetId(secretName)

	//delete default version info
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr = ssmService.DeleteSecretVersion(ctx, secretName, "default")
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

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

	var secretInfo *SecretInfo
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

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	secretName := d.Id()

	var outErr, inErr error
	var secretInfo *SecretInfo
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	d.Partial(true)
	secretName := d.Id()

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
		d.SetPartial("description")
	}

	if d.HasChange("is_enabled") {
		isEnabled := d.Get("is_enabled").(bool)
		err := updateSecretIsEnabled(ctx, ssmService, secretName, isEnabled)
		if err != nil {
			log.Printf("[CRITAL]%s modify SSM secret status failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("is_enabled")
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
		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudSsmSecretRead(d, meta)
}

func resourceTencentCloudSsmSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	secretName := d.Id()
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
