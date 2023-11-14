/*
Use this data source to query detailed information of kms get_public_key

Example Usage

```hcl
data "tencentcloud_kms_get_public_key" "get_public_key" {
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

func dataSourceTencentCloudKmsGetPublicKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsGetPublicKeyRead,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "CMK unique identifier.",
			},

			"public_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Base64-encoded public key content.",
			},

			"public_key_pem": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Public key content in PEM format.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKmsGetPublicKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_get_public_key.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsGetPublicKeyByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		keyId = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(keyId))
	if keyId != nil {
		_ = d.Set("key_id", keyId)
	}

	if publicKey != nil {
		_ = d.Set("public_key", publicKey)
	}

	if publicKeyPem != nil {
		_ = d.Set("public_key_pem", publicKeyPem)
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
