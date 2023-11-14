/*
Use this data source to query detailed information of kms describe_white_box_decrypt_key

Example Usage

```hcl
data "tencentcloud_kms_describe_white_box_decrypt_key" "describe_white_box_decrypt_key" {
  key_id = "244dab8c-6dad-11ea-80c6-5254006d0810"
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

func dataSourceTencentCloudKmsDescribeWhiteBoxDecryptKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsDescribeWhiteBoxDecryptKeyRead,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Globally unique identifier for the white box key.",
			},

			"decrypt_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "White box decryption key, base64 encoded.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKmsDescribeWhiteBoxDecryptKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_describe_white_box_decrypt_key.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_id"); ok {
		paramMap["KeyId"] = helper.String(v.(string))
	}

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsDescribeWhiteBoxDecryptKeyByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		decryptKey = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(decryptKey))
	if decryptKey != nil {
		_ = d.Set("decrypt_key", decryptKey)
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
