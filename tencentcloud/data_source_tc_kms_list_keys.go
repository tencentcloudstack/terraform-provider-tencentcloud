/*
Use this data source to query detailed information of kms list_keys

Example Usage

```hcl
data "tencentcloud_kms_list_keys" "list_keys" {
  offset = 0
  limit = 2
  role =
  hsm_cluster_id = "0"
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

func dataSourceTencentCloudKmsListKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsListKeysRead,
		Schema: map[string]*schema.Schema{
			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The meaning is the same as the Offset of SQL query, which means that this acquisition starts from the Offset-th element of the array arranged in a certain order. The default is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The meaning is the same as Limit in SQL query, which means that up to Limit elements can be obtained this time. The default value is 10, the maximum value is 200.",
			},

			"role": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Filter based on the creator role. The default value is 0, which indicates the cmk created by the user himself, and 1, which indicates the cmk automatically created by authorizing other cloud products.",
			},

			"hsm_cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "HSM cluster ID (only valid for KMS exclusive/managed service instances).",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKmsListKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_list_keys.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("role"); v != nil {
		paramMap["Role"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("hsm_cluster_id"); ok {
		paramMap["HsmClusterId"] = helper.String(v.(string))
	}

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsListKeysByFilter(ctx, paramMap)
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
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
