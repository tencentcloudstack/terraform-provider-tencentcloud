/*
Use this data source to query detailed information of kms key_lists

Example Usage

```hcl
data "tencentcloud_kms_describe_keys" "example" {
  key_ids = [
    "9ffacc8b-6461-11ee-a54e-525400dd8a7d",
    "bffae4ed-6465-11ee-90b2-5254000ef00e"
  ]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKmsDescribeKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsDescribeKeysRead,
		Schema: map[string]*schema.Schema{
			"key_ids": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Query the ID list of CMK, batch query supports up to 100 KeyIds at a time.",
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
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKmsDescribeKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_describe_keys.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		service     = KmsService{client: meta.(*TencentCloudClient).apiV3Conn}
		keyMetadata []*kms.KeyMetadata
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_ids"); ok {
		keyIdsSet := v.(*schema.Set).List()
		paramMap["KeyIds"] = helper.InterfacesStringsPoint(keyIdsSet)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsKeyListsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		keyMetadata = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(keyMetadata))
	tmpList := make([]map[string]interface{}, 0, len(keyMetadata))

	if keyMetadata != nil {
		for _, key := range keyMetadata {
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

			tmpList = append(tmpList, mapping)
			ids = append(ids, *key.KeyId)
		}

		_ = d.Set("key_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
