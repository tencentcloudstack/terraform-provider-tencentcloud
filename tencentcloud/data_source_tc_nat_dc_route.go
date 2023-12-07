package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudNatDcRoute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudNatDcRouteRead,
		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique identifier of Nat Gateway.",
			},

			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique identifier of Vpc.",
			},

			"nat_direct_connect_gateway_route_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data of route.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv4 CIDR of subnet.",
						},
						"gateway_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of next-hop gateway, valid values: DIRECTCONNECT.",
						},
						"gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of next-hop gateway.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of route.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time of route.",
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

func dataSourceTencentCloudNatDcRouteRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_nat_dc_route.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("nat_gateway_id"); ok {
		paramMap["NatGatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["VpcId"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var natDirectConnectGatewayRouteSet []*vpc.NatDirectConnectGatewayRoute

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeNatDcRouteByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		natDirectConnectGatewayRouteSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(natDirectConnectGatewayRouteSet))
	tmpList := make([]map[string]interface{}, 0, len(natDirectConnectGatewayRouteSet))

	if natDirectConnectGatewayRouteSet != nil {
		for _, natDirectConnectGatewayRoute := range natDirectConnectGatewayRouteSet {
			natDirectConnectGatewayRouteMap := map[string]interface{}{}

			if natDirectConnectGatewayRoute.DestinationCidrBlock != nil {
				natDirectConnectGatewayRouteMap["destination_cidr_block"] = natDirectConnectGatewayRoute.DestinationCidrBlock
			}

			if natDirectConnectGatewayRoute.GatewayType != nil {
				natDirectConnectGatewayRouteMap["gateway_type"] = natDirectConnectGatewayRoute.GatewayType
			}

			if natDirectConnectGatewayRoute.GatewayId != nil {
				natDirectConnectGatewayRouteMap["gateway_id"] = natDirectConnectGatewayRoute.GatewayId
			}

			if natDirectConnectGatewayRoute.CreateTime != nil {
				natDirectConnectGatewayRouteMap["create_time"] = natDirectConnectGatewayRoute.CreateTime
			}

			if natDirectConnectGatewayRoute.UpdateTime != nil {
				natDirectConnectGatewayRouteMap["update_time"] = natDirectConnectGatewayRoute.UpdateTime
			}

			ids = append(ids, *natDirectConnectGatewayRoute.GatewayId)
			tmpList = append(tmpList, natDirectConnectGatewayRouteMap)
		}

		_ = d.Set("nat_direct_connect_gateway_route_set", tmpList)
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
