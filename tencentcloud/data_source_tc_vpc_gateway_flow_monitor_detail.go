/*
Use this data source to query detailed information of vpc gateway_flow_monitor_detail

Example Usage

```hcl
data "tencentcloud_vpc_gateway_flow_monitor_detail" "gateway_flow_monitor_detail" {
  time_point      = "2023-06-02 12:15:20"
  vpn_id          = "vpngw-gt8bianl"
  order_field     = "OutTraffic"
  order_direction = "DESC"
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

func dataSourceTencentCloudVpcGatewayFlowMonitorDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcGatewayFlowMonitorDetailRead,
		Schema: map[string]*schema.Schema{
			"time_point": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The point in time. This indicates details of this minute will be queried. For example, in `2019-02-28 18:15:20`, details at `18:15` will be queried.",
			},

			"vpn_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The instance ID of the VPN gateway, such as `vpn-ltjahce6`.",
			},

			"direct_connect_gateway_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The instance ID of the Direct Connect gateway, such as `dcg-ltjahce6`.",
			},

			"peering_connection_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The instance ID of the peering connection, such as `pcx-ltjahce6`.",
			},

			"nat_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The instance ID of the NAT gateway, such as `nat-ltjahce6`.",
			},

			"order_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The order field supports `InPkg`, `OutPkg`, `InTraffic`, and `OutTraffic`.",
			},

			"order_direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Order methods. Ascending: `ASC`, Descending: `DESC`.",
			},

			"gateway_flow_monitor_detail_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The gateway traffic monitoring details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin `IP`.",
						},
						"in_pkg": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Inbound packets.",
						},
						"out_pkg": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Outbound packets.",
						},
						"in_traffic": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Inbound traffic, in Byte.",
						},
						"out_traffic": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Outbound traffic, in Byte.",
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

func dataSourceTencentCloudVpcGatewayFlowMonitorDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_gateway_flow_monitor_detail.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("time_point"); ok {
		paramMap["TimePoint"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpn_id"); ok {
		paramMap["VpnId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("direct_connect_gateway_id"); ok {
		paramMap["DirectConnectGatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("peering_connection_id"); ok {
		paramMap["PeeringConnectionId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("nat_id"); ok {
		paramMap["NatId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_field"); ok {
		paramMap["OrderField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_direction"); ok {
		paramMap["OrderDirection"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var gatewayFlowMonitorDetailSet []*vpc.GatewayFlowMonitorDetail

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcGatewayFlowMonitorDetailByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		gatewayFlowMonitorDetailSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(gatewayFlowMonitorDetailSet))
	tmpList := make([]map[string]interface{}, 0, len(gatewayFlowMonitorDetailSet))

	if gatewayFlowMonitorDetailSet != nil {
		for _, gatewayFlowMonitorDetail := range gatewayFlowMonitorDetailSet {
			gatewayFlowMonitorDetailMap := map[string]interface{}{}

			if gatewayFlowMonitorDetail.PrivateIpAddress != nil {
				gatewayFlowMonitorDetailMap["private_ip_address"] = gatewayFlowMonitorDetail.PrivateIpAddress
			}

			if gatewayFlowMonitorDetail.InPkg != nil {
				gatewayFlowMonitorDetailMap["in_pkg"] = gatewayFlowMonitorDetail.InPkg
			}

			if gatewayFlowMonitorDetail.OutPkg != nil {
				gatewayFlowMonitorDetailMap["out_pkg"] = gatewayFlowMonitorDetail.OutPkg
			}

			if gatewayFlowMonitorDetail.InTraffic != nil {
				gatewayFlowMonitorDetailMap["in_traffic"] = gatewayFlowMonitorDetail.InTraffic
			}

			if gatewayFlowMonitorDetail.OutTraffic != nil {
				gatewayFlowMonitorDetailMap["out_traffic"] = gatewayFlowMonitorDetail.OutTraffic
			}

			ids = append(ids, *gatewayFlowMonitorDetail.PrivateIpAddress)
			tmpList = append(tmpList, gatewayFlowMonitorDetailMap)
		}

		_ = d.Set("gateway_flow_monitor_detail_set", tmpList)
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
