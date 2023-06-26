/*
Use this data source to query detailed information of vpn default_health_check_ip

Example Usage

```hcl
data "tencentcloud_vpn_default_health_check_ip" "default_health_check_ip" {
  vpn_gateway_id = "vpngw-gt8bianl"
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

func dataSourceTencentCloudVpnDefaultHealthCheckIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpnDefaultHealthCheckIpRead,
		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "vpn gateway id.",
			},

			"health_check_local_ip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "local ip of health check.",
			},

			"health_check_remote_ip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "remote ip for health check.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudVpnDefaultHealthCheckIpRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpn_default_health_check_ip.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var vpnGwId string
	res := make(map[string]interface{})

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		vpnGwId = v.(string)
		paramMap["VpnGatewayId"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var defaultHealthCheck *vpc.GenerateVpnConnectionDefaultHealthCheckIpResponseParams

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpnDefaultHealthCheckIp(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		defaultHealthCheck = result
		return nil
	})
	if err != nil {
		return err
	}

	if defaultHealthCheck.HealthCheckLocalIp != nil {
		_ = d.Set("health_check_local_ip", defaultHealthCheck.HealthCheckLocalIp)
		res["health_check_local_ip"] = defaultHealthCheck.HealthCheckLocalIp
	}

	if defaultHealthCheck.HealthCheckRemoteIp != nil {
		_ = d.Set("health_check_remote_ip", defaultHealthCheck.HealthCheckRemoteIp)
		res["health_check_remote_ip"] = defaultHealthCheck.HealthCheckRemoteIp
	}

	d.SetId(vpnGwId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), res); e != nil {
			return e
		}
	}
	return nil
}
