/*
Provides a resource to create a kms key

Example Usage

```hcl
resource "tencentcloud_kms_key" "key" {
  alias = "test"
  description = "test"
  key_usage = "test"
  type = 11
  hsm_cluster_id = "ss"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

kms key can be imported using the id, e.g.

```
terraform import tencentcloud_kms_key.key key_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudKmsKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsKeyCreate,
		Read:   resourceTencentCloudKmsKeyRead,
		Update: resourceTencentCloudKmsKeyUpdate,
		Delete: resourceTencentCloudKmsKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alias": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Key alias.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cmk description.",
			},

			"key_usage": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Test.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Test.",
			},

			"hsm_cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Test.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_key.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = kms.NewCreateKeyRequest()
		response = kms.NewCreateKeyResponse()
		keyId    string
	)
	if v, ok := d.GetOk("alias"); ok {
		request.Alias = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("key_usage"); ok {
		request.KeyUsage = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("hsm_cluster_id"); ok {
		request.HsmClusterId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseKmsClient().CreateKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create kms key failed, reason:%+v", logId, err)
		return err
	}

	keyId = *response.Response.KeyId
	d.SetId(keyId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::kms:%s:uin/:key/%s", region, d.Id())
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

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	keyId := d.Id()

	key, err := service.DescribeKmsKeyById(ctx, keyId)
	if err != nil {
		return err
	}

	if key == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `KmsKey` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if key.Alias != nil {
		_ = d.Set("alias", key.Alias)
	}

	if key.Description != nil {
		_ = d.Set("description", key.Description)
	}

	if key.KeyUsage != nil {
		_ = d.Set("key_usage", key.KeyUsage)
	}

	if key.Type != nil {
		_ = d.Set("type", key.Type)
	}

	if key.HsmClusterId != nil {
		_ = d.Set("hsm_cluster_id", key.HsmClusterId)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "kms", "key", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_key.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		updateAliasRequest  = kms.NewUpdateAliasRequest()
		updateAliasResponse = kms.NewUpdateAliasResponse()
	)

	keyId := d.Id()

	request.KeyId = &keyId

	immutableArgs := []string{"alias", "description", "key_usage", "type", "hsm_cluster_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("alias") {
		if v, ok := d.GetOk("alias"); ok {
			request.Alias = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseKmsClient().UpdateAlias(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update kms key failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("kms", "key", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudKmsKeyRead(d, meta)
}

func resourceTencentCloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_key.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}
	keyId := d.Id()

	if err := service.DeleteKmsKeyById(ctx, keyId); err != nil {
		return err
	}

	return nil
}
