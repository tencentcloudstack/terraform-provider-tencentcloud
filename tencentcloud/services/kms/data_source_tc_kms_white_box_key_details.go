package kms

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKmsWhiteBoxKeyDetails() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsWhiteBoxKeyDetailsRead,
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

func dataSourceTencentCloudKmsWhiteBoxKeyDetailsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_white_box_key_details.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		whiteBoxKeyInfo []*kms.WhiteboxKeyInfo
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("key_status"); v != nil {
		paramMap["KeyStatus"] = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsWhiteBoxKeyDetailsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		whiteBoxKeyInfo = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(whiteBoxKeyInfo))
	tmpList := make([]map[string]interface{}, 0, len(whiteBoxKeyInfo))

	if whiteBoxKeyInfo != nil {
		for _, whiteBoxKey := range whiteBoxKeyInfo {
			whiteBoxKeyInfoMap := map[string]interface{}{}

			if whiteBoxKey.Algorithm != nil {
				whiteBoxKeyInfoMap["algorithm"] = whiteBoxKey.Algorithm
			}

			if whiteBoxKey.CreateTime != nil {
				whiteBoxKeyInfoMap["create_time"] = whiteBoxKey.CreateTime
			}

			if whiteBoxKey.DecryptKey != nil {
				whiteBoxKeyInfoMap["decrypt_key"] = whiteBoxKey.DecryptKey
			}

			if whiteBoxKey.ResourceId != nil {
				whiteBoxKeyInfoMap["resource_id"] = whiteBoxKey.ResourceId
			}

			if whiteBoxKey.KeyId != nil {
				whiteBoxKeyInfoMap["key_id"] = whiteBoxKey.KeyId
			}

			if whiteBoxKey.CreatorUin != nil {
				whiteBoxKeyInfoMap["creator_uin"] = whiteBoxKey.CreatorUin
			}

			if whiteBoxKey.Alias != nil {
				whiteBoxKeyInfoMap["alias"] = whiteBoxKey.Alias
			}

			if whiteBoxKey.Description != nil {
				whiteBoxKeyInfoMap["description"] = whiteBoxKey.Description
			}

			if whiteBoxKey.EncryptKey != nil {
				whiteBoxKeyInfoMap["encrypt_key"] = whiteBoxKey.EncryptKey
			}

			if whiteBoxKey.OwnerUin != nil {
				whiteBoxKeyInfoMap["owner_uin"] = whiteBoxKey.OwnerUin
			}

			if whiteBoxKey.Status != nil {
				whiteBoxKeyInfoMap["status"] = whiteBoxKey.Status
			}

			if whiteBoxKey.DeviceFingerprintBind != nil {
				whiteBoxKeyInfoMap["device_fingerprint_bind"] = whiteBoxKey.DeviceFingerprintBind
			}

			ids = append(ids, *whiteBoxKey.KeyId)
			tmpList = append(tmpList, whiteBoxKeyInfoMap)
		}

		_ = d.Set("key_infos", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
