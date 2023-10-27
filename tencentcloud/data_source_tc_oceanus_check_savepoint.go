/*
Use this data source to query detailed information of oceanus check_savepoint

Example Usage

```hcl
data "tencentcloud_oceanus_check_savepoint" "example" {
  job_id         = "cql-314rw6w0"
  serial_id      = "svp-52xkpymp"
  record_type    = 1
  savepoint_path = "cosn://52xkpymp-12345/12345/10000/cql-12345/2/flink-savepoints/savepoint-000000-12334"
  work_space_id  = "space-2idq8wbr"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
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
				Required:    true,
				Type:        schema.TypeString,
				Description: "Snapshot resource ID.",
			},
			"record_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue(RECORD_TYPE),
				Description:  "Snapshot type. 1:savepoint; 2:checkpoint; 3:cancelWithSavepoint.",
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

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		service        = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		checkSavepoint *oceanus.CheckSavepointResponseParams
		serialId       string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("serial_id"); ok {
		paramMap["SerialId"] = helper.String(v.(string))
		serialId = v.(string)
	}

	if v, ok := d.GetOkExists("record_type"); ok {
		paramMap["RecordType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("savepoint_path"); ok {
		paramMap["SavepointPath"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusCheckSavepointByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		checkSavepoint = result
		return nil
	})

	if err != nil {
		return err
	}

	if checkSavepoint.SerialId != nil {
		_ = d.Set("serial_id", checkSavepoint.SerialId)
	}

	if checkSavepoint.SavepointStatus != nil {
		_ = d.Set("savepoint_status", checkSavepoint.SavepointStatus)
	}

	d.SetId(serialId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
