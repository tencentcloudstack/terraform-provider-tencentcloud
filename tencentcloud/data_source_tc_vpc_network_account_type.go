/*
Use this data source to query detailed information of vpc network_account_type

Example Usage

```hcl
data "tencentcloud_vpc_network_account_type" "network_account_type" {
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

func dataSourceTencentCloudVpcNetworkAccountType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcNetworkAccountTypeRead,
		Schema: map[string]*schema.Schema{
			"network_account_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The network type of the user account, STANDARD is a standard user, LEGACY is a traditional user.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudVpcNetworkAccountTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_network_account_type.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcNetworkAccountTypeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		networkAccountType = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(networkAccountType))
	if networkAccountType != nil {
		_ = d.Set("network_account_type", networkAccountType)
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
