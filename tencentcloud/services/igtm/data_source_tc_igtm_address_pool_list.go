package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmAddressPoolList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmAddressPoolListRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Alert filter conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name, supported list as follows:\n- PoolName: Address pool name.\n- MonitorId: Monitor ID. This is a required parameter, failure to provide will cause interface query failure.",
						},
						"value": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field value.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable fuzzy query, only supports filter field name as domain.\nWhen fuzzy query is enabled, maximum Value length is 1, otherwise maximum Value length is 5. (Reserved field, currently not used).",
						},
					},
				},
			},

			"address_pool_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Resource group list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pool_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Address pool ID.",
						},
						"pool_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address pool name.",
						},
						"addr_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address pool address type: IPV4, IPV6, DOMAIN.",
						},
						"traffic_strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Traffic strategy: WEIGHT load balancing, ALL resolve all.",
						},
						"monitor_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Monitor ID.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OK normal, DOWN failure, WARN risk, UNKNOWN unknown.",
						},
						"address_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Address count.",
						},
						"monitor_group_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Probe point count.",
						},
						"monitor_task_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Detection task count.",
						},
						"instance_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance related information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance name.",
									},
								},
							},
						},
						"address_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Address pool address information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Address value: only supports IPv4, IPv6 and domain name formats;\nLoopback addresses, reserved addresses, internal network addresses and Tencent reserved network segments are not supported.",
									},
									"is_enable": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to enable: DISABLED disabled; ENABLED enabled.",
									},
									"address_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Address ID.",
									},
									"location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Address name.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OK normal, DOWN failure, WARN risk, UNKNOWN detecting, UNMONITORED unknown.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Weight, required when traffic strategy is WEIGHT; range 1-100.",
									},
									"created_on": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time.",
									},
									"updated_on": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Modification time.",
									},
								},
							},
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
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

func dataSourceTencentCloudIgtmAddressPoolListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_address_pool_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*igtmv20231024.ResourceFilter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			resourceFilter := igtmv20231024.ResourceFilter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				resourceFilter.Name = helper.String(v)
			}

			if v, ok := filtersMap["value"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					resourceFilter.Value = append(resourceFilter.Value, helper.String(value))
				}
			}

			if v, ok := filtersMap["fuzzy"].(bool); ok {
				resourceFilter.Fuzzy = helper.Bool(v)
			}

			tmpSet = append(tmpSet, &resourceFilter)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData []*igtmv20231024.AddressPool
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIgtmAddressPoolListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	addressPoolSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, addressPoolSet := range respData {
			addressPoolSetMap := map[string]interface{}{}
			if addressPoolSet.PoolId != nil {
				addressPoolSetMap["pool_id"] = addressPoolSet.PoolId
			}

			if addressPoolSet.PoolName != nil {
				addressPoolSetMap["pool_name"] = addressPoolSet.PoolName
			}

			if addressPoolSet.AddrType != nil {
				addressPoolSetMap["addr_type"] = addressPoolSet.AddrType
			}

			if addressPoolSet.TrafficStrategy != nil {
				addressPoolSetMap["traffic_strategy"] = addressPoolSet.TrafficStrategy
			}

			if addressPoolSet.MonitorId != nil {
				addressPoolSetMap["monitor_id"] = addressPoolSet.MonitorId
			}

			if addressPoolSet.Status != nil {
				addressPoolSetMap["status"] = addressPoolSet.Status
			}

			if addressPoolSet.AddressNum != nil {
				addressPoolSetMap["address_num"] = addressPoolSet.AddressNum
			}

			if addressPoolSet.MonitorGroupNum != nil {
				addressPoolSetMap["monitor_group_num"] = addressPoolSet.MonitorGroupNum
			}

			if addressPoolSet.MonitorTaskNum != nil {
				addressPoolSetMap["monitor_task_num"] = addressPoolSet.MonitorTaskNum
			}

			instanceInfoList := make([]map[string]interface{}, 0, len(addressPoolSet.InstanceInfo))
			if addressPoolSet.InstanceInfo != nil {
				for _, instanceInfo := range addressPoolSet.InstanceInfo {
					instanceInfoMap := map[string]interface{}{}
					if instanceInfo.InstanceId != nil {
						instanceInfoMap["instance_id"] = instanceInfo.InstanceId
					}

					if instanceInfo.InstanceName != nil {
						instanceInfoMap["instance_name"] = instanceInfo.InstanceName
					}

					instanceInfoList = append(instanceInfoList, instanceInfoMap)
				}

				addressPoolSetMap["instance_info"] = instanceInfoList
			}

			addressSetList := make([]map[string]interface{}, 0, len(addressPoolSet.AddressSet))
			if addressPoolSet.AddressSet != nil {
				for _, addressSet := range addressPoolSet.AddressSet {
					addressSetMap := map[string]interface{}{}
					if addressSet.Addr != nil {
						addressSetMap["addr"] = addressSet.Addr
					}

					if addressSet.IsEnable != nil {
						addressSetMap["is_enable"] = addressSet.IsEnable
					}

					if addressSet.AddressId != nil {
						addressSetMap["address_id"] = addressSet.AddressId
					}

					if addressSet.Location != nil {
						addressSetMap["location"] = addressSet.Location
					}

					if addressSet.Status != nil {
						addressSetMap["status"] = addressSet.Status
					}

					if addressSet.Weight != nil {
						addressSetMap["weight"] = addressSet.Weight
					}

					if addressSet.CreatedOn != nil {
						addressSetMap["created_on"] = addressSet.CreatedOn
					}

					if addressSet.UpdatedOn != nil {
						addressSetMap["updated_on"] = addressSet.UpdatedOn
					}

					addressSetList = append(addressSetList, addressSetMap)
				}

				addressPoolSetMap["address_set"] = addressSetList
			}

			if addressPoolSet.CreatedOn != nil {
				addressPoolSetMap["created_on"] = addressPoolSet.CreatedOn
			}

			if addressPoolSet.UpdatedOn != nil {
				addressPoolSetMap["updated_on"] = addressPoolSet.UpdatedOn
			}

			addressPoolSetList = append(addressPoolSetList, addressPoolSetMap)
		}

		_ = d.Set("address_pool_set", addressPoolSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), addressPoolSetList); e != nil {
			return e
		}
	}

	return nil
}
