/*
Use this data source to query detailed information of kms white_box_device_fingerprints

Example Usage

```hcl
data "tencentcloud_kms_white_box_device_fingerprints" "example" {
  key_id = "244dab8c-6dad-11ea-80c6-5254006d0810"
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

func dataSourceTencentCloudKmsWhiteBoxDeviceFingerprints() *schema.Resource {
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
	defer logElapsed("data_source.tencentcloud_kms_white_box_device_fingerprints.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId              = getLogId(contextNil)
		ctx                = context.WithValue(context.TODO(), logIdKey, logId)
		service            = KmsService{client: meta.(*TencentCloudClient).apiV3Conn}
		deviceFingerprints []*kms.DeviceFingerprint
		keyId              string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_id"); ok {
		paramMap["KeyId"] = helper.String(v.(string))
		keyId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsWhiteBoxDeviceFingerprintsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
