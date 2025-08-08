package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDescribeDataEngineImageVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeDataEngineImageVersionsRead,
		Schema: map[string]*schema.Schema{
			"engine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine type only support: SparkSQL/PrestoSQL/SparkBatch.",
			},

			"sort": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort fields: InsertTime (insert time, default), UpdateTime (update time).",
			},

			"asc": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Sort by: false (descending, default), true (ascending).",
			},

			"image_parent_versions": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Major version of the image information list of clusters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the major version of the image.",
						},
						"image_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the major version of the image.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the major version of the image.",
						},
						"is_public": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is a public version: 1: public version; 2: private version.",
						},
						"engine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster types: SparkSQL, PrestoSQL, and SparkBatch.",
						},
						"is_shared_engine": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Version status. 1: initializing; 2: online; 3: offline.",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Version status. 1: initializing; 2: online; 3: offline.",
						},
						"insert_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Insert time.",
						},
						"update_time": {
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

func dataSourceTencentCloudDlcDescribeDataEngineImageVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_data_engine_image_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("engine_type"); ok {
		paramMap["EngineType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort"); ok {
		paramMap["Sort"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("asc"); ok {
		paramMap["Asc"] = helper.Bool(v.(bool))
	}

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var imageParentVersions []*dlc.DataEngineImageVersion

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDataEngineImageVersionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		imageParentVersions = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(imageParentVersions))
	tmpList := make([]map[string]interface{}, 0, len(imageParentVersions))

	if imageParentVersions != nil {
		for _, dataEngineImageVersion := range imageParentVersions {
			dataEngineImageVersionMap := map[string]interface{}{}

			if dataEngineImageVersion.ImageVersionId != nil {
				dataEngineImageVersionMap["image_version_id"] = dataEngineImageVersion.ImageVersionId
			}

			if dataEngineImageVersion.ImageVersion != nil {
				dataEngineImageVersionMap["image_version"] = dataEngineImageVersion.ImageVersion
			}

			if dataEngineImageVersion.Description != nil {
				dataEngineImageVersionMap["description"] = dataEngineImageVersion.Description
			}

			if dataEngineImageVersion.IsPublic != nil {
				dataEngineImageVersionMap["is_public"] = dataEngineImageVersion.IsPublic
			}

			if dataEngineImageVersion.EngineType != nil {
				dataEngineImageVersionMap["engine_type"] = dataEngineImageVersion.EngineType
			}

			if dataEngineImageVersion.IsSharedEngine != nil {
				dataEngineImageVersionMap["is_shared_engine"] = dataEngineImageVersion.IsSharedEngine
			}

			if dataEngineImageVersion.State != nil {
				dataEngineImageVersionMap["state"] = dataEngineImageVersion.State
			}

			if dataEngineImageVersion.InsertTime != nil {
				dataEngineImageVersionMap["insert_time"] = dataEngineImageVersion.InsertTime
			}

			if dataEngineImageVersion.UpdateTime != nil {
				dataEngineImageVersionMap["update_time"] = dataEngineImageVersion.UpdateTime
			}

			ids = append(ids, *dataEngineImageVersion.ImageVersionId)
			tmpList = append(tmpList, dataEngineImageVersionMap)
		}

		_ = d.Set("image_parent_versions", tmpList)
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
