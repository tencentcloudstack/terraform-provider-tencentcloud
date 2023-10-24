/*
Use this data source to query detailed information of dlc check_data_engine_image_can_be_upgrade

Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_image_can_be_upgrade" "check_data_engine_image_can_be_upgrade" {
  data_engine_id = "DataEngine-cgkvbas6"
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

func dataSourceTencentCloudDlcCheckDataEngineImageCanBeUpgrade() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcCheckDataEngineImageCanBeUpgradeRead,
		Schema: map[string]*schema.Schema{
			"data_engine_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine unique id.",
			},

			"child_image_version_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The latest image version id that can be upgraded.",
			},

			"is_upgrade": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Is it possible to upgrade.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDlcCheckDataEngineImageCanBeUpgradeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dlc_check_data_engine_image_can_be_upgrade.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var dataEngineId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("data_engine_id"); ok {
		dataEngineId = v.(string)
		paramMap["DataEngineId"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}
	var checkResult *dlc.CheckDataEngineImageCanBeUpgradeResponseParams

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcCheckDataEngineImageCanBeUpgradeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		checkResult = result
		return nil
	})
	if err != nil {
		return err
	}
	var data = make(map[string]interface{}, 0)

	if checkResult.ChildImageVersionId != nil {
		_ = d.Set("child_image_version_id", checkResult.ChildImageVersionId)
		data["child_image_version_id"] = checkResult.ChildImageVersionId
	}

	if checkResult.IsUpgrade != nil {
		_ = d.Set("is_upgrade", checkResult.IsUpgrade)
		data["is_upgrade"] = checkResult.IsUpgrade

	}

	d.SetId(dataEngineId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), data); e != nil {
			return e
		}
	}
	return nil
}
