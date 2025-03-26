package emr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLiteHbaseInstances() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source will been deprecated in Terraform TencentCloud provider later version. Please use `tencentcloud_serverless_hbase_instances` instead.",
		Read:               dataSourceTencentCloudLiteHbaseInstancesRead,
		Schema: map[string]*schema.Schema{
			"display_strategy": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Cluster filtering policy. Value range:\n" +
					"	* clusterList: Query the list of clusters except the destroyed cluster;\n" +
					"	* monitorManage: Queries the list of clusters except those destroyed, being created, and failed to create.",
			},

			"order_field": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Sorting field. Value range:\n" +
					"	* clusterId: Sorting by instance ID;\n" +
					"	* addTime: sorted by instance creation time;\n" +
					"	* status: sorted by the status code of the instance.",
			},

			"asc": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Sort by OrderField in ascending or descending order. Value range:\n" +
					"	* 0: indicates the descending order;\n" +
					"	* 1: indicates the ascending order;\n" +
					"	The default value is 0.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field value.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster Instance String ID.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster Instance Digital ID.",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State description.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster Instance name.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region ID.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Primary Availability Zone ID.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Primary Availability Zone Name.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User APP ID.",
						},
						"vpc_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Primary Availability Vpc ID.",
						},
						"subnet_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Primary Availability Subnet ID.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status code, please refer to the StatusDesc.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster charging type. 0 means charging by volume, 1 means annual and monthly.",
						},
						"zone_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed configuration of the instance availability zone, including the availability zone name, VPC information, and the total number of nodes, where the total number of nodes must be greater than or equal to 3 and less than or equal to 50.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The availability zone to which the instance belongs, such as ap-guangzhou-1.",
									},
									"vpc_settings": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Private network related information configuration. This parameter can be used to specify the ID of the private network, subnet ID, and other information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "VPC ID.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Subnet ID.",
												},
											},
										},
									},
									"node_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of nodes.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tag value.",
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

func dataSourceTencentCloudLiteHbaseInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lite_hbase_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("display_strategy"); ok {
		paramMap["DisplayStrategy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_field"); ok {
		paramMap["OrderField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("asc"); ok {
		paramMap["Asc"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*emr.Filters, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filters := emr.Filters{}
			if v, ok := filtersMap["name"]; ok {
				filters.Name = helper.String(v.(string))
			}
			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filters.Values = append(filters.Values, helper.String(values))
				}
			}
			tmpSet = append(tmpSet, &filters)
		}
		paramMap["Filters"] = tmpSet
	}

	instances, err := service.DescribeLiteHbaseInstancesByFilter(ctx, paramMap)
	if err != nil {
		return err
	}
	ids := make([]string, 0, len(instances))
	instanceList := make([]map[string]interface{}, 0, len(instances))
	if len(instances) > 0 {
		for _, instance := range instances {
			instanceMap := map[string]interface{}{}

			if instance.ClusterId != nil {
				instanceMap["cluster_id"] = instance.ClusterId
			}

			if instance.Id != nil {
				instanceMap["id"] = instance.Id
				ids = append(ids, (string)(*instance.Id))
			}

			if instance.StatusDesc != nil {
				instanceMap["status_desc"] = instance.StatusDesc
			}

			if instance.ClusterName != nil {
				instanceMap["cluster_name"] = instance.ClusterName
			}

			if instance.RegionId != nil {
				instanceMap["region_id"] = instance.RegionId
			}

			if instance.ZoneId != nil {
				instanceMap["zone_id"] = instance.ZoneId
			}

			if instance.Zone != nil {
				instanceMap["zone"] = instance.Zone
			}

			if instance.AppId != nil {
				instanceMap["app_id"] = instance.AppId
			}

			if instance.VpcId != nil {
				instanceMap["vpc_id"] = instance.VpcId
			}

			if instance.SubnetId != nil {
				instanceMap["subnet_id"] = instance.SubnetId
			}

			if instance.Status != nil {
				instanceMap["status"] = instance.Status
			}

			if instance.AddTime != nil {
				instanceMap["add_time"] = instance.AddTime
			}

			if instance.PayMode != nil {
				instanceMap["pay_mode"] = instance.PayMode
			}

			zoneSettingsList := make([]map[string]interface{}, 0, len(instance.ZoneSettings))
			if instance.ZoneSettings != nil {
				for _, zoneSettings := range instance.ZoneSettings {
					zoneSettingsMap := map[string]interface{}{}

					if zoneSettings.Zone != nil {
						zoneSettingsMap["zone"] = zoneSettings.Zone
					}

					vPCSettingsMap := map[string]interface{}{}

					if zoneSettings.VPCSettings != nil {
						if zoneSettings.VPCSettings.VpcId != nil {
							vPCSettingsMap["vpc_id"] = zoneSettings.VPCSettings.VpcId
						}

						if zoneSettings.VPCSettings.SubnetId != nil {
							vPCSettingsMap["subnet_id"] = zoneSettings.VPCSettings.SubnetId
						}

						zoneSettingsMap["vpc_settings"] = []interface{}{vPCSettingsMap}
					}

					if zoneSettings.NodeNum != nil {
						zoneSettingsMap["node_num"] = zoneSettings.NodeNum
					}

					zoneSettingsList = append(zoneSettingsList, zoneSettingsMap)
				}

				instanceMap["zone_settings"] = zoneSettingsList
			}
			tagsList := make([]map[string]interface{}, 0, len(instance.Tags))
			if instance.Tags != nil {
				for _, tags := range instance.Tags {
					tagsMap := map[string]interface{}{}

					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				instanceMap["tags"] = tagsList
			}
			instanceList = append(instanceList, instanceMap)
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("instance_list", instanceList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceList); e != nil {
			return e
		}
	}

	return nil
}
