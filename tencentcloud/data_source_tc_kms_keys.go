package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKmsKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsKeysRead,
		Schema: map[string]*schema.Schema{
			"role": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Filter by role of the CMK creator. `0` - created by user, `1` - created by cloud product. Default value is `0`.",
			},
			"order_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Order to sort the CMK create time. `0` - desc, `1` - asc. Default value is `0`.",
			},
			"key_state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Filter by state of CMK. `0` - all CMKs are queried, `1` - only Enabled CMKs are queried, `2` - only Disabled CMKs are queried, `3` - only PendingDelete CMKs are queried, `4` - only PendingImport CMKs are queried, `5` - only Archived CMKs are queried.",
			},
			"search_key_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Words used to match the results, and the words can be: key_id and alias.",
			},
			"origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     KMS_ORIGIN_ALL,
				Description: "Filter by origin of CMK. `TENCENT_KMS` - CMK created by KMS, `EXTERNAL` - CMK imported by user, `ALL` - all CMKs. Default value is `ALL`.",
			},
			"key_usage": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     KMS_KEY_USAGE_ENCRYPT_DECRYPT,
				Description: "Filter by usage of CMK. Available values include `ALL`, `ENCRYPT_DECRYPT`, `ASYMMETRIC_DECRYPT_RSA_2048`, `ASYMMETRIC_DECRYPT_SM2`, `ASYMMETRIC_SIGN_VERIFY_SM2`, `ASYMMETRIC_SIGN_VERIFY_RSA_2048`, `ASYMMETRIC_SIGN_VERIFY_ECC`. Default value is `ENCRYPT_DECRYPT`.",
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
							Type:        schema.TypeInt,
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
							Description: "State of CMK.",
						},
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Usage of CMK.",
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
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Next rotate time of CMK when key_rotation_enabled is true.",
						},
						"deletion_date": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Delete time of CMK.",
						},
						"origin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin of CMK. `TENCENT_KMS` - CMK created by KMS, `EXTERNAL` - CMK imported by user.",
						},
						"valid_to": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Valid when origin is `EXTERNAL`, it means the effective date of the key material.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudKmsKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_keys.read")()

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
		keyState := v.(int)
		param["key_state"] = uint64(keyState)
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
			"create_time":          key.CreateTime,
			"description":          key.Description,
			"key_state":            key.KeyState,
			"key_usage":            key.KeyUsage,
			"creator_uin":          key.CreatorUin,
			"key_rotation_enabled": key.KeyRotationEnabled,
			"owner":                key.Owner,
			"next_rotate_time":     key.NextRotateTime,
			"deletion_date":        key.DeletionDate,
			"origin":               key.Origin,
			"valid_to":             key.ValidTo,
		}

		keyList = append(keyList, mapping)
		ids = append(ids, *key.KeyId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("key_list", keyList); e != nil {
		log.Printf("[CRITAL]%s provider set KMS key list fail, reason:%+v", logId, e)
		return e
	}
	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), keyList)
	}
	return nil
}
