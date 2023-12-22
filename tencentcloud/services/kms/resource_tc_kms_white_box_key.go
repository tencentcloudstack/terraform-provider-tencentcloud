package kms

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKmsWhiteBoxKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsWhiteBoxKeyCreate,
		Read:   resourceTencentCloudKmsWhiteBoxKeyRead,
		Update: resourceTencentCloudKmsWhiteBoxKeyUpdate,
		Delete: resourceTencentCloudKmsWhiteBoxKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alias": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "As an alias for the key to be easier to identify and easier to understand, it cannot be empty and is a combination of 1-60 alphanumeric characters - _. The first character must be a letter or number. Alias are not repeatable.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description of the key, up to 1024 bytes.",
			},
			"algorithm": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(WHITE_BOX_KEY_ALGORITHM),
				Description:  "All algorithm types for creating keys, supported values: AES_256, SM4.",
			},
			"status": {
				Optional:     true,
				Type:         schema.TypeString,
				Default:      WHITE_BOX_KEY_STATUS_ENABLED,
				ValidateFunc: tccommon.ValidateAllowedStringValue(WHITE_BOX_KEY_STATUS),
				Description:  "Whether to enable the key. Enabled or Disabled. Default is Enabled.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of Key.",
			},
		},
	}
}

func resourceTencentCloudKmsWhiteBoxKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_white_box_key.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request  = kms.NewCreateWhiteBoxKeyRequest()
		response = kms.NewCreateWhiteBoxKeyResponse()
		keyId    string
	)

	if v, ok := d.GetOk("alias"); ok {
		request.Alias = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("algorithm"); ok {
		request.Algorithm = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseKmsClient().CreateWhiteBoxKey(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("kms whiteBoxKey not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create kms whiteBoxKey failed, reason:%+v", logId, err)
		return err
	}

	keyId = *response.Response.KeyId
	d.SetId(keyId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::kms:%s:uin/:key/%s", region, d.Id())
		if err = tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		if status == WHITE_BOX_KEY_STATUS_DISABLED {
			disableWhiteBoxKeyRequest := kms.NewDisableWhiteBoxKeyRequest()
			disableWhiteBoxKeyRequest.KeyId = &keyId

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseKmsClient().DisableWhiteBoxKey(disableWhiteBoxKeyRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s disable kms whiteBoxKey failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudKmsWhiteBoxKeyRead(d, meta)
}

func resourceTencentCloudKmsWhiteBoxKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_white_box_key.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		keyId   = d.Id()
	)

	whiteBoxKey, err := service.DescribeKmsWhiteBoxKeyById(ctx, keyId)
	if err != nil {
		return err
	}

	if whiteBoxKey == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `KmsWhiteBoxKey` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if whiteBoxKey.Alias != nil {
		_ = d.Set("alias", whiteBoxKey.Alias)
	}

	if whiteBoxKey.Description != nil {
		_ = d.Set("description", whiteBoxKey.Description)
	}

	if whiteBoxKey.Algorithm != nil {
		_ = d.Set("algorithm", whiteBoxKey.Algorithm)
	}

	if whiteBoxKey.Status != nil {
		_ = d.Set("status", whiteBoxKey.Status)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "kms", "key", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudKmsWhiteBoxKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_white_box_key.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                     = tccommon.GetLogId(tccommon.ContextNil)
		enableWhiteBoxKeyRequest  = kms.NewEnableWhiteBoxKeyRequest()
		disableWhiteBoxKeyRequest = kms.NewDisableWhiteBoxKeyRequest()
		keyId                     = d.Id()
	)

	immutableArgs := []string{"alias", "description", "algorithm"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			status := v.(string)
			if status == WHITE_BOX_KEY_STATUS_DISABLED {
				disableWhiteBoxKeyRequest.KeyId = &keyId
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseKmsClient().DisableWhiteBoxKey(disableWhiteBoxKeyRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableWhiteBoxKeyRequest.GetAction(), disableWhiteBoxKeyRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s disable kms whiteBoxKey failed, reason:%+v", logId, err)
					return err
				}
			} else {
				enableWhiteBoxKeyRequest.KeyId = &keyId
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseKmsClient().EnableWhiteBoxKey(enableWhiteBoxKeyRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableWhiteBoxKeyRequest.GetAction(), enableWhiteBoxKeyRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s enable kms whiteBoxKey failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("kms", "key", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudKmsWhiteBoxKeyRead(d, meta)
}

func resourceTencentCloudKmsWhiteBoxKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_white_box_key.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		keyId   = d.Id()
	)

	if err := service.DeleteKmsWhiteBoxKeyById(ctx, keyId); err != nil {
		return err
	}

	return nil
}
