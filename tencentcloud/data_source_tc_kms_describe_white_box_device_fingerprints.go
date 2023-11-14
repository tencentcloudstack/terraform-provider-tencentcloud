/*
Use this data source to query detailed information of kms describe_white_box_device_fingerprints

Example Usage

```hcl
data "tencentcloud_kms_describe_white_box_device_fingerprints" "describe_white_box_device_fingerprints" {
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

func dataSourceTencentCloudKmsDescribeWhiteBoxDeviceFingerprints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsDescribeWhiteBoxDeviceFingerprintsRead,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Globally unique identifier for the white box key.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKmsDescribeWhiteBoxDeviceFingerprintsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_describe_white_box_device_fingerprints.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_id"); ok {
		paramMap["KeyId"] = helper.String(v.(string))
	}

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var deviceFingerprints []*kms.DeviceFingerprint

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsDescribeWhiteBoxDeviceFingerprintsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		deviceFingerprints = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(deviceFingerprints))
	tmpList := make([]map[string]interface{}, 0, len(deviceFingerprints))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
