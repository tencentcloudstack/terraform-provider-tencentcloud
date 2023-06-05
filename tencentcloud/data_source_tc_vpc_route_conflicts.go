/*
Use this data source to query detailed information of vpc route_conflicts

Example Usage

```hcl
data "tencentcloud_vpc_route_conflicts" "route_conflicts" {
  route_table_id = "rtb-6xypllqe"
  destination_cidr_blocks = ["172.18.111.0/24"]
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

func dataSourceTencentCloudVpcRouteConflicts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcRouteConflictsRead,
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Routing table instance ID, for example:rtb-azd4dt1c.",
			},

			"destination_cidr_blocks": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of conflicting destinations to check for.",
			},

			"route_conflict_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "route conflict list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "route table id.",
						},
						"destination_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "destination cidr block.",
						},
						"conflict_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "route conflict list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Destination Cidr Block, like 112.20.51.0/24.",
									},
									"gateway_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "next gateway type.",
									},
									"gateway_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "next hop id.",
									},
									"route_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "route id.",
									},
									"route_description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "route description.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "if enabled.",
									},
									"route_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "routr type.",
									},
									"route_table_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "route table id.",
									},
									"destination_ipv6_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Destination of Ipv6 Cidr Block.",
									},
									"route_item_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "unique policy id.",
									},
									"published_to_vbc": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "if published To ccn.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "create time.",
									},
								},
							},
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

func dataSourceTencentCloudVpcRouteConflictsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_route_conflicts.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("route_table_id"); ok {
		paramMap["RouteTableId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_cidr_blocks"); ok {
		destinationCidrBlocksSet := v.(*schema.Set).List()
		paramMap["DestinationCidrBlocks"] = helper.InterfacesStringsPoint(destinationCidrBlocksSet)
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var routeConflictSet []*vpc.RouteConflict

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcRouteConflicts(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		routeConflictSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(routeConflictSet))
	tmpList := make([]map[string]interface{}, 0, len(routeConflictSet))

	if routeConflictSet != nil {
		for _, routeConflict := range routeConflictSet {
			routeConflictMap := map[string]interface{}{}

			if routeConflict.RouteTableId != nil {
				routeConflictMap["route_table_id"] = routeConflict.RouteTableId
			}

			if routeConflict.DestinationCidrBlock != nil {
				routeConflictMap["destination_cidr_block"] = routeConflict.DestinationCidrBlock
			}

			if routeConflict.ConflictSet != nil {
				conflictSetList := []interface{}{}
				for _, conflictSet := range routeConflict.ConflictSet {
					conflictSetMap := map[string]interface{}{}

					if conflictSet.DestinationCidrBlock != nil {
						conflictSetMap["destination_cidr_block"] = conflictSet.DestinationCidrBlock
					}

					if conflictSet.GatewayType != nil {
						conflictSetMap["gateway_type"] = conflictSet.GatewayType
					}

					if conflictSet.GatewayId != nil {
						conflictSetMap["gateway_id"] = conflictSet.GatewayId
					}

					if conflictSet.RouteId != nil {
						conflictSetMap["route_id"] = conflictSet.RouteId
					}

					if conflictSet.RouteDescription != nil {
						conflictSetMap["route_description"] = conflictSet.RouteDescription
					}

					if conflictSet.Enabled != nil {
						conflictSetMap["enabled"] = conflictSet.Enabled
					}

					if conflictSet.RouteType != nil {
						conflictSetMap["route_type"] = conflictSet.RouteType
					}

					if conflictSet.RouteTableId != nil {
						conflictSetMap["route_table_id"] = conflictSet.RouteTableId
					}

					if conflictSet.DestinationIpv6CidrBlock != nil {
						conflictSetMap["destination_ipv6_cidr_block"] = conflictSet.DestinationIpv6CidrBlock
					}

					if conflictSet.RouteItemId != nil {
						conflictSetMap["route_item_id"] = conflictSet.RouteItemId
					}

					if conflictSet.PublishedToVbc != nil {
						conflictSetMap["published_to_vbc"] = conflictSet.PublishedToVbc
					}

					if conflictSet.CreatedTime != nil {
						conflictSetMap["created_time"] = conflictSet.CreatedTime
					}

					conflictSetList = append(conflictSetList, conflictSetMap)
				}

				routeConflictMap["conflict_set"] = conflictSetList
			}

			ids = append(ids, *routeConflict.RouteTableId)
			tmpList = append(tmpList, routeConflictMap)
		}

		_ = d.Set("route_conflict_set", tmpList)
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
