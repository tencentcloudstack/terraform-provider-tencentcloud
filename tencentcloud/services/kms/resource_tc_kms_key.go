package kms

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			Description:   "Specify whether to enable key. Default value is `false`. This field is conflict with `is_archived`, valid when key_state is `Enabled`, `Disabled`, `Archived`.",
		},
		"is_archived": {
			Type:          schema.TypeBool,
			Optional:      true,
			ConflictsWith: []string{"is_enabled"},
			Description:   "Specify whether to archive key. Default value is `false`. This field is conflict with `is_enabled`, valid when key_state is `Enabled`, `Disabled`, `Archived`.",
		},
		"pending_delete_window_in_days": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      7,
			ValidateFunc: tccommon.ValidateIntegerInRange(7, 30),
			Description:  "Duration in days after which the key is deleted after destruction of the resource, must be between 7 and 30 days. Defaults to 7 days.",
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

func ResourceTencentCloudKmsKey() *schema.Resource {
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
		"hsm_cluster_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The HSM cluster ID corresponding to KMS Advanced Edition (only valid for KMS Exclusive/Managed Edition service instances).",
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
	defer tccommon.LogElapsed("resource.tencentcloud_kms_key.create")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		kmsService = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	keyType := KMS_ORIGIN_TYPE[KMS_ORIGIN_TENCENT_KMS]
	alias := d.Get("alias").(string)
	description := ""
	keyUsage := ""
	hsmClusterId := ""

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}

	if v, ok := d.GetOk("key_usage"); ok {
		keyUsage = v.(string)
	}

	if v, ok := d.GetOk("hsm_cluster_id"); ok {
		hsmClusterId = v.(string)
	}

	var keyId string
	var outErr, inErr error
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		keyId, inErr = kmsService.CreateKey(ctx, keyType, alias, description, keyUsage, hsmClusterId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if outErr != nil {
		log.Printf("[CRITAL]%s create KMS key failed, reason:%+v", logId, outErr)
		return outErr
	}

	d.SetId(keyId)

	if isEnabled := d.Get("is_enabled").(bool); !isEnabled {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKey(ctx, d.Id())
			if e != nil {
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify key state failed, reason:%+v", logId, err)
			return err
		}
	}

	if keyUsage == KMS_KEY_USAGE_ENCRYPT_DECRYPT {
		if keyRotationEnabled := d.Get("key_rotation_enabled").(bool); keyRotationEnabled {
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				e := kmsService.EnableKeyRotation(ctx, d.Id())
				if e != nil {
					return tccommon.RetryError(e)
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify KMS key rotation status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if isArchived := d.Get("is_archived").(bool); isArchived {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.ArchiveKey(ctx, d.Id())
			if e != nil {
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify key state failed, reason:%+v", logId, err)
			return err
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		keyMetaData, err := kmsService.DescribeKeyById(ctx, keyId)
		if err != nil {
			return err
		}

		resourceName := tccommon.BuildTagResourceName("kms", "key", tcClient.Region, *keyMetaData.ResourceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudKmsKeyRead(d, meta)
}

func resourceTencentCloudKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_key.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		kmsService = &KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		keyId      = d.Id()
		key        *kms.KeyMetadata
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := kmsService.DescribeKeyById(ctx, keyId)
		if e != nil {
			return tccommon.RetryError(e)
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
	if key.HsmClusterId != nil {
		_ = d.Set("hsm_cluster_id", key.HsmClusterId)
	}
	transformKeyState(d)

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "kms", "key", tcClient.Region, *key.ResourceId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_key.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		kmsService = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		keyId      = d.Id()
	)

	immutableArgs := []string{"hsm_cluster_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	d.Partial(true)
	if d.HasChange("description") {
		description := d.Get("description").(string)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.UpdateKeyDescription(ctx, keyId, description)
			if e != nil {
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify KMS key description failed, reason:%+v", logId, err)
			return err
		}

	}

	if d.HasChange("alias") {
		alias := d.Get("alias").(string)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.UpdateKeyAlias(ctx, keyId, alias)
			if e != nil {
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify KMS key alias failed, reason:%+v", logId, err)
			return err
		}

	}

	if keyState := d.Get("key_state").(string); keyState == KMS_KEY_STATE_ENABLED || keyState == KMS_KEY_STATE_DISABLED || keyState == KMS_KEY_STATE_ARCHIVED {
		if isArchived, ok := d.GetOk("is_archived"); ok {
			err := updateIsArchived(ctx, kmsService, keyId, isArchived.(bool))
			if err != nil {
				log.Printf("[CRITAL]%s modify key state failed, reason:%+v", logId, err)
				return err
			}

		} else {
			isEnabled := d.Get("is_enabled").(bool)
			err := updateIsEnabled(ctx, kmsService, keyId, isEnabled)
			if err != nil {
				log.Printf("[CRITAL]%s modify key state failed, reason:%+v", logId, err)
				return err
			}

		}
	}

	if v := d.Get("key_usage").(string); v == KMS_KEY_USAGE_ENCRYPT_DECRYPT {
		if d.HasChange("key_rotation_enabled") {
			keyRotationEnabled := d.Get("key_rotation_enabled").(bool)
			err := updateKeyRotationStatus(ctx, kmsService, keyId, keyRotationEnabled)
			if err != nil {
				log.Printf("[CRITAL]%s modify KMS key rotation status failed, reason:%+v", logId, err)
				return err
			}

		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		keyMetaData, err := kmsService.DescribeKeyById(ctx, keyId)
		if err != nil {
			return err
		}
		resourceName := tccommon.BuildTagResourceName("kms", "key", tcClient.Region, *keyMetaData.ResourceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudKmsKeyRead(d, meta)
}

func resourceTencentCloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_key.delete")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		kmsService = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		keyId      = d.Id()
	)

	pendingDeleteWindowInDays := d.Get("pending_delete_window_in_days").(int)
	isEnabled := d.Get("is_enabled").(bool)
	if isEnabled {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKey(ctx, keyId)
			if e != nil {
				ee, ok := e.(*sdkErrors.TencentCloudSDKError)
				if ok && tccommon.IsContains(KMS_RETRYABLE_ERROR, ee.Code) {
					return resource.RetryableError(fmt.Errorf("kms key disable error: %s, retrying", e.Error()))
				}

				return resource.NonRetryableError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify KMS key state failed, reason:%+v", logId, err)
			return err
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := kmsService.DeleteKey(ctx, keyId, uint64(pendingDeleteWindowInDays))
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(KMS_RETRYABLE_ERROR, ee.Code) {
				return resource.RetryableError(fmt.Errorf("kms key delete error: %s, retrying", e.Error()))
			}

			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete KMS key failed, reason:%+v", logId, err)
		return err
	}

	return resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		key, e := kmsService.DescribeKeyById(ctx, keyId)
		if e != nil {
			return tccommon.RetryError(e)
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
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.EnableKeyRotation(ctx, keyId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
	} else {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKeyRotation(ctx, keyId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
	}
	return err
}

func updateIsEnabled(ctx context.Context, kmsService KmsService, keyId string, isEnabled bool) error {
	var err error
	if isEnabled {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.EnableKey(ctx, keyId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})

	} else {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKey(ctx, keyId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
	}
	return err
}

func updateIsArchived(ctx context.Context, kmsService KmsService, keyId string, isArchived bool) error {
	var err error
	if isArchived {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.ArchiveKey(ctx, keyId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
	} else {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.CancelKeyArchive(ctx, keyId)
			if e != nil {
				return tccommon.RetryError(e)
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
