/*
Use this data source to query detailed information of dlc describe_data_engine_image_versions

Example Usage

```hcl
data "tencentcloud_dlc_describe_data_engine_image_versions" "describe_data_engine_image_versions" {
  engine_type = "SparkBatch"
  }
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDlcDescribeDataEngineImageVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeDataEngineImageVersionsRead,
		Schema: map[string]*schema.Schema{
			"engine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine type only support: SparkSQL/PrestoSQL/SparkBatch.",
			},

			"image_parent_versions": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster large version image information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine major version id.",
						},
						"image_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine major version name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image major version description.",
						},
						"is_public": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is a public version, only support: 1: public;/2: private.",
						},
						"engine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine type only support: SparkSQL/PrestoSQL/SparkBatch.",
						},
						"is_shared_engine": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is shared engine, only support: 1:yes/2:no.",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Version status, only support: 1: initialized/2: online/3: offline.",
						},
						"insert_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
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
	defer logElapsed("data_source.tencentcloud_dlc_describe_data_engine_image_versions.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("engine_type"); ok {
		paramMap["EngineType"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var imageParentVersions []*dlc.DataEngineImageVersion

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDataEngineImageVersionsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
