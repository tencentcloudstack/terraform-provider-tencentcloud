/*
Provides a resource to create a vpc vpn_gateway_ssl_client_cert

Example Usage

```hcl
resource "tencentcloud_vpc_vpn_gateway_ssl_client_cert" "vpn_gateway_ssl_client_cert" {
  ssl_vpn_client_id = "vpnc-123456"
}
```

Import

vpc vpn_gateway_ssl_client_cert can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert vpn_gateway_ssl_client_cert_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudVpcVpnGatewaySslClientCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcVpnGatewaySslClientCertCreate,
		Read:   resourceTencentCloudVpcVpnGatewaySslClientCertRead,
		Delete: resourceTencentCloudVpcVpnGatewaySslClientCertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ssl_vpn_client_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "SSL-VPN-CLIENT Instance ID.",
			},
		},
	}
}

func resourceTencentCloudVpcVpnGatewaySslClientCertCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_gateway_ssl_client_cert.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = vpc.NewEnableVpnGatewaySslClientCertRequest()
		response       = vpc.NewEnableVpnGatewaySslClientCertResponse()
		sslVpnClientId string
	)
	if v, ok := d.GetOk("ssl_vpn_client_id"); ok {
		sslVpnClientId = v.(string)
		request.SslVpnClientId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().EnableVpnGatewaySslClientCert(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc vpnGatewaySslClientCert failed, reason:%+v", logId, err)
		return err
	}

	sslVpnClientId = *response.Response.SslVpnClientId
	d.SetId(sslVpnClientId)

	return resourceTencentCloudVpcVpnGatewaySslClientCertRead(d, meta)
}

func resourceTencentCloudVpcVpnGatewaySslClientCertRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_gateway_ssl_client_cert.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	vpnGatewaySslClientCertId := d.Id()

	vpnGatewaySslClientCert, err := service.DescribeVpcVpnGatewaySslClientCertById(ctx, sslVpnClientId)
	if err != nil {
		return err
	}

	if vpnGatewaySslClientCert == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcVpnGatewaySslClientCert` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if vpnGatewaySslClientCert.SslVpnClientId != nil {
		_ = d.Set("ssl_vpn_client_id", vpnGatewaySslClientCert.SslVpnClientId)
	}

	return nil
}

func resourceTencentCloudVpcVpnGatewaySslClientCertDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_gateway_ssl_client_cert.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	vpnGatewaySslClientCertId := d.Id()

	if err := service.DeleteVpcVpnGatewaySslClientCertById(ctx, sslVpnClientId); err != nil {
		return err
	}

	return nil
}
