/*
Use this data source to query detailed information of oceanus tree_jobs

Example Usage

```hcl
data "tencentcloud_oceanus_tree_jobs" "example" {
  work_space_id = "space-2idq8wbr"
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

func dataSourceTencentCloudOceanusTreeJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusTreeJobsRead,
		Schema: map[string]*schema.Schema{
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudOceanusTreeJobsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_oceanus_tree_jobs.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		//treeJobs    *oceanus.DescribeTreeJobsResponseParams
		workSpaceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
		workSpaceId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, e := service.DescribeOceanusTreeJobsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		//treeJobs = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(workSpaceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
