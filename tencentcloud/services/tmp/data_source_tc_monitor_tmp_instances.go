package tmp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
)

func DataSourceTencentCloudMonitorTmpInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorTmpInstancesRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Query according to one or more instance IDs. The instance ID is like: prom-xxxx. The maximum number of instances requested is 100.",
			},

			"instance_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Filter according to instance status.\n" +
					"- 1: Creating;\n" +
					"- 2: In operation;\n" +
					"- 3: Abnormal;\n" +
					"- 4: Reconstruction;\n" +
					"- 5: Destruction;\n" +
					"- 6: Stopped taking;\n" +
					"- 8: Suspension of service due to arrears;\n" +
					"- 9: Service has been suspended due to arrears.",
			},

			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter according to instance name.",
			},

			"zones": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter according to availability area. The availability area is shaped like: ap-Guangzhou-1.",
			},

			"tag_filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter according to tag Key-Value pair. The tag-key is replaced with a specific label key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The key of the tag.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the tag.",
						},
					},
				},
			},

			"ipv4_address": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter according to ipv4 address.",
			},

			"instance_charge_type": {
				Optional: true,
				Type:     schema.TypeInt,
				Description: "Filter according to instance charge type.\n" +
					"- 2: Prepaid;\n" +
					"- 3: Postpaid by hour.",
			},

			"instance_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"instance_charge_type": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Instance charge type.\n" +
								"- 2: Prepaid;\n" +
								"- 3: Postpaid by hour.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region id.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC id.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet id.",
						},
						"data_retention_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data retention time.",
						},
						"instance_status": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Filter according to instance status.\n" +
								"- 1: Creating;\n" +
								"- 2: In operation;\n" +
								"- 3: Abnormal;\n" +
								"- 4: Reconstruction;\n" +
								"- 5: Destruction;\n" +
								"- 6: Stopped taking;\n" +
								"- 8: Suspension of service due to arrears;\n" +
								"- 9: Service has been suspended due to arrears.",
						},
						"grafana_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Grafana panel url.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created_at.",
						},
						"enable_grafana": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Whether to enable grafana.\n" +
								"- 0: closed;\n" +
								"- 1: open.",
						},
						"ipv4_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPV4 address.",
						},
						"tag_specification": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of tags associated with the instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The key of the tag.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the tag.",
									},
								},
							},
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expires for purchased instances.",
						},
						"charge_status": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Charge status.\n" +
								"- 1: Normal;\n" +
								"- 2: Expires;\n" +
								"- 3: Destruction;\n" +
								"- 4: Allocation;\n" +
								"- 5: Allocation failed.",
						},
						"spec_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification name.",
						},
						"auto_renew_flag": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Automatic renewal flag.\n" +
								"- 0: No automatic renewal;\n" +
								"- 1: Enable automatic renewal;\n" +
								"- 2: Automatic renewal is prohibited;\n" +
								"- -1: Invalid.",
						},
						"is_near_expire": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Whether it is about to expire.\n" +
								"- 0: No;\n" +
								"- 1: Expiring soon.",
						},
						"auth_token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Token required for data writing.",
						},
						"remote_write": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address of prometheus remote write.",
						},
						"api_root_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Prometheus http api root address.",
						},
						"proxy_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proxy address.",
						},
						"grafana_status": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Grafana status.\n" +
								"- 1: Creating;\n" +
								"- 2: In operation;\n" +
								"- 3: Abnormal;\n" +
								"- 4: Rebooting;\n" +
								"- 5: Destruction;\n" +
								"- 6: Shutdown;\n" +
								"- 7: Deleted.",
						},
						"grafana_ip_white_list": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Grafana IP whitelist list.",
						},
						"grant": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authorization information for the instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"has_charge_operation": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether you have charging operation authority (1=yes, 2=no).",
									},
									"has_vpc_display": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether to display VPC information (1=yes, 2=no).",
									},
									"has_grafana_status_change": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether the status of Grafana can be modified (1=yes, 2=no).",
									},
									"has_agent_manage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether you have permission to manage the agent (1=yes, 2=no).",
									},
									"has_tke_manage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether you have permission to manage TKE integration (1=yes, 2=no).",
									},
									"has_api_operation": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether to display API and other information (1=yes, 2=no).",
									},
								},
							},
						},
						"grafana_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Binding grafana instance id.",
						},
						"alert_rule_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Alert rule limit.",
						},
						"recording_rule_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Pre-aggregation rule limitations.",
						},
						"migration_type": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Migration status.\n" +
								"- 0: Not in migration;\n+" +
								"- 1: Migrating, original instance;\n+" +
								"- 2: Migrating, target instance.",
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

func dataSourceTencentCloudMonitorTmpInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_tmp_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("instance_status"); ok {
		instanceStatusSet := v.(*schema.Set).List()

		instanceStatusList := make([]*int64, 0)
		for i := range instanceStatusSet {
			instanceStatus := instanceStatusSet[i].(int)
			instanceStatusList = append(instanceStatusList, helper.IntInt64(instanceStatus))
		}
		paramMap["InstanceStatus"] = instanceStatusList
	}

	if v, ok := d.GetOk("instance_name"); ok {
		paramMap["InstanceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		paramMap["Zones"] = helper.InterfacesStringsPoint(zonesSet)
	}

	if v, ok := d.GetOk("tag_filters"); ok {
		tagFiltersSet := v.([]interface{})
		tmpSet := make([]*monitor.PrometheusTag, 0, len(tagFiltersSet))

		for _, item := range tagFiltersSet {
			prometheusTag := monitor.PrometheusTag{}
			prometheusTagMap := item.(map[string]interface{})

			if v, ok := prometheusTagMap["key"]; ok {
				prometheusTag.Key = helper.String(v.(string))
			}
			if v, ok := prometheusTagMap["value"]; ok {
				prometheusTag.Value = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &prometheusTag)
		}
		paramMap["TagFilters"] = tmpSet
	}

	if v, ok := d.GetOk("ipv4_address"); ok {
		iPv4AddressSet := v.(*schema.Set).List()
		paramMap["IPv4Address"] = helper.InterfacesStringsPoint(iPv4AddressSet)
	}

	if v, _ := d.GetOk("instance_charge_type"); v != nil {
		paramMap["InstanceChargeType"] = helper.IntInt64(v.(int))
	}

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var instanceSet []*monitor.PrometheusInstancesItem

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorTmpInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceSet))
	tmpList := make([]map[string]interface{}, 0, len(instanceSet))

	if instanceSet != nil {
		for _, prometheusInstancesItem := range instanceSet {
			prometheusInstancesItemMap := map[string]interface{}{}

			if prometheusInstancesItem.InstanceId != nil {
				prometheusInstancesItemMap["instance_id"] = prometheusInstancesItem.InstanceId
			}

			if prometheusInstancesItem.InstanceName != nil {
				prometheusInstancesItemMap["instance_name"] = prometheusInstancesItem.InstanceName
			}

			if prometheusInstancesItem.InstanceChargeType != nil {
				prometheusInstancesItemMap["instance_charge_type"] = prometheusInstancesItem.InstanceChargeType
			}

			if prometheusInstancesItem.RegionId != nil {
				prometheusInstancesItemMap["region_id"] = prometheusInstancesItem.RegionId
			}

			if prometheusInstancesItem.Zone != nil {
				prometheusInstancesItemMap["zone"] = prometheusInstancesItem.Zone
			}

			if prometheusInstancesItem.VpcId != nil {
				prometheusInstancesItemMap["vpc_id"] = prometheusInstancesItem.VpcId
			}

			if prometheusInstancesItem.SubnetId != nil {
				prometheusInstancesItemMap["subnet_id"] = prometheusInstancesItem.SubnetId
			}

			if prometheusInstancesItem.DataRetentionTime != nil {
				prometheusInstancesItemMap["data_retention_time"] = prometheusInstancesItem.DataRetentionTime
			}

			if prometheusInstancesItem.InstanceStatus != nil {
				prometheusInstancesItemMap["instance_status"] = prometheusInstancesItem.InstanceStatus
			}

			if prometheusInstancesItem.GrafanaURL != nil {
				prometheusInstancesItemMap["grafana_url"] = prometheusInstancesItem.GrafanaURL
			}

			if prometheusInstancesItem.CreatedAt != nil {
				prometheusInstancesItemMap["created_at"] = prometheusInstancesItem.CreatedAt
			}

			if prometheusInstancesItem.EnableGrafana != nil {
				prometheusInstancesItemMap["enable_grafana"] = prometheusInstancesItem.EnableGrafana
			}

			if prometheusInstancesItem.IPv4Address != nil {
				prometheusInstancesItemMap["ipv4_address"] = prometheusInstancesItem.IPv4Address
			}

			if prometheusInstancesItem.TagSpecification != nil {
				tagSpecificationList := []interface{}{}
				for _, tagSpecification := range prometheusInstancesItem.TagSpecification {
					tagSpecificationMap := map[string]interface{}{}

					if tagSpecification.Key != nil {
						tagSpecificationMap["key"] = tagSpecification.Key
					}

					if tagSpecification.Value != nil {
						tagSpecificationMap["value"] = tagSpecification.Value
					}
					tagSpecificationList = append(tagSpecificationList, tagSpecificationMap)
				}
				prometheusInstancesItemMap["tag_specification"] = tagSpecificationList
			}

			if prometheusInstancesItem.ExpireTime != nil {
				prometheusInstancesItemMap["expire_time"] = prometheusInstancesItem.ExpireTime
			}

			if prometheusInstancesItem.ChargeStatus != nil {
				prometheusInstancesItemMap["charge_status"] = prometheusInstancesItem.ChargeStatus
			}

			if prometheusInstancesItem.SpecName != nil {
				prometheusInstancesItemMap["spec_name"] = prometheusInstancesItem.SpecName
			}

			if prometheusInstancesItem.AutoRenewFlag != nil {
				prometheusInstancesItemMap["auto_renew_flag"] = prometheusInstancesItem.AutoRenewFlag
			}

			if prometheusInstancesItem.IsNearExpire != nil {
				prometheusInstancesItemMap["is_near_expire"] = prometheusInstancesItem.IsNearExpire
			}

			if prometheusInstancesItem.AuthToken != nil {
				prometheusInstancesItemMap["auth_token"] = prometheusInstancesItem.AuthToken
			}

			if prometheusInstancesItem.RemoteWrite != nil {
				prometheusInstancesItemMap["remote_write"] = prometheusInstancesItem.RemoteWrite
			}

			if prometheusInstancesItem.ApiRootPath != nil {
				prometheusInstancesItemMap["api_root_path"] = prometheusInstancesItem.ApiRootPath
			}

			if prometheusInstancesItem.ProxyAddress != nil {
				prometheusInstancesItemMap["proxy_address"] = prometheusInstancesItem.ProxyAddress
			}

			if prometheusInstancesItem.GrafanaStatus != nil {
				prometheusInstancesItemMap["grafana_status"] = prometheusInstancesItem.GrafanaStatus
			}

			if prometheusInstancesItem.GrafanaIpWhiteList != nil {
				prometheusInstancesItemMap["grafana_ip_white_list"] = prometheusInstancesItem.GrafanaIpWhiteList
			}

			if prometheusInstancesItem.Grant != nil {
				grantMap := map[string]interface{}{}

				if prometheusInstancesItem.Grant.HasChargeOperation != nil {
					grantMap["has_charge_operation"] = prometheusInstancesItem.Grant.HasChargeOperation
				}

				if prometheusInstancesItem.Grant.HasVpcDisplay != nil {
					grantMap["has_vpc_display"] = prometheusInstancesItem.Grant.HasVpcDisplay
				}

				if prometheusInstancesItem.Grant.HasGrafanaStatusChange != nil {
					grantMap["has_grafana_status_change"] = prometheusInstancesItem.Grant.HasGrafanaStatusChange
				}

				if prometheusInstancesItem.Grant.HasAgentManage != nil {
					grantMap["has_agent_manage"] = prometheusInstancesItem.Grant.HasAgentManage
				}

				if prometheusInstancesItem.Grant.HasTkeManage != nil {
					grantMap["has_tke_manage"] = prometheusInstancesItem.Grant.HasTkeManage
				}

				if prometheusInstancesItem.Grant.HasApiOperation != nil {
					grantMap["has_api_operation"] = prometheusInstancesItem.Grant.HasApiOperation
				}

				prometheusInstancesItemMap["grant"] = []interface{}{grantMap}
			}

			if prometheusInstancesItem.GrafanaInstanceId != nil {
				prometheusInstancesItemMap["grafana_instance_id"] = prometheusInstancesItem.GrafanaInstanceId
			}

			if prometheusInstancesItem.AlertRuleLimit != nil {
				prometheusInstancesItemMap["alert_rule_limit"] = prometheusInstancesItem.AlertRuleLimit
			}

			if prometheusInstancesItem.RecordingRuleLimit != nil {
				prometheusInstancesItemMap["recording_rule_limit"] = prometheusInstancesItem.RecordingRuleLimit
			}

			if prometheusInstancesItem.MigrationType != nil {
				prometheusInstancesItemMap["migration_type"] = prometheusInstancesItem.MigrationType
			}

			ids = append(ids, *prometheusInstancesItem.InstanceId)
			tmpList = append(tmpList, prometheusInstancesItemMap)
		}

		_ = d.Set("instance_set", tmpList)
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
