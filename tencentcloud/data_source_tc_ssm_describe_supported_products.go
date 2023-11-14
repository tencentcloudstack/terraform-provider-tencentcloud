/*
Use this data source to query detailed information of ssm describe_supported_products

Example Usage

```hcl
data "tencentcloud_ssm_describe_supported_products" "describe_supported_products" {
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

func dataSourceTencentCloudSsmDescribeSupportedProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmDescribeSupportedProductsRead,
		Schema: map[string]*schema.Schema{
			"products": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of supported products.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSsmDescribeSupportedProductsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssm_describe_supported_products.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := SsmService{client: meta.(*TencentCloudClient).apiV3Conn}

	var products []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmDescribeSupportedProductsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		products = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(products))
	if products != nil {
		_ = d.Set("products", products)
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
