/*
Use this data source to query detailed information of vpc vpn_customer_gateway_vendors

Example Usage

```hcl
data "tencentcloud_vpn_customer_gateway_vendors" "vpn_customer_gateway_vendors" {}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpnCustomerGatewayVendors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcVpnCustomerGatewayVendorsRead,
		Schema: map[string]*schema.Schema{
			"customer_gateway_vendor_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Customer Gateway Vendor Set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform.",
						},
						"software_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SoftwareVersion.",
						},
						"vendor_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VendorName.",
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

func dataSourceTencentCloudVpcVpnCustomerGatewayVendorsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_vpn_customer_gateway_vendors.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var customerGatewayVendorSet []*vpc.CustomerGatewayVendor

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpnCustomerGatewayVendors(ctx)
		if e != nil {
			return retryError(e)
		}
		customerGatewayVendorSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(customerGatewayVendorSet))
	tmpList := make([]map[string]interface{}, 0, len(customerGatewayVendorSet))

	if customerGatewayVendorSet != nil {
		for _, customerGatewayVendor := range customerGatewayVendorSet {
			customerGatewayVendorMap := map[string]interface{}{}

			if customerGatewayVendor.Platform != nil {
				customerGatewayVendorMap["platform"] = customerGatewayVendor.Platform
			}

			if customerGatewayVendor.SoftwareVersion != nil {
				customerGatewayVendorMap["software_version"] = customerGatewayVendor.SoftwareVersion
			}

			if customerGatewayVendor.VendorName != nil {
				customerGatewayVendorMap["vendor_name"] = customerGatewayVendor.VendorName
			}

			ids = append(ids, *customerGatewayVendor.VendorName)
			tmpList = append(tmpList, customerGatewayVendorMap)
		}

		_ = d.Set("customer_gateway_vendor_set", tmpList)
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
