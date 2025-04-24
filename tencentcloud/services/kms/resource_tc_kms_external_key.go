package kms

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKmsExternalKey() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"wrapping_algorithm": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     KMS_WRAPPING_ALGORITHM_RSAES_PKCS1_V1_5,
			Description: "The algorithm for encrypting key material. Available values include `RSAES_PKCS1_V1_5`, `RSAES_OAEP_SHA_1` and `RSAES_OAEP_SHA_256`. Default value is `RSAES_PKCS1_V1_5`.",
		},
		"key_material_base64": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The base64-encoded key material encrypted with the public_key. For regions using the national secret version, the length of the imported key material is required to be 128 bits, and for regions using the FIPS version, the length of the imported key material is required to be 256 bits.",
		},
		"valid_to": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "This value means the effective timestamp of the key material, 0 means it does not expire. Need to be greater than the current timestamp, the maximum support is 2147443200.",
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
		Create: resourceTencentCloudKmsExternalKeyCreate,
		Read:   resourceTencentCloudKmsExternalKeyRead,
		Update: resourceTencentCloudKmsExternalKeyUpdate,
		Delete: resourceTencentCloudKmsExternalKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: specialInfo,
	}
}

func resourceTencentCloudKmsExternalKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_external_key.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	kmsService := KmsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	keyType := KMS_ORIGIN_TYPE[KMS_ORIGIN_EXTERNAL]
	alias := d.Get("alias").(string)
	description := ""
	hsmClusterId := ""
	keyUsage := KMS_KEY_USAGE_ENCRYPT_DECRYPT
	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
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

	if v, ok := d.GetOk("key_material_base64"); ok {
		param := make(map[string]interface{})
		param["key_id"] = keyId
		param["algorithm"] = d.Get("wrapping_algorithm").(string)
		param["key_spec"] = KMS_WRAPPING_KEY_SPEC_RSA_2048
		param["key_material_base64"] = v.(string)
		param["valid_to"] = d.Get("valid_to").(int)

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.ImportKeyMaterial(ctx, param)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s Create KMS external key failed, reason:%+v", logId, err)
			return err
		}

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

	return resourceTencentCloudKmsExternalKeyRead(d, meta)
}

func resourceTencentCloudKmsExternalKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_external_key.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	keyId := d.Id()
	kmsService := &KmsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var key *kms.KeyMetadata
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := kmsService.DescribeKeyById(ctx, keyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		key = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read KMS external key failed, reason:%+v", logId, err)
		return err
	}

	if key == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("alias", key.Alias)
	_ = d.Set("description", key.Description)
	_ = d.Set("valid_to", key.ValidTo)
	_ = d.Set("key_state", key.KeyState)
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

func resourceTencentCloudKmsExternalKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_external_key.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	keyId := d.Id()
	kmsService := KmsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

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
			log.Printf("[CRITAL]%s modify KMS external key description failed, reason:%+v", logId, err)
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
			log.Printf("[CRITAL]%s modify KMS external key alias failed, reason:%+v", logId, err)
			return err
		}

	}

	if d.HasChange("key_material_base64") || d.HasChange("valid_to") {
		err := updateKeyMaterial(ctx, kmsService, d)
		if err != nil {
			log.Printf("[CRITAL]%s import KMS external key material key failed, reason:%+v", logId, err)
			return err
		}
	}

	var key *kms.KeyMetadata
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := kmsService.DescribeKeyById(ctx, keyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		key = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read KMS external key failed, reason:%+v", logId, err)
		return err
	}

	if *key.KeyState == KMS_KEY_STATE_ENABLED || *key.KeyState == KMS_KEY_STATE_DISABLED || *key.KeyState == KMS_KEY_STATE_ARCHIVED {
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

	return resourceTencentCloudKmsExternalKeyRead(d, meta)
}

func resourceTencentCloudKmsExternalKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_external_key.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	kmsService := KmsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	keyId := d.Id()
	pendingDeleteWindowInDays := d.Get("pending_delete_window_in_days").(int)
	isEnabled := d.Get("is_enabled").(bool)
	if isEnabled {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKey(ctx, keyId)
			if e != nil {
				return tccommon.RetryError(e)
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
			return tccommon.RetryError(e)
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

func updateKeyMaterial(ctx context.Context, kmsService KmsService, d *schema.ResourceData) error {
	param := make(map[string]interface{})
	param["key_id"] = d.Id()
	param["algorithm"] = d.Get("wrapping_algorithm").(string)
	param["key_spec"] = KMS_WRAPPING_KEY_SPEC_RSA_2048
	param["key_material_base64"] = d.Get("key_material_base64").(string)
	param["valid_to"] = d.Get("valid_to").(int)

	var err error
	if d.HasChange("key_material_base64") && param["key_material_base64"] == "" {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.DeleteImportKeyMaterial(ctx, d.Id())
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
	} else {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := kmsService.ImportKeyMaterial(ctx, param)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
	}

	return err
}
