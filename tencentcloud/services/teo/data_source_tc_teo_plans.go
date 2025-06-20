package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoPlans() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceTencentCloudTeoPlansRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions, the upper limit of Filters. Values is 20. The detailed filtering conditions are as follows: <li>plan-type<br>Filter according to [<strong>Package Type</strong>]. <br>Optional types are: <br>plan-trial: Trial Package; <br>plan-personal: Personal Package; <br>plan-basic: Basic Package; <br>plan-standard: Standard Package; <br>plan-enterprise: Enterprise Package. </li><li>plan-id<br>Filter according to [<strong>Package ID</strong>]. The package ID is in the form of: edgeone-268z103ob0sx.</li><li>area<br>Filter according to [<strong>Package Acceleration Region</strong>]. </li>Service area, optional types are: <br>mainland: Mainland China; <br>overseas: Global (excluding Mainland China); <br>global: Global (including Mainland China).<br><li>status<br>Filter by [<strong>Package Status</strong>].<br>The available statuses are:<br>normal: normal status;<br>expiring-soon: about to expire;<br>expired: expired;<br>isolated: isolated.</li>.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Filter value.",
						},
					},
				},
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field, the values are: <li> enable-time: effective time; </li><li> expire-time: expiration time. </li> If not filled in, the default value enable-time will be used.",
			},

			"direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting direction, the possible values are: <li>asc: sort from small to large; </li><li>desc: sort from large to small. </li>If not filled in, the default value desc will be used.",
			},

			// computed
			"plans": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Plan list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plan type. Possible values are: <li>plan-trial: Trial plan; </li><li>plan-personal: Personal plan; </li><li>plan-basic: Basic plan; </li><li>plan-standard: Standard plan; </li><li>plan-enterprise-v2: Enterprise plan; </li><li>plan-enterprise-model-a: Enterprise Model A plan. </li><li>plan-enterprise: Old Enterprise plan. </li>.",
						},
						"plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plan ID.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service area, the values are: <li>mainland: Mainland China; </li><li>overseas: Worldwide (excluding Mainland China); </li><li>global: Worldwide (including Mainland China).</li>.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package status, the values are: <li>normal: normal status; </li><li>expiring-soon: about to expire; </li><li>expired: expired; </li><li>isolated: isolated; </li><li>overdue-isolated: overdue isolated.</li>.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Payment type, possible values: <li>0: post-payment; </li><li>1: pre-payment.</li>.",
						},
						"zones_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Site information bound to the package, including site ID, site name, and site status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone ID.",
									},
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone name.",
									},
									"paused": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the site is disabled. The possible values are: <li>false: not disabled; </li><li>true: disabled.</li>.",
									},
								},
							},
						},
						"smart_request_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of intelligent acceleration requests in the package, unit: times.",
						},
						"vau_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VAU specifications in the package, unit: piece.",
						},
						"acc_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The content acceleration traffic specifications in the package, unit: byte.",
						},
						"smart_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Smart acceleration traffic specifications within the package, unit: byte.",
						},
						"ddos_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DDoS protection traffic specifications within the package, unit: bytes.",
						},
						"sec_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The security flow specification within the package, unit: byte.",
						},
						"sec_request_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of secure requests in the package, unit: times.",
						},
						"l4_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Layer 4 acceleration traffic specifications within the package, unit: byte.",
						},
						"cross_mlc_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The optimized traffic specifications of the Chinese mainland network in the package, unit: bytes.",
						},
						"bindable": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the package allows binding of new sites, the values are: <li>true: allows binding of new sites; </li><li>false: does not allow binding of new sites.</li>.",
						},
						"enabled_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The package effective time.",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiration date of the package.",
						},
						"features": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The functions supported by the package have the following values: <li>ContentAcceleration: content acceleration function; </li><li>SmartAcceleration: smart acceleration function; </li><li>L4: four-layer acceleration function; </li><li>Waf: advanced web protection; </li><li>QUIC: QUIC function; </li><li>CrossMLC: Chinese mainland network optimization function; </li><li>ProcessMedia: media processing function; </li><li>L4DDoS: four-layer DDoS protection function; </li>L7DDoS function will only have one of the following specifications <li>L7DDoS.CM30G; seven-layer DDoS protection function - Chinese mainland 30G minimum bandwidth specification; </li><li>L7DDoS.CM60G; seven-layer DDoS protection function - Chinese mainland 60G minimum bandwidth specification; </li> <li>L7DDoS.CM100G; Layer 7 DDoS protection function - 100G guaranteed bandwidth for mainland China;</li><li>L7DDoS.Anycast300G; Layer 7 DDoS protection function - 300G guaranteed bandwidth for Anycast outside mainland China;</li><li>L7DDoS.AnycastUnlimited; Layer 7 DDoS protection function - unlimited full protection for Anycast outside mainland China;</li><li>L7DDoS.CM30G_Anycast300G; Layer 7 DDoS protection function - 30G guaranteed bandwidth for mainland China </li><li>L7DDoS.CM60G_Anycast300G; Layer 7 DDoS protection function - 60G guaranteed bandwidth in mainland China, 300G guaranteed bandwidth in anycast outside mainland China; </li><li>L7DDoS.CM100G_Anycast300G; Layer 7 DDoS protection function - 100G guaranteed bandwidth in mainland China, 300G guaranteed bandwidth in anycast outside mainland China; </li><li>L7DDoS.CM30G_AnycastUnlimited d; Layer 7 DDoS protection function - 30G guaranteed bandwidth in mainland China, unlimited Anycast protection outside mainland China; </li><li>L7DDoS.CM60G_AnycastUnlimited; Layer 7 DDoS protection function - 60G guaranteed bandwidth in mainland China, unlimited Anycast protection outside mainland China; </li><li>L7DDoS.CM100G_AnycastUnlimited; Layer 7 DDoS protection function - 100G guaranteed bandwidth in mainland China, unlimited Anycast protection outside mainland China; </li>.",
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

func DataSourceTencentCloudTeoPlansRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_plans.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*teo.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filter := teo.Filter{}
			filterMap := item.(map[string]interface{})
			if v, ok := filterMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("direction"); ok {
		paramMap["Direction"] = helper.String(v.(string))
	}

	var plans []*teo.Plan
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoPlansByFilters(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		plans = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(plans))
	tmpList := make([]map[string]interface{}, 0, len(plans))
	if plans != nil {
		for _, plan := range plans {
			planMap := map[string]interface{}{}
			if plan.PlanType != nil {
				planMap["plan_type"] = plan.PlanType
			}

			if plan.PlanId != nil {
				planMap["plan_id"] = plan.PlanId
				ids = append(ids, *plan.PlanId)
			}

			if plan.Area != nil {
				planMap["area"] = plan.Area
			}

			if plan.Status != nil {
				planMap["status"] = plan.Status
			}

			if plan.PayMode != nil {
				planMap["pay_mode"] = plan.PayMode
			}

			if plan.ZonesInfo != nil {
				zonesInfoList := []interface{}{}
				for _, zonesInfo := range plan.ZonesInfo {
					zonesInfoMap := map[string]interface{}{}
					if zonesInfo.ZoneId != nil {
						zonesInfoMap["zone_id"] = zonesInfo.ZoneId
					}

					if zonesInfo.ZoneName != nil {
						zonesInfoMap["zone_name"] = zonesInfo.ZoneName
					}

					if zonesInfo.Paused != nil {
						zonesInfoMap["paused"] = zonesInfo.Paused
					}

					zonesInfoList = append(zonesInfoList, zonesInfoMap)
				}

				planMap["zones_info"] = zonesInfoList
			}

			if plan.SmartRequestCapacity != nil {
				planMap["smart_request_capacity"] = plan.SmartRequestCapacity
			}

			if plan.VAUCapacity != nil {
				planMap["vau_capacity"] = plan.VAUCapacity
			}

			if plan.AccTrafficCapacity != nil {
				planMap["acc_traffic_capacity"] = plan.AccTrafficCapacity
			}

			if plan.SmartTrafficCapacity != nil {
				planMap["smart_traffic_capacity"] = plan.SmartTrafficCapacity
			}

			if plan.DDoSTrafficCapacity != nil {
				planMap["ddos_traffic_capacity"] = plan.DDoSTrafficCapacity
			}

			if plan.SecTrafficCapacity != nil {
				planMap["sec_traffic_capacity"] = plan.SecTrafficCapacity
			}

			if plan.SecRequestCapacity != nil {
				planMap["sec_request_capacity"] = plan.SecRequestCapacity
			}

			if plan.L4TrafficCapacity != nil {
				planMap["l4_traffic_capacity"] = plan.L4TrafficCapacity
			}

			if plan.CrossMLCTrafficCapacity != nil {
				planMap["cross_mlc_traffic_capacity"] = plan.CrossMLCTrafficCapacity
			}

			if plan.Bindable != nil {
				planMap["bindable"] = plan.Bindable
			}

			if plan.EnabledTime != nil {
				planMap["enabled_time"] = plan.EnabledTime
			}

			if plan.ExpiredTime != nil {
				planMap["expired_time"] = plan.ExpiredTime
			}

			if plan.Features != nil {
				planMap["features"] = plan.Features
			}

			tmpList = append(tmpList, planMap)
		}

		_ = d.Set("plans", tmpList)
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
