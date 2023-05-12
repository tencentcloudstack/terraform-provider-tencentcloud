/*
Use this data source to query detailed information of lighthouse reset_instance_blueprint

Example Usage

```hcl
data "tencentcloud_lighthouse_reset_instance_blueprint" "reset_instance_blueprint" {
  instance_id = "lhins-123456"
  offset = 0
  limit = 20
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseResetInstanceBlueprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseResetInstanceBlueprintRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset. Default value is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of returned results. Default value is 20. Maximum value is 100.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter listblueprint-idFilter by image ID.Type: StringRequired: noblueprint-typeFilter by image type.Valid values: APP_OS: application image; PURE_OS: system image; PRIVATE: custom imageType: StringRequired: noplatform-typeFilter by image platform type.Valid values: LINUX_UNIX: Linux or Unix; WINDOWS: WindowsType: StringRequired: noblueprint-nameFilter by image name.Type: StringRequired: noblueprint-stateFilter by image status.Type: StringRequired: noEach request can contain up to 10 Filters and 5 Filter.Values. BlueprintIds and Filters cannot be specified at the same time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter value of field.",
						},
					},
				},
			},

			"reset_instance_blueprint_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of scene info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blueprint_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Mirror details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"blueprint_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image ID, which is the unique identity of Blueprint.",
									},
									"display_title": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The mirror image shows the title to the public.",
									},
									"display_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The image shows the version to the public.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mirror description information.",
									},
									"os_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operating system name.",
									},
									"platform": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operating system platform.",
									},
									"platform_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operating system platform type, such as LINUX_UNIX, WINDOWS.",
									},
									"blueprint_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image type, such as APP_OS, PURE_OS, PRIVATE.",
									},
									"image_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mirror image URL.",
									},
									"required_system_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The size of the system disk required for image (in GB).",
									},
									"blueprint_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mirror status.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time. Expressed according to the ISO8601 standard, and using UTC time. The format is YYYY-MM-DDThh:mm:ssZ.",
									},
									"blueprint_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mirror name.",
									},
									"support_automation_tools": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the image supports automation helper.",
									},
									"required_memory_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Memory required for mirroring (in GB).",
									},
									"image_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CVM image ID after sharing the CVM image to the lightweight application server.",
									},
									"community_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The official website Url.",
									},
									"guide_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Guide article Url.",
									},
									"scene_id_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The mirror association uses the scene Id list.",
									},
									"docker_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Docker version number.",
									},
								},
							},
						},
						"is_resettable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the instance image can be reset to the target image.",
						},
						"non_resettable_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The information cannot be reset. when the mirror can be reset ''.",
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

func dataSourceTencentCloudLighthouseResetInstanceBlueprintRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_reset_instance_blueprint.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["instance_id"] = instanceId
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*lighthouse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := lighthouse.Filter{}
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

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var resetInstanceBlueprintSet []*lighthouse.ResetInstanceBlueprint

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseResetInstanceBlueprintByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		resetInstanceBlueprintSet = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(resetInstanceBlueprintSet))
	for _, resetInstanceBlueprint := range resetInstanceBlueprintSet {
		resetInstanceBlueprintMap := make(map[string]interface{})
		if resetInstanceBlueprint.BlueprintInfo != nil {
			blueprintInfo := make(map[string]interface{})

			if resetInstanceBlueprint.BlueprintInfo.BlueprintId != nil {
				blueprintInfo["blueprint_id"] = *resetInstanceBlueprint.BlueprintInfo.BlueprintId
			}
			if resetInstanceBlueprint.BlueprintInfo.DisplayTitle != nil {
				blueprintInfo["display_title"] = *resetInstanceBlueprint.BlueprintInfo.DisplayTitle
			}
			if resetInstanceBlueprint.BlueprintInfo.DisplayVersion != nil {
				blueprintInfo["display_version"] = *resetInstanceBlueprint.BlueprintInfo.DisplayVersion
			}
			if resetInstanceBlueprint.BlueprintInfo.Description != nil {
				blueprintInfo["description"] = *resetInstanceBlueprint.BlueprintInfo.Description
			}
			if resetInstanceBlueprint.BlueprintInfo.OsName != nil {
				blueprintInfo["os_name"] = *resetInstanceBlueprint.BlueprintInfo.OsName
			}
			if resetInstanceBlueprint.BlueprintInfo.Platform != nil {
				blueprintInfo["platform"] = *resetInstanceBlueprint.BlueprintInfo.Platform
			}
			if resetInstanceBlueprint.BlueprintInfo.PlatformType != nil {
				blueprintInfo["platform_type"] = *resetInstanceBlueprint.BlueprintInfo.PlatformType
			}
			if resetInstanceBlueprint.BlueprintInfo.BlueprintType != nil {
				blueprintInfo["blueprint_type"] = *resetInstanceBlueprint.BlueprintInfo.BlueprintType
			}
			if resetInstanceBlueprint.BlueprintInfo.ImageUrl != nil {
				blueprintInfo["image_url"] = *resetInstanceBlueprint.BlueprintInfo.ImageUrl
			}
			if resetInstanceBlueprint.BlueprintInfo.RequiredSystemDiskSize != nil {
				blueprintInfo["required_system_disk_size"] = *resetInstanceBlueprint.BlueprintInfo.RequiredSystemDiskSize
			}
			if resetInstanceBlueprint.BlueprintInfo.BlueprintState != nil {
				blueprintInfo["blueprint_state"] = *resetInstanceBlueprint.BlueprintInfo.BlueprintState
			}
			if resetInstanceBlueprint.BlueprintInfo.CreatedTime != nil {
				blueprintInfo["created_time"] = *resetInstanceBlueprint.BlueprintInfo.CreatedTime
			}
			if resetInstanceBlueprint.BlueprintInfo.BlueprintName != nil {
				blueprintInfo["blueprint_name"] = *resetInstanceBlueprint.BlueprintInfo.BlueprintName
			}
			if resetInstanceBlueprint.BlueprintInfo.SupportAutomationTools != nil {
				blueprintInfo["support_automation_tools"] = *resetInstanceBlueprint.BlueprintInfo.SupportAutomationTools
			}
			if resetInstanceBlueprint.BlueprintInfo.RequiredMemorySize != nil {
				blueprintInfo["required_memory_size"] = *resetInstanceBlueprint.BlueprintInfo.RequiredMemorySize
			}
			if resetInstanceBlueprint.BlueprintInfo.ImageId != nil {
				blueprintInfo["image_id"] = *resetInstanceBlueprint.BlueprintInfo.ImageId
			}
			if resetInstanceBlueprint.BlueprintInfo.CommunityUrl != nil {
				blueprintInfo["community_url"] = *resetInstanceBlueprint.BlueprintInfo.CommunityUrl
			}
			if resetInstanceBlueprint.BlueprintInfo.GuideUrl != nil {
				blueprintInfo["guide_url"] = *resetInstanceBlueprint.BlueprintInfo.GuideUrl
			}
			if resetInstanceBlueprint.BlueprintInfo.SceneIdSet != nil {
				sceneIds := make([]string, 0)
				for _, sceneId := range resetInstanceBlueprint.BlueprintInfo.SceneIdSet {
					sceneIds = append(sceneIds, *sceneId)
				}
				blueprintInfo["scene_id_set"] = sceneIds
			}
			if resetInstanceBlueprint.BlueprintInfo.DockerVersion != nil {
				blueprintInfo["docker_version"] = *resetInstanceBlueprint.BlueprintInfo.DockerVersion
			}
			resetInstanceBlueprintMap["blueprint_info"] = []map[string]interface{}{blueprintInfo}
		}
		if resetInstanceBlueprint.IsResettable != nil {
			resetInstanceBlueprintMap["is_resettable"] = *resetInstanceBlueprint.IsResettable
		}
		if resetInstanceBlueprint.NonResettableMessage != nil {
			resetInstanceBlueprintMap["non_resettable_message"] = *resetInstanceBlueprint.NonResettableMessage
		}
		tmpList = append(tmpList, resetInstanceBlueprintMap)
	}

	d.SetId(instanceId)
	_ = d.Set("reset_instance_blueprint_set", tmpList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
