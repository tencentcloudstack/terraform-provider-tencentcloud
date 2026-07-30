package cls

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsMachineGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsMachineGroupsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. Maximum 10 filters, each with up to 5 values. Supported keys: `machineGroupName`, `machineGroupId`, `osType`, `tagKey`, `tag:tagKey`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"machine_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of cls machine groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine group ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine group name.",
						},
						"machine_group_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Machine group type.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Machine group type. Valid values: `ip`, `label`.",
									},
									"values": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Machine description list.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"auto_update": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether machine group auto update is enabled.",
						},
						"update_start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Upgrade start time.",
						},
						"update_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Upgrade end time.",
						},
						"service_logging": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether service logging is enabled.",
						},
						"delay_cleanup_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Machine offline periodic cleanup time, in days.",
						},
						"meta_tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Machine group metadata tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Metadata tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Metadata tag value.",
									},
								},
							},
						},
						"os_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Operating system type. 0: Linux, 1: Windows.",
						},
					},
				},
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count of cls machine groups.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudClsMachineGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cls_machine_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cls.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := cls.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Key = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					filter.Values = append(filter.Values, helper.String(value))
				}
			}
			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var (
		respData  []*cls.MachineGroupInfo
		respTotal int64
	)

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, total, e := service.DescribeClsMachineGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		respTotal = total
		return nil
	})

	if reqErr != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId, cls_machine_groups id=%s", d.Id())
		return reqErr
	}

	ids := make([]string, 0, len(respData))
	tmpList := make([]map[string]interface{}, 0, len(respData))

	if respData != nil {
		for _, machineGroupInfo := range respData {
			machineGroupInfoMap := map[string]interface{}{}

			if machineGroupInfo.GroupId != nil {
				machineGroupInfoMap["group_id"] = machineGroupInfo.GroupId
			}

			if machineGroupInfo.GroupName != nil {
				machineGroupInfoMap["group_name"] = machineGroupInfo.GroupName
			}

			if machineGroupInfo.MachineGroupType != nil {
				machineGroupTypeMap := map[string]interface{}{}

				if machineGroupInfo.MachineGroupType.Type != nil {
					machineGroupTypeMap["type"] = machineGroupInfo.MachineGroupType.Type
				}

				if machineGroupInfo.MachineGroupType.Values != nil {
					machineGroupTypeMap["values"] = machineGroupInfo.MachineGroupType.Values
				}

				machineGroupInfoMap["machine_group_type"] = []interface{}{machineGroupTypeMap}
			}

			if machineGroupInfo.CreateTime != nil {
				machineGroupInfoMap["create_time"] = machineGroupInfo.CreateTime
			}

			if machineGroupInfo.Tags != nil {
				tagsList := []interface{}{}
				for _, tag := range machineGroupInfo.Tags {
					tagMap := map[string]interface{}{}

					if tag.Key != nil {
						tagMap["key"] = tag.Key
					}

					if tag.Value != nil {
						tagMap["value"] = tag.Value
					}

					tagsList = append(tagsList, tagMap)
				}

				machineGroupInfoMap["tags"] = tagsList
			}

			if machineGroupInfo.AutoUpdate != nil {
				machineGroupInfoMap["auto_update"] = machineGroupInfo.AutoUpdate
			}

			if machineGroupInfo.UpdateStartTime != nil {
				machineGroupInfoMap["update_start_time"] = machineGroupInfo.UpdateStartTime
			}

			if machineGroupInfo.UpdateEndTime != nil {
				machineGroupInfoMap["update_end_time"] = machineGroupInfo.UpdateEndTime
			}

			if machineGroupInfo.ServiceLogging != nil {
				machineGroupInfoMap["service_logging"] = machineGroupInfo.ServiceLogging
			}

			if machineGroupInfo.DelayCleanupTime != nil {
				machineGroupInfoMap["delay_cleanup_time"] = machineGroupInfo.DelayCleanupTime
			}

			if machineGroupInfo.MetaTags != nil {
				metaTagsList := []interface{}{}
				for _, metaTags := range machineGroupInfo.MetaTags {
					metaTagsMap := map[string]interface{}{}

					if metaTags.Key != nil {
						metaTagsMap["key"] = metaTags.Key
					}

					if metaTags.Value != nil {
						metaTagsMap["value"] = metaTags.Value
					}

					metaTagsList = append(metaTagsList, metaTagsMap)
				}

				machineGroupInfoMap["meta_tags"] = metaTagsList
			}

			if machineGroupInfo.OSType != nil {
				machineGroupInfoMap["os_type"] = machineGroupInfo.OSType
			}

			if machineGroupInfo.GroupId != nil {
				ids = append(ids, *machineGroupInfo.GroupId)
			}

			tmpList = append(tmpList, machineGroupInfoMap)
		}

		_ = d.Set("machine_groups", tmpList)
	}

	_ = d.Set("total_count", respTotal)

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
