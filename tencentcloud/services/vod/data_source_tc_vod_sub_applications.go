package vod

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVodSubApplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVodSubApplicationsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application name for exact match filtering.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag key-value pairs for filtering applications. Applications matching all specified tags will be returned.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results in JSON format.",
			},
			"sub_application_info_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of sub-application information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sub_app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sub-application ID.",
						},
						"sub_app_id_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sub-application name.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Legacy name field (for backward compatibility).",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sub-application description.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in ISO 8601 format.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application status. Valid values: On, Off, Destroying, Destroyed.",
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application mode. Valid values: fileid, fileid+path.",
						},
						"storage_regions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of enabled storage regions.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Resource tags bound to the sub-application.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVodSubApplicationsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vod_sub_applications.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	filters := make(map[string]interface{})

	// Build filter map
	if v, ok := d.GetOk("name"); ok {
		filters["name"] = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := v.(map[string]interface{})
		tags := make([]*vod.ResourceTag, 0, len(tagsMap))
		for key, value := range tagsMap {
			tag := &vod.ResourceTag{
				TagKey:   helper.String(key),
				TagValue: helper.String(value.(string)),
			}
			tags = append(tags, tag)
		}
		filters["tags"] = tags
	}

	// Call service layer
	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	subAppInfos, err := vodService.DescribeSubApplicationsByFilter(ctx, filters)
	if err != nil {
		return err
	}

	// Map results to Terraform state
	subAppInfoList := make([]map[string]interface{}, 0, len(subAppInfos))
	ids := make([]string, 0, len(subAppInfos))

	for _, info := range subAppInfos {
		subAppInfoMap := map[string]interface{}{}

		if info.SubAppId != nil {
			subAppInfoMap["sub_app_id"] = int(*info.SubAppId)
			ids = append(ids, helper.UInt64ToStr(*info.SubAppId))
		}

		if info.SubAppIdName != nil {
			subAppInfoMap["sub_app_id_name"] = *info.SubAppIdName
		}

		if info.Name != nil {
			subAppInfoMap["name"] = *info.Name
		}

		if info.Description != nil {
			subAppInfoMap["description"] = *info.Description
		}

		if info.CreateTime != nil {
			subAppInfoMap["create_time"] = *info.CreateTime
		}

		if info.Status != nil {
			subAppInfoMap["status"] = *info.Status
		}

		if info.Mode != nil {
			subAppInfoMap["mode"] = *info.Mode
		}

		// Convert storage regions
		if info.StorageRegions != nil {
			regions := make([]string, 0, len(info.StorageRegions))
			for _, region := range info.StorageRegions {
				if region != nil {
					regions = append(regions, *region)
				}
			}
			subAppInfoMap["storage_regions"] = regions
		}

		// Convert tags
		if info.Tags != nil {
			tagsMap := make(map[string]string, len(info.Tags))
			for _, tag := range info.Tags {
				if tag.TagKey != nil && tag.TagValue != nil {
					tagsMap[*tag.TagKey] = *tag.TagValue
				}
			}
			subAppInfoMap["tags"] = tagsMap
		}

		subAppInfoList = append(subAppInfoList, subAppInfoMap)
	}

	// Set data source ID
	d.SetId(helper.DataResourceIdsHash(ids))

	// Set results
	if err := d.Set("sub_application_info_set", subAppInfoList); err != nil {
		log.Printf("[CRITAL]%s provider set sub_application_info_set fail, reason:%s", logId, err.Error())
		return err
	}

	// Export to file if requested
	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), subAppInfoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]", logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
