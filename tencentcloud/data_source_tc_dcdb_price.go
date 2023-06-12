/*
Use this data source to query detailed information of dcdb price

Example Usage

```hcl
data "tencentcloud_dcdb_price" "price" {
	instance_count   = 1
	zone             = var.default_az
	period           = 1
	shard_node_count = 2
	shard_memory     = 2
	shard_storage    = 10
	shard_count      = 2
	paymode          = "postpaid"
	amount_unit      = "pent"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbPrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbPriceRead,
		Schema: map[string]*schema.Schema{
			"instance_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The count of instances wants to buy.",
			},

			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "AZ ID of the purchased instance.",
			},

			"period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Purchase period in months.",
			},

			"shard_node_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of instance shard nodes.",
			},

			"shard_memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Shard memory size in GB.",
			},

			"shard_storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Shard storage capacity in GB.",
			},

			"shard_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of instance shards.",
			},

			"paymode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Billing type. Valid values: `postpaid` (pay-as-you-go), `prepaid` (monthly subscription).",
			},

			"amount_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Price unit. Valid values: `pent` (cent), `microPent` (microcent).",
			},

			"original_price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Original price. Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. Currency: CNY (Chinese site), USD (international site).",
			},

			"price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The actual price may be different from the original price due to discounts. Unit: Cent (default). If the request parameter contains `AmountUnit`, see `AmountUnit` description. Currency: CNY (Chinese site), USD (international site).",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDcdbPriceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_price.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		ids []string
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("instance_count"); v != nil {
		paramMap["InstanceCount"] = helper.IntInt64(v.(int))
		ids = append(ids, helper.IntToStr(v.(int)))
	}

	if v, ok := d.GetOk("zone"); ok {
		paramMap["Zone"] = helper.String(v.(string))
		ids = append(ids, v.(string))
	}

	if v, _ := d.GetOk("shard_count"); v != nil {
		paramMap["ShardCount"] = helper.IntInt64(v.(int))
		ids = append(ids, helper.IntToStr(v.(int)))
	}

	if v, _ := d.GetOk("shard_node_count"); v != nil {
		paramMap["ShardNodeCount"] = helper.IntInt64(v.(int))
		ids = append(ids, helper.IntToStr(v.(int)))
	}

	if v, _ := d.GetOk("shard_memory"); v != nil {
		paramMap["ShardMemory"] = helper.IntInt64(v.(int))
		ids = append(ids, helper.IntToStr(v.(int)))
	}

	if v, _ := d.GetOk("shard_storage"); v != nil {
		paramMap["ShardStorage"] = helper.IntInt64(v.(int))
		ids = append(ids, helper.IntToStr(v.(int)))
	}

	if v, _ := d.GetOk("period"); v != nil {
		paramMap["Period"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("paymode"); ok {
		paramMap["Paymode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("amount_unit"); ok {
		paramMap["AmountUnit"] = helper.String(v.(string))
	}

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *dcdb.DescribeDCDBPriceResponseParams
	var e error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e = service.DescribeDcdbPriceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result != nil {
		if result.OriginalPrice != nil {
			_ = d.Set("original_price", result.OriginalPrice)
		}

		if result.Price != nil {
			_ = d.Set("price", result.Price)
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
