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

func DataSourceTencentCloudVpnDefaultHealthCheckIp() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_vpn_default_health_check_ip.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var vpnGwId string
	res := make(map[string]interface{})

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		vpnGwId = v.(string)
		paramMap["VpnGatewayId"] = helper.String(v.(string))
	}

	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var defaultHealthCheck *vpc.GenerateVpnConnectionDefaultHealthCheckIpResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpnDefaultHealthCheckIp(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
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
		if e := tccommon.WriteToFile(output.(string), res); e != nil {
			return e
		}
	}
	return nil
}
