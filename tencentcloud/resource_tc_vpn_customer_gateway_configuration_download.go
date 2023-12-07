package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpnCustomerGatewayConfigurationDownload() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnCustomerGatewayConfigurationDownloadCreate,
		Read:   resourceTencentCloudVpnCustomerGatewayConfigurationDownloadRead,
		Delete: resourceTencentCloudVpnCustomerGatewayConfigurationDownloadDelete,
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

			"customer_gateway_configuration": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "xml configuration.",
			},
		},
	}
}

func resourceTencentCloudVpnCustomerGatewayConfigurationDownloadCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpn_customer_gateway_configuration_download.read")()
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
		return nil
	}

	d.SetId(vpnGatewayId + FILED_SP + vpnConnectionId)

	_ = d.Set("customer_gateway_configuration", response.Response.CustomerGatewayConfiguration)

	return resourceTencentCloudVpnCustomerGatewayConfigurationDownloadRead(d, meta)
}

func resourceTencentCloudVpnCustomerGatewayConfigurationDownloadRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_customer_gateway_configuration_download.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpnCustomerGatewayConfigurationDownloadDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_customer_gateway_configuration_download.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
