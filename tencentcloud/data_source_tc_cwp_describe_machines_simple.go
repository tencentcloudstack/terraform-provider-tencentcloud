/*
Use this data source to query detailed information of cwp describe_machines_simple

Example Usage

```hcl
data "tencentcloud_cwp_describe_machines_simple" "describe_machines_simple" {
  machine_type = ""
  machine_region = ""
  filters {
		name = ""
		values =

  }
  project_ids =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCwpDescribeMachinesSimple() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCwpDescribeMachinesSimpleRead,
		Schema: map[string]*schema.Schema{
			"machine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service types- CVM: Cloud Virtual Machine- ECM: Edge Computing Machine- LH: encentCloud Lighthouse- Other: Mixed cloud- ALL: All server types.",
			},

			"machine_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The area where the machine belongs,Such as：ap-guangzhou，ap-shanghaiall-regions: All server region types.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Keywords filtrationVersion current Protection Vsersion  - PRO_VERSION：Pro Version  - BASIC_VERSION：Basic Version  - Flagship : Flagship Version  - ProtectedMachines: Pro Version And Flagship Version  - UnFlagship : Not Flagship Version  - PRO_POST_PAY：Pro Version Pay as you go  - PRO_PRE_PAY：Pro Version Monthly subscription.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key name。.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "One or more filter values。.",
						},
					},
				},
			},

			"project_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
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
							Description: "Machine name。.",
						},
						"machine_os": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine OS System。.",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cwp client sole UUID。.",
						},
						"quuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud server sole UUID。.",
						},
						"machine_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine Internal net IP。.",
						},
						"is_pro_version": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Paid version or not。true： yesfalse：no.",
						},
						"machine_wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine Outer net IP。.",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payment model。POSTPAY: Pay as you go PREPAY: Monthly subscription.",
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
										Description: "Region，Such as ap-guangzhou，ap-shanghai，ap-beijing.",
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
										Description: "Region Code，如 gz，sh，bj.",
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
							Description: "Instance status TERMINATED_PRO_VERSION 已销毁.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"machine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service types- CVM: Cloud Virtual Machine- ECM: Edge Computing Machine- LH: encentCloud Lighthouse- Other: Mixed cloud- ALL: All server types.",
						},
						"kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Core Version.",
						},
						"protect_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protection Version- BASIC_VERSION Basic Version- PRO_VERSION Pro Version- Flagship Flagship Version- GENERAL_DISCOUNT CWP-LH Version.",
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

func dataSourceTencentCloudCwpDescribeMachinesSimpleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cwp_describe_machines_simple.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsSet := v.(*schema.Set).List()
		for i := range projectIdsSet {
			projectIds := projectIdsSet[i].(int)
			paramMap["ProjectIds"] = append(paramMap["ProjectIds"], helper.IntUint64(projectIds))
		}
	}

	service := CwpService{client: meta.(*TencentCloudClient).apiV3Conn}

	var machines []*cwp.MachineSimple

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCwpDescribeMachinesSimpleByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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

				machineSimpleMap["tag"] = []interface{}{tagList}
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

				machineSimpleMap["cloud_tags"] = []interface{}{cloudTagsList}
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
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
