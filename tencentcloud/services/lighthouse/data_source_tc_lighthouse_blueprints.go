package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseBlueprints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseBlueprintsRead,
		Schema: map[string]*schema.Schema{
			"blueprint_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Blueprint ID list.",
			},

			"filters": {
				Optional: true,
				Type:     schema.TypeList,
				Description: "Filter list.\n" +
					"- `blueprint-id`: Filter by blueprint ID.\n" +
					"- `blueprint-type`: Filter by blueprint type. Values: `APP_OS`, `PURE_OS`, `DOCKER`, `PRIVATE`, `SHARED`.\n" +
					"- `platform-type`: Filter by platform type. Values: `LINUX_UNIX`, `WINDOWS`.\n" +
					"- `blueprint-name`: Filter by blueprint name.\n" +
					"- `blueprint-state`: Filter by blueprint state.\n" +
					"- `scene-id`: Filter by scene ID.\n" +
					"NOTE: The upper limit of Filters per request is 10. The upper limit of Filter.Values is 100. Parameter does not support specifying both BlueprintIds and Filters.",
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

			"blueprint_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of blueprint details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blueprint_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blueprint ID, which is the unique identifier of Blueprint.",
						},
						"display_title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blueprint display title.",
						},
						"display_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blueprint display version.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blueprint description.",
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
							Description: "Platform type, such as LINUX_UNIX and WINDOWS.",
						},
						"blueprint_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blueprint type, such as APP_OS, PURE_OS, DOCKER, PRIVATE, and SHARED.",
						},
						"image_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blueprint image URL.",
						},
						"required_system_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "System disk size required by blueprint in GB.",
						},
						"blueprint_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blueprint state.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time according to ISO 8601 standard. UTC time is used. Format is YYYY-MM-DDThh:mm:ssZ.",
						},
						"blueprint_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blueprint name.",
						},
						"support_automation_tools": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the blueprint supports automation tools.",
						},
						"required_memory_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size required by blueprint in GB.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Lighthouse blueprint shared from a CVM image. Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"community_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of official website of the open-source project.",
						},
						"guide_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Guide documentation URL.",
						},
						"scene_id_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Array of IDs of scenes associated with a blueprint. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"docker_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Docker version. Note: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudLighthouseBlueprintsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_blueprints.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})

	// Parse blueprint_ids
	if v, ok := d.GetOk("blueprint_ids"); ok {
		blueprintIdsSet := v.(*schema.Set).List()
		paramMap["BlueprintIds"] = helper.InterfacesStringsPoint(blueprintIdsSet)
	}

	// Parse filters
	if v, ok := d.GetOk("filters"); ok {
		filtersList := v.([]interface{})
		filters := make([]*lighthouse.Filter, 0, len(filtersList))
		for _, item := range filtersList {
			filterMap := item.(map[string]interface{})
			filter := lighthouse.Filter{}
			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			filters = append(filters, &filter)
		}
		paramMap["Filters"] = filters
	}

	// Call service layer
	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var blueprintSet []*lighthouse.Blueprint

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseBlueprintsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		blueprintSet = result
		return nil
	})

	if err != nil {
		return err
	}

	// Map response to schema
	ids := make([]string, 0, len(blueprintSet))
	tmpList := make([]map[string]interface{}, 0, len(blueprintSet))

	if blueprintSet != nil {
		for _, blueprint := range blueprintSet {
			blueprintMap := map[string]interface{}{}

			if blueprint.BlueprintId != nil {
				blueprintMap["blueprint_id"] = blueprint.BlueprintId
				ids = append(ids, *blueprint.BlueprintId)
			}

			if blueprint.DisplayTitle != nil {
				blueprintMap["display_title"] = blueprint.DisplayTitle
			}

			if blueprint.DisplayVersion != nil {
				blueprintMap["display_version"] = blueprint.DisplayVersion
			}

			if blueprint.Description != nil {
				blueprintMap["description"] = blueprint.Description
			}

			if blueprint.OsName != nil {
				blueprintMap["os_name"] = blueprint.OsName
			}

			if blueprint.Platform != nil {
				blueprintMap["platform"] = blueprint.Platform
			}

			if blueprint.PlatformType != nil {
				blueprintMap["platform_type"] = blueprint.PlatformType
			}

			if blueprint.BlueprintType != nil {
				blueprintMap["blueprint_type"] = blueprint.BlueprintType
			}

			if blueprint.ImageUrl != nil {
				blueprintMap["image_url"] = blueprint.ImageUrl
			}

			if blueprint.RequiredSystemDiskSize != nil {
				blueprintMap["required_system_disk_size"] = blueprint.RequiredSystemDiskSize
			}

			if blueprint.BlueprintState != nil {
				blueprintMap["blueprint_state"] = blueprint.BlueprintState
			}

			if blueprint.CreatedTime != nil {
				blueprintMap["created_time"] = blueprint.CreatedTime
			}

			if blueprint.BlueprintName != nil {
				blueprintMap["blueprint_name"] = blueprint.BlueprintName
			}

			if blueprint.SupportAutomationTools != nil {
				blueprintMap["support_automation_tools"] = blueprint.SupportAutomationTools
			}

			if blueprint.RequiredMemorySize != nil {
				blueprintMap["required_memory_size"] = blueprint.RequiredMemorySize
			}

			if blueprint.ImageId != nil {
				blueprintMap["image_id"] = blueprint.ImageId
			}

			if blueprint.CommunityUrl != nil {
				blueprintMap["community_url"] = blueprint.CommunityUrl
			}

			if blueprint.GuideUrl != nil {
				blueprintMap["guide_url"] = blueprint.GuideUrl
			}

			if blueprint.SceneIdSet != nil {
				blueprintMap["scene_id_set"] = blueprint.SceneIdSet
			}

			if blueprint.DockerVersion != nil {
				blueprintMap["docker_version"] = blueprint.DockerVersion
			}

			tmpList = append(tmpList, blueprintMap)
		}

		_ = d.Set("blueprint_set", tmpList)
	}

	// Set resource ID
	d.SetId(helper.DataResourceIdsHash(ids))

	// Handle result_output_file
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
