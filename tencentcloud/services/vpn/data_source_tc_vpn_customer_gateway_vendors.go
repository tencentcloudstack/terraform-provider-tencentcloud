package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpnCustomerGatewayVendors() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_vpn_customer_gateway_vendors.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var customerGatewayVendorSet []*vpc.CustomerGatewayVendor

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpnCustomerGatewayVendors(ctx)
		if e != nil {
			return tccommon.RetryError(e)
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
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
