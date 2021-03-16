/*
Provide a resource to create a KMS key.

Example Usage

```hcl
resource "tencentcloud_kms_key" "foo" {
	alias = "test"
	description = "describe key test message."
	key_rotation_enabled = true
	is_enabled = true

	tags = {
		"test-tag":"key-test"
	}
}
```

Import

KMS keys can be imported using the id, e.g.

```
$ terraform import tencentcloud_kms_key.foo 287e8f40-7cbb-11eb-9a3a-5254004f7f94
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TencentKmsBasicInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alias": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of CMK. The name can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of CMK. The maximum is 1024 bytes.",
		},
		"is_enabled": {
			Type:          schema.TypeBool,
			Optional:      true,
			ConflictsWith: []string{"is_archived"},
			Description:   "Specify whether to enable key. Default value is `false`.",
		},
		"is_archived": {
			Type:          schema.TypeBool,
			Optional:      true,
			ConflictsWith: []string{"is_enabled"},
			Description:   "Specify whether to archive key. Default value is `false`.",
		},
		"pending_delete_window_in_days": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     7,
			Description: "Duration in days after which the key is deleted after destruction of the resource, must be between 7 and 30 days. Defaults to 7 days.",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Tags of CMK.",
		},
		"key_state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "State of CMK.",
		},
	}
}

func resourceTencentCloudKmsKey() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"key_usage": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     KMS_KEY_USAGE_ENCRYPT_DECRYPT,
			Description: "Usage of CMK. Available values include `ENCRYPT_DECRYPT`, `ASYMMETRIC_DECRYPT_RSA_2048`, `ASYMMETRIC_DECRYPT_SM2`, `ASYMMETRIC_SIGN_VERIFY_SM2`, `ASYMMETRIC_SIGN_VERIFY_RSA_2048`, `ASYMMETRIC_SIGN_VERIFY_ECC`. Default value is `ENCRYPT_DECRYPT`.",
		},
		"key_rotation_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Specify whether to enable key rotation, valid when key_usage is `ENCRYPT_DECRYPT`. Default value is `false`.",
		},
	}

	basic := TencentKmsBasicInfo()
	for k, v := range basic {
		specialInfo[k] = v
	}

	return &schema.Resource{
		Create: resourceTencentCloudKmsKeyCreate,
		Read:   resourceTencentCloudKmsKeyRead,
		Update: resourceTencentCloudKmsKeyUpdate,
		Delete: resourceTencentCloudKmsKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: specialInfo,
	}
}

func resourceTencentCloudKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_key.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	kmsService := KmsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	keyType := KMS_ORIGIN_TYPE[KMS_ORIGIN_TENCENT_KMS]
	alias := d.Get("alias").(string)
	description := ""
	keyUsage := ""
	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}
	if v, ok := d.GetOk("key_usage"); ok {
		keyUsage = v.(string)
	}

	var keyId string
	var outErr, inErr error
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		keyId, inErr = kmsService.CreateKey(ctx, keyType, alias, description, keyUsage)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		log.Printf("[CRITAL]%s create KMS key failed, reason:%+v", logId, outErr)
		return outErr
	}
	d.SetId(keyId)

	if isEnabled := d.Get("is_enabled").(bool); !isEnabled {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKey(ctx, d.Id())
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify key state failed, reason:%+v", logId, err)
			return err
		}
	}

	if keyRotationEnabled := d.Get("key_rotation_enabled").(bool); keyRotationEnabled {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.EnableKeyRotation(ctx, d.Id())
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify KMS key rotation status failed, reason:%+v", logId, err)
			return err
		}
	}

	if isArchived := d.Get("is_archived").(bool); isArchived {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.ArchiveKey(ctx, d.Id())
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify key state failed, reason:%+v", logId, err)
			return err
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		keyMetaData, err := kmsService.DescribeKeyById(ctx, keyId)
		if err != nil {
			return err
		}
		resourceName := BuildTagResourceName("kms", "key", tcClient.Region, *keyMetaData.ResourceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudKmsKeyRead(d, meta)

}

func resourceTencentCloudKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_key.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	keyId := d.Id()
	kmsService := &KmsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var key *kms.KeyMetadata
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := kmsService.DescribeKeyById(ctx, keyId)
		if e != nil {
			return retryError(e)
		}
		key = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read KMS key failed, reason:%+v", logId, err)
		return err
	}

	if key == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("alias", key.Alias)
	_ = d.Set("description", key.Description)
	_ = d.Set("key_state", key.KeyState)
	_ = d.Set("key_usage", key.KeyUsage)
	_ = d.Set("key_rotation_enabled", key.KeyRotationEnabled)
	transformKeyState(d)

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "kms", "key", tcClient.Region, *key.ResourceId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_key.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	keyId := d.Id()
	kmsService := KmsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	d.Partial(true)

	if d.HasChange("description") {
		description := d.Get("description").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.UpdateKeyDescription(ctx, keyId, description)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify KMS key description failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("description")
	}

	if d.HasChange("alias") {
		alias := d.Get("alias").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.UpdateKeyAlias(ctx, keyId, alias)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify KMS key alias failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("alias")
	}

	if isArchived, ok := d.GetOk("is_archived"); ok {
		err := updateIsArchived(ctx, kmsService, keyId, isArchived.(bool))
		if err != nil {
			log.Printf("[CRITAL]%s modify key state failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("is_archived")
	} else {
		isEnabled := d.Get("is_enabled").(bool)
		err := updateIsEnabled(ctx, kmsService, keyId, isEnabled)
		if err != nil {
			log.Printf("[CRITAL]%s modify key state failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("is_enabled")
	}

	if d.HasChange("key_rotation_enabled") {
		keyRotationEnabled := d.Get("key_rotation_enabled").(bool)
		err := updateKeyRotationStatus(ctx, kmsService, keyId, keyRotationEnabled)
		if err != nil {
			log.Printf("[CRITAL]%s modify KMS key rotation status failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("key_rotation_enabled")
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		keyMetaData, err := kmsService.DescribeKeyById(ctx, keyId)
		if err != nil {
			return err
		}
		resourceName := BuildTagResourceName("kms", "key", tcClient.Region, *keyMetaData.ResourceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudKmsKeyRead(d, meta)
}

func resourceTencentCloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_key.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	kmsService := KmsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	keyId := d.Id()
	pendingDeleteWindowInDays := d.Get("pending_delete_window_in_days").(int)
	isEnabled := d.Get("is_enabled").(bool)
	if isEnabled {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKey(ctx, keyId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify KMS key state failed, reason:%+v", logId, err)
			return err
		}
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := kmsService.DeleteKey(ctx, keyId, uint64(pendingDeleteWindowInDays))
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete KMS key failed, reason:%+v", logId, err)
		return err
	}

	return resource.Retry(readRetryTimeout, func() *resource.RetryError {
		key, e := kmsService.DescribeKeyById(ctx, keyId)
		if e != nil {
			return retryError(e)
		}
		if *key.KeyState == KMS_KEY_STATE_PENDINGDELETE {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}

func updateKeyRotationStatus(ctx context.Context, kmsService KmsService, keyId string, keyRotationEnabled bool) error {
	var err error
	if keyRotationEnabled {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.EnableKeyRotation(ctx, keyId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	} else {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKeyRotation(ctx, keyId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	}
	return err
}

func updateIsEnabled(ctx context.Context, kmsService KmsService, keyId string, isEnabled bool) error {
	var err error
	if isEnabled {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.EnableKey(ctx, keyId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})

	} else {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKey(ctx, keyId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	}
	return err
}

func updateIsArchived(ctx context.Context, kmsService KmsService, keyId string, isArchived bool) error {
	var err error
	if isArchived {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.ArchiveKey(ctx, keyId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	} else {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.CancelKeyArchive(ctx, keyId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	}
	return err
}

func transformKeyState(d *schema.ResourceData) {
	keyState := d.Get("key_state").(string)
	switch keyState {
	case KMS_KEY_STATE_ENABLED:
		_ = d.Set("is_enabled", true)
	case KMS_KEY_STATE_DISABLED:
		_ = d.Set("is_enabled", false)
	case KMS_KEY_STATE_ARCHIVED:
		_ = d.Set("is_archived", true)
	}
}
