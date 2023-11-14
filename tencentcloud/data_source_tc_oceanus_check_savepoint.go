/*
Use this data source to query detailed information of oceanus check_savepoint

Example Usage

```hcl
data "tencentcloud_oceanus_check_savepoint" "check_savepoint" {
  job_id = "cql-52xkpymp"
    record_type = 1
  savepoint_path = "cosn://xxx/xxx"
  work_space_id = "space-1327"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudOceanusCheckSavepoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusCheckSavepointRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job id.",
			},

			"serial_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Snapshot resource ID.",
			},

			"record_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Snapshot type. 1:savepoint; 2:checkpoint; 3:cancelWithSavepoint.",
			},

			"savepoint_path": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Snapshot path, currently only supports COS path.",
			},

			"work_space_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Workspace ID.",
			},

			"savepoint_status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "1=available, 2=unavailable.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudOceanusCheckSavepointRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_oceanus_check_savepoint.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("record_type"); v != nil {
		paramMap["RecordType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("savepoint_path"); ok {
		paramMap["SavepointPath"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
	}

	service := OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusCheckSavepointByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		serialId = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(serialId))
	if serialId != nil {
		_ = d.Set("serial_id", serialId)
	}

	if savepointStatus != nil {
		_ = d.Set("savepoint_status", savepointStatus)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
