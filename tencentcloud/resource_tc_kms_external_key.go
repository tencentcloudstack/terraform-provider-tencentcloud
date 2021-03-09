/*
Provide a resource to create a KMS external import key.

Example Usage

```hcl
resource "tencentcloud_kms_external_key" "foo" {
	alias = "test"
	description = "describe key test message."
	wrapping_algorithm = "RSAES_PKCS1_V1_5"
	key_material_base64 = "MTIzMTIzMTIzMTIzMTIzQQ=="
	valid_to = 2147443200
}
```

Import

KMS external keys can be imported using the id, e.g.

```
$ terraform import tencentcloud_kms_external_key.foo 287e8f40-7cbb-11eb-9a3a-5254004f7f94
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudKmsExternalKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsExternalKeyCreate,
		Read:   resourceTencentCloudKmsExternalKeyRead,
		Update: resourceTencentCloudKmsExternalKeyUpdate,
		Delete: resourceTencentCloudKmsExternalKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of CMK.",
			},
			"alias": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of CMK.The name can only contain English letters, numbers, underscore and hyphen '-'.The first character must be a letter or number.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Description of CMK.The maximum is 1024 bytes.",
			},
			"key_state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(KMS_KEY_STATE),
				Computed:     true,
				Description:  "State of CMK.Available values include `Enabled`, `Disabled`, `PendingDelete`, `PendingImport`, `Archived`.",
			},
			"pending_delete_window_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validateIntegerInRange(7, 30),
				Description:  "Duration in days after which the key is deleted after destruction of the resource, must be between 7 and 30 days. Defaults to 7 days.",
			},
			"wrapping_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      KMS_WRAPPING_ALGORITHM_RSAES_PKCS1_V1_5,
				ValidateFunc: validateAllowedStringValue(KMS_WRAPPING_ALGORITHM),
				Description:  "The algorithm for encrypting key material.Available values include `RSAES_PKCS1_V1_5`, `RSAES_OAEP_SHA_1` and `RSAES_OAEP_SHA_256`.",
			},
			"key_material_base64": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The base64-encoded key material encrypted with the public_key.For regions using the national secret version, the length of the imported key material is required to be 128 bits, and for regions using the FIPS version, the length of the imported key material is required to be 256 bits。",
			},
			"valid_to": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "this value means the effective timestamp of the key material, 0 means it does not expire.Need to be greater than the current time point, the maximum support is 2147443200.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Tags of CMK.",
			},
		},
	}
}

func resourceTencentCloudKmsExternalKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_external_key.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	kmsService := KmsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	keyType := KMS_ORIGIN_TYPE[KMS_ORIGIN_EXTERNAL]
	alias := d.Get("alias").(string)
	description := ""
	keyUsage := KMS_KEY_USAGE_ENCRYPT_DECRYPT
	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
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
	_ = d.Set("key_id", helper.String(keyId))

	if v, ok := d.GetOk("key_material_base64"); ok {
		param := make(map[string]interface{})
		param["key_id"] = keyId
		param["algorithm"] = d.Get("wrapping_algorithm").(string)
		param["key_spec"] = KMS_WRAPPING_KEY_SPEC_RSA_2048
		param["key_material_base64"] = v.(string)
		param["valid_to"] = d.Get("valid_to").(int)

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.ImportKeyMaterial(ctx, param)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s Create KMS external key failed, reason:%+v", logId, err)
			return err
		}
	}

	if keyState := d.Get("key_state").(string); keyState == KMS_KEY_STATE_DISABLED {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.DisableKey(ctx, d.Id())
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

	return resourceTencentCloudKmsExternalKeyRead(d, meta)
}

func resourceTencentCloudKmsExternalKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_external_key.read")()
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
		log.Printf("[CRITAL]%s read KMS external key failed, reason:%+v", logId, err)
		return err
	}

	if key == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("key_id", key.KeyId)
	_ = d.Set("alias", key.Alias)
	_ = d.Set("description", key.Description)
	_ = d.Set("key_state", key.KeyState)
	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "kms", "key", tcClient.Region, *key.ResourceId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudKmsExternalKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_external_key.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	keyId := d.Id()
	kmsService := KmsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

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
			log.Printf("[CRITAL]%s modify KMS external key description failed, reason:%+v", logId, err)
			return err
		}
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

	if d.HasChange("key_state") {
		oldKeyState, newKeyState := d.GetChange("key_state")
		err := updateKeyState(ctx, kmsService, keyId, oldKeyState.(string), newKeyState.(string))
		if err != nil {
			log.Printf("[CRITAL]%s modify KMS key state failed, reason:%+v", logId, err)
			return err
		}
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
	}

	return resourceTencentCloudKmsKeyRead(d, meta)
}

func resourceTencentCloudKmsExternalKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_external_key.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	kmsService := KmsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	keyId := d.Id()
	pendingDeleteWindowInDays := d.Get("pending_delete_window_in_days").(int)
	keyState := d.Get("key_state").(string)
	if keyState == KMS_KEY_STATE_ENABLED {
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

	return resourceTencentCloudKmsKeyRead(d, meta)
}

func updateKeyMaterial(ctx context.Context, kmsService KmsService, d *schema.ResourceData) error {
	param := make(map[string]interface{})
	param["key_id"] = d.Id()
	param["algorithm"] = d.Get("wrapping_algorithm").(string)
	param["key_spec"] = KMS_WRAPPING_KEY_SPEC_RSA_2048
	param["key_material_base64"] = d.Get("key_material_base64")
	param["valid_to"] = d.Get("valid_to").(int)

	var err error
	if param["key_material_base64"] == "" {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.DeleteImportKeyMaterial(ctx, d.Id())
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	} else {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := kmsService.ImportKeyMaterial(ctx, param)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	}

	return err
}
