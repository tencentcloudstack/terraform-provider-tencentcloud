/*
Use this data source to query detailed information of vpc tenant_ccn

Example Usage

```hcl
data "tencentcloud_vpc_tenant_ccn" "tenant_ccn" {
  name = "ccn-ids"
  values =
  offset = 0
  limit = 20
      }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcTenantCcn() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcTenantCcnRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Attribute name, if there are multiple Filters, the relationship between Filters is logical AND (AND) relationship.",
			},

			"values": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Attribute value, if there are multiple Values in the same Filter, the relationship between Values under the same Filter is a logical OR (OR) relationship.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset, default 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The amount of data returned in a single page, the optional value is an integer between 0 and 100, and the default is 20.",
			},

			"ccn_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "CCN objects.Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of objects eligible.",
			},

			"request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Unique request ID, returned for every request. The RequestId of the request needs to be provided when locating the problem.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudVpcTenantCcnRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_tenant_ccn.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		paramMap["Name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("values"); ok {
		valuesSet := v.(*schema.Set).List()
		paramMap["Values"] = helper.InterfacesStringsPoint(valuesSet)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var ccnSet []*vpc.CcnInstanceInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcTenantCcnByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		ccnSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ccnSet))
	tmpList := make([]map[string]interface{}, 0, len(ccnSet))

	if ccnSet != nil {
		_ = d.Set("ccn_set", ccnSet)
	}

	if totalCount != nil {
		_ = d.Set("total_count", totalCount)
	}

	if requestId != nil {
		_ = d.Set("request_id", requestId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
