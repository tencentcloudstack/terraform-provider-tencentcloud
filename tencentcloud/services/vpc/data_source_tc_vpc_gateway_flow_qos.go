package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcGatewayFlowQos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcGatewayFlowQosRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Network instance ID, the network instance types we currently support are:Private line gateway instance ID, in the form of `dcg-ltjahce6`;Nat gateway instance ID, in the form of `nat-ltjahce6`;VPN gateway instance ID, in the form of `vpn-ltjahce6`.",
			},

			"ip_addresses": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Intranet IP of the cloud server with traffic limitation.",
			},

			"gateway_qos_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "instance detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vpc id.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cvm ip address.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "bandwidth value.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
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

func dataSourceTencentCloudVpcGatewayFlowQosRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_gateway_flow_qos.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_addresses"); ok {
		ipAddressesSet := v.(*schema.Set).List()
		paramMap["IpAddresses"] = helper.InterfacesStringsPoint(ipAddressesSet)
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var gatewayQosSet []*vpc.GatewayQos

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcGatewayFlowQosByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		gatewayQosSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(gatewayQosSet))
	tmpList := make([]map[string]interface{}, 0, len(gatewayQosSet))

	if gatewayQosSet != nil {
		for _, gatewayQos := range gatewayQosSet {
			gatewayQosMap := map[string]interface{}{}

			if gatewayQos.VpcId != nil {
				gatewayQosMap["vpc_id"] = gatewayQos.VpcId
			}

			if gatewayQos.IpAddress != nil {
				gatewayQosMap["ip_address"] = gatewayQos.IpAddress
			}

			if gatewayQos.Bandwidth != nil {
				gatewayQosMap["bandwidth"] = gatewayQos.Bandwidth
			}

			if gatewayQos.CreateTime != nil {
				gatewayQosMap["create_time"] = gatewayQos.CreateTime
			}

			ids = append(ids, *gatewayQos.IpAddress)
			tmpList = append(tmpList, gatewayQosMap)
		}

		_ = d.Set("gateway_qos_set", tmpList)
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
