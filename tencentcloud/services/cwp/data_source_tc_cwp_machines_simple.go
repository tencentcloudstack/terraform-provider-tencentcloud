package cwp

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCwpMachinesSimple() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCwpMachinesSimpleRead,
		Schema: map[string]*schema.Schema{
			"machine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service types. -CVM: Cloud Virtual Machine; -ECM: Edge Computing Machine; -LH: Lighthouse; -Other: Mixed cloud; -ALL: All server types.",
			},
			"machine_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The area where the machine belongs,Such as: ap-guangzhou, ap-shanghai, all-regions: All server region types.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "filter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Only supported Keywords, Version and TagId.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "If `name` is `Keywords`: enter keyword query; If `name` is `Version`: enter PRO_VERSION: Professional Edition | BASIC_VERSION: Basic | Flagship: Flagship | ProtectedMachines: Professional+Flagship | UnFlagship: Non Flagship | PRO_POST_PAY: Professional Edition Pay by Volume | PRO_PRE_PAY: Professional Edition Monthly Package query; If `name` is `TagId`: enter tag ID query.",
						},
						"exact_match": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "exact match. true or false.",
						},
					},
				},
			},
			"project_ids": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Project id list.",
			},
			"machines": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Machine list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"machine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine name.",
						},
						"machine_os": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine OS System.",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cwp client sole UUID.",
						},
						"quuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud server sole UUID.",
						},
						"machine_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine Internal net IP.",
						},
						"is_pro_version": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Paid version or not. true: yes; false: no.",
						},
						"machine_wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine Outer net IP.",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payment model. POSTPAY: Pay as you go; PREPAY: Monthly subscription.",
						},
						"tag": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relevance tag id.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag name.",
									},
									"tag_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Tag ID.",
									},
								},
							},
						},
						"region_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Region detail.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region, Such as ap-guangzhou, ap-shanghai, ap-beijing.",
									},
									"region_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regional Chinese name.",
									},
									"region_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Region ID.",
									},
									"region_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region Code.",
									},
									"region_name_en": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regional English name.",
									},
								},
							},
						},
						"instance_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"machine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service types. -CVM: Cloud Virtual Machine; -ECM: Edge Computing Machine -LH: Lighthouse; -Other: Mixed cloud; -ALL: All server types.",
						},
						"kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Core Version.",
						},
						"protect_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protection Version. -BASIC_VERSION: Basic Version; -PRO_VERSION: Pro Version -Flagship: Flagship Version; -GENERAL_DISCOUNT: CWP-LH Version.",
						},
						"license_order": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "License Order ObjectNote: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"license_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "License ID.",
									},
									"license_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "License Types.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "License Order Status.",
									},
									"source_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Order types.",
									},
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource ID.",
									},
								},
							},
						},
						"cloud_tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cloud tags detailNote: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance IDNote: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudCwpMachinesSimpleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cwp_machines_simple.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = CwpService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		machines []*cwp.MachineSimple
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("machine_type"); ok {
		paramMap["MachineType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("machine_region"); ok {
		paramMap["MachineRegion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cwp.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := cwp.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			if v, ok := filterMap["exact_match"]; ok {
				filter.ExactMatch = helper.Bool(v.(bool))
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsSet := v.(*schema.Set).List()
		paramMap["ProjectIds"] = helper.InterfacesIntUInt64Point(projectIdsSet)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCwpMachinesSimpleByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		machines = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(machines))
	tmpList := make([]map[string]interface{}, 0, len(machines))

	if machines != nil {
		for _, machineSimple := range machines {
			machineSimpleMap := map[string]interface{}{}

			if machineSimple.MachineName != nil {
				machineSimpleMap["machine_name"] = machineSimple.MachineName
			}

			if machineSimple.MachineOs != nil {
				machineSimpleMap["machine_os"] = machineSimple.MachineOs
			}

			if machineSimple.Uuid != nil {
				machineSimpleMap["uuid"] = machineSimple.Uuid
			}

			if machineSimple.Quuid != nil {
				machineSimpleMap["quuid"] = machineSimple.Quuid
			}

			if machineSimple.MachineIp != nil {
				machineSimpleMap["machine_ip"] = machineSimple.MachineIp
			}

			if machineSimple.IsProVersion != nil {
				machineSimpleMap["is_pro_version"] = machineSimple.IsProVersion
			}

			if machineSimple.MachineWanIp != nil {
				machineSimpleMap["machine_wan_ip"] = machineSimple.MachineWanIp
			}

			if machineSimple.PayMode != nil {
				machineSimpleMap["pay_mode"] = machineSimple.PayMode
			}

			if machineSimple.Tag != nil {
				tagList := []interface{}{}
				for _, tag := range machineSimple.Tag {
					tagMap := map[string]interface{}{}

					if tag.Rid != nil {
						tagMap["rid"] = tag.Rid
					}

					if tag.Name != nil {
						tagMap["name"] = tag.Name
					}

					if tag.TagId != nil {
						tagMap["tag_id"] = tag.TagId
					}

					tagList = append(tagList, tagMap)
				}

				machineSimpleMap["tag"] = tagList
			}

			if machineSimple.RegionInfo != nil {
				regionInfoMap := map[string]interface{}{}

				if machineSimple.RegionInfo.Region != nil {
					regionInfoMap["region"] = machineSimple.RegionInfo.Region
				}

				if machineSimple.RegionInfo.RegionName != nil {
					regionInfoMap["region_name"] = machineSimple.RegionInfo.RegionName
				}

				if machineSimple.RegionInfo.RegionId != nil {
					regionInfoMap["region_id"] = machineSimple.RegionInfo.RegionId
				}

				if machineSimple.RegionInfo.RegionCode != nil {
					regionInfoMap["region_code"] = machineSimple.RegionInfo.RegionCode
				}

				if machineSimple.RegionInfo.RegionNameEn != nil {
					regionInfoMap["region_name_en"] = machineSimple.RegionInfo.RegionNameEn
				}

				machineSimpleMap["region_info"] = []interface{}{regionInfoMap}
			}

			if machineSimple.InstanceState != nil {
				machineSimpleMap["instance_state"] = machineSimple.InstanceState
			}

			if machineSimple.ProjectId != nil {
				machineSimpleMap["project_id"] = machineSimple.ProjectId
			}

			if machineSimple.MachineType != nil {
				machineSimpleMap["machine_type"] = machineSimple.MachineType
			}

			if machineSimple.KernelVersion != nil {
				machineSimpleMap["kernel_version"] = machineSimple.KernelVersion
			}

			if machineSimple.ProtectType != nil {
				machineSimpleMap["protect_type"] = machineSimple.ProtectType
			}

			if machineSimple.LicenseOrder != nil {
				licenseOrderMap := map[string]interface{}{}

				if machineSimple.LicenseOrder.LicenseId != nil {
					licenseOrderMap["license_id"] = machineSimple.LicenseOrder.LicenseId
				}

				if machineSimple.LicenseOrder.LicenseType != nil {
					licenseOrderMap["license_type"] = machineSimple.LicenseOrder.LicenseType
				}

				if machineSimple.LicenseOrder.Status != nil {
					licenseOrderMap["status"] = machineSimple.LicenseOrder.Status
				}

				if machineSimple.LicenseOrder.SourceType != nil {
					licenseOrderMap["source_type"] = machineSimple.LicenseOrder.SourceType
				}

				if machineSimple.LicenseOrder.ResourceId != nil {
					licenseOrderMap["resource_id"] = machineSimple.LicenseOrder.ResourceId
				}

				machineSimpleMap["license_order"] = []interface{}{licenseOrderMap}
			}

			if machineSimple.CloudTags != nil {
				cloudTagsList := []interface{}{}
				for _, cloudTags := range machineSimple.CloudTags {
					cloudTagsMap := map[string]interface{}{}

					if cloudTags.TagKey != nil {
						cloudTagsMap["tag_key"] = cloudTags.TagKey
					}

					if cloudTags.TagValue != nil {
						cloudTagsMap["tag_value"] = cloudTags.TagValue
					}

					cloudTagsList = append(cloudTagsList, cloudTagsMap)
				}

				machineSimpleMap["cloud_tags"] = cloudTagsList
			}

			if machineSimple.InstanceId != nil {
				machineSimpleMap["instance_id"] = machineSimple.InstanceId
			}

			ids = append(ids, *machineSimple.Quuid)
			tmpList = append(tmpList, machineSimpleMap)
		}

		_ = d.Set("machines", tmpList)
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
