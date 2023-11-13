/*
Use this data source to query detailed information of mariadb renewal_price

Example Usage

```hcl
data "tencentcloud_mariadb_renewal_price" "renewal_price" {
  instance_id = ""
  period =
  amount_unit = ""
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

func dataSourceTencentCloudMariadbRenewalPrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbRenewalPriceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Renewal duration, default: 1 month.",
			},

			"amount_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Price unit. Valid values: `* pent` (cent), `* microPent` (microcent).",
			},

			"original_price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Original price * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).",
			},

			"price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The actual price may be different from the original price due to discounts. * Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. * Currency: CNY (Chinese site), USD (international site).",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMariadbRenewalPriceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_renewal_price.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("period"); v != nil {
		paramMap["Period"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("amount_unit"); ok {
		paramMap["AmountUnit"] = helper.String(v.(string))
	}

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbRenewalPriceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		originalPrice = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(originalPrice))
	if originalPrice != nil {
		_ = d.Set("original_price", originalPrice)
	}

	if price != nil {
		_ = d.Set("price", price)
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
