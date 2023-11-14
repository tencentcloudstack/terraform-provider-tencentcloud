/*
Use this data source to query detailed information of kms describe_keys

Example Usage

```hcl
data "tencentcloud_kms_describe_keys" "describe_keys" {
  key_ids =
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
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Query the ID list of CMK, batch query supports up to 100 KeyIds at a time.",
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

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_ids"); ok {
		keyIdsSet := v.(*schema.Set).List()
		paramMap["KeyIds"] = helper.InterfacesStringsPoint(keyIdsSet)
	}

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var keyMetadata []*kms.KeyMetadata

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsDescribeKeysByFilter(ctx, paramMap)
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

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
