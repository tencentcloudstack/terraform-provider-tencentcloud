package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseInstanceTrafficPackage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseInstanceTrafficPackageRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list.",
			},

			"offset": {
				Optional:    true,
				Default:     0,
				Type:        schema.TypeInt,
				Description: "Offset. Default value is 0.",
			},

			"limit": {
				Optional:    true,
				Default:     20,
				Type:        schema.TypeInt,
				Description: "Number of returned results. Default value is 20. Maximum value is 100.",
			},

			"instance_traffic_package_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of details of instance traffic packages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"traffic_package_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of traffic package details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"traffic_package_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Traffic packet ID.",
									},
									"traffic_used": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Traffic has been used during the effective period of the traffic packet, in bytes.",
									},
									"traffic_package_total": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total traffic in bytes during the effective period of the traffic packet.",
									},
									"traffic_package_remaining": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The remaining traffic during the effective period of the traffic packet, in bytes.",
									},
									"traffic_overflow": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The amount of traffic that exceeds the quota of the traffic packet during the effective period of the traffic packet, in bytes.",
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The start time of the effective cycle of the traffic packet. Expressed according to the ISO8601 standard, and using UTC time. The format is YYYY-MM-DDThh:mm:ssZ.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The end time of the effective period of the traffic packet. Expressed according to the ISO8601 standard, and using UTC time. The format is YYYY-MM-DDThh:mm:ssZ.",
									},
									"deadline": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The expiration time of the traffic package. Expressed according to the ISO8601 standard, and using UTC time. The format is YYYY-MM-DDThh:mm:ssZ..",
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "Traffic packet status:" +
											"- `NETWORK_NORMAL`: normal." +
											"- `OVERDUE_NETWORK_DISABLED`: network disconnection due to arrears.",
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

func dataSourceTencentCloudLighthouseInstanceTrafficPackageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_instance_traffic_package.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		instanceIds := make([]string, 0)
		for _, instanceId := range instanceIdsSet {
			instanceIds = append(instanceIds, instanceId.(string))
		}
		paramMap["instance_ids"] = instanceIds
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceTrafficPackageSet []*lighthouse.InstanceTrafficPackage

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseInstanceTrafficPackageByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceTrafficPackageSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceTrafficPackageSet))
	tmpList := make([]map[string]interface{}, 0, len(instanceTrafficPackageSet))

	if instanceTrafficPackageSet != nil {
		for _, instanceTrafficPackage := range instanceTrafficPackageSet {
			instanceTrafficPackageMap := map[string]interface{}{}

			if instanceTrafficPackage.InstanceId != nil {
				instanceTrafficPackageMap["instance_id"] = instanceTrafficPackage.InstanceId
			}

			if instanceTrafficPackage.TrafficPackageSet != nil {
				trafficPackageSetList := []map[string]interface{}{}
				for _, trafficPackageSet := range instanceTrafficPackage.TrafficPackageSet {
					trafficPackageSetMap := map[string]interface{}{}

					if trafficPackageSet.TrafficPackageId != nil {
						trafficPackageSetMap["traffic_package_id"] = trafficPackageSet.TrafficPackageId
					}

					if trafficPackageSet.TrafficUsed != nil {
						trafficPackageSetMap["traffic_used"] = trafficPackageSet.TrafficUsed
					}

					if trafficPackageSet.TrafficPackageTotal != nil {
						trafficPackageSetMap["traffic_package_total"] = trafficPackageSet.TrafficPackageTotal
					}

					if trafficPackageSet.TrafficPackageRemaining != nil {
						trafficPackageSetMap["traffic_package_remaining"] = trafficPackageSet.TrafficPackageRemaining
					}

					if trafficPackageSet.TrafficOverflow != nil {
						trafficPackageSetMap["traffic_overflow"] = trafficPackageSet.TrafficOverflow
					}

					if trafficPackageSet.StartTime != nil {
						trafficPackageSetMap["start_time"] = trafficPackageSet.StartTime
					}

					if trafficPackageSet.EndTime != nil {
						trafficPackageSetMap["end_time"] = trafficPackageSet.EndTime
					}

					if trafficPackageSet.Deadline != nil {
						trafficPackageSetMap["deadline"] = trafficPackageSet.Deadline
					}

					if trafficPackageSet.Status != nil {
						trafficPackageSetMap["status"] = trafficPackageSet.Status
					}

					trafficPackageSetList = append(trafficPackageSetList, trafficPackageSetMap)
				}

				instanceTrafficPackageMap["traffic_package_set"] = trafficPackageSetList
			}

			ids = append(ids, *instanceTrafficPackage.InstanceId)
			tmpList = append(tmpList, instanceTrafficPackageMap)
		}

		_ = d.Set("instance_traffic_package_set", tmpList)
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
