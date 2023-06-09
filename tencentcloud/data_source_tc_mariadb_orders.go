/*
Use this data source to query detailed information of mariadb orders

Example Usage

```hcl
data "tencentcloud_mariadb_orders" "orders" {
  deal_name = "20230607164033835942781"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
)

func dataSourceTencentCloudMariadbOrders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbOrdersRead,
		Schema: map[string]*schema.Schema{
			"deal_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "List of long order numbers to be queried, which are returned for the APIs for creating, renewing, or scaling instances.",
			},
			"deals": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Order information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deal_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Order number.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of items.",
						},
						"flow_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the associated process, which can be used to query the process execution status.",
						},
						"instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The ID of the created instance, which is required only for the order that creates an instance.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Payment mode. Valid values: 0 (postpaid), 1 (prepaid).",
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

func dataSourceTencentCloudMariadbOrdersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_orders.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		deals    []*mariadb.Deal
		dealName string
	)

	if v, ok := d.GetOk("deal_name"); ok {
		dealName = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbOrdersByFilter(ctx, dealName)
		if e != nil {
			return retryError(e)
		}

		deals = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(deals))

	if deals != nil {
		for _, deal := range deals {
			dealMap := map[string]interface{}{}

			if deal.DealName != nil {
				dealMap["deal_name"] = deal.DealName
			}

			if deal.OwnerUin != nil {
				dealMap["owner_uin"] = deal.OwnerUin
			}

			if deal.Count != nil {
				dealMap["count"] = deal.Count
			}

			if deal.FlowId != nil {
				dealMap["flow_id"] = deal.FlowId
			}

			if deal.InstanceIds != nil {
				dealMap["instance_ids"] = deal.InstanceIds
			}

			if deal.PayMode != nil {
				dealMap["pay_mode"] = deal.PayMode
			}

			tmpList = append(tmpList, dealMap)
		}

		_ = d.Set("deals", tmpList)
	}

	d.SetId(dealName)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
