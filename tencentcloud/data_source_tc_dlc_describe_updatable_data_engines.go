/*
Use this data source to query detailed information of dlc describe_updatable_data_engines

Example Usage

```hcl
data "tencentcloud_dlc_describe_updatable_data_engines" "describe_updatable_data_engines" {
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
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

func dataSourceTencentCloudDlcDescribeUpdatableDataEngines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeUpdatableDataEnginesRead,
		Schema: map[string]*schema.Schema{
			"data_engine_config_command": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine configuration operation command, UpdateSparkSQLLakefsPath updates the managed table path, UpdateSparkSQLResultPath updates the result bucket path.",
			},

			"data_engine_basic_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Engine basic information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_engine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine name.",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Engine state, only support: 0:Init/-1:Failed/-2:Deleted/1:Pause/2:Running/3:ToBeDelete/4:Deleting.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Create time.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Update time.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Returned messages.",
						},
						"data_engine_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine unique id.",
						},
						"data_engine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine type, valid values: PrestoSQL/SparkSQL/SparkBatch.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User unique ID.",
						},
						"user_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User unique uin.",
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

func dataSourceTencentCloudDlcDescribeUpdatableDataEnginesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dlc_describe_updatable_data_engines.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("data_engine_config_command"); ok {
		paramMap["DataEngineConfigCommand"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var dataEngineBasicInfos []*dlc.DataEngineBasicInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeUpdatableDataEnginesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		dataEngineBasicInfos = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dataEngineBasicInfos))
	tmpList := make([]map[string]interface{}, 0, len(dataEngineBasicInfos))

	if dataEngineBasicInfos != nil {
		for _, dataEngineBasicInfo := range dataEngineBasicInfos {
			dataEngineBasicInfoMap := map[string]interface{}{}

			if dataEngineBasicInfo.DataEngineName != nil {
				dataEngineBasicInfoMap["data_engine_name"] = dataEngineBasicInfo.DataEngineName
			}

			if dataEngineBasicInfo.State != nil {
				dataEngineBasicInfoMap["state"] = dataEngineBasicInfo.State
			}

			if dataEngineBasicInfo.CreateTime != nil {
				dataEngineBasicInfoMap["create_time"] = dataEngineBasicInfo.CreateTime
			}

			if dataEngineBasicInfo.UpdateTime != nil {
				dataEngineBasicInfoMap["update_time"] = dataEngineBasicInfo.UpdateTime
			}

			if dataEngineBasicInfo.Message != nil {
				dataEngineBasicInfoMap["message"] = dataEngineBasicInfo.Message
			}

			if dataEngineBasicInfo.DataEngineId != nil {
				dataEngineBasicInfoMap["data_engine_id"] = dataEngineBasicInfo.DataEngineId
			}

			if dataEngineBasicInfo.DataEngineType != nil {
				dataEngineBasicInfoMap["data_engine_type"] = dataEngineBasicInfo.DataEngineType
			}

			if dataEngineBasicInfo.AppId != nil {
				dataEngineBasicInfoMap["app_id"] = dataEngineBasicInfo.AppId
			}

			if dataEngineBasicInfo.UserUin != nil {
				dataEngineBasicInfoMap["user_uin"] = dataEngineBasicInfo.UserUin
			}

			ids = append(ids, *dataEngineBasicInfo.DataEngineId)
			tmpList = append(tmpList, dataEngineBasicInfoMap)
		}

		_ = d.Set("data_engine_basic_infos", tmpList)
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
