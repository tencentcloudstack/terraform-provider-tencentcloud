/*
Use this data source to query detailed information of dc internet_address_quota

Example Usage

```hcl
data "tencentcloud_dc_internet_address_quota" "internet_address_quota" {
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

func dataSourceTencentCloudDcInternetAddressQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcInternetAddressQuotaRead,
		Schema: map[string]*schema.Schema{
			"ipv6_prefix_len": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The minimum prefix length allowed on the IPv6 Internet public network.",
			},

			"ipv4_bgp_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "BGP type IPv4 Internet address quota.",
			},

			"ipv4_other_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Non-BGP type IPv4 Internet address quota.",
			},

			"ipv4_bgp_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of used BGP type IPv4 Internet addresses.",
			},

			"ipv4_other_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of non-BGP Internet addresses used.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDcInternetAddressQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dc_internet_address_quota.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcInternetAddressQuotaByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		ipv6PrefixLen = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ipv6PrefixLen))
	if ipv6PrefixLen != nil {
		_ = d.Set("ipv6_prefix_len", ipv6PrefixLen)
	}

	if ipv4BgpQuota != nil {
		_ = d.Set("ipv4_bgp_quota", ipv4BgpQuota)
	}

	if ipv4OtherQuota != nil {
		_ = d.Set("ipv4_other_quota", ipv4OtherQuota)
	}

	if ipv4BgpNum != nil {
		_ = d.Set("ipv4_bgp_num", ipv4BgpNum)
	}

	if ipv4OtherNum != nil {
		_ = d.Set("ipv4_other_num", ipv4OtherNum)
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
