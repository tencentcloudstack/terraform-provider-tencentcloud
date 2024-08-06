package ccn

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCcnRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnRoutesRead,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the CCN to be queried.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered. Support `route-id`, `cidr-block`, `instance-type`, `instance-region`, `instance-id`, `route-table-id`.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Filter value of the field.",
						},
					},
				},
			},
			// Computed
			"route_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "CCN route list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "route ID.",
						},
						"destination_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next hop type (associated instance type), all types: VPC, DIRECTCONNECT.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next jump (associated instance ID).",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next jump (associated instance name).",
						},
						"instance_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next jump (associated instance region).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "update time.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is routing enabled.",
						},
						"instance_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UIN (root account) to which the associated instance belongs.",
						},
						"extra_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Extension status of routing.",
						},
						"is_bgp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is it dynamic routing.",
						},
						"route_priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Routing priority.",
						},
						"instance_extra_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next hop extension name (associated instance extension name).",
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

func dataSourceTencentCloudCcnRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ccn_routes.read")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ccn_id"); ok {
		paramMap["CcnId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filter := vpc.Filter{}
			filterMap := item.(map[string]interface{})
			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var routeSet []*vpc.CcnRoute
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcDescribeCcnRoutesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		routeSet = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(routeSet))
	tmpList := make([]map[string]interface{}, 0, len(routeSet))

	if routeSet != nil {
		for _, route := range routeSet {
			tmpMap := make(map[string]interface{})
			if route.RouteId != nil {
				tmpMap["route_id"] = route.RouteId
			}

			if route.DestinationCidrBlock != nil {
				tmpMap["destination_cidr_block"] = route.DestinationCidrBlock
			}

			if route.InstanceType != nil {
				tmpMap["instance_type"] = route.InstanceType
			}

			if route.InstanceId != nil {
				tmpMap["instance_id"] = route.InstanceId
			}

			if route.InstanceName != nil {
				tmpMap["instance_name"] = route.InstanceName
			}

			if route.InstanceRegion != nil {
				tmpMap["instance_region"] = route.InstanceRegion
			}

			if route.UpdateTime != nil {
				tmpMap["update_time"] = route.UpdateTime
			}

			if route.Enabled != nil {
				tmpMap["enabled"] = route.Enabled
			}

			if route.InstanceUin != nil {
				tmpMap["instance_uin"] = route.InstanceUin
			}

			if route.ExtraState != nil {
				tmpMap["extra_state"] = route.ExtraState
			}

			if route.IsBgp != nil {
				tmpMap["is_bgp"] = route.IsBgp
			}

			if route.RoutePriority != nil {
				tmpMap["route_priority"] = route.RoutePriority
			}

			if route.InstanceExtraName != nil {
				tmpMap["instance_extra_name"] = route.InstanceExtraName
			}

			ids = append(ids, *route.RouteId)
			tmpList = append(tmpList, tmpMap)
		}

		_ = d.Set("route_list", tmpList)
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
