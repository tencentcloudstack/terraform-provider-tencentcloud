package kms

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKmsWhiteBoxDeviceFingerprints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsWhiteBoxDeviceFingerprintsRead,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Globally unique identifier for the white box key.",
			},
			"list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Device fingerprint list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "identity.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
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

func dataSourceTencentCloudKmsWhiteBoxDeviceFingerprintsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_white_box_device_fingerprints.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		deviceFingerprints []*kms.DeviceFingerprint
		keyId              string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_id"); ok {
		paramMap["KeyId"] = helper.String(v.(string))
		keyId = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsWhiteBoxDeviceFingerprintsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		deviceFingerprints = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(deviceFingerprints))
	if deviceFingerprints != nil {
		for _, item := range deviceFingerprints {
			itemMap := map[string]interface{}{}

			if item.Identity != nil {
				itemMap["identity"] = item.Identity
			}

			if item.Description != nil {
				itemMap["description"] = item.Description
			}

			tmpList = append(tmpList, itemMap)
		}

		_ = d.Set("list", tmpList)
	}

	d.SetId(keyId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
