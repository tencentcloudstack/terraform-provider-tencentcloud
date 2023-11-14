/*
Use this data source to query detailed information of kms describe_white_box_key_details

Example Usage

```hcl
data "tencentcloud_kms_describe_white_box_key_details" "describe_white_box_key_details" {
  key_status =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKmsDescribeWhiteBoxKeyDetails() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsDescribeWhiteBoxKeyDetailsRead,
		Schema: map[string]*schema.Schema{
			"key_status": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Filter condition: status of the key, 0: disabled, 1: enabled.",
			},

			"key_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of white box key information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of algorithm used by the key.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Key creation time, Unix timestamp.",
						},
						"decrypt_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "White box decryption key, base64 encoded.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID, format: creatorUin/$creatorUin/$keyId.",
						},
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Globally unique identifier for the white box key.",
						},
						"creator_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creator.",
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "As an alias for a key that is easier to identify and easier to understand, it cannot be empty and is a combination of 1-60 alphanumeric characters - _. The first character must be a letter or number. It cannot be repeated.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the key.",
						},
						"encrypt_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "White box encryption key, base64 encoded.",
						},
						"owner_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creator.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the white box key, the value is: Enabled | Disabled.",
						},
						"device_fingerprint_bind": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is there a device fingerprint bound to the current key?.",
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

func dataSourceTencentCloudKmsDescribeWhiteBoxKeyDetailsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_describe_white_box_key_details.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("key_status"); v != nil {
		paramMap["KeyStatus"] = helper.IntInt64(v.(int))
	}

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsDescribeWhiteBoxKeyDetailsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		totalCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(totalCount))
	if keyInfos != nil {
		for _, whiteboxKeyInfo := range keyInfos {
			whiteboxKeyInfoMap := map[string]interface{}{}

			if whiteboxKeyInfo.Algorithm != nil {
				whiteboxKeyInfoMap["algorithm"] = whiteboxKeyInfo.Algorithm
			}

			if whiteboxKeyInfo.CreateTime != nil {
				whiteboxKeyInfoMap["create_time"] = whiteboxKeyInfo.CreateTime
			}

			if whiteboxKeyInfo.DecryptKey != nil {
				whiteboxKeyInfoMap["decrypt_key"] = whiteboxKeyInfo.DecryptKey
			}

			if whiteboxKeyInfo.ResourceId != nil {
				whiteboxKeyInfoMap["resource_id"] = whiteboxKeyInfo.ResourceId
			}

			if whiteboxKeyInfo.KeyId != nil {
				whiteboxKeyInfoMap["key_id"] = whiteboxKeyInfo.KeyId
			}

			if whiteboxKeyInfo.CreatorUin != nil {
				whiteboxKeyInfoMap["creator_uin"] = whiteboxKeyInfo.CreatorUin
			}

			if whiteboxKeyInfo.Alias != nil {
				whiteboxKeyInfoMap["alias"] = whiteboxKeyInfo.Alias
			}

			if whiteboxKeyInfo.Description != nil {
				whiteboxKeyInfoMap["description"] = whiteboxKeyInfo.Description
			}

			if whiteboxKeyInfo.EncryptKey != nil {
				whiteboxKeyInfoMap["encrypt_key"] = whiteboxKeyInfo.EncryptKey
			}

			if whiteboxKeyInfo.OwnerUin != nil {
				whiteboxKeyInfoMap["owner_uin"] = whiteboxKeyInfo.OwnerUin
			}

			if whiteboxKeyInfo.Status != nil {
				whiteboxKeyInfoMap["status"] = whiteboxKeyInfo.Status
			}

			if whiteboxKeyInfo.DeviceFingerprintBind != nil {
				whiteboxKeyInfoMap["device_fingerprint_bind"] = whiteboxKeyInfo.DeviceFingerprintBind
			}

			ids = append(ids, *whiteboxKeyInfo.KeyId)
			tmpList = append(tmpList, whiteboxKeyInfoMap)
		}

		_ = d.Set("key_infos", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
