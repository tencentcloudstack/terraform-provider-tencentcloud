/*
Use this data source to query detailed information of KMS key

Example Usage

```hcl
data "tencentcloud_kms_key" "foo" {
	search_key_alias = "test"
	key_state = "All"
	origin = "TENCENT_KMS"
	key_usage = "ALL"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
)

func dataSourceTencentCloudKmsKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsKeyRead,
		Schema: map[string]*schema.Schema{
			"role": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Default:      0,
				Description:  "Role of the CMK creator.`0` - created by user, `1` - created by cloud product.Default value is `0`.",
			},
			"order_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Default:      0,
				Description:  "Order to sort the CMK create time.`0` - desc, `1` - asc.Default value is `0`.",
			},
			"key_state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(KMS_KEY_STATE_FILTER),
				Default:      KMS_KEY_STATE_ALL,
				Description:  "State of CMK.Available values include `All`, `Enabled`, `Disabled`, `PendingDelete`, `PendingImport`, `Archived`.",
			},
			"search_key_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Words used to match the results,and the words can be: key_id and alias.",
			},
			"origin": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(KMS_ORIGIN_FILTER),
				Default:      KMS_ORIGIN_ALL,
				Description:  "Origin of CMK.`TENCENT_KMS` - CMK created by KMS, `EXTERNAL` - CMK imported by user, `ALL` - All CMK.Default value is `ALL`.",
			},
			"key_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(KMS_KEY_USAGE_FILTER),
				Default:      KMS_KEY_USAGE_ENCRYPT_DECRYPT,
				Description:  "Usage of CMK.Available values include `ALL`, `ENCRYPT_DECRYPT`, `ASYMMETRIC_DECRYPT_RSA_2048`, `ASYMMETRIC_DECRYPT_SM2`, `ASYMMETRIC_SIGN_VERIFY_SM2`, `ASYMMETRIC_SIGN_VERIFY_RSA_2048`, `ASYMMETRIC_SIGN_VERIFY_ECC`.Default value is `ENCRYPT_DECRYPT`.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags to filter CMK.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"key_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of KMS keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of CMK.",
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CMK.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of CMK.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of CMK.",
						},
						"key_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of CMK.Available values include `Enabled`, `Disabled`, `PendingDelete`, `PendingImport`, `Archived`.",
						},
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Usage of CMK.Available values include `ENCRYPT_DECRYPT`, `ASYMMETRIC_DECRYPT_RSA_2048`, `ASYMMETRIC_DECRYPT_SM2`, `ASYMMETRIC_SIGN_VERIFY_SM2`, `ASYMMETRIC_SIGN_VERIFY_RSA_2048`, `ASYMMETRIC_SIGN_VERIFY_ECC`.",
						},
						"creator_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Uin of CMK Creator.",
						},
						"key_rotation_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specify whether to enable key rotation.",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creator of CMK.",
						},
						"next_rotate_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next rotate time of CMK when key_rotation_enabled is true.",
						},
						"deletion_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Delete time of CMK.",
						},
						"origin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin of CMK.`TENCENT_KMS` - CMK created by KMS, `EXTERNAL` - CMK imported by user.",
						},
						"valid_to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Valid when Origin is EXTERNAL, it means the effective date of the key material.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_key.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	param := make(map[string]interface{})
	if v, ok := d.GetOk("role"); ok {
		param["role"] = v.(int)
	}
	if v, ok := d.GetOk("order_type"); ok {
		param["order_type"] = v.(int)
	}
	if v, ok := d.GetOk("key_state"); ok {
		keyState := v.(string)
		param["key_state"] = KMS_KEY_STATE_MAP[keyState]
	}
	if v, ok := d.GetOk("search_key_alias"); ok {
		param["search_key_alias"] = v.(string)
	}
	if v, ok := d.GetOk("origin"); ok {
		param["origin"] = v.(string)
	}
	if v, ok := d.GetOk("key_usage"); ok {
		param["key_usage"] = v.(string)
	}
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		param["tag_filter"] = tags
	}

	kmsService := KmsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var keys []*kms.KeyMetadata
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := kmsService.DescribeKeysByFilter(ctx, param)
		if e != nil {
			return retryError(e)
		}
		keys = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read KMS keys failed, reason:%+v", logId, err)
		return err
	}
	keyList := make([]map[string]interface{}, 0, len(keys))
	ids := make([]string, 0, len(keys))
	for _, key := range keys {
		mapping := map[string]interface{}{
			"key_id":               key.KeyId,
			"alias":                key.Alias,
			"create_time":          helper.FormatUnixTime(*key.CreateTime),
			"description":          key.Description,
			"key_state":            key.KeyState,
			"key_usage":            key.KeyUsage,
			"creator_uin":          key.CreatorUin,
			"key_rotation_enabled": key.KeyRotationEnabled,
			"owner":                key.Owner,
			"origin":               key.Origin,
		}
		if *key.KeyRotationEnabled {
			mapping["next_rotate_time"] = helper.FormatUnixTime(*key.NextRotateTime)
		}
		if *key.KeyState == KMS_KEY_STATE_PENDINGDELETE {
			mapping["deletion_date"] = helper.FormatUnixTime(*key.DeletionDate)
		}
		if *key.Origin == KMS_ORIGIN_EXTERNAL {
			if *key.ValidTo != 0 {
				mapping["valid_to"] = helper.FormatUnixTime(*key.ValidTo)
			} else {
				mapping["valid_to"] = "never expire"
			}
		}
		keyList = append(keyList, mapping)
		ids = append(ids, *key.KeyId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("key_list", keyList); e != nil {
		log.Printf("[CRITAL]%s provider set KMS key list fail, reason:%+v", logId, e)
		return e
	}
	return nil
}
