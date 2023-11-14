/*
Provides a resource to create a vpc vpn_customer_gateway_configuration_download

Example Usage

```hcl
resource "tencentcloud_vpc_vpn_customer_gateway_configuration_download" "vpn_customer_gateway_configuration_download" {
  vpn_gateway_id = "vpngw-c6orbuv7"
  vpn_connection_id = "vpnx-osftvdea"
  customer_gateway_vendor {
		platform = "comware"
		software_version = "V1.0"
		vendor_name = "h3c"

  }
  interface_name = ""
}
```

Import

vpc vpn_customer_gateway_configuration_download can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_vpn_customer_gateway_configuration_download.vpn_customer_gateway_configuration_download vpn_customer_gateway_configuration_download_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudVpcVpnCustomerGatewayConfigurationDownload() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcVpnCustomerGatewayConfigurationDownloadCreate,
		Read:   resourceTencentCloudVpcVpnCustomerGatewayConfigurationDownloadRead,
		Delete: resourceTencentCloudVpcVpnCustomerGatewayConfigurationDownloadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPN Gateway Instance ID.",
			},

			"vpn_connection_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPN Connection Instance id.",
			},

			"customer_gateway_vendor": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Customer Gateway Vendor Info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"platform": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Platform.",
						},
						"software_version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "SoftwareVersion.",
						},
						"vendor_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VendorName.",
						},
					},
				},
			},

			"interface_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPN connection access device physical interface name.",
			},
		},
	}
}

func resourceTencentCloudVpcVpnCustomerGatewayConfigurationDownloadCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_customer_gateway_configuration_download.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = vpc.NewDownloadCustomerGatewayConfigurationRequest()
		response        = vpc.NewDownloadCustomerGatewayConfigurationResponse()
		vpnGatewayId    string
		vpnConnectionId string
	)
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		vpnGatewayId = v.(string)
		request.VpnGatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpn_connection_id"); ok {
		vpnConnectionId = v.(string)
		request.VpnConnectionId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "customer_gateway_vendor"); ok {
		customerGatewayVendor := vpc.CustomerGatewayVendor{}
		if v, ok := dMap["platform"]; ok {
			customerGatewayVendor.Platform = helper.String(v.(string))
		}
		if v, ok := dMap["software_version"]; ok {
			customerGatewayVendor.SoftwareVersion = helper.String(v.(string))
		}
		if v, ok := dMap["vendor_name"]; ok {
			customerGatewayVendor.VendorName = helper.String(v.(string))
		}
		request.CustomerGatewayVendor = &customerGatewayVendor
	}

	if v, ok := d.GetOk("interface_name"); ok {
		request.InterfaceName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DownloadCustomerGatewayConfiguration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc vpnCustomerGatewayConfigurationDownload failed, reason:%+v", logId, err)
		return err
	}

	vpnGatewayId = *response.Response.VpnGatewayId
	d.SetId(strings.Join([]string{vpnGatewayId, vpnConnectionId}, FILED_SP))

	return resourceTencentCloudVpcVpnCustomerGatewayConfigurationDownloadRead(d, meta)
}

func resourceTencentCloudVpcVpnCustomerGatewayConfigurationDownloadRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_customer_gateway_configuration_download.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcVpnCustomerGatewayConfigurationDownloadDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_customer_gateway_configuration_download.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
