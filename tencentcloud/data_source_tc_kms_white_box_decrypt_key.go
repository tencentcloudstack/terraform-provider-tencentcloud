/*
Use this data source to query detailed information of kms white_box_decrypt_key

Example Usage

```hcl
data "tencentcloud_kms_white_box_decrypt_key" "example" {
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

func dataSourceTencentCloudKmsWhiteBoxDecryptKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsWhiteBoxDecryptKeyRead,
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

func dataSourceTencentCloudKmsWhiteBoxDecryptKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_white_box_decrypt_key.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId              = getLogId(contextNil)
		ctx                = context.WithValue(context.TODO(), logIdKey, logId)
		service            = KmsService{client: meta.(*TencentCloudClient).apiV3Conn}
		whiteBoxDecryptKey *kms.DescribeWhiteBoxDecryptKeyResponseParams
		keyId              string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_id"); ok {
		paramMap["KeyId"] = helper.String(v.(string))
		keyId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsWhiteBoxDecryptKeyByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		whiteBoxDecryptKey = result
		return nil
	})

	if err != nil {
		return err
	}

	if whiteBoxDecryptKey.DecryptKey != nil {
		_ = d.Set("decrypt_key", whiteBoxDecryptKey.DecryptKey)
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
